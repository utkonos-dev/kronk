package stan

type StanCfg struct {
	ClusterID string `mapstructure:"clusterId"`
	ClientID  string `mapstructure:"clientId"`
	NatsURL   string `mapstructure:"natsUrl"`
}
