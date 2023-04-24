package service

import (
	"context"
	"fmt"
	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
	"github.com/logrusorgru/aurora"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"git.selly.red/Selly-Server/affiliate/external/utils/mgquery"
	responsemodel "git.selly.red/Selly-Server/affiliate/pkg/admin/model/response"
)

// ClickInterface ...
type ClickInterface interface {
	// All ...
	All(ctx context.Context, q mgquery.AppQuery) (result responsemodel.ResponseClickAll)

	// Statistic ...
	Statistic(ctx context.Context, q mgquery.Affiliate) (result responsemodel.ResponseClickStatistic)

	// FindByID ...
	FindByID(ctx context.Context, id primitive.ObjectID) mgaffiliate.Click

	// UpdateStatusByID ...
	UpdateStatusByID(ctx context.Context, id primitive.ObjectID, status string) error

	// MigrationSearchString ...
	MigrationSearchString()
}

// clickImplement ...
type clickImplement struct{}

// Click ...
func Click() ClickInterface {
	return &clickImplement{}
}

func (clickImplement) MigrationSearchString() {
	var (
		ctx    = context.Background()
		limit  = 500
		page   = 0 // for logger
		lastID primitive.ObjectID
		q      = mgquery.AppQuery{
			Limit: int64(limit),
			SortInterface: bson.D{
				{"_id", 1},
			},
		}
		cond = bson.M{}
	)
	for {
		if !lastID.IsZero() {
			cond["_id"] = bson.M{
				"$gt": lastID,
			}
		}
		data := dao.Click().FindByCondition(ctx, cond, q.GetFindOptionsWithPage())
		if len(data) == 0 {
			break
		}
		fmt.Println(aurora.Green(fmt.Sprintf("MigrationSearchString page : %d", page)))
		var (
			w []mongo.WriteModel
		)
		for _, click := range data {
			condition := bson.M{
				"_id": click.ID,
			}
			update := bson.M{
				"searchString": click.ID.Hex(),
			}
			w = append(w, mongo.NewUpdateOneModel().SetFilter(condition).SetUpdate(bson.M{
				"$set": update,
			}))
		}
		if len(w) > 0 {
			if err := dao.Click().BulkWrite(ctx, w); err != nil {
				fmt.Println("Bulk write err: ", err)
			}
		}
		lastID = data[len(data)-1].ID
		page++
	}
}
