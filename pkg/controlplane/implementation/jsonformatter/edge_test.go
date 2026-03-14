package jsonformatter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wernsiet/morchy/pkg/controlplane/domain/workload"
)

func TestNewEdgeFromDomain(t *testing.T) {
	t.Run("edge with all fields", func(t *testing.T) {
		domainEdge := &workload.Edge{
			UpstreamAddress: "http://localhost:8080",
			ProxyPath:       "/api/v1",
		}

		response := NewEdgeFromDomain(domainEdge)

		require.NotNil(t, response)
		assert.Equal(t, "http://localhost:8080", response.UpstreamAddress)
		assert.Equal(t, "/api/v1", response.ProxyPath)
	})

	t.Run("edge with empty proxy path", func(t *testing.T) {
		domainEdge := &workload.Edge{
			UpstreamAddress: "http://localhost:8080",
			ProxyPath:       "",
		}

		response := NewEdgeFromDomain(domainEdge)

		require.NotNil(t, response)
		assert.Equal(t, "http://localhost:8080", response.UpstreamAddress)
		assert.Empty(t, response.ProxyPath)
	})

	t.Run("edge with root path", func(t *testing.T) {
		domainEdge := &workload.Edge{
			UpstreamAddress: "http://localhost:8080",
			ProxyPath:       "/",
		}

		response := NewEdgeFromDomain(domainEdge)

		require.NotNil(t, response)
		assert.Equal(t, "/", response.ProxyPath)
	})

	t.Run("edge with complex path", func(t *testing.T) {
		domainEdge := &workload.Edge{
			UpstreamAddress: "http://localhost:8080",
			ProxyPath:       "/api/v1/service/subresource",
		}

		response := NewEdgeFromDomain(domainEdge)

		require.NotNil(t, response)
		assert.Equal(t, "/api/v1/service/subresource", response.ProxyPath)
	})

	t.Run("edge with various upstream addresses", func(t *testing.T) {
		tests := []struct {
			name            string
			upstreamAddress string
		}{
			{"localhost", "http://localhost:8080"},
			{"IP address", "http://192.168.1.1:8080"},
			{"domain", "http://example.com:8080"},
			{"HTTPS", "https://secure.example.com:443"},
			{"with path", "http://localhost:8080/service"},
			{"with query", "http://localhost:8080?query=value"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				domainEdge := &workload.Edge{
					UpstreamAddress: tt.upstreamAddress,
					ProxyPath:       "/proxy",
				}

				response := NewEdgeFromDomain(domainEdge)

				assert.Equal(t, tt.upstreamAddress, response.UpstreamAddress)
			})
		}
	})

	t.Run("edge with trailing slash in path", func(t *testing.T) {
		domainEdge := &workload.Edge{
			UpstreamAddress: "http://localhost:8080",
			ProxyPath:       "/api/v1/",
		}

		response := NewEdgeFromDomain(domainEdge)

		assert.Equal(t, "/api/v1/", response.ProxyPath)
	})
}
