package hasher

import "github.com/spaolacci/murmur3"

// Murmur64 implements Hasher using murmur3 64-bit.
type Murmur64 struct{}

func (Murmur64) Sum(data []byte) uint64 {
	return murmur3.Sum64(data)
}
