package cchainatomictx

import (
	"github.com/MetalBlockchain/metalgo/utils/formatting/address"
	"github.com/coinbase/rosetta-sdk-go/types"

	"github.com/MetalBlockchain/metal-rosetta/constants"
	"github.com/MetalBlockchain/metal-rosetta/mapper"
)

// IsCChainBech32Address checks whether a given account identifier contains a C-chain Bech32 type address
func IsCChainBech32Address(accountIdentifier *types.AccountIdentifier) bool {
	if chainID, _, _, err := address.Parse(accountIdentifier.Address); err == nil {
		return chainID == constants.CChain.String()
	}
	return false
}

// IsAtomicOpType determines whether a given C-chain operation is an atomic one
func IsAtomicOpType(t string) bool {
	return t == mapper.OpExport || t == mapper.OpImport
}
