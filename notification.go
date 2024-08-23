package gviper

type Notification interface {
	Notify(configName string, err error)
}
