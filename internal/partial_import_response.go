package internal

type PartialImportPrepareResponse struct {
	ImportID      string         `json:"importId"`
	ListItems     []interface{}  `json:"listItems"`
	Lists         []List         `json:"lists"`
	Models        []interface{}  `json:"models"`
	Plans         []Plans        `json:"plans"`
	RulesProjects []RulesProject `json:"rulesProjects"`
}
