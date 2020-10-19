// Package ini 包实现了读取简单 ini 文件的功能，
// 并且提供了对应的接口函数，用来监听配置文件是否被修改等一些特定的时间
// 让开发者自己决定处理配置变化，最终返回一个新的配置文件数据结构。
package ini

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

// Listener interface used in watch function
//
// Listener 接口使用在 watch 函数中，其中的 listen 函数用于监听特定事件的发生
type Listener interface{ listen(inifile string) }

// ListenFunc a default type to implement Listener interface
//
// ListenFunc 是 Listener 接口的一个默认实现结构，可以将实现的 listen 函数转换成该类型传入 watch 函数
type ListenFunc func(string)

// Watch ini-file if some event happened then it returns a new Config pointer and error
//
// Watch 函数接收一个 Listener 接口，函数会调用其中的 listen 函数，用于监听特定事件发生。
// 当事件发生后，线程从 listen 函数返回，重新执行 Load 函数返回新的 Config 指针或可能产生的错误。
func Watch(filename string, listener Listener) (*Config, error) {
	listener.listen(filename)
	return Load(filename)
}

func (l ListenFunc) listen(filename string) {
	l(filename)
}

// Config represents a INI files in memory
//
// Config 结构将 ini 文件，以特定的形式将文件内容保存在内存中
type Config struct {
	// To keep data in order.
	SecList []string
	// Actual data is stored here.
	Sections map[string]*Section
}

// Load ini file, return error if exists
//
// Load 函数读取 ini 文件，根据文件内容返回 Config 结构指针或可能产生的错误。
func Load(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := bufio.NewReader(f)
	cfg, err := parse(buf)
	return cfg, err
}

func (cfg *Config) init() {
	cfg.Sections = make(map[string]*Section)
}

func (cfg *Config) newSection(name string) error {
	if len(name) == 0 {
		return errors.New("empty section name")
	}
	if cfg.Sections[name] != nil {
		return fmt.Errorf("section(%s) name already exists", name)
	}
	cfg.SecList = append(cfg.SecList, name)
	cfg.Sections[name] = &Section{
		KeyVal: make(map[string]string),
	}
	return nil
}
