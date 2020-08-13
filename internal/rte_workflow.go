package internal

type RteWorkflow struct {
	CollectionSize int    `json:"collectionSize"`
	Items          []Item `json:"items"`
	LastPage       bool   `json:"lastPage"`
	Offset         int    `json:"offset"`
	Type           string `json:"type"`
}

type Item struct {
	App               string `json:"app"`
	BaseInputSchemaID string `json:"baseInputSchemaId"`
	Config            Config `json:"config,omitempty"`
	CreatedAt         int64  `json:"createdAt"`
	CreatedBy         string `json:"createdBy"`
	Desc              string `json:"desc"`
	ID                string `json:"id"`
	Name              string `json:"name"`
	OutcomeConfig     struct {
		Outcomes []struct {
			Type              string      `json:"@type"`
			DefaultValue      interface{} `json:"defaultValue"`
			ID                string      `json:"id"`
			Label             string      `json:"label"`
			CategoricalValues []string    `json:"categoricalValues,omitempty"`
		} `json:"outcomes"`
	} `json:"outcomeConfig"`
	Ownership struct {
		CurrentUserRights int `json:"currentUserRights"`
		Groups            struct {
			JrmgtUwhFazNKB7BVLZFJw int `json:"jrmgtUwhFazNKB7bVLZFJw"`
		} `json:"groups"`
		UserVisibleGroups []interface{} `json:"userVisibleGroups"`
		Visibility        string        `json:"visibility"`
	} `json:"ownership"`
	UpdatedAt int64  `json:"updatedAt"`
	UpdatedBy string `json:"updatedBy"`
}

type Config struct {
	ActionHandlers      []interface{}     `json:"actionHandlers"`
	Actions             []Action          `json:"actions"`
	AdvancedProperties  struct{}          `json:"advancedProperties"`
	Connections         []interface{}     `json:"connections"`
	Elements            []WorkflowElement `json:"elements"`
	EventStorageEnabled bool              `json:"eventStorageEnabled"`
	PartitionKeys       []interface{}     `json:"partitionKeys"`
}

type WorkflowElement struct {
	Type          string        `json:"@type"`
	Configuration struct{}      `json:"configuration"`
	DependencyID  string        `json:"dependencyId"`
	Desc          string        `json:"desc"`
	ID            string        `json:"id"`
	InputMapping  struct{}      `json:"inputMapping"`
	Metadata      string        `json:"metadata"`
	OutcomeIds    []interface{} `json:"outcomeIds"`
	OutputMapping struct{}      `json:"outputMapping"`
}

type Action struct {
	Name   string        `json:"name"`
	Params []interface{} `json:"params"`
}
