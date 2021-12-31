package server

type serverConfig struct {
	Addr string `long:"addr" default:":8080" env:"SERVICE_ADDR" description:"service address"`

	TargetHashing     int  `short:"t" long:"target-num-hash-goroutines" default:"0" env:"TASK_COMPLEXITY " description:"Sets target number of concurrently run hash functions. If set to 0, will be set to half of NumCPU."`
	DefaultComplexity uint `short:"c" long:"default-complexity" default:"20" env:"TASK_COMPLEXITY" description:"sets default task complexity (leading zeros in hash)"`
	TaskCacheSize     uint `long:"task-cache-size" default:"128" env:"TASK_CACHE_SIZE" description:"sets task cache size in Mb"`

	Console  bool   `long:"console" env:"CONSOLE" description:"extended debug mode; adapts logs output for console"`
	LogLevel string `long:"log-level" default:"info" env:"LOG_LEVEL" description:"set log level (debug|info|warn |error)"`
}
