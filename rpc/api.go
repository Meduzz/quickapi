package rpc

type (
	ReadRequest struct {
		ID      string            `json:"id"`
		Preload map[string]string `json:"preload,omitempty"`
	}

	UpdateRequest struct {
		ID   string `json:"id,omitempty"`
		Data any    `json:"data"`
	}

	DeleteRequest struct {
		ID string `json:"id"`
	}

	SearchRequest struct {
		Skip    int                          `json:"skip"`
		Take    int                          `json:"take"`
		Where   map[string]string            `json:"where,omitempty"`
		Sort    map[string]string            `json:"sort,omitempty"`
		Preload map[string]string            `json:"preload,omitempty"`
		Scopes  map[string]map[string]string `json:"scopes,omitempty"`
	}

	PatchRequest struct {
		ID      string            `json:"id"`
		Preload map[string]string `json:"preload,omitempty"`
		Data    map[string]any    `json:"data"`
	}
)
