package internal

type PartialImportSchemasRequest struct {
	Lists         []List         `json:"lists"`
	Models        []interface{}  `json:"models"`
	Plans         []Plans        `json:"plans"`
	RulesProjects []RulesProject `json:"rulesProjects"`
}

type List struct {
	Desc         string `json:"desc"`
	ExistingList bool   `json:"existingList"`
	ID           string `json:"id"`
	MatchingType string `json:"matchingType"`
	TenancyScope string `json:"tenancyScope"`
	Tokenized    bool   `json:"tokenized"`
	Type         string `json:"type"`
}

type Snapshot struct {
	Desc             string            `json:"desc,omitempty"`
	ID               string            `json:"id,omitempty"`
	WorkflowMappings []WorkflowMapping `json:"workflowMapping,omitempty"`
}

type WorkflowMapping struct {
	WorkflowElementId string `json:"workflowElementId,omitempty"`
	WorkflowId        string `json:"workflowId,omitempty"`
}

type RulesProject struct {
	ID              string     `json:"id,omitempty"`
	Desc            string     `json:"desc,omitempty"`
	DestinationDesc string     `json:"destinationDesc,omitempty"`
	DestinationId   string     `json:"destinationId,omitempty"`
	Snapshots       []Snapshot `json:"snapshots,omitempty"`
	Type            string     `json:"type,omitempty"`
}

type Plans struct {
	Desc            string       `json:"desc,omitempty"`
	Executions      []Executions `json:"executions,omitempty"`
	ID              string       `json:"id,omitempty"`
	DestinationId   string       `json:"destinationId,omitempty"`
	DestinationDesc string       `json:"destinationDesc,omitempty"`
}

type Executions struct {
	Desc string `json:"desc,omitempty"`
	ID   string `json:"id"`
}
