package config

type Config struct {
	Interval      int
	HomeAssistant HomeAssistantConfig
	Git           GitConfig
	Kubernetes    KubernetesConfig
	Metrics       MetricsConfig
}

type HomeAssistantConfig struct {
	ConfigPath string
	Url        string
	Token      string
}

type GitConfig struct {
	SshKeyPath string
	Token      string
}

type KubernetesConfig struct {
	KubeconfigPath string
	Namespace      string
	Deployment     string
}

type MetricsConfig struct {
	Enabled bool
	Port    int
}
