package linkanalyzer

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mikitabablo/webpages/internal/domain"
	"github.com/rs/zerolog"
)

type (
	IUseCase interface {
		Analyze(ctx context.Context, baseUrl string, links []string) (domain.Links, error)
	}

	Usecase struct {
		logger zerolog.Logger

		httpClient *http.Client
	}

	AvailabilityResult struct {
		url        string
		accessible bool
		err        error
	}
)

func NewUsecase(logger zerolog.Logger) *Usecase {
	return &Usecase{
		logger: logger,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (u *Usecase) Analyze(ctx context.Context, baseUrl string, links []string) (domain.Links, error) {
	resultsChannel := make(chan AvailabilityResult, len(links))
	for _, link := range links {
		go u.checkLinkAccessibility(link, resultsChannel)
	}

	var analyzedLinks domain.Links

	for range links {
		result := <-resultsChannel
		if result.err != nil {
			u.logger.Err(result.err).Msg("failed to check accessibility")
		}

		switch isLinkInternal(baseUrl, result.url) {
		case true:
			analyzedLinks.Internal.Count++
			analyzedLinks.Internal.Count++
			if !result.accessible {
				analyzedLinks.Inaccessible++
			}

		case false:
			analyzedLinks.External.Count++
			if !result.accessible {
				analyzedLinks.Inaccessible++
			}
		}
	}

	return analyzedLinks, nil
}

func (u *Usecase) checkLinkAccessibility(url string, results chan<- AvailabilityResult) {
	resp, err := u.httpClient.Head(url)
	if err != nil {
		results <- AvailabilityResult{url: url, accessible: false, err: err}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		results <- AvailabilityResult{url: url, accessible: true, err: nil}
	} else {
		results <- AvailabilityResult{url: url, accessible: false, err: fmt.Errorf("received non-2xx status code: %d", resp.StatusCode)}
	}
}

func isLinkInternal(baseURL, hrefURL string) bool {
	base, err := url.Parse(baseURL)
	if err != nil {
		return false
	}

	target, err := url.Parse(hrefURL)
	if err != nil {
		return false
	}

	if base.Scheme != target.Scheme {
		return false
	}

	return getMainDomain(base.Host) == getMainDomain(target.Host)
}

func getMainDomain(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "." + parts[len(parts)-1]
	}

	return host
}
