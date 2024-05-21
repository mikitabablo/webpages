package app

import (
	"github.com/mikitabablo/webpages/internal/usecase/linkanalyzer"
	"github.com/mikitabablo/webpages/internal/usecase/webpageanalyzer"
	"github.com/mikitabablo/webpages/pkg/scrapper"
	"github.com/rs/zerolog"
)

type Factory struct {
	log zerolog.Logger

	scrapper *scrapper.Scrapper

	linkAnalyzerUsecase    *linkanalyzer.Usecase
	webpageAnalyzerUsecase *webpageanalyzer.Usecase
}

func NewFactory(log zerolog.Logger) *Factory {
	return &Factory{
		log: log,
	}
}

func (f *Factory) GetScrapper() *scrapper.Scrapper {
	if f.webpageAnalyzerUsecase == nil {
		f.scrapper = scrapper.NewScrapper()
	}

	return f.scrapper
}

func (f *Factory) GetLinkAnalyzerUsecase() *linkanalyzer.Usecase {
	if f.linkAnalyzerUsecase == nil {
		f.linkAnalyzerUsecase = linkanalyzer.NewUsecase(f.GetLogger())
	}

	return f.linkAnalyzerUsecase
}

func (f *Factory) GetWebpageAnalyzerUsecase() *webpageanalyzer.Usecase {
	if f.webpageAnalyzerUsecase == nil {
		f.webpageAnalyzerUsecase = webpageanalyzer.NewUsecase(
			f.GetLogger(),
			f.GetScrapper(),
			f.GetLinkAnalyzerUsecase(),
		)
	}

	return f.webpageAnalyzerUsecase
}

func (f *Factory) GetLogger() zerolog.Logger {
	return f.log
}
