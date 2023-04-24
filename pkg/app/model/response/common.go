package responsemodel

// ResponseDetail ...
type ResponseDetail struct {
	Data interface{} `json:"data"`
}

// ResponseList ...
type ResponseList struct {
	List          interface{} `json:"list"`
	EndData       bool        `json:"endData"`
	NextPageToken string      `json:"nextPageToken"`
}

// Upsert ...
type Upsert struct {
	ID string `json:"_id"`
}

// LabelValue ...
type LabelValue struct {
	Label string `json:"label"`
	Value string `json:"value"`
}
