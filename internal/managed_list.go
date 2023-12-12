package internal

type ManagedList struct {
	CollectionSize int        `json:"collectionSize"`
	ListItems      []ListItem `json:"items"`
	LastPage       bool       `json:"lastPage"`
	Offset         int        `json:"offset"`
	Type           string     `json:"type"`
}

type ListItem struct {
	Type           string        `json:"@type"`
	App            string        `json:"app"`
	Color          string        `json:"color"`
	Comment        string        `json:"comment"`
	CreatedAt      int64         `json:"createdAt"`
	CreatedBy      string        `json:"createdBy"`
	Desc           string        `json:"desc"`
	ID             string        `json:"id"`
	ItemValuesType string        `json:"itemValuesType"`
	Items          []interface{} `json:"items"`
	KeyLabel       string        `json:"keyLabel"`
	MatchingType   string        `json:"matchingType"`
	MultiTenant    bool          `json:"multiTenant"`
	Name           string        `json:"name"`
	Ownership      Ownership     `json:"ownership"`
	Tags           []interface{} `json:"tags"`
	TenancyScope   string        `json:"tenancyScope"`
	Tokenized      bool          `json:"tokenized"`
	UpdatedAt      int64         `json:"updatedAt"`
	UpdatedBy      string        `json:"updatedBy"`
}
