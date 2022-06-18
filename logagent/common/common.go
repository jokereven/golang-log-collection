package common

// Config 要收集日志配置项的结构体
type Config struct {
	Path  string `json:"path"`
	Topic string `json:"topic"`
}
