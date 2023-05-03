package client

import (
	"context"

	"github.com/MetalBlockchain/metalgo/api/info"
	"github.com/MetalBlockchain/metalgo/ids"
	"github.com/MetalBlockchain/metalgo/utils/rpc"
)

// InfoClient collects all Avalanchego info.Client methods common
// to Rosetta Clients
type InfoClient interface {
	GetBlockchainID(context.Context, string, ...rpc.Option) (ids.ID, error)
	IsBootstrapped(context.Context, string, ...rpc.Option) (bool, error)
	Peers(context.Context, ...rpc.Option) ([]info.Peer, error)
}
