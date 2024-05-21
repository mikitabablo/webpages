package app

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/mikitabablo/webpages/config"
	httpHandler "github.com/mikitabablo/webpages/internal/handler/http"
	"github.com/rs/zerolog"
)

type App struct {
	server *gin.Engine

	logger zerolog.Logger
}

func NewApp(cfg *config.Config, l zerolog.Logger) *App {
	r := gin.Default()

	factory := NewFactory(l)

	apiRoutes := r.Group("/api")
	httpHandler.InitHandler(
		l,
		apiRoutes,
		factory.GetParserUsecase(),
	)

	// health route
	r.GET("/health", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	return &App{
		server: r,
		logger: l,
	}
}

func (a *App) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := a.server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Err(err).Ctx(ctx).Msg("Server start Failed")
		}
	}()
}
