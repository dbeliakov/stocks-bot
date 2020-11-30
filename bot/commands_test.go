package bot

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"

	mockprovider "github.com/dbeliakov/stocks-bot/stocks/mock"
	mockstorage "github.com/dbeliakov/stocks-bot/storage/mock"
)

const (
	testChatID int64   = 42
	testSymbol string  = "AAA"
	testPrice  float64 = 4.2
)

var (
	currentState string
	addedSymbols []string
)

func TestCurrentPrice(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mp := mockprovider.NewMockProvider(ctrl)
	mp.EXPECT().CurrentPrice(testSymbol).Return(testPrice, nil).AnyTimes()

	ms := mockstorage.NewMockStorage(ctrl)
	ms.EXPECT().GetState(testChatID).DoAndReturn(func(_ int64) (string, error) {
		return currentState, nil
	}).AnyTimes()
	ms.EXPECT().SetState(testChatID, gomock.AssignableToTypeOf("")).DoAndReturn(func(_ int64, state string) error {
		currentState = state
		return nil
	}).AnyTimes()
	replies := make(chan Reply, 1)

	p := NewProcessor(mp, ms, testChatID, replies, newTestCounter(), newTestCounter())
	p.Process(IncomingMessage{
		Command: "price",
		Message: "",
	})
	r := getReply(t, replies)
	require.Equal(t, EnterSymbolText, r.Message, "Unexpected reply on /price command")

	p.Process(IncomingMessage{
		Command: "",
		Message: testSymbol,
	})
	r = getReply(t, replies)
	require.Equal(t, fmt.Sprintf(PriceText, testSymbol, testPrice), r.Message, "Unexpected price reply")
}

func TestSymbolsList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mp := mockprovider.NewMockProvider(ctrl)
	mp.EXPECT().CurrentPrice(testSymbol).Return(testPrice, nil).AnyTimes()

	ms := mockstorage.NewMockStorage(ctrl)
	ms.EXPECT().GetState(testChatID).DoAndReturn(func(_ int64) (string, error) {
		return currentState, nil
	}).AnyTimes()
	ms.EXPECT().SetState(testChatID, gomock.AssignableToTypeOf("")).DoAndReturn(func(_ int64, state string) error {
		currentState = state
		return nil
	}).AnyTimes()
	ms.EXPECT().AddSymbol(testChatID, gomock.AssignableToTypeOf("")).DoAndReturn(func(_ int64, s string) error {
		addedSymbols = append(addedSymbols, s)
		return nil
	}).AnyTimes()
	ms.EXPECT().RemoveSymbol(testChatID, gomock.AssignableToTypeOf("")).DoAndReturn(func(_ int64, s string) error {
		// Clean all list in this test cases
		addedSymbols = nil
		return nil
	}).AnyTimes()
	ms.EXPECT().Symbols(testChatID).DoAndReturn(func(_ int64) ([]string, error) {
		return addedSymbols, nil
	}).AnyTimes()

	replies := make(chan Reply, 1)

	p := NewProcessor(mp, ms, testChatID, replies, newTestCounter(), newTestCounter())
	p.Process(IncomingMessage{
		Command: "my",
		Message: "",
	})
	r := getReply(t, replies)
	require.Equal(t, NoSymbolsText, r.Message, "Unexpected reply on /my command")

	p.Process(IncomingMessage{
		Command: "add",
		Message: "",
	})
	r = getReply(t, replies)
	require.Equal(t, EnterSymbolText, r.Message, "Unexpected reply on /add command")

	p.Process(IncomingMessage{
		Command: "",
		Message: testSymbol,
	})
	r = getReply(t, replies)
	require.Equal(t, SuccessText, r.Message, "Not a success message")

	p.Process(IncomingMessage{
		Command: "my",
		Message: "",
	})
	r = getReply(t, replies)
	require.Equal(t, fmt.Sprintf(PriceText+"\n", testSymbol, testPrice), r.Message, "Unexpected my symbols reply")

	p.Process(IncomingMessage{
		Command: "remove",
		Message: "",
	})
	r = getReply(t, replies)
	require.Equal(t, EnterSymbolText, r.Message, "Unexpected reply in /remove command")

	p.Process(IncomingMessage{
		Command: "",
		Message: testSymbol,
	})
	r = getReply(t, replies)
	require.Equal(t, SuccessText, r.Message, "Not a success message")

	p.Process(IncomingMessage{
		Command: "my",
		Message: "",
	})
	r = getReply(t, replies)
	require.Equal(t, NoSymbolsText, r.Message, "Unexpected reply on /my command")
}

func getReply(t *testing.T, m <-chan Reply) Reply {
	select {
	case r := <-m:
		return r
	default:
		require.Fail(t, "No reply in chan")
		return Reply{}
	}
}

func newTestCounter() prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{})
}
