package gviper

import "github.com/spf13/viper"

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

func WithAutomaticEnv() Option {
	return func(config *Config) {
		config.AutomaticEnv()
	}
}

func WithAllowEmptyEnv(allowEmptyEnv bool) Option {
	return func(config *Config) {
		config.AllowEmptyEnv(allowEmptyEnv)
	}
}

func WithDecoderConfigOptions(options ...viper.DecoderConfigOption) Option {
	return func(config *Config) {
		config.decoderConfigOptions = append(config.decoderConfigOptions, options...)
	}
}
