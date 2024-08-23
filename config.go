package gviper

import (
	"github.com/ace-zhaoy/errors"
	"github.com/ace-zhaoy/go-utils/uslice"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"path/filepath"
)

type Listener func(viper *viper.Viper) error

type configParam struct {
	configName string
	configType string
	configFile string
	tagName    string
	viper      *viper.Viper
	onChange   Listener
	data       any
}

type Config struct {
	viper             *viper.Viper
	configPath        string
	defaultConfigType string
	configs           []*configParam
	notifications     []Notification
}

func NewConfig(configPath string, names ...string) *Config {
	c := &Config{
		viper:             viper.New(),
		configPath:        configPath,
		defaultConfigType: "yaml",
	}
	c.Register(names...)
	return c
}

func NewConfigWithOptions(options ...Option) *Config {
	c := &Config{
		viper:             viper.New(),
		defaultConfigType: "yaml",
	}
	for _, option := range options {
		option(c)
	}
	return c
}

func (c *Config) AutomaticEnv() {
	c.viper.AutomaticEnv()
}

func (c *Config) AllowEmptyEnv(allowEmptyEnv bool) {
	c.viper.AllowEmptyEnv(allowEmptyEnv)
}

func (c *Config) parseName(name string) (configName string, configType string, configFile string) {
	fileName := filepath.Base(name)
	fileExt := filepath.Ext(fileName)
	if fileExt == "" {
		return name, c.defaultConfigType, filepath.Join(c.configPath, name+"."+c.defaultConfigType)
	}
	return fileName[:len(fileName)-len(fileExt)], fileExt[1:], filepath.Join(c.configPath, name)
}

func (c *Config) find(configName string) *configParam {
	for _, v := range c.configs {
		if v.configName == configName {
			return v
		}
	}
	return nil
}

func (c *Config) add(configName string, configType string, configFile string) *configParam {
	if cp := c.find(configName); cp != nil {
		return cp
	}
	v := viper.New()
	v.SetConfigName(configName)
	v.SetConfigType(configType)
	v.SetConfigFile(configFile)
	cp := &configParam{
		configName: configName,
		configType: configType,
		configFile: configFile,
		viper:      v,
	}

	c.configs = append(c.configs, cp)

	return cp
}

func (c *Config) resolveConfigParam(name string) *configParam {
	configName, fileType, filePath := c.parseName(name)
	return c.add(configName, fileType, filePath)
}

func (c *Config) Register(names ...string) {
	uslice.ForEach(names, func(name string) { c.resolveConfigParam(name) })
}

func (c *Config) OnChange(name string, listener Listener) {
	cp := c.resolveConfigParam(name)
	cp.onChange = listener
}

func (c *Config) Bind(name string, data any) {
	c.BindWithTag(name, data, "json")
}

func (c *Config) BindWithTag(name string, data any, tagName string) {
	cp := c.resolveConfigParam(name)
	cp.data = data
	cp.tagName = tagName
}

func (c *Config) RegisterNotification(notifications ...Notification) {
	c.notifications = append(c.notifications, notifications...)
}

func (c *Config) notify(configName string, err error) {
	uslice.ForEach(c.notifications, func(n Notification) { n.Notify(configName, err) })
}

func (c *Config) buildChangeFunc(cp *configParam) func() error {
	return func() (err error) {
		defer errors.Recover(func(e error) { err = e })
		c.viper.Set(cp.configName, cp.viper.AllSettings())
		if cp.data != nil {
			err = c.viper.Unmarshal(cp.data, func(dc *mapstructure.DecoderConfig) { dc.TagName = cp.tagName })
			errors.Check(errors.Wrap(err, "unmarshal config [%s] failed", cp.configName))
		}
		if cp.onChange != nil {
			err = cp.onChange(cp.viper)
			errors.Check(errors.Wrap(err, "onchange config [%s] failed", cp.configName))
		}
		return
	}
}

func (c *Config) Load() (err error) {
	defer errors.Recover(func(e error) { err = e })
	uslice.ForEach(c.configs, func(cp *configParam) {
		errors.Check(errors.Wrap(cp.viper.ReadInConfig(), "read config [%s] error", cp.configName))
		errors.Check(c.buildChangeFunc(cp)())
	})
	return nil
}

func (c *Config) Watch() {
	uslice.ForEach(c.configs, func(cp *configParam) {
		c.viper.OnConfigChange(func(_ fsnotify.Event) {
			defer errors.Recover(func(e error) {
				c.notify(cp.configName, e)
			})
			errors.Check(c.buildChangeFunc(cp)())
		})
		c.viper.WatchConfig()
	})
}
