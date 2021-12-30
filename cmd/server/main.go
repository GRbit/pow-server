package main

import (
	"github.com/rs/zerolog/log"

	"github.com/GRbit/pow-server/inernal/app/server"
)

func main() {
	if err := server.Serve(); err != nil {
		log.Error().Err(err).Msg("server error")
		panic(err)
	}
}
