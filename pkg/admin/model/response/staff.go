package responsemodel

import "go.mongodb.org/mongo-driver/bson/primitive"

// StaffRaw ...
type StaffRaw struct {
	ID   primitive.ObjectID `bson:"_id" json:"_id"`
	Name string             `bson:"name" json:"name"`
}

// ResponseStaff ...
type ResponseStaff struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}
