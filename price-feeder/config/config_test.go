package config_test

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/cosmos/cosmos-sdk/telemetry"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"github.com/tessornetwork/nebula/price-feeder/config"
	"github.com/tessornetwork/nebula/price-feeder/oracle/provider"
)

func TestValidate(t *testing.T) {
	validConfig := func() config.Config {
		return config.Config{
			Server: config.Server{
				ListenAddr:     "0.0.0.0:7171",
				VerboseCORS:    false,
				AllowedOrigins: []string{},
			},
			CurrencyPairs: []config.CurrencyPair{
				{Base: "ATOM", Quote: "USDT", Providers: []provider.Name{provider.ProviderKraken}},
			},
			Account: config.Account{
				Address:   "fromaddr",
				Validator: "valaddr",
				ChainID:   "chain-id",
			},
			Keyring: config.Keyring{
				Backend: "test",
				Dir:     "/Users/username/.nebula",
			},
			RPC: config.RPC{
				TMRPCEndpoint: "http://localhost:26657",
				GRPCEndpoint:  "localhost:9090",
				RPCTimeout:    "100ms",
			},
			Telemetry: telemetry.Config{
				ServiceName:             "price-feeder",
				Enabled:                 true,
				EnableHostname:          true,
				EnableHostnameLabel:     true,
				EnableServiceLabel:      true,
				GlobalLabels:            make([][]string, 1),
				PrometheusRetentionTime: 120,
			},
			GasAdjustment: 1.5,
		}
	}
	emptyPairs := validConfig()
	emptyPairs.CurrencyPairs = []config.CurrencyPair{}

	invalidBase := validConfig()
	invalidBase.CurrencyPairs = []config.CurrencyPair{
		{Base: "", Quote: "USDT", Providers: []provider.Name{provider.ProviderKraken}},
	}

	invalidQuote := validConfig()
	invalidQuote.CurrencyPairs = []config.CurrencyPair{
		{Base: "ATOM", Quote: "", Providers: []provider.Name{provider.ProviderKraken}},
	}

	emptyProviders := validConfig()
	emptyProviders.CurrencyPairs = []config.CurrencyPair{
		{Base: "ATOM", Quote: "USDT", Providers: []provider.Name{}},
	}

	invalidEndpoints := validConfig()
	invalidEndpoints.ProviderEndpoints = []provider.Endpoint{
		{
			Name: provider.ProviderBinance,
		},
	}

	invalidEndpointsProvider := validConfig()
	invalidEndpointsProvider.ProviderEndpoints = []provider.Endpoint{
		{
			Name:      "foo",
			Rest:      "bar",
			Websocket: "baz",
		},
	}

	testCases := []struct {
		name      string
		cfg       config.Config
		expectErr bool
	}{
		{
			"valid config",
			validConfig(),
			false,
		},
		{
			"empty pairs",
			emptyPairs,
			true,
		},
		{
			"invalid base",
			invalidBase,
			true,
		},
		{
			"invalid quote",
			invalidQuote,
			true,
		},
		{
			"empty providers",
			emptyProviders,
			true,
		},
		{
			"invalid endpoints",
			invalidEndpoints,
			true,
		},
		{
			"invalid endpoint provider",
			invalidEndpointsProvider,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.cfg.Validate() != nil, tc.expectErr)
		})
	}
}

