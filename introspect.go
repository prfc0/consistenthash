package consistenthash

// Nodes returns a snapshot slice of all physical nodes currently in the ring.
// The order is not guaranteed.
func (r *Ring) Nodes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]string, 0, len(r.nodes))
	for n := range r.nodes {
		out = append(out, n)
	}
	return out
}

// Slots returns a copy of all virtual slots (vnodes) currently in the ring.
// Each Slot contains the 64-bit hash and the associated physical node id.
func (r *Ring) Slots() []Slot {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Slot, len(r.slots))
	copy(out, r.slots)
	return out
}
