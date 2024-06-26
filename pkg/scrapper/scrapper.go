package scrapper

import (
	"net/http"
	"net/url"
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

	//var errResult ErrorResult
	var parseResult ParseResult

	_, err := url.Parse(scrapURL)
	if err != nil {
		return parseResult, ErrParseURLFailed
	}

	// Find the DOCTYPE declaration
	c.OnHTML("doctype", func(e *colly.HTMLElement) {
		parseResult.HTMLVersion = strings.TrimSpace(e.Text)
	})

	// Title
	c.OnHTML("title", func(e *colly.HTMLElement) {
		parseResult.Title = e.Text
	})

	// Body Elements
	c.OnHTML("body", func(e *colly.HTMLElement) {
		// Forms
		e.ForEach("form", func(_ int, fe *colly.HTMLElement) {
			var form Form
			form.Action = strings.ToLower(fe.Attr("action"))
			fe.ForEach("input", func(_ int, ie *colly.HTMLElement) {
				form.Inputs = append(form.Inputs, Input{
					Type: ie.Attr("type"),
				})
			})

			parseResult.Forms = append(parseResult.Forms, form)
		})

		// Links
		e.ForEach("a[href]", func(_ int, ae *colly.HTMLElement) {
			hrefURL := ae.Attr("href")
			hrefURL = e.Request.AbsoluteURL(hrefURL)

			if hrefURL == "" {
				hrefURL = e.Request.URL.String()
			}

			// Skip javascript:void(0) links
			if strings.HasPrefix(hrefURL, "javascript:void(0)") {
				return
			}

			parseResult.Links = append(parseResult.Links, e.Request.AbsoluteURL(hrefURL))
		})

		// Headings
		e.ForEach("h1", func(_ int, el *colly.HTMLElement) {
			parseResult.Headings.H1 = append(parseResult.Headings.H1, el.Text)
		})
		e.ForEach("h2", func(_ int, el *colly.HTMLElement) {
			parseResult.Headings.H2 = append(parseResult.Headings.H2, el.Text)
		})
		e.ForEach("h3", func(_ int, el *colly.HTMLElement) {
			parseResult.Headings.H3 = append(parseResult.Headings.H3, el.Text)
		})
		e.ForEach("h4", func(_ int, el *colly.HTMLElement) {
			parseResult.Headings.H4 = append(parseResult.Headings.H4, el.Text)
		})
		e.ForEach("h5", func(_ int, el *colly.HTMLElement) {
			parseResult.Headings.H5 = append(parseResult.Headings.H5, el.Text)
		})
		e.ForEach("h6", func(_ int, el *colly.HTMLElement) {
			parseResult.Headings.H6 = append(parseResult.Headings.H6, el.Text)
		})
	})

	err = c.Visit(scrapURL)
	if err != nil {
		return parseResult, err
	}

	// If DOCTYPE declaration was not found, return a default value
	if parseResult.HTMLVersion == "" {
		parseResult.HTMLVersion = "HTML5"
	}

	return parseResult, nil
}
