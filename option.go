package gviper

type Option func(*Config)

func WithConfigPath(configPath string) Option {
	return func(config *Config) {
		config.configPath = configPath
	}
}

func WithDefaultConfigType(configType string) Option {
	return func(config *Config) {
		config.defaultConfigType = configType
	}
}

func WithNotification(notifications ...Notification) Option {
	return func(config *Config) {
		config.notifications = append(config.notifications, notifications...)
	}
}
