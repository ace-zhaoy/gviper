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
