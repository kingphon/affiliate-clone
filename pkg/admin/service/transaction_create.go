package service

import (
	"context"

	mgaffiliate "git.selly.red/Selly-Server/affiliate/external/model/mg/affiliate"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/dao"
)

// InsertOne ...
func (s transactionImplement) InsertOne(ctx context.Context, doc mgaffiliate.Transaction) (err error) {
	var d = dao.Transaction()
	err = d.InsertOne(ctx, doc)
	return
}
