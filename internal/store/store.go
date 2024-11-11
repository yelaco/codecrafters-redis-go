package store

import "sync"

var (
	store = make(map[string]string)
	mu    = &sync.RWMutex{}
)

func Set(key, value string) {
	mu.Lock()
	defer mu.Unlock()
	store[key] = value
}

func Get(key string) (string, bool) {
	mu.RLock()
	defer mu.RUnlock()
	v, ok := store[key]
	return v, ok
}

func Delete(key string) {
	mu.Lock()
	defer mu.Unlock()
	delete(store, key)
}
