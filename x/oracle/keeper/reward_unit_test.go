package keeper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrependNebulaIfUnique(t *testing.T) {
	require := require.New(t)
	tcs := []struct {
		in  []string
		out []string
	}{
		// Should prepend "unebula" to a slice of denoms, unless it is already present.
		{[]string{}, []string{"unebula"}},
		{[]string{"a"}, []string{"unebula", "a"}},
		{[]string{"x", "a", "heeeyyy"}, []string{"unebula", "x", "a", "heeeyyy"}},
		{[]string{"x", "a", "unebula"}, []string{"x", "a", "unebula"}},
	}
	for i, tc := range tcs {
		require.Equal(tc.out, prependNebulaIfUnique(tc.in), i)
	}

}
