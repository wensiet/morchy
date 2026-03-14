package workload

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEdge_Fields(t *testing.T) {
	t.Run("edge with all fields", func(t *testing.T) {
		edge := &Edge{
			UpstreamAddress: "http://localhost:8080",
			ProxyPath:       "/api/v1",
		}

		assert.Equal(t, "http://localhost:8080", edge.UpstreamAddress)
		assert.Equal(t, "/api/v1", edge.ProxyPath)
	})

	t.Run("edge with minimal fields", func(t *testing.T) {
		edge := &Edge{
			UpstreamAddress: "http://localhost:8080",
			ProxyPath:       "",
		}

		assert.Equal(t, "http://localhost:8080", edge.UpstreamAddress)
		assert.Empty(t, edge.ProxyPath)
	})
}

func TestEdge_UpstreamAddress(t *testing.T) {
	t.Run("various upstream address formats", func(t *testing.T) {
		tests := []struct {
			name            string
			upstreamAddress string
		}{
			{"localhost", "http://localhost:8080"},
			{"IP address", "http://192.168.1.1:8080"},
			{"domain name", "http://example.com:8080"},
			{"HTTPS", "https://secure.example.com:443"},
			{"with path", "http://localhost:8080/service"},
			{"with query", "http://localhost:8080?query=value"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				edge := &Edge{
					UpstreamAddress: tt.upstreamAddress,
					ProxyPath:       "/proxy",
				}

				assert.Equal(t, tt.upstreamAddress, edge.UpstreamAddress)
			})
		}
	})
}

func TestEdge_ProxyPath(t *testing.T) {
	t.Run("various proxy path formats", func(t *testing.T) {
		tests := []struct {
			name      string
			proxyPath string
		}{
			{"root path", "/"},
			{"single segment", "/api"},
			{"multiple segments", "/api/v1/service"},
			{"with trailing slash", "/api/v1/"},
			{"empty path", ""},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				edge := &Edge{
					UpstreamAddress: "http://localhost:8080",
					ProxyPath:       tt.proxyPath,
				}

				assert.Equal(t, tt.proxyPath, edge.ProxyPath)
			})
		}
	})
}

func TestEdge_EmptyFields(t *testing.T) {
	t.Run("empty upstream address", func(t *testing.T) {
		edge := &Edge{
			UpstreamAddress: "",
			ProxyPath:       "/api",
		}

		assert.Empty(t, edge.UpstreamAddress)
		assert.Equal(t, "/api", edge.ProxyPath)
	})

	t.Run("empty proxy path", func(t *testing.T) {
		edge := &Edge{
			UpstreamAddress: "http://localhost:8080",
			ProxyPath:       "",
		}

		assert.Equal(t, "http://localhost:8080", edge.UpstreamAddress)
		assert.Empty(t, edge.ProxyPath)
	})

	t.Run("both fields empty", func(t *testing.T) {
		edge := &Edge{
			UpstreamAddress: "",
			ProxyPath:       "",
		}

		assert.Empty(t, edge.UpstreamAddress)
		assert.Empty(t, edge.ProxyPath)
	})
}
