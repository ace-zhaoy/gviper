package gviper

import (
	"github.com/ace-zhaoy/errors"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDefault(t *testing.T) {
	d := t.TempDir()
	t.Logf("tmpdir: %s", d)
	serverConfigFile := filepath.Join(d, "server.yaml")
	err := os.WriteFile(serverConfigFile, []byte("name: gviper\ndate: 2021-11-17T16:25:15+08:00"), 0644)
	if err != nil {
		t.Fatalf("Failed to create server.yaml: %v", err)
	}

	type MyServer struct {
		Name string    `json:"name"`
		Date time.Time `json:"date"`
	}
	var myServer MyServer

	config := Default(d)
	config.Bind("server", &myServer)
	err = config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	td, _ := time.Parse(time.RFC3339, "2021-11-17T16:25:15+08:00")
	assert.Equal(t, MyServer{
		Name: "gviper",
		Date: td,
	}, myServer)
}

func TestNewConfig(t *testing.T) {
	config := NewConfig("/test/path")

	if config.configPath != "/test/path" {
		t.Errorf("Expected config path to be /test/path, got %s", config.configPath)
	}
}

func TestConfig_AutomaticEnv(t *testing.T) {
	config := NewConfig(".")
	config.AutomaticEnv()
	t.Setenv("TEST_ENV", "test")
	t.Setenv("TEST_ENV_2", "")

	if config.viper.Get("test_env") != "test" {
		t.Errorf("Expected test_env to be test, got %s", config.viper.Get("test_env"))
	}
	if config.viper.Get("test_env_2") != nil {
		t.Errorf("Expected test_env_2 to be nil, got %s", config.viper.Get("test_env_2"))
	}

	if config.viper.Get("test_env_3") != nil {
		t.Errorf("Expected test_env_3 to be nil, got %s", config.viper.Get("test_env_3"))
	}
}

func TestConfig_AllowEmptyEnv(t *testing.T) {
	config := NewConfig(".")
	config.AutomaticEnv()
	config.AllowEmptyEnv(true)
	t.Setenv("TEST_ENV", "test")
	t.Setenv("TEST_ENV_2", "")

	if config.viper.Get("test_env") != "test" {
		t.Errorf("Expected test_env to be test, got %s", config.viper.Get("test_env"))
	}
	if config.viper.Get("test_env_2") != "" {
		t.Errorf("Expected test_env_2 to be '', got %s", config.viper.Get("test_env_2"))
	}

	if config.viper.Get("test_env_3") != nil {
		t.Errorf("Expected test_env_3 to be nil, got %s", config.viper.Get("test_env_3"))
	}
}

func TestConfig_Register(t *testing.T) {
	config := NewConfig("/test/path", "config1")
	config.Register("config2")
	config.Register("config3.json")
	wants := []configParam{
		{
			configName: "config1",
			configType: "yaml",
			configFile: "/test/path/config1.yaml",
		},
		{
			configName: "config2",
			configType: "yaml",
			configFile: "/test/path/config2.yaml",
		},
		{
			configName: "config3",
			configType: "json",
			configFile: "/test/path/config3.json",
		},
	}

	if len(config.configs) != len(wants) {
		t.Errorf("Expected %d configs, got %d", len(wants), len(config.configs))
	}

	for i, want := range wants {
		if config.configs[i].configName != want.configName {
			t.Errorf("Expected config name to be %s, got %s", want.configName, config.configs[i].configName)
		}
		if config.configs[i].configType != want.configType {
			t.Errorf("Expected config type to be %s, got %s", want.configType, config.configs[i].configType)
		}
		if config.configs[i].configFile != want.configFile {
			t.Errorf("Expected config file to be %s, got %s", want.configFile, config.configs[i].configFile)
		}
	}
}

func TestDefaultConfigType(t *testing.T) {
	type args struct {
		defaultConfigType string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test json",
			args: args{defaultConfigType: "json"},
			want: "json",
		},
		{
			name: "test ini",
			args: args{defaultConfigType: "ini"},
			want: "ini",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewConfig(".")
			got.defaultConfigType = tt.args.defaultConfigType
			got.Register("config1")
			if got.configs[0].configType != tt.want {
				t.Errorf("Expected config type to be %s, got %s", tt.want, got.configs[0].configType)
			}
		})
	}
}

