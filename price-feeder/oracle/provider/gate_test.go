package provider

import (
	"context"
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/tessornetwork/nebula/price-feeder/oracle/types"
)

func TestGateProvider_GetTickerPrices(t *testing.T) {
	p, err := NewGateProvider(
		context.TODO(),
		zerolog.Nop(),
		Endpoint{},
		types.CurrencyPair{Base: "ATOM", Quote: "USDT"},
	)
	require.NoError(t, err)

	t.Run("valid_request_single_ticker", func(t *testing.T) {
		lastPrice := "34.69000000"
		volume := "2396974.02000000"

		tickerMap := map[string]GateTicker{}
		tickerMap["ATOM_USDT"] = GateTicker{
			Symbol: "ATOM_USDT",
			Last:   lastPrice,
			Vol:    volume,
		}

		p.tickers = tickerMap

		prices, err := p.GetTickerPrices(types.CurrencyPair{Base: "ATOM", Quote: "USDT"})
		require.NoError(t, err)
		require.Len(t, prices, 1)
		require.Equal(t, sdk.MustNewDecFromStr(lastPrice), prices["ATOMUSDT"].Price)
		require.Equal(t, sdk.MustNewDecFromStr(volume), prices["ATOMUSDT"].Volume)
	})

	t.Run("valid_request_multi_ticker", func(t *testing.T) {
		lastPriceAtom := "34.69000000"
		lastPriceNEBULA := "41.35000000"
		volume := "2396974.02000000"

		tickerMap := map[string]GateTicker{}
		tickerMap["ATOM_USDT"] = GateTicker{
			Symbol: "ATOM_USDT",
			Last:   lastPriceAtom,
			Vol:    volume,
		}

		tickerMap["NEBULA_USDT"] = GateTicker{
			Symbol: "NEBULA_USDT",
			Last:   lastPriceNEBULA,
			Vol:    volume,
		}

		p.tickers = tickerMap
		prices, err := p.GetTickerPrices(
			types.CurrencyPair{Base: "ATOM", Quote: "USDT"},
			types.CurrencyPair{Base: "NEBULA", Quote: "USDT"},
		)
		require.NoError(t, err)
		require.Len(t, prices, 2)
		require.Equal(t, sdk.MustNewDecFromStr(lastPriceAtom), prices["ATOMUSDT"].Price)
		require.Equal(t, sdk.MustNewDecFromStr(volume), prices["ATOMUSDT"].Volume)
		require.Equal(t, sdk.MustNewDecFromStr(lastPriceNEBULA), prices["NEBULAUSDT"].Price)
		require.Equal(t, sdk.MustNewDecFromStr(volume), prices["NEBULAUSDT"].Volume)
	})

	t.Run("invalid_request_invalid_ticker", func(t *testing.T) {
		prices, err := p.GetTickerPrices(types.CurrencyPair{Base: "FOO", Quote: "BAR"})
		require.EqualError(t, err, "gate failed to get ticker price for FOO_BAR")
		require.Nil(t, prices)
	})
}

func TestGateCurrencyPairToGatePair(t *testing.T) {
	cp := types.CurrencyPair{Base: "ATOM", Quote: "USDT"}
	GateSymbol := currencyPairToGatePair(cp)
	require.Equal(t, GateSymbol, "ATOM_USDT")
}

func TestGateProvider_getSubscriptionMsgs(t *testing.T) {
	provider := &GateProvider{
		subscribedPairs: map[string]types.CurrencyPair{},
	}
	cps := []types.CurrencyPair{
		{Base: "ATOM", Quote: "USDT"},
	}
	subMsgs := provider.getSubscriptionMsgs(cps...)

	msg, _ := json.Marshal(subMsgs[0])
	require.Equal(t, "{\"method\":\"ticker.subscribe\",\"params\":[\"ATOM_USDT\"],\"id\":1}", string(msg))

	msg, _ = json.Marshal(subMsgs[1])
	require.Equal(t, "{\"method\":\"kline.subscribe\",\"params\":[\"ATOM_USDT\",60],\"id\":2}", string(msg))
}
