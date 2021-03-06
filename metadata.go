package tileinspect

//go:generate counterfeiter MetadataCmd
type MetadataCmd interface {
	LoadMetadata(target interface{}) error
}

type TileProperty struct {
	Name               string      `json:"name"`
	Type               string      `json:"type"`
	Configurable       bool        `json:"configurable"`
	Default            interface{} `json:"default"`
	Optional           bool        `json:"optional"`
	Options            []Option
	ChildProperties    []TileProperties `json:"option_templates"`
	PropertyBlueprints []TileProperty   `json:"property_blueprints"`
}

type JobType struct {
	Name               string         `json:"name"`
	PropertyBlueprints []TileProperty `json:"property_blueprints"`
}

type TileProperties struct {
	Name               string                 `json:"name"`
	PropertyBlueprints []TileProperty         `json:"property_blueprints"`
	SelectValue        string                 `json:"select_value"`
	StemcellCriteria   map[string]interface{} `json:"stemcell_criteria"`
	JobTypes           []JobType              `json:"job_types"`
}

type Option struct {
	Name  interface{} `json:"name"`
	Label interface{} `json:"label"`
}
