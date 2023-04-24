package initialize

import (
	"git.selly.red/Selly-Server/affiliate/pkg/app/database"
)

// mongoDB ...
func mongoDB() {
	database.ConnectMongoDBAffiliate()
}
