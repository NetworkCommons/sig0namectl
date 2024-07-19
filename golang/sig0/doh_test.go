package sig0

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindDOHEndpoint(t *testing.T) {
	endpoint, err := FindDOHEndpoint("zenr.io")
	require.NoError(t, err)

	want := "https://doh.zenr.io/dns-query"
	require.Equal(t, want, endpoint.String())
}
