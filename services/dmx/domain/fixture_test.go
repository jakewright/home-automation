package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateFixtures(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		fixtures []Fixture
		want     bool
	}{
		{
			name:     "No fixtures",
			fixtures: []Fixture{},
			want:     true,
		},
		{
			name: "One fixture",
			fixtures: []Fixture{
				&MockFixture{Ofs: 0, Len: 5},
			},
			want: true,
		},
		{
			name: "Three non-overlapping fixtures",
			fixtures: []Fixture{
				&MockFixture{Ofs: 0, Len: 7},
				&MockFixture{Ofs: 7, Len: 5},
				&MockFixture{Ofs: 12, Len: 10},
			},
			want: true,
		},
		{
			name: "Overlapping fixtures",
			fixtures: []Fixture{
				&MockFixture{Ofs: 0, Len: 7},
				&MockFixture{Ofs: 6, Len: 5},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateFixtures(tt.fixtures)
			if tt.want {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}
