package hasher

// Hasher defines the interface for hashing algorithms.
// It allows plugging in different hash functions.
type Hasher interface {
	Sum(data []byte) uint64
}
