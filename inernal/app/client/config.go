package client

type clientConfig struct {
	Addr          string `long:"addr" default:"localhost:1444" env:"SERVICE_ADDR" description:"service address"`
	MaxComplexity uint64 `long:"max-complexity" default:"4294967296" env:"MAX_COMPLEXITY"`

	Console  bool   `long:"console" env:"CONSOLE" description:"extended debug mode; adapts logs output for console"`
	LogLevel string `long:"log-level" default:"info" env:"LOG_LEVEL" description:"set log level (debug|info|warn |error)"`
}
