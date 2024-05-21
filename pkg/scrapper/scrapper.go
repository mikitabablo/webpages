package scrapper

import (
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/gocolly/colly/v2"
)

type (
	IScrapper interface {
		Scrap(url string) (ParseResult, error)
	}

	Scrapper struct {
		scrapper   *colly.Collector
		httpClient *http.Client
	}
)

func NewScrapper() *Scrapper {
	return &Scrapper{
		scrapper:   colly.NewCollector(),
		httpClient: &http.Client{},
	}
}

func (s *Scrapper) Scrap(scrapURL string) (ParseResult, error) {
	c := colly.NewCollector()

	var errResult ErrorResult
	var parseResult ParseResult

	_, err := url.Parse(scrapURL)
	if err != nil {
		return parseResult, ErrParseURLFailed
	}

	c.OnError(func(_ *colly.Response, err error) {
		errResult.Error = err
	})

	c.OnHTML("form", func(e *colly.HTMLElement) {
		action := strings.ToLower(e.Attr("action"))
		hasUsername := false
		hasPassword := false

		e.ForEach("input", func(_ int, el *colly.HTMLElement) {
			inputType := el.Attr("type")
			if slices.Contains([]string{"text", "email", "username", "login"}, inputType) {
				hasUsername = true
			}
			if inputType == "password" {
				hasPassword = true
			}
		})

		if hasPassword && (hasUsername ||
			strings.Contains(action, "signin") ||
			strings.Contains(action, "sign-in") ||
			strings.Contains(action, "login") ||
			strings.Contains(action, "log-in")) {
			parseResult.ContainsLoginForm = true
		}
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		doc := e.DOM.ParentsUntil("html")
		if doc.Length() > 0 {
			docType := strings.ToUpper(doc.Nodes[0].FirstChild.Data)
			switch {
			case strings.Contains(docType, "HTML 5"):
				parseResult.HTMLVersion = "HTML5"
			case strings.Contains(docType, "XHTML"):
				parseResult.HTMLVersion = "XHTML"
			case strings.Contains(docType, "HTML 4.01"):
				parseResult.HTMLVersion = "HTML 4.01"
			default:
				parseResult.HTMLVersion = "Unknown or not a valid HTML document"
			}
		} else {
			parseResult.HTMLVersion = "DOCTYPE not found or unable to determine HTML version"
		}
	})

	c.OnHTML("title", func(e *colly.HTMLElement) {
		parseResult.Title = e.Text
	})

	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
			hrefURL := el.Attr("href")
			if slices.Contains([]string{"", "#"}, hrefURL) {
				return
			}

			// Skip javascript:void(0) links
			if strings.HasPrefix(hrefURL, "javascript:void(0)") {
				return
			}

			hrefURL = e.Request.AbsoluteURL(hrefURL)
			isURLAccessible := s.isLinkAccessible(hrefURL)

			switch s.isLinkInternal(scrapURL, hrefURL) {
			case true:
				parseResult.Links.Internal.Count++
				if isURLAccessible {
					parseResult.Links.Internal.Accessible.Count++
					parseResult.Links.Internal.Accessible.Links = append(
						parseResult.Links.Internal.Accessible.Links,
						hrefURL,
					)
				} else {
					parseResult.Links.Internal.Inaccessible.Count++
					parseResult.Links.Internal.Inaccessible.Links = append(
						parseResult.Links.Internal.Inaccessible.Links,
						hrefURL,
					)
				}

			case false:
				if isURLAccessible {
					parseResult.Links.External.Count++
					parseResult.Links.External.Accessible.Count++
					parseResult.Links.External.Accessible.Links = append(
						parseResult.Links.External.Accessible.Links,
						hrefURL,
					)
				} else {
					parseResult.Links.External.Count++
					parseResult.Links.External.Inaccessible.Count++
					parseResult.Links.External.Inaccessible.Links = append(
						parseResult.Links.External.Inaccessible.Links,
						hrefURL,
					)
				}
			}
		})

		e.ForEach("h1", func(_ int, el *colly.HTMLElement) {
			parseResult.Headings.H1++
		})
		e.ForEach("h2", func(_ int, el *colly.HTMLElement) {
			parseResult.Headings.H2++
		})
		e.ForEach("h3", func(_ int, el *colly.HTMLElement) {
			parseResult.Headings.H3++
		})
		e.ForEach("h4", func(_ int, el *colly.HTMLElement) {
			parseResult.Headings.H4++
		})
		e.ForEach("h5", func(_ int, el *colly.HTMLElement) {
			parseResult.Headings.H5++
		})
		e.ForEach("h6", func(_ int, el *colly.HTMLElement) {
			parseResult.Headings.H6++
		})
	})

	_ = c.Visit(scrapURL)
	if errResult.Error != nil {
		return parseResult, errResult.Error
	}

	return parseResult, nil
}

func (s *Scrapper) isLinkInternal(baseURL, hrefURL string) bool {
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

	return s.getMainDomain(base.Host) == s.getMainDomain(target.Host)
}

func (s *Scrapper) isLinkAccessible(url string) bool {
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return false
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (s *Scrapper) getMainDomain(host string) string {
	parts := strings.Split(host, ".")
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "." + parts[len(parts)-1]
	}

	return host
}
