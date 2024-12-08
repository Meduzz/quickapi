package model

type (
	jsonEntity struct {
		name         string
		factory      func() any
		arrayFactory func() any
		filters      []*NamedFilter
	}
)

func NewJsonEntity[T any](name string, filters ...*NamedFilter) Entity {
	if name == "" {
		panic("Entity must have a name!")
	}

	return &jsonEntity{
		name:         name,
		factory:      func() any { return new(T) },
		arrayFactory: func() any { return make([]*T, 0) },
		filters:      filters,
	}
}

const (
	KindJson = "json"
)

func (j *jsonEntity) Name() string {
	return j.name
}

func (j *jsonEntity) Filters() []*NamedFilter {
	return j.filters
}

func (j *jsonEntity) Create() any {
	return j.factory()
}

func (j *jsonEntity) CreateArray() any {
	return j.arrayFactory()
}

func (j *jsonEntity) Preload(string) map[string]*PreloadConfig {
	return nil
}

func (j *jsonEntity) Kind() string {
	return KindJson
}
