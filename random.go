package gotil

type LCGState struct {
	a uint32
	c uint32
	m uint32
	seed int64
}

func NewLCGState(s int64) *LCGState {
	return &LCGState{m: 1 << 31, a: 1103515245, c: 12345, seed: s}
}

func (r *LCGState) SetSeed(s int64) {
	r.seed = s
}

func (r *LCGState) Random(max int) int {
	r.seed = ((int64(r.a) * r.seed) + int64(r.c)) % int64(r.m)
	return int(r.seed % int64(max))
}

func (r *LCGState) URandom(max int) uint64 {
	r.seed = ((int64(r.a) * r.seed) + int64(r.c)) % int64(r.m)
	return uint64(uint64(r.seed) % uint64(max))
}
