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
		// Kind tells us what kind of entity we're deling with normal|json
		Kind() string
	}
)
