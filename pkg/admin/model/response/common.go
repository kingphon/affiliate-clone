package responsemodel

// ResponseCreate ...
type ResponseCreate struct {
	ID string `json:"_id"`
}

// ResponseUpdate ...
type ResponseUpdate struct {
	ID string `json:"_id"`
}

// ResponseChangeStatus ...
type ResponseChangeStatus struct {
	ID     string `json:"_id"`
	Status string `json:"status"`
}
