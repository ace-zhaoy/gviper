package gviper

import (
	"github.com/ace-zhaoy/errors"
	"github.com/ace-zhaoy/go-utils/uslice"
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"path/filepath"
	"time"
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
		configPath:        ".",
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
			err = cp.viper.Unmarshal(cp.data, func(dc *mapstructure.DecoderConfig) { dc.TagName = cp.tagName })
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
		cp.viper.OnConfigChange(func(_ fsnotify.Event) {
			defer errors.Recover(func(e error) {
				c.notify(cp.configName, e)
			})
			errors.Check(errors.Wrap(cp.viper.ReadInConfig(), "read config [%s] error", cp.configName))
			errors.Check(c.buildChangeFunc(cp)())
		})
		cp.viper.WatchConfig()
	})
}

func (c *Config) Has(key string) bool {
	return c.viper.IsSet(key)
}

func (c *Config) Get(key string) any {
	return c.viper.Get(key)
}

func (c *Config) IsSet(key string) bool {
	return c.viper.IsSet(key)
}

func (c *Config) GetString(key string) string {
	return c.viper.GetString(key)
}

func (c *Config) GetBool(key string) bool {
	return c.viper.GetBool(key)
}

func (c *Config) GetInt(key string) int {
	return c.viper.GetInt(key)
}

func (c *Config) GetInt32(key string) int32 {
	return c.viper.GetInt32(key)
}

func (c *Config) GetInt64(key string) int64 {
	return c.viper.GetInt64(key)
}

func (c *Config) GetUint(key string) uint {
	return c.viper.GetUint(key)
}

func (c *Config) GetUint32(key string) uint32 {
	return c.viper.GetUint32(key)
}

func (c *Config) GetUint64(key string) uint64 {
	return c.viper.GetUint64(key)
}

func (c *Config) GetFloat64(key string) float64 {
	return c.viper.GetFloat64(key)
}

func (c *Config) GetTime(key string) time.Time {
	return c.viper.GetTime(key)
}

func (c *Config) GetDuration(key string) time.Duration {
	return c.viper.GetDuration(key)
}

func (c *Config) GetIntSlice(key string) []int {
	return c.viper.GetIntSlice(key)
}

func (c *Config) GetStringSlice(key string) []string {
	return c.viper.GetStringSlice(key)
}

func (c *Config) GetStringMap(key string) map[string]any {
	return c.viper.GetStringMap(key)
}

func (c *Config) GetStringMapString(key string) map[string]string {
	return c.viper.GetStringMapString(key)
}

func (c *Config) GetStringMapStringSlice(key string) map[string][]string {
	return c.viper.GetStringMapStringSlice(key)
}

func (c *Config) GetSizeInBytes(key string) uint {
	return c.viper.GetSizeInBytes(key)
}

func (c *Config) Sub(key string) *viper.Viper {
	return c.viper.Sub(key)
}

func (c *Config) AllSettings() map[string]any {
	return c.viper.AllSettings()
}

func (c *Config) Default(key string, defaultValue any) any {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.Get(key)
}

func (c *Config) DefaultString(key string, defaultValue string) string {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetString(key)
}

func (c *Config) DefaultBool(key string, defaultValue bool) bool {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetBool(key)
}

func (c *Config) DefaultInt(key string, defaultValue int) int {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetInt(key)
}

func (c *Config) DefaultInt32(key string, defaultValue int32) int32 {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetInt32(key)
}

func (c *Config) DefaultInt64(key string, defaultValue int64) int64 {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetInt64(key)
}

func (c *Config) DefaultUint(key string, defaultValue uint) uint {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetUint(key)
}

func (c *Config) DefaultUint32(key string, defaultValue uint32) uint32 {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetUint32(key)
}

func (c *Config) DefaultUint64(key string, defaultValue uint64) uint64 {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetUint64(key)
}

func (c *Config) DefaultFloat64(key string, defaultValue float64) float64 {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetFloat64(key)
}

func (c *Config) DefaultTime(key string, defaultValue time.Time) time.Time {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetTime(key)
}

func (c *Config) DefaultDuration(key string, defaultValue time.Duration) time.Duration {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetDuration(key)
}

func (c *Config) DefaultIntSlice(key string, defaultValue []int) []int {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetIntSlice(key)
}

func (c *Config) DefaultStringSlice(key string, defaultValue []string) []string {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetStringSlice(key)
}

func (c *Config) DefaultStringMap(key string, defaultValue map[string]any) map[string]any {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetStringMap(key)
}

func (c *Config) DefaultStringMapString(key string, defaultValue map[string]string) map[string]string {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetStringMapString(key)
}

func (c *Config) DefaultStringMapStringSlice(key string, defaultValue map[string][]string) map[string][]string {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetStringMapStringSlice(key)
}

func (c *Config) DefaultSizeInBytes(key string, defaultValue uint) uint {
	if !c.viper.IsSet(key) {
		return defaultValue
	}
	return c.viper.GetSizeInBytes(key)
}
