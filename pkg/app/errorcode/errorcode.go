package errorcode

import "git.selly.red/Selly-Server/affiliate/external/response"

// Init list error codes
// Code from 200-299
// campaign: 201-299
// platform: 301-399
// transaction: 401-499
func Init() {
	// Init common code first
	response.Init()
	response.AddListCodes(campaign)
	response.AddListCodes(platform)
	response.AddListCodes(transaction)
}
