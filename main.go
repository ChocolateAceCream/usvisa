package main

import (
	"context"
	"time"
	"visa/global"
	"visa/job"
	"visa/utils"

	"go.uber.org/zap"
)

func main() {

	// dir, _ := os.Getwd()
	// fmt.Println(dir)
	global.VIPER = utils.ViperInit("/tmp")
	// global.VIPER = utils.ViperInit(dir)
	// must first load config
	global.LOGGER = utils.LoggerInit()
	global.LOGGER.Info("Init Success")

	count := 0
	for {
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(global.CONFIG.Visa.TimeInterval)*time.Minute)
		go func(ctx context.Context) {
			usVisaJob := job.JobGroupInstance.USVisa
			usVisaJob.Init()
			usVisaJob.Run(ctx)
		}(ctx)
		count++
		time.Sleep(time.Duration(global.CONFIG.Visa.TimeInterval) * time.Minute)
		time.Sleep(5 * time.Second)
		global.LOGGER.Info("--current loop#: ", zap.Int("Count: ", count))
	}
}
