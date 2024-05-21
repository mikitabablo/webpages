package http

import "github.com/mikitabablo/webpages/internal/domain"

type (
	ErrResponse struct {
		Error string `json:"error"`
	}

	AnalyzeResponse struct {
		HTMLVersion       string          `json:"html_version"`
		Title             string          `json:"title"`
		Headings          HeadingsCounter `json:"headings"`
		Links             Links           `json:"links"`
		ContainsLoginForm bool            `json:"contains_login_form"`
	}

	ErrorResult struct {
		Error error `json:"error"`
	}

	HeadingsCounter struct {
		H1 int `json:"h1"`
		H2 int `json:"h2"`
		H3 int `json:"h3"`
		H4 int `json:"h4"`
		H5 int `json:"h5"`
		H6 int `json:"h6"`
	}

	Links struct {
		Internal     LinksInfo `json:"internal"`
		External     LinksInfo `json:"external"`
		Inaccessible int       `json:"inaccessible"`
	}

	LinksInfo struct {
		Count int `json:"count"`
	}
)

func fromDomain(domain *domain.AnalyzeResult) *AnalyzeResponse {
	return &AnalyzeResponse{
		HTMLVersion: domain.HTMLVersion,
		Title:       domain.Title,
		Headings: HeadingsCounter{
			H1: domain.Headings.H1,
			H2: domain.Headings.H2,
			H3: domain.Headings.H3,
			H4: domain.Headings.H4,
			H5: domain.Headings.H5,
			H6: domain.Headings.H6,
		},
		Links: Links{
			Internal:     LinksInfo{Count: domain.Links.Internal.Count},
			External:     LinksInfo{Count: domain.Links.External.Count},
			Inaccessible: domain.Links.Inaccessible,
		},
		ContainsLoginForm: domain.ContainsLoginForm,
	}
}
