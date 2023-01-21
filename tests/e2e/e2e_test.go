package e2e

import (
	"fmt"
	"strings"
	"time"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	appparams "github.com/tessornetwork/nebula/app/params"
)

func (s *IntegrationTestSuite) TestIBCTokenTransfer() {
	var ibcStakeDenom string

	valAddr, err := s.chain.validators[0].keyInfo.GetAddress()
	s.Require().NoError(err)

	s.Run("send_stake_to_nebula", func() {
		recipient := valAddr.String()
		token := sdk.NewInt64Coin("stake", 3300000000) // 3300stake
		s.sendIBC(gaiaChainID, s.chain.id, recipient, token)

		nebulaAPIEndpoint := fmt.Sprintf("http://%s", s.valResources[0].GetHostPort("1317/tcp"))

		// require the recipient account receives the IBC tokens (IBC packets ACKd)
		var (
			balances sdk.Coins
			err      error
		)
		s.Require().Eventually(
			func() bool {
				balances, err = queryNebulaAllBalances(nebulaAPIEndpoint, recipient)
				s.Require().NoError(err)

				return balances.Len() == 3
			},
			time.Minute,
			5*time.Second,
		)

		for _, c := range balances {
			if strings.Contains(c.Denom, "ibc/") {
				ibcStakeDenom = c.Denom
				s.Require().Equal(token.Amount.Int64(), c.Amount.Int64())
				break
			}
		}

		s.Require().NotEmpty(ibcStakeDenom)
	})

	var ibcStakeERC20Addr string
	s.Run("deploy_stake_erc20 ibcStakeERC20Addr", func() {
		s.T().Skip("paused due to Ethereum PoS migration and PoW fork")
		s.Require().NotEmpty(ibcStakeDenom)
		ibcStakeERC20Addr = s.deployERC20Token(ibcStakeDenom)
	})

	// send 300 stake tokens from Nebula to Ethereum
	s.Run("send_stake_tokens_to_eth", func() {
		s.T().Skip("paused due to Ethereum PoS migration and PoW fork")
		nebulaValIdxSender := 0
		orchestratorIdxReceiver := 1
		amount := sdk.NewCoin(ibcStakeDenom, math.NewInt(300))
		nebulaFee := sdk.NewCoin(appparams.BondDenom, math.NewInt(10000))
		gravityFee := sdk.NewCoin(ibcStakeDenom, math.NewInt(7))

		s.sendFromNebulaToEthCheck(nebulaValIdxSender, orchestratorIdxReceiver, ibcStakeERC20Addr, amount, nebulaFee, gravityFee)
	})

	// send 300 stake tokens from Ethereum back to Nebula
	s.Run("send_stake_tokens_from_eth", func() {
		s.T().Skip("paused due to Ethereum PoS migration and PoW fork")
		nebulaValIdxReceiver := 0
		orchestratorIdxSender := 1
		amount := uint64(300)

		s.sendFromEthToNebulaCheck(orchestratorIdxSender, nebulaValIdxReceiver, ibcStakeERC20Addr, ibcStakeDenom, amount)
	})
}

func (s *IntegrationTestSuite) TestPhotonTokenTransfers() {
	s.T().Skip("paused due to Ethereum PoS migration and PoW fork")
	// deploy photon ERC20 token contact
	var photonERC20Addr string
	s.Run("deploy_photon_erc20", func() {
		photonERC20Addr = s.deployERC20Token(photonDenom)
	})

	// send 100 photon tokens from Nebula to Ethereum
	s.Run("send_photon_tokens_to_eth", func() {
		nebulaValIdxSender := 0
		orchestratorIdxReceiver := 1
		amount := sdk.NewCoin(photonDenom, math.NewInt(100))
		nebulaFee := sdk.NewCoin(appparams.BondDenom, math.NewInt(10000))
		gravityFee := sdk.NewCoin(photonDenom, math.NewInt(3))

		s.sendFromNebulaToEthCheck(nebulaValIdxSender, orchestratorIdxReceiver, photonERC20Addr, amount, nebulaFee, gravityFee)
	})

	// send 100 photon tokens from Ethereum back to Nebula
	s.Run("send_photon_tokens_from_eth", func() {
		s.T().Skip("paused due to Ethereum PoS migration and PoW fork")
		nebulaValIdxReceiver := 0
		orchestratorIdxSender := 1
		amount := uint64(100)

		s.sendFromEthToNebulaCheck(orchestratorIdxSender, nebulaValIdxReceiver, photonERC20Addr, photonDenom, amount)
	})
}

func (s *IntegrationTestSuite) TestNebulaTokenTransfers() {
	s.T().Skip("paused due to Ethereum PoS migration and PoW fork")
	// deploy nebula ERC20 token contract
	var nebulaERC20Addr string
	s.Run("deploy_nebula_erc20", func() {
		nebulaERC20Addr = s.deployERC20Token(appparams.BondDenom)
	})

	// send 300 nebula tokens from Nebula to Ethereum
	s.Run("send_unebula_tokens_to_eth", func() {
		nebulaValIdxSender := 0
		orchestratorIdxReceiver := 1
		amount := sdk.NewCoin(appparams.BondDenom, math.NewInt(300))
		nebulaFee := sdk.NewCoin(appparams.BondDenom, math.NewInt(10000))
		gravityFee := sdk.NewCoin(appparams.BondDenom, math.NewInt(7))

		s.sendFromNebulaToEthCheck(nebulaValIdxSender, orchestratorIdxReceiver, nebulaERC20Addr, amount, nebulaFee, gravityFee)
	})

	// send 300 nebula tokens from Ethereum back to Nebula
	s.Run("send_unebula_tokens_from_eth", func() {
		nebulaValIdxReceiver := 0
		orchestratorIdxSender := 1
		amount := uint64(300)

		s.sendFromEthToNebulaCheck(orchestratorIdxSender, nebulaValIdxReceiver, nebulaERC20Addr, appparams.BondDenom, amount)
	})
}
