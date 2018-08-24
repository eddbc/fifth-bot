package fifth

import (
	"github.com/robfig/cron"
	"context"
)

var C *cron.Cron

var ctx context.Context

func init() {
	ctx = context.Background()
	C = cron.New()

	//C.AddFunc("@every 10s", func(){getAllContracts()})
	C.AddFunc("0 0 20 * * ?", func(){circlejerk()})

	C.Start()
}

func circlejerk(){
	SendImportantMsg("Daily reminder: Don't forget to circlejerk! <https://www.pandemic-legion.pl/forums/topic/16750-circlejerk/?do=getNewComment>")
}