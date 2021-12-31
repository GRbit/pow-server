package server

import (
	"net/http"
	"runtime/debug"
	"strconv"

	"github.com/GRbit/pow-server/inernal/pow"
	"github.com/GRbit/pow-server/inernal/quotes"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"github.com/valyala/fasthttp"
)

const jsonContentType = "application/json"

func requestHandler(p pow.PoW) func(ctx *fasthttp.RequestCtx) {
	task := takeTask(p)
	word := receiveWord(p)

	return func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if rvr := recover(); rvr != nil {
				log.Error().
					Interface("recover", rvr).
					Str("stacktrace", string(debug.Stack())).
					Msg("panic")
				ctxError(ctx, "panic", http.StatusInternalServerError)
			}
		}()

		switch string(ctx.Path()) {
		case "/task":
			task(ctx)
		case "/word":
			word(ctx)
		default:
			ctxError(ctx, `{"error":"Unsupported path"}`, fasthttp.StatusNotFound)
		}
	}
}

func takeTask(p pow.PoW) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType(jsonContentType)

		key, c, err := p.CreateTask()
		if err != nil {
			ctx.Error(`{"error":"can't create task'"}`, fasthttp.StatusInternalServerError)
			log.Error().Err(err).Msg("creating task")

			return
		}

		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody([]byte(`{"key":"` + key + `","complexity":` + strconv.Itoa(c) + `}`))
	}
}

type solvedTask struct {
	Key   string
	Nonce uint64
}

func receiveWord(p pow.PoW) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		var t solvedTask
		ctx.SetContentType(jsonContentType)

		if err := jsoniter.Unmarshal(ctx.PostBody(), &t); err != nil {
			ctx.Error(`{"error":"decoding request"}`, fasthttp.StatusBadRequest)
			log.Error().Err(err).Msg("decoding receive word request")

			return
		}

		if err := p.ValidateTask(t.Key, t.Nonce); err != nil {
			ctx.Error(`{"error":"incorrect task solution"}`, fasthttp.StatusBadRequest)
			log.Error().Err(err).Msg("validating task")

			return
		}

		resp, err := jsoniter.Marshal(quotes.RandomQuote())
		if err != nil {
			ctxError(ctx, `{"error":"encoding word of wisdom"}`, fasthttp.StatusInternalServerError)
			log.Error().Err(err).Msg("encoding word of wisdom")

			return
		}

		ctx.SetStatusCode(http.StatusOK)
		ctx.SetBody(resp)
	}
}

func ctxError(ctx *fasthttp.RequestCtx, msg string, statusCode int) {
	ctx.Response.Reset()
	ctx.SetStatusCode(statusCode)
	ctx.SetContentType(jsonContentType)
	ctx.SetBodyString(msg)
}
