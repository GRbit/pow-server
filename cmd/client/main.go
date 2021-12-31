package main

import (
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/GRbit/pow-server/inernal/app/client"
)

func main() {
	app, err := client.New()
	if err != nil {
		log.Error().Err(err).Msg("creating client error")
		panic(err)
	}

	gr := errgroup.Group{}
	for i := uint64(0); i < app.NumOfRequests; i++ {
		gr.Go(app.GetWord)
	}

	if err = gr.Wait(); err != nil {
		panic(err)
	}
}
