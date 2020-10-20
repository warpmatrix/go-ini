<!-- omit in toc -->
# ini

ini 包实现了读取简单 ini 文件的功能，并且提供了对应的接口函数，用来监听配置文件是否被修改，让开发者自己决定处理配置变化，最终返回一个新的配置文件数据结构。

<!-- omit in toc -->
## Table of Contents

- [Getting Started](#getting-started)
  - [`Load` 函数的使用](#load-函数的使用)
  - [`Watch` 函数的使用](#watch-函数的使用)
- [使用示例](#使用示例)
  - [`Load` 函数使用示例](#load-函数使用示例)
  - [`Watch` 函数使用示例](#watch-函数使用示例)

## Getting Started

使用指令 `go get gitee.com/warpmatrix/go-ini` 或指令 `go get github.com/warpmatrix/go-ini` 可以安装该 `ini` 包。

该包的主要用法可以参见 [api 文档](godoc.html)，下面提供简单说明主要函数的用法。并在最后提供相关函数的的使用示例。

### `Load` 函数的使用

`Load` 函数可以将 ini 文件内容，填充到 `Config` 数据结构中。通过 `Sections[secName]` 的形式可以访问 `Config` 变量中对应节的内容。每一个 `Section` 变量也可以通过 `KeyVal[keyName]` 的形式访问对应键的值。

### `Watch` 函数的使用

`Watch` 函数接收一个 `listener` 接口，用于侦听特定事件的发生，如：配置文件是否有被修改。当事件发生后，`Watch` 函数将重新读取配置文件返回相应的 `Config` 变量。

> 需要注意的是，`Watch` 函数使用的是阻塞的方式实现事件的侦听。要用非阻塞的方式需要调用者自行使用 go 程并行执行。

关于 `listener` 接口，包内提供了一种默认的实现方式 `ListenFunc`。可以将编写好的函数（该函数的功能为阻塞线程，直到指定事件发生后直接返回）赋值给该类型的变量，传递给 `Watch` 函数即可。一个简单的代码片段如下：

```go
var listener ini.ListenFunc = fileChange
cfg, err = ini.Watch("my.ini", listener)
```

## 使用示例

### `Load` 函数使用示例

可以执行以下指令建立实验环境：

```bash
mkdir -p /tmp/ini
cd /tmp/ini
touch my.ini main.go
```

在 `my.ini` 文件中填充以下内容：

```ini
# possible values : production, development
app_mode = development

[paths]
# Path to where grafana can store temp files, sessions, and the sqlite3 db (if that is used)
data = /home/git/grafana

[server]
# Protocol (http or https)
protocol = http

# The http port  to use
http_port = 9999

# Redirect to correct domain if host header does not match domain
# Prevents DNS rebinding attacks
enforce_domain = true
```

在 `main.go` 文件中填充以下内容：

```go
package main

import (
    "fmt"

    // ini "domain-name/user/repo"
    ini "github.com/warpmatrix/go-ini"
)

func main() {
    cfg, err := ini.Load("my.ini")
    if err != nil {
        fmt.Println(err)
        return
    }
    printCfg(cfg)
}

func printCfg(cfg *ini.Config) {
    for _, secName := range cfg.SecList {
        sec := cfg.Sections[secName]
        fmt.Printf("[%s]\n", secName)
        for _, key := range sec.KeyList {
            fmt.Printf("%s = %s\n", key, sec.KeyVal[key])
        }
    }
}
```

使用指令 `go run main.go`，运行上述程序得到以下输出：

```plaintext
[DEFAULT]
app_mode = development
[paths]
data = /home/git/grafana
[server]
protocol = http
http_port = 9999
enforce_domain = true
```

### `Watch` 函数使用示例

以上面的程序作为基础，我们可以进一步改动使用 `Watch` 函数。

将 `main.go` 文件改为以下内容，该程序可以实现侦听文件修改的事件，将变化前后的配置文件内容输出到屏幕上：

```go
package main

import (
    "fmt"
    "log"

    "github.com/fsnotify/fsnotify"
    // ini "domain-name/user/repo"
    ini "github.com/warpmatrix/go-ini"
)

func main() {
    cfg, err := ini.Load("my.ini")
    if err != nil {
        fmt.Println(err)
        return
    }
    printCfg(cfg)

    var listener ini.ListenFunc = fileChange
    cfg, err = ini.Watch("my.ini", listener)
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("")
    fmt.Println("after changed:")
    printCfg(cfg)
}

func fileChange(filename string) {
    watcher, err := fsnotify.NewWatcher()
    if err != nil {
        log.Fatal(err)
    }
    defer watcher.Close()
    done := make(chan bool)
    go func() {
        for {
            select {
            case event, ok := <-watcher.Events:
                if !ok {
                    return
                }
                if event.Op&fsnotify.Write == fsnotify.Write {
                    done <- true
                    return
                }
            case err, ok := <-watcher.Errors:
                if !ok {
                    return
                }
                log.Println("error:", err)
            }
        }
    }()
    watcher.Add(filename)
    <-done
}

func printCfg(cfg *ini.Config) {
    for _, secName := range cfg.SecList {
        sec := cfg.Sections[secName]
        fmt.Printf("[%s]\n", secName)
        for _, key := range sec.KeyList {
            fmt.Printf("%s = %s\n", key, sec.KeyVal[key])
        }
    }
}
```

例如将 `server` 一节注释掉，可以得到以下输出：

```plaintext
[DEFAULT]
app_mode = development
[paths]
data = /home/git/grafana
[server]
protocol = http
http_port = 9999
enforce_domain = true

after changed:
[DEFAULT]
app_mode = development
[paths]
data = /home/git/grafana
```
