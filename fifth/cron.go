package fifth

import (
	"context"
	"github.com/robfig/cron"
)

var c *cron.Cron

var ctx context.Context

func init() {
	ctx = context.Background()
	c = cron.New()

	c.AddFunc("@every 10s", func() { timerCron() })
	c.AddFunc("@every 30s", func() { theraCron() })
	//C.AddFunc("@every 3m", func(){ getAllContracts() })
	//c.AddFunc("0 0 19 * * *", func() { circlejerk() })

	c.Start()
}

func circlejerk() {
	SendImportantMsg("Daily reminder: Don't forget to circlejerk! <https://forums.nog8s.space/topic/3-circlejerk/?do=getNewComment>")
}
