package gviper

import (
	"testing"
)

func TestWithConfigPath(t *testing.T) {
	config := &Config{}
	option := WithConfigPath("/test/path")
	option(config)

	if config.configPath != "/test/path" {
		t.Errorf("Expected config path to be /test/path, got %s", config.configPath)
	}
}

func TestWithDefaultConfigType(t *testing.T) {
	config := &Config{}
	option := WithDefaultConfigType("json")
	option(config)

	if config.defaultConfigType != "json" {
		t.Errorf("Expected config type to be json, got %s", config.defaultConfigType)
	}
}

func TestWithNotification(t *testing.T) {
	config := &Config{}
	mockNotification := &MockNotification{}
	option := WithNotification(mockNotification)
	option(config)

	if len(config.notifications) != 1 {
		t.Errorf("Expected 1 notification, got %d", len(config.notifications))
	}
}

func TestWithAutomaticEnv(t *testing.T) {
	t.Setenv("TEST_ENV", "test")

	config1 := NewConfigWithOptions()
	config2 := NewConfigWithOptions(WithAutomaticEnv())

	if config1.viper.Get("test_env") != nil {
		t.Errorf("Expected test_env to be nil, got %s", config1.viper.Get("test_env"))
	}

	if config2.viper.Get("test_env") != "test" {
		t.Errorf("Expected test_env to be test, got %s", config2.viper.Get("test_env"))
	}

}

func TestWithAllowEmptyEnv(t *testing.T) {
	t.Setenv("TEST_ENV_2", "")

	config1 := NewConfigWithOptions(WithAutomaticEnv())
	config2 := NewConfigWithOptions(
		WithAutomaticEnv(),
		WithAllowEmptyEnv(true),
	)

	if config1.viper.Get("test_env_2") != nil {
		t.Errorf("Expected test_env_2 to be nil, got %s", config1.viper.Get("test_env_2"))
	}

	if config2.viper.Get("test_env_2") != "" {
		t.Errorf("Expected test_env_2 to be '', got %s", config2.viper.Get("test_env_2"))
	}
}
