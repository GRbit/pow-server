package main

import (
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"

	"github.com/GRbit/pow-server/pkg/client"
)

func main() {
	app, err := client.New()
	if err != nil {
		log.Error().Err(err).Msg("creating client error")
		panic(err)
	}

	for {
		time.Sleep(time.Millisecond * 100)
		if err = app.GetWord(); err == nil {
			break
		} else {
			log.Error().Err(err).Msg("waiting till server is running")
		}
	}

	gr := errgroup.Group{}
	for i := uint64(0); i < app.NumOfRequests; i++ {
		gr.Go(app.GetWord)
	}

	if err = gr.Wait(); err != nil {
		panic(err)
	}
}
