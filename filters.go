package quickapi

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	// Scope is a func that takes a map (gin QueryMap) and creates a gorm.Scope
	// from it. With the gorm.Scope you get full access to the query and can do
	// custom filters, preloading etc.
	Scope func(map[string]string) func(*gorm.DB) *gorm.DB

	// NamedFilter allows you to encode buisnesss rules (and a lot of other things)
	// in code and then activate them with query params.
	NamedFilter struct {
		Name  string
		Scope Scope
	}
)

// NewFilter creates a new filter for you. Name is the name of the queryMap. See examples for ... just that.
func NewFilter(name string, handler func(map[string]string) func(*gorm.DB) *gorm.DB) *NamedFilter {
	return &NamedFilter{
		Name:  name,
		Scope: handler,
	}
}

func createScopes(ctx *gin.Context, filters []*NamedFilter) []func(*gorm.DB) *gorm.DB {
	if len(filters) == 0 {
		return nil
	}

	scopes := make([]func(*gorm.DB) *gorm.DB, 0)

	for _, filter := range filters {
		data, ok := ctx.GetQueryMap(filter.Name)

		if ok {
			scopes = append(scopes, filter.Scope(data))
		}
	}

	return scopes
}