func TestParseConfig_Valid(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "price-feeder*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte(`
gas_adjustment = 1.5

[server]
listen_addr = "0.0.0.0:99999"
read_timeout = "20s"
verbose_cors = true
write_timeout = "20s"

[[currency_pairs]]
base = "ATOM"
quote = "USDT"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[[currency_pairs]]
base = "NEBULA"
quote = "USDT"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[[currency_pairs]]
base = "USDT"
quote = "USD"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[account]
address = "nebula15nejfgcaanqpw25ru4arvfd0fwy6j8clccvwx4"
validator = "nebulavalcons14rjlkfzp56733j5l5nfk6fphjxymgf8mj04d5p"
chain_id = "nebula-local-testnet"

[keyring]
backend = "test"
dir = "/Users/username/.nebula"
pass = "keyringPassword"

[rpc]
tmrpc_endpoint = "http://localhost:26657"
grpc_endpoint = "localhost:9090"
rpc_timeout = "100ms"

[telemetry]
service-name = "price-feeder"
enabled = true
enable-hostname = true
enable-hostname-label = true
enable-service-label = true
prometheus-retention = 120
global-labels = [["chain-id", "nebula-local-testnet"]]
`)
	_, err = tmpFile.Write(content)
	require.NoError(t, err)

	cfg, err := config.ParseConfig(tmpFile.Name())
	require.NoError(t, err)

	require.Equal(t, "0.0.0.0:99999", cfg.Server.ListenAddr)
	require.Equal(t, "20s", cfg.Server.WriteTimeout)
	require.Equal(t, "20s", cfg.Server.ReadTimeout)
	require.True(t, cfg.Server.VerboseCORS)
	require.Len(t, cfg.CurrencyPairs, 3)
	require.Equal(t, "ATOM", cfg.CurrencyPairs[0].Base)
	require.Equal(t, "USDT", cfg.CurrencyPairs[0].Quote)
	require.Len(t, cfg.CurrencyPairs[0].Providers, 3)
	require.Equal(t, provider.ProviderKraken, cfg.CurrencyPairs[0].Providers[0])
	require.Equal(t, provider.ProviderBinance, cfg.CurrencyPairs[0].Providers[1])
}

func TestParseConfig_Valid_NoTelemetry(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "price-feeder*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte(`
gas_adjustment = 1.5

[server]
listen_addr = "0.0.0.0:99999"
read_timeout = "20s"
verbose_cors = true
write_timeout = "20s"

[[currency_pairs]]
base = "ATOM"
quote = "USDT"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[[currency_pairs]]
base = "NEBULA"
quote = "USDT"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[[currency_pairs]]
base = "USDT"
quote = "USD"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[account]
address = "nebula15nejfgcaanqpw25ru4arvfd0fwy6j8clccvwx4"
validator = "nebulavalcons14rjlkfzp56733j5l5nfk6fphjxymgf8mj04d5p"
chain_id = "nebula-local-testnet"

[keyring]
backend = "test"
dir = "/Users/username/.nebula"
pass = "keyringPassword"

[rpc]
tmrpc_endpoint = "http://localhost:26657"
grpc_endpoint = "localhost:9090"
rpc_timeout = "100ms"

[telemetry]
enabled = false
`)
	_, err = tmpFile.Write(content)
	require.NoError(t, err)

	cfg, err := config.ParseConfig(tmpFile.Name())
	require.NoError(t, err)

	require.Equal(t, "0.0.0.0:99999", cfg.Server.ListenAddr)
	require.Equal(t, "20s", cfg.Server.WriteTimeout)
	require.Equal(t, "20s", cfg.Server.ReadTimeout)
	require.True(t, cfg.Server.VerboseCORS)
	require.Len(t, cfg.CurrencyPairs, 3)
	require.Equal(t, "ATOM", cfg.CurrencyPairs[0].Base)
	require.Equal(t, "USDT", cfg.CurrencyPairs[0].Quote)
	require.Len(t, cfg.CurrencyPairs[0].Providers, 3)
	require.Equal(t, provider.ProviderKraken, cfg.CurrencyPairs[0].Providers[0])
	require.Equal(t, provider.ProviderBinance, cfg.CurrencyPairs[0].Providers[1])
	require.Equal(t, cfg.Telemetry.Enabled, false)
}

func TestParseConfig_InvalidProvider(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "price-feeder*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte(`
listen_addr = ""

[[currency_pairs]]
base = "ATOM"
quote = "USD"
providers = [
	"kraken",
	"binance"
]

[[currency_pairs]]
base = "NEBULA"
quote = "USD"
providers = [
	"kraken",
	"foobar"
]
`)
	_, err = tmpFile.Write(content)
	require.NoError(t, err)

	_, err = config.ParseConfig(tmpFile.Name())
	require.Error(t, err)
}

func TestParseConfig_NonUSDQuote(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "price-feeder*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte(`
listen_addr = ""

[[currency_pairs]]
base = "ATOM"
quote = "USDT"
providers = [
	"kraken",
	"binance"
]

[[currency_pairs]]
base = "NEBULA"
quote = "USDT"
providers = [
	"kraken",
	"binance"
]
`)
	_, err = tmpFile.Write(content)
	require.NoError(t, err)

	_, err = config.ParseConfig(tmpFile.Name())
	require.Error(t, err)
}

func TestParseConfig_Valid_Deviations(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "price-feeder*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte(`
gas_adjustment = 1.5

[server]
listen_addr = "0.0.0.0:99999"
read_timeout = "20s"
verbose_cors = true
write_timeout = "20s"

[[deviation_thresholds]]
base = "USDT"
threshold = "2"

[[deviation_thresholds]]
base = "ATOM"
threshold = "1.5"

[[currency_pairs]]
base = "ATOM"
quote = "USDT"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[[currency_pairs]]
base = "NEBULA"
quote = "USDT"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[[currency_pairs]]
base = "USDT"
quote = "USD"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[account]
address = "nebula15nejfgcaanqpw25ru4arvfd0fwy6j8clccvwx4"
validator = "nebulavalcons14rjlkfzp56733j5l5nfk6fphjxymgf8mj04d5p"
chain_id = "nebula-local-testnet"

[keyring]
backend = "test"
dir = "/Users/username/.nebula"
pass = "keyringPassword"

[rpc]
tmrpc_endpoint = "http://localhost:26657"
grpc_endpoint = "localhost:9090"
rpc_timeout = "100ms"

[telemetry]
service-name = "price-feeder"
enabled = true
enable-hostname = true
enable-hostname-label = true
enable-service-label = true
prometheus-retention = 120
global-labels = [["chain-id", "nebula-local-testnet"]]
`)
	_, err = tmpFile.Write(content)
	require.NoError(t, err)

	cfg, err := config.ParseConfig(tmpFile.Name())
	require.NoError(t, err)

	require.Equal(t, "0.0.0.0:99999", cfg.Server.ListenAddr)
	require.Equal(t, "20s", cfg.Server.WriteTimeout)
	require.Equal(t, "20s", cfg.Server.ReadTimeout)
	require.True(t, cfg.Server.VerboseCORS)
	require.Len(t, cfg.CurrencyPairs, 3)
	require.Equal(t, "ATOM", cfg.CurrencyPairs[0].Base)
	require.Equal(t, "USDT", cfg.CurrencyPairs[0].Quote)
	require.Len(t, cfg.CurrencyPairs[0].Providers, 3)
	require.Equal(t, provider.ProviderKraken, cfg.CurrencyPairs[0].Providers[0])
	require.Equal(t, provider.ProviderBinance, cfg.CurrencyPairs[0].Providers[1])
	require.Equal(t, "2", cfg.Deviations[0].Threshold)
	require.Equal(t, "USDT", cfg.Deviations[0].Base)
	require.Equal(t, "1.5", cfg.Deviations[1].Threshold)
	require.Equal(t, "ATOM", cfg.Deviations[1].Base)
}

func TestParseConfig_Invalid_Deviations(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "price-feeder*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte(`
gas_adjustment = 1.5

[server]
listen_addr = "0.0.0.0:99999"
read_timeout = "20s"
verbose_cors = true
write_timeout = "20s"

[[deviation_thresholds]]
base = "USDT"
threshold = "4.0"

[[deviation_thresholds]]
base = "ATOM"
threshold = "1.5"

[[currency_pairs]]
base = "ATOM"
quote = "USDT"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[[currency_pairs]]
base = "NEBULA"
quote = "USDT"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[[currency_pairs]]
base = "USDT"
quote = "USD"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[account]
address = "nebula15nejfgcaanqpw25ru4arvfd0fwy6j8clccvwx4"
validator = "nebulavalcons14rjlkfzp56733j5l5nfk6fphjxymgf8mj04d5p"
chain_id = "nebula-local-testnet"

[keyring]
backend = "test"
dir = "/Users/username/.nebula"
pass = "keyringPassword"

[rpc]
tmrpc_endpoint = "http://localhost:26657"
grpc_endpoint = "localhost:9090"
rpc_timeout = "100ms"

[telemetry]
service-name = "price-feeder"
enabled = true
enable-hostname = true
enable-hostname_label = true
enable-service_label = true
prometheus-retention = 120
global-labels = [["chain-id", "nebula-local-testnet"]]
`)
	_, err = tmpFile.Write(content)
	require.NoError(t, err)

	_, err = config.ParseConfig(tmpFile.Name())
	require.Error(t, err)
}

func TestParseConfig_Env_Vars(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "price-feeder*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte(`
gas_adjustment = 1.5

[server]
listen_addr = "0.0.0.0:99999"
read_timeout = "20s"
verbose_cors = true
write_timeout = "20s"

[[currency_pairs]]
base = "ATOM"
quote = "USDT"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[[currency_pairs]]
base = "NEBULA"
quote = "USDT"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[[currency_pairs]]
base = "USDT"
quote = "USD"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[account]
address = "nebula15nejfgcaanqpw25ru4arvfd0fwy6j8clccvwx4"
validator = "nebulavalcons14rjlkfzp56733j5l5nfk6fphjxymgf8mj04d5p"
chain_id = "nebula-local-testnet"

[keyring]
backend = "test"
dir = "/Users/username/.nebula"
pass = "keyringPassword"

[rpc]
tmrpc_endpoint = "http://localhost:26657"
grpc_endpoint = "localhost:9090"
rpc_timeout = "100ms"

[telemetry]
service-name = "price-feeder"
enabled = true
enable-hostname = true
enable-hostname_label = true
enable-service_label = true
prometheus-retention = 120
global-labels = [["chain-id", "nebula-local-testnet"]]
`)
	_, err = tmpFile.Write(content)
	require.NoError(t, err)

	// Set env variables to overwrite config files
	os.Setenv("SERVER.LISTEN_ADDR", "0.0.0.0:888888")
	os.Setenv("SERVER.WRITE_TIMEOUT", "10s")
	os.Setenv("SERVER.READ_TIMEOUT", "10s")
	os.Setenv("SERVER.VERBOSE_CORS", "false")

	cfg, err := config.ParseConfig(tmpFile.Name())
	require.NoError(t, err)

	require.Equal(t, "0.0.0.0:888888", cfg.Server.ListenAddr)
	require.Equal(t, "10s", cfg.Server.WriteTimeout)
	require.Equal(t, "10s", cfg.Server.ReadTimeout)
	require.False(t, cfg.Server.VerboseCORS)
	require.Len(t, cfg.CurrencyPairs, 3)
	require.Equal(t, "ATOM", cfg.CurrencyPairs[0].Base)
	require.Equal(t, "USDT", cfg.CurrencyPairs[0].Quote)
	require.Len(t, cfg.CurrencyPairs[0].Providers, 3)
	require.Equal(t, provider.ProviderKraken, cfg.CurrencyPairs[0].Providers[0])
	require.Equal(t, provider.ProviderBinance, cfg.CurrencyPairs[0].Providers[1])
}

func TestCheckProviderMins_Valid(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "price-feeder*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte(`
gas_adjustment = 1.5

[server]
listen_addr = "0.0.0.0:99999"
read_timeout = "20s"
verbose_cors = true
write_timeout = "20s"

[[currency_pairs]]
base = "ATOM"
quote = "USDT"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[[currency_pairs]]
base = "NEBULA"
quote = "USDT"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[[currency_pairs]]
base = "USDT"
quote = "USD"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[account]
address = "nebula15nejfgcaanqpw25ru4arvfd0fwy6j8clccvwx4"
validator = "nebulavalcons14rjlkfzp56733j5l5nfk6fphjxymgf8mj04d5p"
chain_id = "nebula-local-testnet"

[keyring]
backend = "test"
dir = "/Users/username/.nebula"
pass = "keyringPassword"

[rpc]
tmrpc_endpoint = "http://localhost:26657"
grpc_endpoint = "localhost:9090"
rpc_timeout = "100ms"

[telemetry]
enabled = false
`)
	_, err = tmpFile.Write(content)
	require.NoError(t, err)

	cfg, err := config.ParseConfig(tmpFile.Name())
	require.NoError(t, err)

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	err = config.CheckProviderMins(context.TODO(), logger, cfg)
	require.NoError(t, err)
}

func TestCheckProviderMins_Invalid(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "price-feeder*.toml")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	content := []byte(`
gas_adjustment = 1.5

[server]
listen_addr = "0.0.0.0:99999"
read_timeout = "20s"
verbose_cors = true
write_timeout = "20s"

[[currency_pairs]]
base = "ATOM"
quote = "USDT"
providers = [
	"kraken",
	"binance",
]

[[currency_pairs]]
base = "NEBULA"
quote = "USDT"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[[currency_pairs]]
base = "USDT"
quote = "USD"
providers = [
	"kraken",
	"binance",
	"huobi"
]

[account]
address = "nebula15nejfgcaanqpw25ru4arvfd0fwy6j8clccvwx4"
validator = "nebulavalcons14rjlkfzp56733j5l5nfk6fphjxymgf8mj04d5p"
chain_id = "nebula-local-testnet"

[keyring]
backend = "test"
dir = "/Users/username/.nebula"
pass = "keyringPassword"

[rpc]
tmrpc_endpoint = "http://localhost:26657"
grpc_endpoint = "localhost:9090"
rpc_timeout = "100ms"

[telemetry]
enabled = false
`)
	_, err = tmpFile.Write(content)
	require.NoError(t, err)

	cfg, err := config.ParseConfig(tmpFile.Name())
	require.NoError(t, err)

	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.InfoLevel).With().Timestamp().Logger()
	err = config.CheckProviderMins(context.TODO(), logger, cfg)
	require.EqualError(t, err, "must have at least 3 providers for ATOM")
}
