package model

import (
	"gorm.io/gorm"
)

type (
	Hook func(*gorm.DB) *gorm.DB

	// Scope is a func that takes a map (gin QueryMap) and creates a gorm.Scope
	// from it. With the gorm.Scope you get full access to the query and can do
	// custom filters, preloading etc.
	Scope func(map[string]string) Hook

	// NamedFilter allows you to encode buisnesss rules (and a lot of other things)
	// in code and then activate them with query params.
	NamedFilter struct {
		Name  string
		Scope Scope
	}
)

// NewFilter creates a new filter for you. Name is the name of the queryMap. See examples for ... just that.
func NewFilter(name string, handler Scope) *NamedFilter {
	return &NamedFilter{
		Name:  name,
		Scope: handler,
	}
}
