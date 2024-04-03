package config

type Server struct {
	Zap      Zap      `mapstructure:"zap" json:"zap" yaml:"zap"`
	Visa     Visa     `mapstructure:"visa" json:"visa" yaml:"visa"`
	Chromedp Chromedp `mapstructure:"chromedp" json:"chromedp" yaml:"chromedp"`
}
