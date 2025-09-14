package consistenthash

import (
	"fmt"
	"sort"
	"sync"

	"github.com/prfc0/consistenthash/hasher"
)

type Slot struct {
	Hash uint64
	Node string
}

type Ring struct {
	mu     sync.RWMutex
	hasher hasher.Hasher

	replicas int
	slots    []Slot
	nodes    map[string]struct{}
}

func New(replicas int, h hasher.Hasher) *Ring {
	if replicas <= 0 {
		replicas = 20
	}
	if h == nil {
		h = hasher.Murmur64{}
	}
	return &Ring{
		hasher:   h,
		replicas: replicas,
		nodes:    make(map[string]struct{}),
	}
}

// AddNode adds a new physical node and its virtual slots to the ring.
func (r *Ring) AddNode(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[node]; exists {
		return // already present
	}

	for i := 0; i < r.replicas; i++ {
		// vnode key = node + index
		data := []byte(fmt.Sprintf("%s|%d", node, i))
		h := r.hasher.Sum(data)
		r.slots = append(r.slots, Slot{Hash: h, Node: node})
	}

	r.nodes[node] = struct{}{}
	r.rebuild()
}

// RemoveNode removes a physical node and all its virtual slots from the ring.
func (r *Ring) RemoveNode(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.nodes[node]; !exists {
		return // not present
	}

	filtered := r.slots[:0]
	for _, s := range r.slots {
		if s.Node != node {
			filtered = append(filtered, s)
		}
	}
	r.slots = filtered
	delete(r.nodes, node)

	// no rebuild() needed — slots remains sorted
	// r.rebuild()
}

// rebuild sorts the slot hashes for binary search lookups.
func (r *Ring) rebuild() {
	sort.Slice(r.slots, func(i, j int) bool {
		return r.slots[i].Hash < r.slots[j].Hash
	})
}
