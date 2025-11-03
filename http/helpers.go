package http

import (
	"fmt"
	"strconv"

	"github.com/Meduzz/quickapi/model"
	"github.com/gin-gonic/gin"
)

func ExtractID(param string) func(*gin.Context) string {
	return func(ctx *gin.Context) string {
		return ctx.Param(param)
	}
}

func ExtractQueryMap(param string) func(*gin.Context) map[string]string {
	return func(ctx *gin.Context) map[string]string {
		return ctx.QueryMap(param)
	}
}

func ExtractBody(entity any, ctx *gin.Context) (any, error) {
	err := ctx.BindJSON(entity)
	return entity, err
}

func CreateHooks(entity model.Entity, ctx *gin.Context) []model.Hook {
	hooks := make([]model.Hook, 0)
	scopeSupport, ok := entity.(model.ScopeSupport)

	if ok {
		hooks = createScopes(ctx, scopeSupport.Scopes())
	}

	return hooks
}

func ExtractQueryInt(param string, defaultValue int) func(*gin.Context) int {
	return func(ctx *gin.Context) int {
		sSkip := ctx.DefaultQuery(param, fmt.Sprintf("%d", defaultValue))
		iSkip, err := strconv.Atoi(sSkip)

		if err != nil {
			println("parsing query parameter 'skip' threw error", err.Error())
			return defaultValue
		}

		return iSkip
	}
}
