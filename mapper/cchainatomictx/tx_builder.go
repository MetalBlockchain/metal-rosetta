package cchainatomictx

import (
	"errors"
	"fmt"

	"github.com/MetalBlockchain/metalgo/codec"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils"
	"github.com/MetalBlockchain/metalgo/utils/crypto/secp256k1"
	"github.com/MetalBlockchain/metalgo/utils/formatting/address"
	"github.com/MetalBlockchain/metalgo/vms/components/avax"
	"github.com/MetalBlockchain/metalgo/vms/secp256k1fx"
	"github.com/MetalBlockchain/coreth/plugin/evm"
	"github.com/coinbase/rosetta-sdk-go/parser"
	"github.com/coinbase/rosetta-sdk-go/types"
	ethcommon "github.com/ethereum/go-ethereum/common"

	"github.com/MetalBlockchain/metal-rosetta/mapper"
)

var errMissingCoinIdentifier = errors.New("input operation does not have coin identifier")

// BuildTx constructs an evm tx based on the provided operation type, Rosetta matches and metadata
// This method is only used during construction.
func BuildTx(opType string, matches []*parser.Match, metadata Metadata, codec codec.Manager, avaxAssetID ids.ID) (*evm.Tx, []*types.AccountIdentifier, error) {
	switch opType {
	case mapper.OpExport:
		return buildExportTx(matches, metadata, codec, avaxAssetID)
	case mapper.OpImport:
		return buildImportTx(matches, metadata, codec, avaxAssetID)
	default:
		return nil, nil, fmt.Errorf("unsupported atomic operation type %s", opType)
	}
}

// [buildExportTx] returns a duly initialized tx if it does not err
func buildExportTx(
	matches []*parser.Match,
	metadata Metadata,
	codec codec.Manager,
	avaxAssetID ids.ID,
) (*evm.Tx, []*types.AccountIdentifier, error) {
	ins, signers := buildIns(matches, metadata, avaxAssetID)

	exportedOutputs, err := buildExportedOutputs(matches, codec, avaxAssetID)
	if err != nil {
		return nil, nil, err
	}

	tx := &evm.Tx{UnsignedAtomicTx: &evm.UnsignedExportTx{
		NetworkID:        metadata.NetworkID,
		BlockchainID:     metadata.CChainID,
		DestinationChain: *metadata.DestinationChainID,
		Ins:              ins,
		ExportedOutputs:  exportedOutputs,
	}}
	return tx, signers, tx.Sign(codec, nil)
}

// [buildImportTx] returns a duly initialized tx if it does not err
func buildImportTx(
	matches []*parser.Match,
	metadata Metadata,
	codec codec.Manager,
	avaxAssetID ids.ID,
) (*evm.Tx, []*types.AccountIdentifier, error) {
	importedInputs, signers, err := buildImportedInputs(matches, avaxAssetID)
	if err != nil {
		return nil, nil, err
	}

	outs := buildOuts(matches, avaxAssetID)

	tx := &evm.Tx{UnsignedAtomicTx: &evm.UnsignedImportTx{
		NetworkID:      metadata.NetworkID,
		BlockchainID:   metadata.CChainID,
		SourceChain:    *metadata.SourceChainID,
		ImportedInputs: importedInputs,
		Outs:           outs,
	}}
	return tx, signers, tx.Sign(codec, nil)
}

func buildIns(matches []*parser.Match, metadata Metadata, avaxAssetID ids.ID) ([]evm.EVMInput, []*types.AccountIdentifier) {
	inputMatch := matches[0]

	ins := []evm.EVMInput{}
	signers := []*types.AccountIdentifier{}
	for i, op := range inputMatch.Operations {
		ins = append(ins, evm.EVMInput{
			Address: ethcommon.HexToAddress(op.Account.Address),
			Amount:  inputMatch.Amounts[i].Uint64(),
			AssetID: avaxAssetID,
			Nonce:   metadata.Nonce,
		})
		signers = append(signers, op.Account)
	}

	// we do not use the signers as signing is performed externally to Rosetta
	// instead we are using a dummy array with the same length as ins
	evmSigners := make([][]*secp256k1.PrivateKey, len(ins))
	evm.SortEVMInputsAndSigners(ins, evmSigners)

	return ins, signers
}

func buildImportedInputs(matches []*parser.Match, avaxAssetID ids.ID) ([]*avax.TransferableInput, []*types.AccountIdentifier, error) {
	inputMatch := matches[0]

	importedInputs := []*avax.TransferableInput{}
	signers := []*types.AccountIdentifier{}
	for i, op := range inputMatch.Operations {
		if op.CoinChange == nil || op.CoinChange.CoinIdentifier == nil {
			return nil, nil, errMissingCoinIdentifier
		}
		utxoID, err := mapper.DecodeUTXOID(op.CoinChange.CoinIdentifier.Identifier)
		if err != nil {
			return nil, nil, err
		}

		importedInputs = append(importedInputs, &avax.TransferableInput{
			UTXOID: *utxoID,
			Asset:  avax.Asset{ID: avaxAssetID},
			In: &secp256k1fx.TransferInput{
				Amt: inputMatch.Amounts[i].Uint64(),
				Input: secp256k1fx.Input{
					SigIndices: []uint32{0},
				},
			},
		})
		signers = append(signers, op.Account)
	}
	utils.Sort(importedInputs)

	return importedInputs, signers, nil
}

func buildOuts(matches []*parser.Match, avaxAssetID ids.ID) []evm.EVMOutput {
	outputMatch := matches[1]

	outs := []evm.EVMOutput{}
	for i, op := range outputMatch.Operations {
		outs = append(outs, evm.EVMOutput{
			Address: ethcommon.HexToAddress(op.Account.Address),
			Amount:  outputMatch.Amounts[i].Uint64(),
			AssetID: avaxAssetID,
		})
	}

	evm.SortEVMOutputs(outs)

	return outs
}

func buildExportedOutputs(matches []*parser.Match, codec codec.Manager, avaxAssetID ids.ID) ([]*avax.TransferableOutput, error) {
	outputMatch := matches[1]

	outs := []*avax.TransferableOutput{}
	for i, op := range outputMatch.Operations {
		destinationAddress, err := address.ParseToID(op.Account.Address)
		if err != nil {
			return nil, err
		}

		outs = append(outs, &avax.TransferableOutput{
			Asset: avax.Asset{ID: avaxAssetID},
			Out: &secp256k1fx.TransferOutput{
				Amt: outputMatch.Amounts[i].Uint64(),
				OutputOwners: secp256k1fx.OutputOwners{
					Locktime:  0,
					Threshold: 1,
					Addrs:     []ids.ShortID{destinationAddress},
				},
			},
		})
	}

	avax.SortTransferableOutputs(outs, codec)

	return outs, nil
}
