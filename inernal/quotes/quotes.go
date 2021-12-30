package quotes

import (
	_ "embed"
	"encoding/json"
	"math/rand"

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
	quotes []Quote
	l      float64
)

func init() {
	if err := json.Unmarshal(quotesJSON, &quotes); err != nil {
		panic(xerrors.Errorf("can't unmarshal assets/quotes.json: %w", err))
	}

	l = float64(len(quotes))
}

func RandomQuote() Quote {
	return quotes[int(rand.Float64()*l)]
}
