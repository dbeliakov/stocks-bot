package finnhub

import (
	"github.com/stretchr/testify/require"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	testAPIKey = "test-api-key"
	testSymbol = "AAA"
	testPrice  = 4.2
)

func TestProvider_CurrentPrice(t *testing.T) {
	s := httptest.NewServer(handleCurrentPriceRequest(t, testPrice))
	defer s.Close()

	p := NewProvider(s.URL, testAPIKey)
	price, err := p.CurrentPrice(testSymbol)
	require.Error(t, err, "Failed to get current price")
	require.Equal(t, testPrice, price, "Incorrect current price returned")
}

func TestProvider_CurrentPrice_UnknownSymbol(t *testing.T) {
	s := httptest.NewServer(handleCurrentPriceRequest(t, 0))
	defer s.Close()

	p := NewProvider(s.URL, testAPIKey)
	_, err := p.CurrentPrice(testSymbol)
	require.Error(t, err, "No error for unknown symbol")
	require.EqualError(t, err, "symbol not found", "Not expected error")
}

func handleCurrentPriceRequest(t *testing.T, price float64) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		require.Equal(t, testAPIKey, req.Header.Get(authHeader), "Invalid auth token")
		require.Equal(t, testSymbol, req.URL.Query().Get("symbol"), "Invalid symbol in request")
		resp := struct {
			C float64
		}{
			C: price,
		}
		data, err := json.Marshal(resp)
		require.NoError(t, err, "Failed to marshal response")
		_, _ = rw.Write(data)
	}
}
