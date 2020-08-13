package internal


type ImportRequest struct {
	ImportID                string `json:"importId"`
	OwnershipGroupsMatching struct {
		ImportID string `json:"importId"`
	} `json:"ownershipGroupsMatching"`
}
