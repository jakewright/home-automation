package universe

import (
	"sort"
	"sync"

	"github.com/jakewright/home-automation/service.dmx/domain"
)

// Universe represents a set of fixtures in a 512 channel space
type Universe struct {
	Number int

	fixtures map[string]domain.Fixture
	mux      *sync.RWMutex
}

// New returns an initialised universe
func New(number int) *Universe {
	return &Universe{
		Number:   number,
		fixtures: make(map[string]domain.Fixture),
		mux:      &sync.RWMutex{},
	}
}

// Find returns the fixture with the given ID
func (u *Universe) Find(id string) domain.Fixture {
	u.mux.RLock()
	defer u.mux.RUnlock()

	return u.fixtures[id]
}

// AddFixture will add the given fixture to
// the universe if it does not already exist
func (u *Universe) AddFixture(f domain.Fixture) {
	u.mux.Lock()
	defer u.mux.Unlock()

	if _, ok := u.fixtures[f.ID()]; ok {
		return
	}

	u.fixtures[f.ID()] = f
}

// Valid returns false if any fixtures have overlapping channel ranges
func (u *Universe) Valid() bool {
	u.mux.RLock()
	defer u.mux.RUnlock()

	var f []domain.Fixture
	for _, fixture := range u.fixtures {
		f = append(f, fixture)
	}

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
