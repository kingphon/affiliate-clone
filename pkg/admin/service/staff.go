package service

import (
	"encoding/json"
	"fmt"

	"git.selly.red/Selly-Modules/authentication"
	"git.selly.red/Selly-Modules/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	externalauth "git.selly.red/Selly-Server/affiliate/external/auth"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
)

// StaffInterface ...
type StaffInterface interface {
	// GetListStaffInfoByIDs ...
	GetListStaffInfoByIDs(staffIds []primitive.ObjectID) (result []responsemodel.ResponseStaff)
}

// staffImplement ...
type staffImplement struct{}

// Staff ...
func Staff() StaffInterface {
	return staffImplement{}
}

// GetListStaffInfoByIDs ...
func (s staffImplement) GetListStaffInfoByIDs(staffIds []primitive.ObjectID) (result []responsemodel.ResponseStaff) {
	var (
		payload authentication.GetInfoStaff
	)

	payload.Condition = bson.M{"_id": bson.M{"$in": staffIds}}
	response, err := externalauth.RequestGetListStaff(payload)

	if err != nil {
		fmt.Println(err)
	}

	data := make([]responsemodel.StaffRaw, 0)
	if err := json.Unmarshal(response.Data, &data); err != nil {
		logger.Error("authentication.Staff - gte List staff by IDs", logger.LogData{
			Data: bson.M{
				"error":  err.Error(),
				"data: ": response.Data,
			},
		})

		fmt.Println("error: ", err)
		return nil
	}

	// Convert to staff response
	result = make([]responsemodel.ResponseStaff, 0)
	for _, r := range data {
		result = append(result, responsemodel.ResponseStaff{
			ID:   r.ID.Hex(),
			Name: r.Name,
		})
	}

	return
}
