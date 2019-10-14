package timemap

import (
	"sort"
	"sync"
	"time"
)

// tuple is used to hold the value and the timestamp.
type tuple struct {
	timestamp time.Time
	value     interface{}
}

// Map is a thread-safe fixed size time-based map.
type Map struct {
	timemap map[string][]*tuple
	lock    sync.RWMutex
}

// New creates a time-based map.
func New() *Map {
	return &Map{
		timemap: make(map[string][]*tuple),
	}
}

// Set sets a new value and timestamp to a specific key.
func (t *Map) Set(key string, value interface{}, timestamp time.Time) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	if key == "" {
		return false
	}
	pair := &tuple{timestamp, value}
	t.timemap[key] = append(t.timemap[key], pair)
	return true
}

// Get looks up a key's value based on a timestamp.
func (t *Map) Get(key string, timestamp time.Time) (value interface{}, ok bool) {
	t.lock.Lock()
	defer t.lock.Unlock()
	if !t.Contains(key) {
		return nil, false
	}
	ent, _ := t.timemap[key]
	i := sort.Search(len(ent), func(i int) bool { return ent[i].timestamp.After(timestamp) })
	return ent[i-1].value, true
}

// Contains checks if a key exists in the map.
func (t *Map) Contains(key string) bool {
	_, ok := t.timemap[key]
	return ok
}

// Remove removes the provided key from the map.
func (t *Map) Remove(key string) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	if !t.Contains(key) {
		return false
	}
	delete(t.timemap, key)
	return true
}

// Keys returns a slice of the keys from the map.
func (t *Map) Keys() []string {
	t.lock.RLock()
	defer t.lock.RUnlock()
	keys := []string{}
	for k := range t.timemap {
		keys = append(keys, k)
	}
	return keys
}
