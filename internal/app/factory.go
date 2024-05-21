package app

import (
	"github.com/mikitabablo/webpages/internal/usecase/webpageparser"
	"github.com/mikitabablo/webpages/pkg/scrapper"
	"github.com/rs/zerolog"
)

type Factory struct {
	log zerolog.Logger

	scrapper *scrapper.Scrapper

	parserUsecase *webpageparser.Usecase
}

func NewFactory(log zerolog.Logger) *Factory {
	return &Factory{
		log: log,
	}
}

func (f *Factory) GetScrapper() *scrapper.Scrapper {
	if f.parserUsecase == nil {
		f.scrapper = scrapper.NewScrapper()
	}

	return f.scrapper
}

func (f *Factory) GetParserUsecase() *webpageparser.Usecase {
	if f.parserUsecase == nil {
		f.parserUsecase = webpageparser.NewUsecase(
			f.GetLogger(),
			f.GetScrapper(),
		)
	}

	return f.parserUsecase
}

func (f *Factory) GetLogger() zerolog.Logger {
	return f.log
}