func TestConfig_RegisterNotification(t *testing.T) {
	config := NewConfig("/test/path", "config.yaml")
	notification1 := &MockNotification{}
	notification2 := &MockNotification{}

	config.RegisterNotification(notification1, notification2)

	if len(config.notifications) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(config.notifications))
	}

	if config.notifications[0] != notification1 {
		t.Error("Expected first notification to be notification1")
	}

	if config.notifications[1] != notification2 {
		t.Error("Expected second notification to be notification2")
	}
}

func TestConfig_Load(t *testing.T) {
	d := t.TempDir()
	t.Logf("tmpdir: %s", d)
	err := os.WriteFile(filepath.Join(d, "server.yaml"), []byte("name: gviper\nenv: test\nhttp:\n  port: 8080"), 0644)
	if err != nil {
		t.Fatalf("Failed to create server.yaml: %v", err)
	}

	err = os.WriteFile(filepath.Join(d, "log.toml"), []byte("type = \"console\"\nlevel = \"ACCESS\""), 0644)
	if err != nil {
		t.Fatalf("Failed to create log.toml: %v", err)
	}

	config := NewConfig(d, "server", "log.toml")
	err = config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want any
	}{
		{
			name: "test1",
			args: args{"server.name"},
			want: "gviper",
		},
		{
			name: "test2",
			args: args{"server.env"},
			want: "test",
		},
		{
			name: "test3",
			args: args{"server.http.port"},
			want: 8080,
		},
		{
			name: "test4",
			args: args{"log.type"},
			want: "console",
		},
		{
			name: "test5",
			args: args{"log.level"},
			want: "ACCESS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if config.Get(tt.args.key) != tt.want {
				t.Errorf("Expected config value to be %v, got %v", tt.want, config.Get(tt.args.key))
			}
		})
	}
}

func TestConfig_Load_error(t *testing.T) {
	config := NewConfig(".", "server", "log.toml")
	err := config.Load()
	assert.Error(t, err)
}

func TestConfig_Watch(t *testing.T) {
	d := t.TempDir()
	t.Logf("tmpdir: %s", d)
	serverConfigFile := filepath.Join(d, "server.yaml")
	err := os.WriteFile(serverConfigFile, []byte("name: gviper"), 0644)
	if err != nil {
		t.Fatalf("Failed to create server.yaml: %v", err)
	}
	logConfigFile := filepath.Join(d, "log.toml")
	err = os.WriteFile(logConfigFile, []byte("type = \"console\"\nlevel = \"ACCESS\""), 0644)
	if err != nil {
		t.Fatalf("Failed to create log.toml: %v", err)
	}

	config := NewConfig(d, "server", "log.toml")
	err = config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	config.Watch()

	assert.Equal(t, "gviper", config.Get("server.name"))
	assert.Equal(t, nil, config.Get("server.env"))
	assert.Equal(t, "ACCESS", config.Get("log.level"))
	assert.Equal(t, nil, config.Get("log.rotate"))

	f1, err := os.OpenFile(serverConfigFile, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("Failed to open server.yaml: %v", err)
	}
	_, err = f1.WriteString("\nenv: test")
	if err != nil {
		t.Fatalf("Failed to write to server.yaml: %v", err)
	}
	_ = f1.Close()

	f2, err := os.OpenFile(logConfigFile, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("Failed to open log.toml: %v", err)
	}
	_, err = f2.WriteString("\nrotate = 1")
	if err != nil {
		t.Fatalf("Failed to write to log.toml: %v", err)
	}
	_ = f2.Close()
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, "gviper", config.Get("server.name"))
	assert.Equal(t, "test", config.Get("server.env"))
	assert.Equal(t, "ACCESS", config.Get("log.level"))
	assert.Equal(t, int64(1), config.Get("log.rotate"))
}

