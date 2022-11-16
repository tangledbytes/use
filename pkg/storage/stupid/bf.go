package stupid

import "sync"

type bf interface {
	Add([]byte)
	Contains([]byte) bool
	Delete([]byte)
}

type bfsync struct {
	bf
	mu *sync.Mutex
}

func newBfSync(bf bf) *bfsync {
	return &bfsync{
		bf: bf,
		mu: &sync.Mutex{},
	}
}

func (b *bfsync) Add(item []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.bf.Add(item)
}

func (b *bfsync) Contains(item []byte) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.bf.Contains(item)
}

func (b *bfsync) Delete(item []byte) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.bf.Delete(item)
}
