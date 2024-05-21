package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type (
	IWebpageParser interface {
		Scrape(ctx context.Context, url string) (interface{}, error)
	}

	Handler struct {
		log          zerolog.Logger
		routerEngine *gin.RouterGroup

		parserUsecase IWebpageParser
	}
)

func InitHandler(
	log zerolog.Logger,
	router *gin.RouterGroup,
	parserUsecase IWebpageParser,
) {
	handler := &Handler{
		log:           log,
		routerEngine:  router,
		parserUsecase: parserUsecase,
	}

	webpagesGroup := router.Group("webpage")
	webpagesGroup.POST("analyze", handler.Analyze)
}

func (h *Handler) Analyze(c *gin.Context) {
	ctx := c.Request.Context()

	var req AnalyzeRequest
	err := c.BindJSON(&req)
	if err != nil {
		h.log.Err(err).Msg(ErrCannotParseRequest.Error())
		c.JSON(http.StatusBadRequest, ErrResponse{
			Error: ErrCannotParseRequest.Error(),
		})
		return
	}

	res, err := h.parserUsecase.Scrape(ctx, req.URL)
	if err != nil {
		h.log.Err(err).Msg(ErrInvalidURL.Error())
		c.JSON(http.StatusNotFound, ErrResponse{
			Error: ErrInvalidURL.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, res)
}