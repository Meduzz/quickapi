package rpc

type (
	ReadRequest struct {
		ID string `json:"id"`
	}

	DeleteRequest struct {
		ID string `json:"id"`
	}

	SearchRequest struct {
		Skip   int                          `json:"skip"`
		Take   int                          `json:"take"`
		Where  map[string]string            `json:"where"`
		Scopes map[string]map[string]string `json:"scopes"`
	}

	PatchRequest struct {
		ID   string         `json:"id"`
		Data map[string]any `json:"data"`
	}
)