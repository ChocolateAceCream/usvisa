package utils

import (
	"math/rand"
	"time"
	"visa/global"

	"go.uber.org/zap"
)

func RandomSleep(min, max int) {
	rand.Seed(time.Now().UnixNano())
	randomSleep := rand.Intn(max-min+1) + min
	duration := time.Duration(randomSleep) * time.Second
	global.LOGGER.Info("Sleeping for %d seconds\n", zap.Int("randomSleep", randomSleep))

	time.Sleep(duration)
}
