package http

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Meduzz/quickapi/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type (
	Spec struct {
		Name   string  `json:"name,omitempty"`   // entity.Name
		Kind   string  `json:"kind,omitempty"`   // quickapi entity.Kind (normal|json)
		Type   string  `json:"type,omitempty"`   // struct.name
		Entity *Entity `json:"entity,omitempty"` // struct.fields
	}

	Entity struct {
		Name   string   `json:"name"`
		Fields []*Field `json:"fields"`
	}

	Field struct {
		Name   string  `json:"name"`
		Type   string  `json:"type"`
		Array  bool    `json:"array,omitempty"` // is array
		Map    bool    `json:"map,omitempty"`   // is map
		Entity *Entity `json:"entity,omitempty"`
	}
)

// For sets up routing for T in the provided router group
// but leaves up to you to deal with the server and run migrations.
func For(db *gorm.DB, e *gin.RouterGroup, entity model.Entity) {
	r := newRouter(db, entity)

	if entity.Name() != "" {
		api := e.Group(fmt.Sprintf("/%s", entity.Name()))

		// setup REST endpoints
		api.POST("/", r.Create)      // create
		api.GET("/:id", r.Read)      // read
		api.PUT("/:id", r.Update)    // update
		api.DELETE("/:id", r.Delete) // delete
		api.GET("/", r.Search)       // list/search
		api.PATCH("/:id", r.Patch)   // patch
		api.GET("/_meta", metaEndpoint(entity))
	} else {
		// setup REST endpoints
		e.POST("/", r.Create)      // create
		e.GET("/:id", r.Read)      // read
		e.PUT("/:id", r.Update)    // update
		e.DELETE("/:id", r.Delete) // delete
		e.GET("/", r.Search)       // list/search
		e.PATCH("/:id", r.Patch)   // patch
		e.GET("/_meta", metaEndpoint(entity))
	}
}

func createScopes(ctx *gin.Context, filters []*model.NamedFilter) []model.Hook {
	if len(filters) == 0 {
		return nil
	}

	scopes := []model.Hook{}

	for _, filter := range filters {
		data, ok := ctx.GetQueryMap(filter.Name)

		if ok {
			scopes = append(scopes, filter.Scope(data))
		}
	}

	return scopes
}

func metaEndpoint(entity model.Entity) func(*gin.Context) {
	data := entity.Create()
	var def any

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Pointer {
		v = v.Elem() // drop pointer
	}

	if v.Kind() == reflect.Struct {
		def = parseStruct(v)
	} else {
		t := v.Type()
		def = t.String()
	}

	return func(ctx *gin.Context) {
		ctx.JSON(200, def)
	}
}

func parseStruct(v reflect.Value) *Entity {
	p := &Entity{}
	p.Name = v.Type().Name()

	fieldCount := v.NumField()

	for i := 0; i < fieldCount; i++ {
		rf := v.Type().Field(i)

		p.Fields = append(p.Fields, parseField(rf))
	}

	return p
}

func parseField(rf reflect.StructField) *Field {
	f := &Field{}

	f.Name = strings.ToLower(rf.Name)
	f.Type = rf.Type.Name()

	raw := rf.Type

	// TODO can we reliably parse tags and use that? validatin:requied|json:optional etc

	if rf.Type.Kind() == reflect.Array || rf.Type.Kind() == reflect.Slice {
		f.Array = true
		raw = raw.Elem()
	}

	if rf.Type.Kind() == reflect.Map {
		f.Map = true
		raw = raw.Elem()
	}

	if raw.Kind() == reflect.Pointer {
		raw = raw.Elem()
	}

	if raw.Kind() == reflect.Struct {
		f.Type = "struct"

		it := reflect.New(raw)

		if it.Kind() == reflect.Pointer {
			it = it.Elem()
		}

		f.Entity = parseStruct(it)
	}

	return f
}
