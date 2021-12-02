package u256

import "github.com/zeebo/blake3"

// Hash returns a Blake3 hash as U256.
func Hash(b []byte) U256 {
	out := blake3.Sum256(b)
	u := U256{}
	copy(u[:], out[:])
	return u
}

// Hasher computes Blake3 hash as U256.
type Hasher struct {
	h *blake3.Hasher
}

// NewHasher returns a new hasher.
func NewHasher() *Hasher {
	return &Hasher{
		h: blake3.New(),
	}
}

// Write hashes bytes.
func (h *Hasher) Write(p []byte) (int, error) {
	h.h.Write(p)
	return len(p), nil
}

// WriteString hashes a string.
func (h *Hasher) WriteString(p string) (int, error) {
	h.h.WriteString(p)
	return len(p), nil
}

// Sum returns a hash as U256.
func (h *Hasher) Sum() U256 {
	sum := [32]byte{}
	h.h.Sum(sum[:0])

	u := U256{}
	copy(u[:], sum[:])
	return u
}
