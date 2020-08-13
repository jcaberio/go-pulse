package internal

type UpdateResponse struct {
	HasFinished bool `json:"hasFinished"`
	Members     []struct {
		MemberDesc string `json:"memberDesc"`
		MemberID   string `json:"memberId"`
		Messages   []struct {
			Status string `json:"status"`
			Task   string `json:"task"`
		} `json:"messages"`
		Status string `json:"status"`
	} `json:"members"`
	OperationID             string `json:"operationId"`
	OperationStartTimestamp int64  `json:"operationStartTimestamp"`
	OperationType           string `json:"operationType"`
	Rolling                 bool   `json:"rolling"`
	Status                  string `json:"status"`
	User                    string `json:"user"`
}
