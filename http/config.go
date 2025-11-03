package http

import (
	"github.com/gin-gonic/gin"
)

type (
	IDExtractor   func(*gin.Context) string
	MapExtractor  func(*gin.Context) map[string]string
	IntExtractor  func(*gin.Context) int
	BodyExtractor func(any, *gin.Context) (any, error)

	Configurer func(*Config)

	Config struct {
		ID      IDExtractor
		Preload MapExtractor
		Sorting MapExtractor
		Where   MapExtractor
		Skip    IntExtractor
		Take    IntExtractor
		Body    BodyExtractor
	}
)

const (
	ID      = "id"
	PRELOAD = "preload"
	TAKE    = "take"
	SKIP    = "skip"
	WHERE   = "where"
	SORT    = "sort"
)

func DefaultConfig() *Config {
	cfg := &Config{}

	WithPathParamIdStrategy(ID)(cfg)
	WithPreloadQueryMapStrategy(PRELOAD)(cfg)
	WithSortingQueryMapStrategy(SORT)(cfg)
	WithWhereQueryMapStrategy(WHERE)(cfg)
	WithSkipQueryIntStrategy(SKIP, 0)(cfg)
	WithTakeQueryIntStrategy(TAKE, 25)(cfg)
	WithJsonBodyExtractor()(cfg)

	return cfg
}

func WithPathParamIdStrategy(param string) Configurer {
	return func(c *Config) {
		c.ID = ExtractID(param)
	}
}

func WithQueryparamIdStrategy(param string) Configurer {
	return func(c *Config) {
		c.ID = func(ctx *gin.Context) string {
			return ctx.Query(param)
		}
	}
}

func WithPreloadQueryMapStrategy(param string) Configurer {
	return func(c *Config) {
		c.Preload = ExtractQueryMap(param)
	}
}

func WithSortingQueryMapStrategy(param string) Configurer {
	return func(c *Config) {
		c.Sorting = ExtractQueryMap(param)
	}
}

func WithWhereQueryMapStrategy(param string) Configurer {
	return func(c *Config) {
		c.Where = ExtractQueryMap(param)
	}
}

func WithSkipQueryIntStrategy(param string, defaultValue int) Configurer {
	return func(c *Config) {
		c.Skip = ExtractQueryInt(param, defaultValue)
	}
}

func WithTakeQueryIntStrategy(param string, defaultValue int) Configurer {
	return func(c *Config) {
		c.Take = ExtractQueryInt(param, defaultValue)
	}
}

func WithJsonBodyExtractor() Configurer {
	return func(c *Config) {
		c.Body = ExtractBody
	}
}
