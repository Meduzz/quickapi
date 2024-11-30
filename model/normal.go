package model

type (
	defaultEntity struct {
		name            string
		filters         []*NamedFilter
		factory         func() any
		arrayFactory    func() any
		preloadDelegate PreloadDelegate
	}
)

const (
	KindNormal = "normal"
)

// NewEntity creates a new entity, that is the base for the api logic.
func NewEntity[T any](name string, preload PreloadDelegate, filters ...*NamedFilter) Entity {
	if preload == nil {
		preload = func(s string) map[string]*PreloadConfig {
			return nil
		}
	}
	return &defaultEntity{
		name:            name,
		filters:         filters,
		factory:         func() any { return new(T) },
		arrayFactory:    func() any { return make([]*T, 0) },
		preloadDelegate: preload,
	}
}

func (d *defaultEntity) Name() string {
	return d.name
}

func (d *defaultEntity) Filters() []*NamedFilter {
	return d.filters
}

func (d *defaultEntity) Create() any {
	return d.factory()
}

func (d *defaultEntity) CreateArray() any {
	return d.arrayFactory()
}

func (d *defaultEntity) Preload(name string) map[string]*PreloadConfig {
	return d.preloadDelegate(name)
}

func (d *defaultEntity) Kind() string {
	return KindNormal
}
