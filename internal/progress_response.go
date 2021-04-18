package internal

type ProgressResponse struct {
	HasFinished             bool      `json:"hasFinished"`
	Members                 []Members `json:"members"`
	OperationId             string    `json:"operationId"`
	OperationStartTimestamp int64     `json:"operationStartTimestamp"`
	OperationType           string    `json:"operationType"`
	Rolling                 bool      `json:"rolling"`
	Status                  string    `json:"status"`
	User                    string    `json:"user"`
}
type Messages struct {
	Status string `json:"status"`
	Task   string `json:"task"`
}
type Members struct {
	MemberDesc string     `json:"memberDesc"`
	MemberId   string     `json:"memberId"`
	Messages   []Messages `json:"messages"`
	Status     string     `json:"status"`
}
