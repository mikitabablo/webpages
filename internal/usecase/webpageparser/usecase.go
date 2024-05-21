package webpageparser

import (
	"context"

	"github.com/mikitabablo/webpages/pkg/scrapper"
	"github.com/rs/zerolog"
)

type (
	Usecase struct {
		log zerolog.Logger

		scrapper scrapper.IScrapper
	}
)

func NewUsecase(
	log zerolog.Logger,
	scrapper scrapper.IScrapper,
) *Usecase {
	return &Usecase{
		log:      log,
		scrapper: scrapper,
	}
}

func (u *Usecase) Scrape(ctx context.Context, url string) (interface{}, error) {
	res, err := u.scrapper.
		Scrap(url)
	if err != nil {
		u.log.Err(err).Msg("Visit failed")
		return nil, err
	}

	return res, nil
}
