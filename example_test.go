package gviper_test

import (
	"github.com/ace-zhaoy/gviper"
	"github.com/ace-zhaoy/gviper/notifications"
	"github.com/mitchellh/mapstructure"
	"time"
)

func ExampleNewConfig() {
	config := gviper.NewConfig(".", "app", "database.json", "log.toml")

	// 加载配置文件
	if err := config.Load(); err != nil {
		panic(err)
	}

	// 监听配置文件变化（非阻塞）
	config.Watch()

	_ = config.GetString("app.name")
	_ = config.GetInt("database.port")
	_ = config.GetString("log.level")
}

func ExampleConfig_Register() {
	config := gviper.NewConfig(".")

	config.Register("app", "database.json", "log.toml")

	// 加载配置文件
	if err := config.Load(); err != nil {
		panic(err)
	}

	// 监听配置文件变化（非阻塞）
	config.Watch()

	_ = config.GetString("app.name")
	_ = config.GetInt("database.port")
	_ = config.GetString("log.level")
}

func ExampleConfig_Bind() {
	type LogConfig struct {
		Type  string `json:"type"`
		Level string `json:"level"`
	}

	type AppConfig struct {
		Name string `json:"name"`
	}

	config := gviper.NewConfig(".")

	var lc LogConfig
	var ac AppConfig

	config.Bind("log.toml", &lc)
	config.Bind("app", &ac)

	_ = config.Load()
	config.Watch()

	_ = lc.Level
	_ = ac.Name
}

func ExampleConfig_Bind_register() {
	type LogConfig struct {
		Type  string `json:"type"`
		Level string `json:"level"`
	}

	type AppConfig struct {
		Name string `json:"name"`
	}

	config := gviper.NewConfig(".", "log.toml", "app")

	var lc LogConfig
	var ac AppConfig

	config.Bind("log", &lc)
	config.Bind("app", &ac)

	_ = config.Load()
	config.Watch()

	_ = lc.Level
	_ = ac.Name
}

func ExampleConfig_BindWithTag() {
	type LogConfig struct {
		Type  string `yaml:"type"`
		Level string `yaml:"level"`
	}

	type AppConfig struct {
		Name string `yaml:"name"`
	}

	type DateConfig struct {
		Date time.Time `yaml:"date"`
	}

	config := gviper.NewConfig(".")

	var lc LogConfig
	var ac AppConfig
	var dc DateConfig

	config.BindWithTag("log.toml", &lc, "yaml")
	config.BindWithTag("app", &ac, "yaml")
	config.BindWithTag("date", &dc, "yaml", func(d *mapstructure.DecoderConfig) {
		d.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			d.DecodeHook,
			mapstructure.StringToTimeDurationHookFunc(),
		)
	})

	_ = config.Load()
	config.Watch()

	_ = lc.Level
	_ = ac.Name
}

func ExampleConfig_RegisterNotification() {
	config := gviper.NewConfig(".", "app", "database.json", "log.toml")
	config.RegisterNotification(
		notifications.NewFeishuBotHook("https://open.feishu.cn/open-apis/bot/v2/hook/xxx"),
	)

	_ = config.Load()
	config.Watch()
}

func ExampleConfig_Load() {
	config := gviper.NewConfig(".", "app", "database.json", "log.toml")

	_ = config.Load()
}

func ExampleConfig_Watch() {
	config := gviper.NewConfig(".", "app", "database.json", "log.toml")

	_ = config.Load()
	config.Watch()
}
