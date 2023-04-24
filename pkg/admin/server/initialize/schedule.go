package initialize

import (
	"git.selly.red/Selly-Server/affiliate/internal/schedule"
	"git.selly.red/Selly-Server/affiliate/pkg/admin/service"
)

// initSchedule ...
func initSchedule() {
	jobs := []*schedule.Job{
		{
			Spec: "0 */1 * * *",
			Name: "Crawl transaction",
			Cmd:  service.TransactionCrawl().CrawlTransactionData,
		},
	}

	s := schedule.New(jobs...)
	s.Start()
}
