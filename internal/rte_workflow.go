package internal

type RteWorkflow struct {
	CollectionSize int    `json:"collectionSize"`
	Items          []Item `json:"items"`
	LastPage       bool   `json:"lastPage"`
	Offset         int    `json:"offset"`
	Type           string `json:"type"`
}

type Item struct {
	App               string        `json:"app"`
	Baseinputschemaid string        `json:"baseInputSchemaId"`
	Config            Config        `json:"config"`
	Createdat         int64         `json:"createdAt"`
	Createdby         string        `json:"createdBy"`
	Desc              string        `json:"desc"`
	ID                string        `json:"id"`
	Name              string        `json:"name"`
	Outcomeconfig     Outcomeconfig `json:"outcomeConfig"`
	Ownership         Ownership     `json:"ownership"`
	Updatedat         int64         `json:"updatedAt"`
	Updatedby         string        `json:"updatedBy"`
}
type Actions struct {
	Name   string        `json:"name"`
	Params []interface{} `json:"params"`
}
type Advancedproperties struct {
	NumberEventStorageSplits string `json:"numberEventStorageSplits"`
}
type Filter struct {
	Fields   []string `json:"fields,omitempty"`
	Template string   `json:"template,omitempty"`
}
type Connections struct {
	Filter   *Filter `json:"filter,omitempty"`
	ID       string  `json:"id"`
	SinkId   string  `json:"sinkId"`
	SourceId string  `json:"sourceId"`
}
type Configuration struct {
}
type Inputmapping struct {
}
type Outputmapping struct {
}
type Elements struct {
	Type          string        `json:"@type"`
	Configuration Configuration `json:"configuration"`
	DependencyId  string        `json:"dependencyId"`
	Desc          string        `json:"desc"`
	ID            string        `json:"id"`
	InputMapping  Inputmapping  `json:"inputMapping"`
	Metadata      string        `json:"metadata"`
	OutcomeIds    []interface{} `json:"outcomeIds"`
	OutputMapping Outputmapping `json:"outputMapping"`
}
type Config struct {
	ActionHandlers      []interface{}      `json:"actionHandlers"`
	Actions             []Actions          `json:"actions"`
	AdvancedProperties  Advancedproperties `json:"advancedProperties"`
	Connections         []Connections      `json:"connections"`
	Elements            []Elements         `json:"elements"`
	EventStorageEnabled bool               `json:"eventStorageEnabled"`
	PartitionKeys       []string           `json:"partitionKeys"`
	RecoveryExpression  string             `json:"recoveryExpression"`
}
type Outcomes struct {
	Type              string      `json:"@type"`
	DefaultValue      interface{} `json:"defaultValue"`
	ID                string      `json:"id"`
	Label             string      `json:"label"`
	CategoricalValues []string    `json:"categoricalValues,omitempty"`
}
type Outcomeconfig struct {
	Outcomes []Outcomes `json:"outcomes"`
}
type Groups struct {
	Jrmgtuwhfaznkb7Bvlzfjw int `json:"jrmgtUwhFazNKB7bVLZFJw"`
}
type Ownership struct {
	CurrentUserRights int           `json:"currentUserRights"`
	Groups            Groups        `json:"groups"`
	UserVisibleGroups []interface{} `json:"userVisibleGroups"`
	Visibility        string        `json:"visibility"`
}

type ValidateRestoreState struct {
	RecoveryExpression string `json:"recoveryExpression"`
}
