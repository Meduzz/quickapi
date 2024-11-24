package model

type (
	PreloadConfig struct {
		Condition string           // sql condition for the preload
		Converter func(string) any // converter from string to "correct" value tyep
	}

	PreloadDelegate func(string) map[string]*PreloadConfig

	// Entity is the base of the api
	Entity interface {
		// Name will be used as a prefix in the path for the api
		Name() string
		// Filters will be used in all parts of the api except create
		Filters() []*NamedFilter
		// Crate creates a new *T (as any)
		Create() any
		// CreateArray created an array of *T so []*T (as any)
		CreateArray() any
		// From a preload alias, return the actual preload data
		Preload(string) map[string]*PreloadConfig
	}

	defaultEntity struct {
		name            string
		filters         []*NamedFilter
		factory         func() any
		arrayFactory    func() any
		preloadDelegate PreloadDelegate
	}
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
