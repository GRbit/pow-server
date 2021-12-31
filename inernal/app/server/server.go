package server

import (
	"github.com/GRbit/pow-server/inernal/logger"
	"github.com/GRbit/pow-server/inernal/pow"
	"github.com/VictoriaMetrics/fastcache"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
	"golang.org/x/xerrors"
)

// Serve runs http.ListenAndServe on go-chi router with address specified in Config
func Serve() error {
	var cfg serverConfig

	if _, err := flags.Parse(&cfg); err != nil {
		return xerrors.Errorf("parsing server config: %w", err)
	}

	p, err := initService(&cfg)
	if err != nil {
		return err
	}

	log.Info().Msg("Everything configured. ListenAndServe.")

	return fasthttp.ListenAndServe(cfg.Addr, requestHandler(p))
}

func initService(cfg *serverConfig) (pow.PoW, error) {
	logger.Config(cfg.LogLevel, cfg.Console)

	c := fastcache.New(int(cfg.TaskCacheSize) * 1024 * 1024)

	if cfg.DefaultComplexity < 1 || cfg.DefaultComplexity > 127 {
		return nil, xerrors.Errorf("can't use complexity out of bounds 1-127: %d", cfg.DefaultComplexity)
	}

	return pow.New(c, byte(cfg.DefaultComplexity), cfg.TargetHashing), nil
}
