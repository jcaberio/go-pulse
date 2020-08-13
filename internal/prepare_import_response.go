package internal

type PrepareImportResponse struct {
	Errors          []interface{} `json:"errors"`
	ImportID        string        `json:"importId"`
	OwnershipGroups struct {
		GroupsUsedInApp []interface{} `json:"groupsUsedInApp"`
		SourceRootGroup struct {
			Children []struct {
				Children    []interface{} `json:"children"`
				CreatedAt   int64         `json:"createdAt"`
				CreatedBy   string        `json:"createdBy"`
				ForTenant   bool          `json:"forTenant"`
				HierarchyID string        `json:"hierarchyId"`
				ID          string        `json:"id"`
				Name        string        `json:"name"`
				ParentID    string        `json:"parentId"`
				Properties  struct {
				} `json:"properties"`
				Roles     []string `json:"roles"`
				UpdatedAt int64    `json:"updatedAt"`
				UpdatedBy string   `json:"updatedBy"`
			} `json:"children"`
			CreatedAt   int64  `json:"createdAt"`
			CreatedBy   string `json:"createdBy"`
			ForTenant   bool   `json:"forTenant"`
			HierarchyID string `json:"hierarchyId"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			ParentID    string `json:"parentId"`
			Properties  struct {
			} `json:"properties"`
			Roles     []interface{} `json:"roles"`
			UpdatedAt int64         `json:"updatedAt"`
			UpdatedBy string        `json:"updatedBy"`
		} `json:"sourceRootGroup"`
		TargetRootGroup struct {
			Children []struct {
				Children    []interface{} `json:"children"`
				CreatedAt   int64         `json:"createdAt"`
				CreatedBy   string        `json:"createdBy"`
				ForTenant   bool          `json:"forTenant"`
				HierarchyID string        `json:"hierarchyId"`
				ID          string        `json:"id"`
				Name        string        `json:"name"`
				ParentID    string        `json:"parentId"`
				Properties  struct {
				} `json:"properties"`
				Roles     []string `json:"roles"`
				UpdatedAt int64    `json:"updatedAt"`
				UpdatedBy string   `json:"updatedBy"`
			} `json:"children"`
			CreatedAt   int64  `json:"createdAt"`
			CreatedBy   string `json:"createdBy"`
			ForTenant   bool   `json:"forTenant"`
			HierarchyID string `json:"hierarchyId"`
			ID          string `json:"id"`
			Name        string `json:"name"`
			ParentID    string `json:"parentId"`
			Properties  struct {
			} `json:"properties"`
			Roles     []interface{} `json:"roles"`
			UpdatedAt int64         `json:"updatedAt"`
			UpdatedBy string        `json:"updatedBy"`
		} `json:"targetRootGroup"`
	} `json:"ownershipGroups"`
}
