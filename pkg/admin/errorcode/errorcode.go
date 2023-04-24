package errorcode

import "git.selly.red/Selly-Server/affiliate/external/response"

// Init list error codes
// Code from 100 -> 1000
// campaign: 100 - 299
// platform: 300 - 499
// transaction: 500 - 699
// reconciliation: 700 - 899
// seller 900-1000
func Init() {
	// Init common code first
	response.Init()
	response.AddListCodes(audit)
	response.AddListCodes(campaign)
	response.AddListCodes(platform)
	response.AddListCodes(transaction)
	response.AddListCodes(reconciliation)
}
