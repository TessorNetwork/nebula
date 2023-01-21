package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/tessornetwork/nebula/v3/x/leverage/types"
)

func (s *IntegrationTestSuite) TestGetCollateralAmount() {
	ctx, require := s.ctx, s.Require()
	uDenom := types.ToUTokenDenom(nebulaDenom)

	// get u/nebula collateral amount of empty account address
	collateral := s.tk.GetCollateral(ctx, sdk.AccAddress{}, uDenom)
	require.Equal(coin(uDenom, 0), collateral)

	// fund an account
	addr := s.newAccount(coin(nebulaDenom, 1000))

	// get u/nebula collateral amount of non-empty account address
	collateral = s.tk.GetCollateral(ctx, addr, uDenom)
	require.Equal(coin(uDenom, 0), collateral)

	// supply 1000 u/unebula but do not collateralize
	s.supply(addr, coin(nebulaDenom, 1000))

	// confirm collateral amount is 0 u/unebula
	collateral = s.tk.GetCollateral(ctx, addr, uDenom)
	require.Equal(coin(uDenom, 0), collateral)

	// enable u/nebula as collateral
	s.collateralize(addr, coin(uDenom, 1000))

	// confirm collateral amount is 1000 u/unebula
	collateral = s.tk.GetCollateral(ctx, addr, uDenom)
	require.Equal(coin(uDenom, 1000), collateral)

	// collateral amount of non-utoken denom (zero)
	collateral = s.tk.GetCollateral(ctx, addr, nebulaDenom)
	require.Equal(coin(nebulaDenom, 0), collateral)

	// collateral amount of unregistered denom (zero)
	collateral = s.tk.GetCollateral(ctx, addr, "abcd")
	require.Equal(coin("abcd", 0), collateral)

	// disable u/nebula as collateral
	s.decollateralize(addr, coin(uDenom, 1000))

	// confirm collateral amount is 0 u/unebula
	collateral = s.tk.GetCollateral(ctx, addr, uDenom)
	require.Equal(coin(uDenom, 0), collateral)

	// we do not test empty denom, as that will cause a panic
}

func (s *IntegrationTestSuite) TestSetCollateralAmount() {
	ctx, require := s.ctx, s.Require()
	uDenom := types.ToUTokenDenom(nebulaDenom)

	// set u/nebula collateral amount of empty account address (error)
	err := s.tk.SetCollateral(ctx, sdk.AccAddress{}, coin(uDenom, 0))
	require.ErrorIs(err, types.ErrEmptyAddress)

	addr := s.newAccount()

	// force invalid denom
	err = s.tk.SetCollateral(ctx, addr, sdk.Coin{Denom: "", Amount: sdk.ZeroInt()})
	require.ErrorContains(err, "invalid denom")

	// base denom
	err = s.tk.SetCollateral(ctx, addr, sdk.Coin{Denom: nebulaDenom, Amount: sdk.ZeroInt()})
	require.ErrorIs(err, types.ErrNotUToken)

	// negative amount
	err = s.tk.SetCollateral(ctx, addr, sdk.Coin{Denom: uDenom, Amount: sdk.NewInt(-1)})
	require.ErrorContains(err, "negative coin amount")

	// set u/nebula collateral amount
	require.NoError(s.tk.SetCollateral(ctx, addr, coin(uDenom, 10)))

	// confirm effect
	collateral := s.tk.GetCollateral(ctx, addr, uDenom)
	require.Equal(coin(uDenom, 10), collateral)

	// set u/nebula collateral amount to zero
	require.NoError(s.tk.SetCollateral(ctx, addr, coin(uDenom, 0)))

	// confirm effect
	collateral = s.tk.GetCollateral(ctx, addr, uDenom)
	require.Equal(coin(uDenom, 0), collateral)

	// set unregistered token collateral amount
	require.NoError(s.tk.SetCollateral(ctx, addr, coin("u/abcd", 10)))

	// confirm effect
	collateral = s.tk.GetCollateral(ctx, addr, "u/abcd")
	require.Equal(coin("u/abcd", 10), collateral)

	// set unregistered token collateral amount to zero
	require.NoError(s.tk.SetCollateral(ctx, addr, coin("u/abcd", 0)))

	// confirm effect
	collateral = s.tk.GetCollateral(ctx, addr, "u/abcd")
	require.Equal(coin("u/abcd", 0), collateral)

	// we do not test empty denom, as that will cause a panic
}

func (s *IntegrationTestSuite) TestTotalCollateral() {
	app, ctx, require := s.app, s.ctx, s.Require()

	// not a uToken
	collateral := app.LeverageKeeper.GetTotalCollateral(ctx, nebulaDenom)
	require.Equal(sdk.Coin{}, collateral, "not a uToken")

	// Test zero collateral
	uDenom := types.ToUTokenDenom(nebulaDenom)
	collateral = app.LeverageKeeper.GetTotalCollateral(ctx, uDenom)
	require.Equal(coin(uDenom, 0), collateral, "zero collateral")

	// create a supplier which will have 100 u/NEBULA collateral
	addr := s.newAccount(coin(nebulaDenom, 100_000000))
	s.supply(addr, coin(nebulaDenom, 100_000000))
	s.collateralize(addr, coin(uDenom, 100_000000))

	// Test nonzero collateral
	collateral = app.LeverageKeeper.GetTotalCollateral(ctx, uDenom)
	require.Equal(coin(uDenom, 100_000000), collateral, "nonzero collateral")
}
