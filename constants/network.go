package constants

import (
	"github.com/MetalBlockchain/coreth/params"
	"github.com/MetalBlockchain/metalgo/utils/constants"
)

const (
	MainnetChainID = 381931
	MainnetAssetID = "FvwEAhmxKfeiG8SnEvq42hc6whRyY3EFYAvebMqDNDGCgxN5Z"
	MainnetNetwork = constants.MainnetName

	TahoeChainID = 381932
	TahoeAssetID = "2QpCJwPk3nzi1VqJEuaFA44WM2UUzraBXQyH6jMGLTLQhqoe4n"
	TahoeNetwork = constants.TahoeName
)

var (
	MainnetAP5Activation = params.MetalMainnetChainConfig.ApricotPhase5BlockTimestamp
	TahoeAP5Activation   = params.MetalTahoeChainConfig.ApricotPhase5BlockTimestamp
)
