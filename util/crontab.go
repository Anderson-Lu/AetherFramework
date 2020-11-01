package util

import "github.com/robfig/cron"

var CronSchaduler JobSchduler

type JobSchduler struct {
	jobSchduler *cron.Cron
}

func init() {
	CronSchaduler = JobSchduler{
		jobSchduler: cron.New(),
	}
}

func (self *JobSchduler) AddAndStartJob(job func(), errHandler func(error), taskName string, cronTab string) {
	go func() {
		c := cron.New()
		c.AddFunc(cronTab, func() {
			defer func() {
				if r := recover(); r != nil {
					errHandler(r.(error))
				}
			}()
			job()
		})
		c.Start()
		select {}
	}()
}
