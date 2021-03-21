package mime_test

import (
	"github.com/partyzanex/s3-http-proxy/internal/mime"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTypes(t *testing.T) {
	types, err := mime.Types()
	require.NoError(t, err)
	require.NotNil(t, types)

	t.Log(types)
}
