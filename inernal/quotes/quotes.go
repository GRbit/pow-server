package quotes

import (
	_ "embed"
	"math/rand"

	jsoniter "github.com/json-iterator/go"
	"golang.org/x/xerrors"
)

//go:embed quotes.json
var quotesJSON []byte

type Quote struct {
	Text     string
	Citation string
	Topics   []string
}

func (q Quote) String() string {
	return q.Text + " [" + q.Citation + "]"
}

var (
	quotes []Quote //nolint:gochecknoglobals // since it's private, so it's ok
	l      float64 //nolint:gochecknoglobals // since it's private, so it's ok
)

func init() {
	if err := jsoniter.Unmarshal(quotesJSON, &quotes); err != nil {
		panic(xerrors.Errorf("can't unmarshal assets/quotes.json: %w", err))
	}

	l = float64(len(quotes))
}

func RandomQuote() Quote {
	return quotes[int(rand.Float64()*l)] //nolint:gosec // it's not for security proposes
}
