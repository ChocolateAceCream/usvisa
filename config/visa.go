package config

type Visa struct {
	Email        string `mapstructure:"email" json:"email" yaml:"email"`
	Password     string `mapstructure:"password" json:"password" yaml:"password"`
	TimeInterval int64  `mapstructure:"time-interval" json:"time-interval" yaml:"time-interval"` // log live age
}
