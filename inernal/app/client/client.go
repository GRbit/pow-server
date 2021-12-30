package client

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/GRbit/pow-server/inernal/logger"
	"github.com/GRbit/pow-server/inernal/quotes"
	"github.com/go-resty/resty/v2"
	"github.com/jessevdk/go-flags"
	jsoniter "github.com/json-iterator/go"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/sha3"
	"golang.org/x/xerrors"
)

func GetWord() error {
	var cfg clientConfig

	if _, err := flags.Parse(&cfg); err != nil {
		return xerrors.Errorf("parsing client config: %w", err)
	}

	logger.Config(cfg.LogLevel, cfg.Console)

	c := resty.New()

	r, err := c.R().Post("http://" + cfg.Addr + "/task")
	if err != nil {
		return xerrors.Errorf("making request for task: %w", err)
	}

	k := struct {
		Key        string
		Complexity int
	}{}
	if err = jsoniter.Unmarshal(r.Body(), &k); err != nil {
		return xerrors.Errorf("unmarshalling task request: %w", err)
	}

	log.Debug().
		Str("key", k.Key).
		Int("complexity", k.Complexity).
		Msg("task received")

	nonce, err := solveTask(k.Key, k.Complexity, cfg.MaxComplexity)
	if err != nil {
		return xerrors.Errorf("solving task: %w", err)
	}

	r, err = c.R().SetBody(map[string]interface{}{
		"key":   k.Key,
		"nonce": nonce,
	}).
		Post("http://" + cfg.Addr + "/word")
	if err != nil {
		return xerrors.Errorf("making request for word: %w", err)
	}

	q := quotes.Quote{}
	if err = jsoniter.Unmarshal(r.Body(), &q); err != nil {
		return xerrors.Errorf("unmarshalling word request: %w", err)
	}

	log.Debug().
		Str("text", q.Text).
		Str("citation", q.Citation).
		Strs("topics", q.Topics).
		Msg("word of wisdom received")

	fmt.Println(q)

	return nil
}

func solveTask(key string, c int, max uint64) (uint64, error) {
	bKey, err := hex.DecodeString(key)
	if err != nil {
		return 0, xerrors.Errorf("decoding key string: %w", err)
	}

	var sum big.Int

	target := big.NewInt(1)
	target.Lsh(target, uint(256-int(c)))

	for i := uint64(0); i < max; i++ {
		n := make([]byte, 8)
		binary.LittleEndian.PutUint64(n[0:], i)
		s := sha3.Sum256(append(bKey, n...))
		sum.SetBytes(s[:])

		if sum.Cmp(target) == -1 {
			return i, nil
		}
	}

	return 0, xerrors.Errorf("reached maximum of tries")
}
