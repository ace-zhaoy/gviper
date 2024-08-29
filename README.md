# gviper

`gviper` 是一个基于 [Viper](https://github.com/spf13/viper) 的 Golang 配置管理库，提供了更加便捷的多配置文件加载、环境变量绑定、配置动态更新等功能。`gviper` 支持 YAML、JSON、TOML 等多种配置文件格式，并允许用户通过监听器在配置变化时触发相应的处理逻辑。

## 功能特性

- **支持多个配置文件**：允许同时加载和监听多个配置文件
- **动态配置加载**：支持配置文件的动态加载与更新。
- **配置解析**:结构化绑定配置到Go结构体
- **配置变更监听**：支持在配置文件发生变化时触发自定义回调。
- **通知机制**：在配置加载或更新失败时，发送通知。
- **多种配置文件格式支持**：支持 YAML、JSON、TOML 等格式。
- **环境变量支持**：自动绑定环境变量，支持自定义环境变量前缀。

## 安装

使用 `go get` 命令安装 `gviper`：

```sh
go get github.com/ace-zhaoy/gviper
```

## 使用示例
### 基本使用
```go
import (
    "fmt"
    "github.com/ace-zhaoy/gviper"
)

func main() {
    // 当前目录下的 app.yaml、database.json、log.toml 三个文件
    config := gviper.NewConfig(".", "app", "database.json", "log.toml")
    
    // 加载配置文件
    if err := config.Load(); err != nil {
        panic(err)
    }

    // 读取配置项
    fmt.Println("App Name:", config.GetString("app.name"))
    fmt.Println("Database Port:", config.GetInt("database.port"))
    fmt.Println("Log Level:", config.GetString("log.level"))
}
```

### 配置文件热加载
```go
// 监听配置文件变化（非阻塞）
config.Watch()
```

### 绑定配置到结构体
```go
import (
	"fmt"
	"github.com/ace-zhaoy/gviper"
)

type LogConfig struct {
	Type  string `json:"type"`
	Level string `json:"level"`
}

func main() {
	// 当前目录下的 app.yaml、database.json、log.toml 三个文件
	config := gviper.NewConfig(".", "app", "database.json", "log.toml")

	// 绑定结构体
	var lc LogConfig
	config.Bind("log", &lc)

	// 加载配置文件
	if err := config.Load(); err != nil {
		panic(err)
	}
	
	// 读取配置项
	fmt.Println("Log Level:", lc.Level)
}
```

### 注册配置变更监听器
```go
import (
	"fmt"
	"github.com/ace-zhaoy/gviper"
	"github.com/spf13/viper"
)

func main() {
	// 当前目录下的 app.yaml、database.json、log.toml 三个文件
	config := gviper.NewConfig(".", "app", "database.json", "log.toml")

	config.OnChange("database", func(viper *viper.Viper) error {
		fmt.Println("database config changed!")
		fmt.Println("database init success!")
		return nil
	})

	// 加载配置文件
	if err := config.Load(); err != nil {
		panic(err)
	}

	// 监听配置文件变化
	config.Watch()
}
```

### 注册通知机制
```go

import (
	"fmt"
	"github.com/ace-zhaoy/gviper"
)

type MyNotification struct{}

func (m *MyNotification) Notify(configName string, err error) {
	// 配置热更新失败
	fmt.Printf(" Config %s reload failed: %v\n", configName, err)
}

func main() {
	// 当前目录下的 app.yaml、database.json、log.toml 三个文件
	config := gviper.NewConfig(".", "app", "database.json", "log.toml")

	notification := &MyNotification{}
	config.RegisterNotification(notification)

	if err := config.Load(); err != nil {
		panic(err)
	}

	// 监听配置文件变化
	config.Watch()
}

```

### 环境变量自动绑定
```go
import (
	"fmt"
	"github.com/ace-zhaoy/gviper"
)

func main() {
	// 当前目录下的 app.yaml、database.json、log.toml 三个文件
	config := gviper.NewConfig(".", "app", "database.json", "log.toml")

	// 自动绑定环境变量
	config.AutomaticEnv()

	if err := config.Load(); err != nil {
		panic(err)
	}
	fmt.Println(config.GetString("HOME"))
}
```

### 通过 Option 初始化
```go
import (
	"github.com/ace-zhaoy/gviper"
)

func main() {
	// 当前目录下的 app.ini、database.json、log.toml 三个文件
	config := gviper.NewConfigWithOptions(
		gviper.WithConfigPath("."),
		gviper.WithDefaultConfigType("ini"),
	)
	
	// app.ini、database.json、log.toml
	config.Register("app", "database.json", "log.toml")

	if err := config.Load(); err != nil {
		panic(err)
	}

	// 监听配置文件变化
	config.Watch()
}
```