package clientjstream

import (
	"git.selly.red/Selly-Modules/natsio"
	jsmodel "git.selly.red/Selly-Modules/natsio/js/model"
	jssubject "git.selly.red/Selly-Modules/natsio/js/subject"
	"git.selly.red/Selly-Server/affiliate/external/utils/pjson"
)

// UpdateSellerAffiliateStatistic ...
func (s ClientJestStreamPull) UpdateSellerAffiliateStatistic(payload jsmodel.PayloadUpdateSellerAffiliateStatistic) (bool, error) {
	var bytesData = pjson.ToBytes(payload)
	return s.publishWithJetStream(natsio.StreamNameSelly, jssubject.Selly.UpdateSellerAffiliateStatistic, bytesData)
}

// InsertCashflowBySeller ...
func (s ClientJestStreamPull) InsertCashflowBySeller(payload jsmodel.PayloadCashflowsBySeller) (bool, error) {
	var bytesData = pjson.ToBytes(payload)
	return s.publishWithJetStream(natsio.StreamNameSelly, jssubject.Selly.CheckAnDInsertCashflowBySeller, bytesData)
}