func TestConfig_OnChange(t *testing.T) {
	d := t.TempDir()
	t.Logf("tmpdir: %s", d)
	serverConfigFile := filepath.Join(d, "server.yaml")
	err := os.WriteFile(serverConfigFile, []byte("name: gviper"), 0644)
	if err != nil {
		t.Fatalf("Failed to create server.yaml: %v", err)
	}

	ch := make(chan any, 1)
	config := NewConfig(d)
	config.OnChange("server", func(v *viper.Viper) error {
		ch <- v.Get("env")
		return nil
	})
	err = config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	config.Watch()

	f1, err := os.OpenFile(serverConfigFile, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("Failed to open server.yaml: %v", err)
	}
	_, err = f1.WriteString("\nenv: test")
	if err != nil {
		t.Fatalf("Failed to write to server.yaml: %v", err)
	}
	_ = f1.Close()

	v := <-ch
	t.Logf("load success")
	assert.Equal(t, nil, v)

	v = <-ch
	t.Logf("change success")
	assert.Equal(t, "test", v)
	close(ch)
}

func TestConfig_Bind(t *testing.T) {
	d := t.TempDir()
	t.Logf("tmpdir: %s", d)
	serverConfigFile := filepath.Join(d, "server.yaml")
	err := os.WriteFile(serverConfigFile, []byte("name: gviper"), 0644)
	if err != nil {
		t.Fatalf("Failed to create server.yaml: %v", err)
	}

	type MyServer struct {
		Name string `json:"name"`
		Env  string `json:"env"`
	}
	var myServer MyServer

	config := NewConfig(d)
	config.Bind("server", &myServer)
	err = config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	config.Watch()

	assert.Equal(t, MyServer{
		Name: "gviper",
		Env:  "",
	}, myServer)

	f1, err := os.OpenFile(serverConfigFile, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("Failed to open server.yaml: %v", err)
	}
	_, err = f1.WriteString("\nenv: test")
	if err != nil {
		t.Fatalf("Failed to write to server.yaml: %v", err)
	}
	_ = f1.Close()

	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, MyServer{
		Name: "gviper",
		Env:  "test",
	}, myServer)
}

func TestConfig_BindWithTag(t *testing.T) {
	d := t.TempDir()
	t.Logf("tmpdir: %s", d)
	serverConfigFile := filepath.Join(d, "server.yaml")
	err := os.WriteFile(serverConfigFile, []byte("name: gviper"), 0644)
	if err != nil {
		t.Fatalf("Failed to create server.yaml: %v", err)
	}

	type MyServer struct {
		Name string `yaml:"name"`
		Env  string `yaml:"env"`
	}
	var myServer MyServer

	config := NewConfig(d)
	config.BindWithTag("server", &myServer, "yaml")
	err = config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	config.Watch()

	assert.Equal(t, MyServer{
		Name: "gviper",
		Env:  "",
	}, myServer)

	f1, err := os.OpenFile(serverConfigFile, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("Failed to open server.yaml: %v", err)
	}
	_, err = f1.WriteString("\nenv: test")
	if err != nil {
		t.Fatalf("Failed to write to server.yaml: %v", err)
	}
	_ = f1.Close()

	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, MyServer{
		Name: "gviper",
		Env:  "test",
	}, myServer)
}

type MyNotification struct {
	configName string
	err        error
}

func (m *MyNotification) Notify(configName string, err error) {
	m.configName = configName
	m.err = err
}

func TestConfig_Watch_RegisterNotification(t *testing.T) {
	d := t.TempDir()
	t.Logf("tmpdir: %s", d)
	serverConfigFile := filepath.Join(d, "server.yaml")
	err := os.WriteFile(serverConfigFile, []byte("name: gviper"), 0644)
	if err != nil {
		t.Fatalf("Failed to create server.yaml: %v", err)
	}

	type MyServer struct {
		Name string `json:"name"`
		Env  string `json:"env"`
		Port int64  `json:"port"`
	}
	var myServer MyServer
	myNotification := &MyNotification{}

	config := NewConfig(d)
	config.Bind("server", &myServer)
	config.RegisterNotification(myNotification)
	err = config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	config.Watch()
	assert.Equal(t, "", myNotification.configName)
	assert.Equal(t, nil, myNotification.err)
	assert.Equal(t, false, errors.Is(myNotification.err, errors.NewWithMessage("read config [%s] error", "server")))
	assert.Equal(t, false, errors.Is(myNotification.err, errors.NewWithMessage("unmarshal config [%s] failed", "server")))

	f2, err := os.OpenFile(serverConfigFile, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("Failed to open server.yaml: %v", err)
	}
	_, err = f2.WriteString("\nport: noport")
	if err != nil {
		t.Fatalf("Failed to write to server.yaml: %v", err)
	}
	_ = f2.Close()
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, "server", myNotification.configName)
	assert.Equal(t, true, errors.Is(myNotification.err, errors.NewWithMessage("unmarshal config [%s] failed", "server")))

	myNotification.configName = ""
	myNotification.err = nil

	f1, err := os.OpenFile(serverConfigFile, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("Failed to open server.yaml: %v", err)
	}
	_, err = f1.WriteString("env: test")
	if err != nil {
		t.Fatalf("Failed to write to server.yaml: %v", err)
	}
	_ = f1.Close()
	time.Sleep(10 * time.Millisecond)

	assert.Equal(t, "server", myNotification.configName)
	assert.Equal(t, true, errors.Is(myNotification.err, errors.NewWithMessage("read config [%s] error", "server")))
}

