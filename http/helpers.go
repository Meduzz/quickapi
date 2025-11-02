package http

import (
	"fmt"
	"strconv"

	"github.com/Meduzz/quickapi/model"
	"github.com/gin-gonic/gin"
)

func ExtractID(param string, ctx *gin.Context) string {
	return ctx.Param(param)
}

func ExtractQueryMap(param string, ctx *gin.Context) map[string]string {
	return ctx.QueryMap(param)
}

func ExtractBody(entity model.Entity, ctx *gin.Context) (any, error) {
	e := entity.Create()
	err := ctx.BindJSON(e)

	return e, err
}

func CreateHooks(entity model.Entity, ctx *gin.Context) []model.Hook {
	hooks := make([]model.Hook, 0)
	scopeSupport, ok := entity.(model.ScopeSupport)

	if ok {
		hooks = createScopes(ctx, scopeSupport.Scopes())
	}

	return hooks
}

func ExtractQueryInt(param string, defaultValue int, ctx *gin.Context) int {
	sSkip := ctx.DefaultQuery(param, fmt.Sprintf("%d", defaultValue))
	iSkip, err := strconv.Atoi(sSkip)

	if err != nil {
		println("parsing query parameter 'skip' threw error", err.Error())
		return defaultValue
	}

	return iSkip
}
