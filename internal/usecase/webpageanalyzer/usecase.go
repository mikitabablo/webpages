package webpageanalyzer

import (
	"context"
	"slices"
	"strings"

	"github.com/mikitabablo/webpages/internal/domain"
	"github.com/mikitabablo/webpages/internal/usecase/linkanalyzer"
	"github.com/mikitabablo/webpages/pkg/scrapper"
	"github.com/rs/zerolog"
)

var (
	loginFormTypeNames    = []string{"text", "email", "username", "login"}
	passwordFormTypeNames = []string{"password"}
)

type (
	Usecase struct {
		log zerolog.Logger

		scrapper scrapper.IScrapper

		linkAnalyzer linkanalyzer.IUseCase
	}
)

func NewUsecase(
	log zerolog.Logger,
	scrapper scrapper.IScrapper,
	linkAnalyzer linkanalyzer.IUseCase,
) *Usecase {
	return &Usecase{
		log:          log,
		scrapper:     scrapper,
		linkAnalyzer: linkAnalyzer,
	}
}

func (u *Usecase) AnalyzeWebpage(ctx context.Context, url string) (interface{}, error) {
	res, err := u.scrapper.Scrap(url)
	if err != nil {
		u.log.Err(err).Msg("Visit failed")
		return nil, err
	}

	links, err := u.linkAnalyzer.Analyze(ctx, url, res.Links)
	if err != nil {
		u.log.Err(err).Msg("Links analyze failed")
		return nil, err
	}

	analyzeRes := domain.AnalyzeResult{
		HTMLVersion: res.HTMLVersion,
		Title:       res.Title,
		Headings: domain.HeadingsCounter{
			H1: len(res.Headings.H1),
			H2: len(res.Headings.H2),
			H3: len(res.Headings.H3),
			H4: len(res.Headings.H4),
			H5: len(res.Headings.H5),
			H6: len(res.Headings.H6),
		},
		Links:             links,
		ContainsLoginForm: u.isLoginFormExistsInScrapResults(res.Forms),
	}

	return analyzeRes, nil
}

func (u *Usecase) isLoginFormExistsInScrapResults(forms []scrapper.Form) bool {
	var hasUsername, hasPassword bool
	var action string

	for _, form := range forms {
		action = strings.ToLower(form.Action)
		hasUsername = false
		hasPassword = false

		for _, input := range form.Inputs {
			if slices.Contains(loginFormTypeNames, input.Type) {
				hasUsername = true
			}

			if slices.Contains(passwordFormTypeNames, input.Type) {
				hasPassword = true
			}
		}

		if hasPassword && (hasUsername || u.isLoginAction(action)) {
			return true
		}
	}

	return false
}

func (u *Usecase) isLoginAction(action string) bool {
	return strings.Contains(action, "signin") ||
		strings.Contains(action, "sign-in") ||
		strings.Contains(action, "login") ||
		strings.Contains(action, "log-in")
}
