package server

import (
	"net/http"
	"runtime/debug"

	"github.com/GRbit/pow-server/inernal/logger"
	"github.com/GRbit/pow-server/inernal/pow"
	"github.com/VictoriaMetrics/fastcache"
	"github.com/go-chi/chi/v5"
	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
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

	log.Info().Msg("Creating router...")

	r := chi.NewRouter()
	r.Use(recoverer)
	r.Post("/task", takeTask(p))
	r.Post("/word", receiveWord(p))

	log.Info().Msg("Everything configured. ListenAndServe.")

	return http.ListenAndServe(cfg.Addr, r)
}

func initService(cfg *serverConfig) (pow.PoW, error) {
	logger.Config(cfg.LogLevel, cfg.Console)

	c := fastcache.New(int(cfg.TaskCacheSize) * 1024 * 1024)

	if cfg.DefaultComplexity < 1 || cfg.DefaultComplexity > 127 {
		return nil, xerrors.Errorf("can't use complexity out of bounds 1-127: %d", cfg.DefaultComplexity)
	}

	return pow.New(c, byte(cfg.DefaultComplexity)), nil
}

func recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				log.Error().
					Interface("recover", rvr).
					Str("stacktrace", string(debug.Stack())).
					Msg("panic")
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
