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
    Numbereventstoragesplits string `json:"numberEventStorageSplits"`
}
type Filter struct {
	Fields   []string `json:"fields"`
	Template string   `json:"template"`
}
type Connections struct {
	Filter   Filter `json:"filter,omitempty"`
	ID       string `json:"id"`
	Sinkid   string `json:"sinkId"`
	Sourceid string `json:"sourceId"`
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
	Dependencyid  string        `json:"dependencyId"`
	Desc          string        `json:"desc"`
	ID            string        `json:"id"`
	Inputmapping  Inputmapping  `json:"inputMapping"`
	Metadata      string        `json:"metadata"`
	Outcomeids    []interface{} `json:"outcomeIds"`
	Outputmapping Outputmapping `json:"outputMapping"`
}
type Config struct {
	Actionhandlers      []interface{}      `json:"actionHandlers"`
	Actions             []Actions          `json:"actions"`
	Advancedproperties  Advancedproperties `json:"advancedProperties"`
	Connections         []Connections      `json:"connections"`
	Elements            []Elements         `json:"elements"`
	Eventstorageenabled bool               `json:"eventStorageEnabled"`
	Partitionkeys       []string           `json:"partitionKeys"`
}
type Outcomes struct {
	Type              string   `json:"@type"`
	Defaultvalue      bool     `json:"defaultValue"`
	ID                string   `json:"id"`
	Label             string   `json:"label"`
	Categoricalvalues []string `json:"categoricalValues,omitempty"`
}
type Outcomeconfig struct {
	Outcomes []Outcomes `json:"outcomes"`
}
type Groups struct {
	Jrmgtuwhfaznkb7Bvlzfjw int `json:"jrmgtUwhFazNKB7bVLZFJw"`
}
type Ownership struct {
	Currentuserrights int           `json:"currentUserRights"`
	Groups            Groups        `json:"groups"`
	Uservisiblegroups []interface{} `json:"userVisibleGroups"`
	Visibility        string        `json:"visibility"`
}