package tileinspect

type ConfigFile struct {
	ProductProperties map[string]*ConfigFileProperty `json:"product-properties"`
}

type ConfigFileProperty struct {
	Type     string      `json:"type"`
	Value    interface{} `json:"value"`
	Required *bool       `json:"required,omitempty"`
}
