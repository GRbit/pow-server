package pow

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"math/big"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/sha3"
	"golang.org/x/xerrors"
)

const maxComplexity = byte(127)

type pow struct {
	cache      *fastcache.Cache
	complexity byte
}

func New(c *fastcache.Cache, defaultComplexity byte) PoW {
	p := &pow{
		cache:      c,
		complexity: defaultComplexity,
	}
	go p.balanceLoad()

	return p
}

type PoW interface {
	CreateTask() (string, int, error)
	ValidateTask(key string, nonce uint64) error
}

func (p *pow) balanceLoad() {}

func (p *pow) CreateTask() (key string, complexity int, err error) {
	k := make([]byte, 16)
	if _, err = rand.Read(k); err != nil {
		return "", 0, xerrors.Errorf("reading random bytes: %w", err)
	}

	c := p.complexity
	p.cache.Set(k, []byte{c})

	key = hex.EncodeToString(k)

	log.Debug().
		Str("key", key).
		Int("complexity", int(c)).
		Msg("key created")

	return key, int(c), nil
}

func (p *pow) ValidateTask(key string, nonce uint64) error {
	bKey, err := hex.DecodeString(key)
	if err != nil {
		return xerrors.Errorf("decoding key string: %w", err)
	}

	c := p.cache.Get(nil, bKey)
	n := make([]byte, 8)
	binary.LittleEndian.PutUint64(n[0:], nonce)
	s := sha3.Sum256(append(bKey, n...))

	log.Debug().
		Str("key", key).
		Uint64("nonce", nonce).
		Str("hash", hex.EncodeToString(s[:])).
		Int("complexity", int(c[0])).
		Msg("validating task")

	if !p.checkHash(s[:], c[0]) {
		return xerrors.New("incorrect solution")
	}

	p.cache.Del(bKey)

	return nil
}

func (p *pow) checkHash(s []byte, c byte) bool {
	var sumI big.Int
	sumI.SetBytes(s[:])

	target := big.NewInt(1)
	target.Lsh(target, uint(256-int(c)))

	return sumI.Cmp(target) == -1
}
