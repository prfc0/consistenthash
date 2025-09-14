package consistenthash

import (
	"sort"
)

// GetOwner returns the physical node responsible for the given key.
// Returns ErrEmptyRing if ring has no slots.
func (r *Ring) GetOwner(key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.slots) == 0 {
		return "", ErrEmptyRing
	}

	h := r.hasher.Sum([]byte(key))
	idx := r.searchSlotIndex(h)
	return r.slots[idx].Node, nil
}

// GetReplicas returns up to n distinct physical nodes responsible for key.
// If n > total nodes, it returns all distinct nodes in the ring.
func (r *Ring) GetReplicas(key string, n int) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if len(r.slots) == 0 {
		return nil, ErrEmptyRing
	}
	totalNodes := len(r.nodes)
	if n <= 0 || n > totalNodes {
		n = totalNodes
	}

	h := r.hasher.Sum([]byte(key))
	start := r.searchSlotIndex(h)

	result := make([]string, 0, n)
	seen := make(map[string]struct{}, n)

	// Walk clockwise over slots until we have n distinct nodes
	for i := 0; len(result) < n && i < len(r.slots); i++ {
		idx := (start + i) % len(r.slots)
		node := r.slots[idx].Node
		if _, exists := seen[node]; exists {
			continue
		}
		seen[node] = struct{}{}
		result = append(result, node)
	}

	return result, nil
}

// searchSlotIndex finds the index in r.Slots for the first hash >= h.
// If all hashes < h, it wraps and returns 0.
func (r *Ring) searchSlotIndex(h uint64) int {
	// caller holds RLock/RWLock
	idx := sort.Search(len(r.slots), func(i int) bool {
		return r.slots[i].Hash >= h
	})
	if idx == len(r.slots) {
		idx = 0 // wrap-around
	}
	return idx
}
