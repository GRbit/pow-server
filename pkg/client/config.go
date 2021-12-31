package client

type clientConfig struct {
	Addr          string `long:"addr" default:"server:8080" env:"SERVER_ADDR" description:"server address"`
	MaxComplexity uint64 `long:"max-complexity" default:"4294967296" env:"MAX_COMPLEXITY"`
	NumOfRequests uint64 `long:"requests-num" default:"1024" env:"NUMBER_OF_REQUESTS" description:"how much requests will be sent to server"`

	Console  bool   `long:"console" env:"CONSOLE" description:"extended debug mode; adapts logs output for console"`
	LogLevel string `long:"log-level" default:"info" env:"LOG_LEVEL" description:"set log level (debug|info|warn |error)"`
}
