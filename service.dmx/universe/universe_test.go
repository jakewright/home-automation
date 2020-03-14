package universe

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	devicedef "github.com/jakewright/home-automation/libraries/go/device/def"
	deviceregistrydef "github.com/jakewright/home-automation/service.device-registry/def"
	"github.com/jakewright/home-automation/service.dmx/domain"
)

var ctr = 0

type mockFixture struct {
	id     string
	offset int
	len    int
}

func (f *mockFixture) DMXValues() []byte                                  { return make([]byte, f.len) }
func (f *mockFixture) Offset() int                                        { return f.offset }
func (f *mockFixture) ID() string                                         { return f.id }
func (f *mockFixture) SetHeader(*deviceregistrydef.DeviceHeader) error    { panic("implement me") }
func (f *mockFixture) ToDef() *devicedef.Device                           { panic("implement me") }
func (f *mockFixture) SetProperties(map[string]interface{}) (bool, error) { panic("implement me") }
func (f *mockFixture) Copy() (domain.Fixture, error)                      { panic("implement me") }

func mf(offset, len int) *mockFixture {
	ctr++
	id := strconv.Itoa(ctr)

	return &mockFixture{
		id:     id,
		offset: offset,
		len:    len,
	}
}

func TestUniverse_Valid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fixtures []domain.Fixture
		want     bool
	}{
		{
			name:     "No fixtures",
			fixtures: []domain.Fixture{},
			want:     true,
		},
		{
			name:     "1 fixture",
			fixtures: []domain.Fixture{mf(0, 5)},
			want:     true,
		},
		{
			name: "3 non-overlapping fixtures",
			fixtures: []domain.Fixture{
				mf(0, 7),
				mf(7, 5),
				mf(12, 10),
			},
			want: true,
		},
		{
			name: "overlapping fixtures",
			fixtures: []domain.Fixture{
				mf(0, 7),
				mf(6, 5),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			u := New(1)
			for _, f := range tt.fixtures {
				u.AddFixture(f)
			}

			got := u.Valid()

			require.Equal(t, tt.want, got)
		})
	}
}
