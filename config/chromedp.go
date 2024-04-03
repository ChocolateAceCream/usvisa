package config

type Chromedp struct {
	Headless                bool `mapstructure:"headless" json:"headless" yaml:"headless"`
	IgnoreCertificateErrors bool `mapstructure:"ignoreCertificateErrors" json:"ignoreCertificateErrors" yaml:"ignore-certificate-errors"`
}
