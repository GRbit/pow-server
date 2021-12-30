package server

type serverConfig struct {
	Addr string `long:"addr" default:"localhost:1444" env:"SERVICE_ADDR" description:"service address"`

	DefaultComplexity uint `short:"c" long:"default-complexity" default:"16" env:"TASK_COMPLEXITY" description:"sets default task complexity (leading zeros in hash)"`
	TaskCacheSize     uint `long:"task-cache-size" default:"128" env:"TASK_CACHE_SIZE" description:"sets task cache size in Mb"`

	Console  bool   `long:"console" env:"CONSOLE" description:"extended debug mode; adapts logs output for console"`
	LogLevel string `long:"log-level" default:"info" env:"LOG_LEVEL" description:"set log level (debug|info|warn |error)"`
}
