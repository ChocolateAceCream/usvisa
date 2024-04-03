package global

import (
	"visa/config"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	CONFIG config.Server
	LOGGER *zap.Logger
	VIPER  *viper.Viper
)
