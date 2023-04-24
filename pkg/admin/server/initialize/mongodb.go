package initialize

import "git.selly.red/Selly-Server/affiliate/pkg/admin/database"

// mongoDB ...
func mongoDB() {
	database.ConnectMongoDBAffiliate()
}
