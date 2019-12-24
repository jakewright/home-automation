package domain

import (
	"sort"
	"sync"
)

// Universe represents a set of fixtures in a 512 channel space
type Universe struct {
	Number int

	fixtures []Fixture
	mux      sync.RWMutex
}

// Find returns the fixture with the given ID
func (u *Universe) Find(id string) Fixture {
	u.mux.RLock()
	defer u.mux.RUnlock()

	for _, f := range u.fixtures {
		if f.ID() == id {
			return f
		}
	}

	return nil
}

// AddFixture will add the given fixture to
// the universe if it does not already exist
func (u *Universe) AddFixture(f Fixture) {
	if existing := u.Find(f.ID()); existing != nil {
		return
	}

	u.mux.Lock()
	defer u.mux.Unlock()

	u.fixtures = append(u.fixtures, f)
}

// Valid returns false if any fixtures have overlapping channel ranges
func (u *Universe) Valid() bool {
	u.mux.RLock()
	defer u.mux.RUnlock()

	// Don't modify the slice
	var f []Fixture
	copy(f, u.fixtures)

	// Sort the fixtures by offset
	sort.Slice(f, func(i, j int) bool {
		return f[i].Offset() < f[j].Offset()
	})

	// Make sure each fixture ends before the next one begins
	for i := 0; i < len(f)-1; i++ {
		if f[i].Offset()+len(f[i].DMXValues()) > f[i+1].Offset() {
			return false
		}
	}

	return true
}

// DMXValues returns the value of all channels in the universe
func (u *Universe) DMXValues() [512]byte {
	u.mux.RLock()
	defer u.mux.RUnlock()

	var v [512]byte
	for _, f := range u.fixtures {
		copy(v[f.Offset():], f.DMXValues())
	}
	return v
}