func TestConfig_find(t *testing.T) {
	config := &Config{
		configs: []*configParam{
			{configName: "config1"},
			{configName: "config2"},
		},
	}

	result := config.find("config1")
	assert.NotNil(t, result)
	assert.Equal(t, "config1", result.configName)

	result = config.find("config3")
	assert.Nil(t, result)
}

func TestConfig_add(t *testing.T) {
	config := &Config{
		configs: []*configParam{},
	}

	configName := "test"
	configType := "yaml"
	configFile := "test.yaml"
	cp := config.add(configName, configType, configFile)
	assert.Equal(t, configName, cp.configName)
	assert.Equal(t, configType, cp.configType)
	assert.Equal(t, configFile, cp.configFile)
	assert.Equal(t, 1, len(config.configs))

	cp2 := config.add(configName, configType, configFile)
	assert.Equal(t, cp, cp2)
	assert.Equal(t, 1, len(config.configs))
}

func TestConfig_viper_method(t *testing.T) {
	config := NewConfig(".")
	m := map[string]any{
		"string":               "hello",
		"bool":                 true,
		"int":                  1,
		"int32":                int32(1),
		"int64":                int64(1),
		"uint":                 uint(1),
		"uint32":               uint32(1),
		"uint64":               uint64(1),
		"float64":              float64(1.0),
		"time":                 time.Unix(1725271117, 0),
		"duration":             time.Second,
		"intSlice":             []int{1, 2, 3},
		"stringSlice":          []string{"a", "b", "c"},
		"stringMap":            map[string]any{"a": "b", "c": "d"},
		"stringMapString":      map[string]string{"a": "b", "c": "d"},
		"stringMapStringSlice": map[string][]string{"a": {"b", "c"}, "d": {"e", "f"}},
		"sizeInBytes":          "1kb",
	}
	config.viper.Set("test", m)

	assert.Equal(t, true, config.Has("test.string"))
	assert.Equal(t, false, config.IsSet("test.string11"))
	assert.Equal(t, m["string"], config.GetString("test.string"))
	assert.Equal(t, m["bool"], config.GetBool("test.bool"))
	assert.Equal(t, m["int"], config.GetInt("test.int"))
	assert.Equal(t, m["int32"], config.GetInt32("test.int32"))
	assert.Equal(t, m["int64"], config.GetInt64("test.int64"))
	assert.Equal(t, m["uint"], config.GetUint("test.uint"))
	assert.Equal(t, m["uint32"], config.GetUint32("test.uint32"))
	assert.Equal(t, m["uint64"], config.GetUint64("test.uint64"))
	assert.Equal(t, m["float64"], config.GetFloat64("test.float64"))
	assert.Equal(t, m["time"], config.GetTime("test.time"))
	assert.Equal(t, m["duration"], config.GetDuration("test.duration"))
	assert.Equal(t, m["intSlice"], config.GetIntSlice("test.intSlice"))
	assert.Equal(t, m["stringSlice"], config.GetStringSlice("test.stringSlice"))
	assert.Equal(t, m["stringMap"], config.GetStringMap("test.stringMap"))
	assert.Equal(t, m["stringMapString"], config.GetStringMapString("test.stringMapString"))
	assert.Equal(t, m["stringMapStringSlice"], config.GetStringMapStringSlice("test.stringMapStringSlice"))
	assert.Equal(t, uint(1024), config.GetSizeInBytes("test.sizeInBytes"))
}
