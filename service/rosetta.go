package service

import (
	"fmt"

	"github.com/MetalBlockchain/metalgo/version"
)

var NodeVersion = fmt.Sprintf(
	"%d.%d.%d",
	version.Current.Major,
	version.Current.Minor,
	version.Current.Patch,
)

const (
	MiddlewareVersion = "0.1.30"
	BlockchainName    = "Metal"
)
