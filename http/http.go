package http

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Meduzz/helper/fp/result"
	"github.com/Meduzz/helper/fp/slice"
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
		Type   string  `json:"type"`            // field.Type | struct
		Array  bool    `json:"array,omitempty"` // is array
		Map    bool    `json:"map,omitempty"`   // is map
		Entity *Entity `json:"entity,omitempty"`
	}

	Discovery struct {
		Entities []string `json:"entities"`
	}
)

// For sets up routing for T in the provided router group
// but leaves up to you to deal with the server and run migrations.
func For(db *gorm.DB, e *gin.RouterGroup, entities ...model.Entity) error {
	discovery := &Discovery{}

	listOfMaybeNames := slice.Map(entities, func(entity model.Entity) *result.Operation[string] {
		aRouter := result.Execute(newRouter(db, entity))

		return result.Map(aRouter, func(r *router) string {
			api := e.Group(fmt.Sprintf("/%s", entity.Name()))

			// setup REST endpoints
			api.POST("/", r.Create)                          // create
			api.GET("/:id", r.Read)                          // read
			api.PUT("/:id", r.Update)                        // update
			api.DELETE("/:id", r.Delete)                     // delete
			api.GET("/", r.Search)                           // list/search
			api.PATCH("/:id", r.Patch)                       // patch
			api.GET("/_meta", serveMeta(entityMeta(entity))) // TODO make this opt-in too?

			return entity.Name()
		})
	})

	agg := result.Execute(make([]string, 0), nil)

	maybeListOfNames := slice.Fold(listOfMaybeNames, agg, func(op *result.Operation[string], agg *result.Operation[[]string]) *result.Operation[[]string] {
		return result.Then(agg, func(n []string) ([]string, error) {
			v, err := op.Get()

			if err != nil {
				return nil, err
			}

			return append(n, v), nil
		})
	})

	names, err := maybeListOfNames.Get()

	if err != nil {
		return err
	}

	discovery.Entities = names

	e.GET("/_discover", func(ctx *gin.Context) {
		ctx.JSON(200, discovery)
	})

	return nil
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

func serveMeta(entity any) func(*gin.Context) {
	return func(ctx *gin.Context) {
		ctx.JSON(200, entity)
	}
}

func entityMeta(entity model.Entity) any {
	data := entity.Create()
	var def any

	v := reflect.ValueOf(data)

	if v.Kind() == reflect.Pointer {
		v = v.Elem() // drop pointer
	}

	if v.Kind() == reflect.Struct {

		if entity.Kind() == model.JsonKind {
			root := &Entity{}
			root.Name = "jsonTable"
			root.Fields = append(root.Fields, generateSimpleField("id", "int64"), generateSimpleField("created", "int64"), generateSimpleField("udpated", "int64"))

			dataField := &Field{}
			dataField.Name = "data"
			dataField.Entity = parseStruct(v)
			dataField.Array = false
			dataField.Map = false
			dataField.Type = "struct"

			root.Fields = append(root.Fields, dataField)

			def = root
		} else {
			def = parseStruct(v)
		}
	} else {
		t := v.Type()
		def = t.String()
	}

	return def
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

	tag := rf.Tag
	jsonData, ok := tag.Lookup("json")

	if ok {
		tagSplit := strings.Split(jsonData, ",")
		if len(tagSplit) > 0 {
			if tagSplit[0] != "" {
				f.Name = strings.TrimSpace(tagSplit[0])
			}
		}
	}

	raw := rf.Type

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

func generateSimpleField(name, typ string) *Field {
	f := &Field{}

	f.Name = name
	f.Type = typ
	f.Array = false
	f.Map = false

	return f
}
