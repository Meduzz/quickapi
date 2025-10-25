package model

type (
	PreloadConfig struct {
		Condition string           // sql condition for the preload
		Converter func(string) any // converter from string to "correct" value type
	}

	EntityKind string

	// Entity is the base of the api
	Entity interface {
		// Name will be used as a prefix in the path for the api
		Name() string
		// Crate creates a new *T (as any)
		Create() any
		// CreateArray created an array of *T so []*T (as any)
		CreateArray() any
	}

	// PreloadSupport allows you to preload a defined collection with optional conditions
	PreloadSupport interface {
		// From a preload alias, return the actual preload data
		Preload(string) map[string]*PreloadConfig
	}

	ScopeSupport interface {
		Scopes() []*NamedFilter
	}
)
