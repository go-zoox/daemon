package model

type Config = MasterConfig

func Run(cfg *Config) error {
	return Master(cfg)
}
