package gviper

import (
	"testing"
)

type MockNotification struct {
	notified bool
}

func (m *MockNotification) Notify(configName string, err error) {
	m.notified = true
}

func TestNotification(t *testing.T) {
	mockNotification := &MockNotification{}
	mockNotification.Notify("test-config", nil)

	if !mockNotification.notified {
		t.Error("Expected notification to be sent, but it was not")
	}
}
