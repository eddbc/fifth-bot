package fifth

import (
	"context"
	"github.com/robfig/cron"
)

var C *cron.Cron

var ctx context.Context

func init() {
	ctx = context.Background()
	C = cron.New()

	C.AddFunc("@every 10s", func() { timerCron() })
	//C.AddFunc("@every 3m", func(){ getAllContracts() })
	C.AddFunc("0 0 19 * * *", func() { circlejerk() })

	C.Start()
}

func circlejerk() {
	SendImportantMsg("Daily reminder: Don't forget to circlejerk! <https://forums.nog8s.space/topic/3-circlejerk/?do=getNewComment>")
}
