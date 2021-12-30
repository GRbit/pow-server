package server

import (
	"net/http"
	"strconv"

	"github.com/GRbit/pow-server/inernal/pow"
	"github.com/GRbit/pow-server/inernal/quotes"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
)

func takeTask(p pow.PoW) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)

		key, c, err := p.CreateTask()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"can't create task'"}`))
			log.Error().Err(err).Msg("creating task")
			return
		}

		t := `{"key":"` + key + `","complexity":` + strconv.Itoa(c) + `}`
		if _, err := w.Write([]byte(t)); err != nil {
			log.Error().Err(err).Msg("HealthCheck: body write err")
		}
	}
}

type solvedTask struct {
	Key   string
	Nonce uint64
}

func receiveWord(p pow.PoW) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var t solvedTask

		if err := jsoniter.NewDecoder(r.Body).Decode(&t); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"decoding request"}`))
			log.Error().Err(err).Msg("decoding receive word request")
			return
		}

		if err := p.ValidateTask(t.Key, t.Nonce); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"error":"incorrect task solution"}`))
			log.Error().Err(err).Msg("validating task")
			return
		}

		resp, err := jsoniter.Marshal(quotes.RandomQuote())
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"encoding word of wisdom"}`))
			log.Error().Err(err).Msg("encoding word of wisdom")
			return
		}

		w.WriteHeader(http.StatusOK)
		if _, err = w.Write(resp); err != nil {
			log.Error().Err(err).Msg("writing word of wisdom")
		}
	}
}
