package cached

import (
	"testing"
	"time"

	"github.com/dbeliakov/stocks-bot/stocks/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

const (
	testSymbol = "AAA"
	testPrice  = 4.2
	testTTL    = time.Minute
)

func TestProvider_CurrentPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mp := mock.NewMockProvider(ctrl)
	mp.EXPECT().CurrentPrice(testSymbol).Return(testPrice, nil)

	p := NewProvider(mp, 1, testTTL)
	price, err := p.CurrentPrice(testSymbol)
	require.NoError(t, err, "Failed to get current price")
	require.Equal(t, testPrice, price, "Incorrect price returned")
}

func TestProvider_CurrentPrice_Cached(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mp := mock.NewMockProvider(ctrl)
	mp.EXPECT().CurrentPrice(testSymbol).Return(testPrice, nil).Times(1)

	p := NewProvider(mp, 1, testTTL)
	p1, err := p.CurrentPrice(testSymbol)
	require.NoError(t, err, "Failed to get current price")
	p2, err := p.CurrentPrice(testSymbol)
	require.NoError(t, err, "Failed to get current price")
	require.Equal(t, p1, p2, "Two values are not equal")
}

func TestProvider_CurrentPrice_Expired(t *testing.T) {
	const smallTTL = time.Second

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mp := mock.NewMockProvider(ctrl)
	mp.EXPECT().CurrentPrice(testSymbol).Return(testPrice, nil).Times(2)

	p := NewProvider(mp, 1, smallTTL)
	p1, err := p.CurrentPrice(testSymbol)
	require.NoError(t, err, "Failed to get current price")
	time.Sleep(smallTTL + 10*time.Millisecond)
	p2, err := p.CurrentPrice(testSymbol)
	require.NoError(t, err, "Failed to get current price")
	require.Equal(t, p1, p2, "Two values are not equal")
}
