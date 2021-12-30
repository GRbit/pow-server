package main

import (
	"github.com/GRbit/pow-server/inernal/app/server"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := server.Serve(); err != nil {
		log.Error().Err(err).Msg("server error")
		panic(err)
	}
}
