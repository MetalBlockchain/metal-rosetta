#!/bin/bash

export METAL_NETWORK=${METAL_NETWORK:-testnet}
export METAL_CHAIN=${METAL_CHAIN:-43113}
export METAL_MODE=${METAL_MODE:-online}
export AVALANCHE_GENESIS_HASH=${AVALANCHE_GENESIS_HASH:-"0x31ced5b9beb7f8782b014660da0cb18cc409f121f408186886e1ca3e8eeca96b"}
export AVALANCHEGO_RPC_BASE_URL=${AVALANCHEGO_RPC_BASE_URL:-"http://localhost:9650"}
export AVALANCHEGO_INDEXER_BASE_URL=${AVALANCHEGO_INDEXER_BASE_URL:-$AVALANCHEGO_RPC_BASE_URL}

cat <<EOF > /app/avalanchego-config.json
{
  "network-id": "$METAL_NETWORK",
  "http-host": "0.0.0.0",
  "api-keystore-enabled": false,
  "api-admin-enabled": false,
  "api-ipcs-enabled": false,
  "api-keystore-enabled": false,
  "db-dir": "/data",
  "chain-config-dir": "/app/configs/chains",
  "index-enabled": true,
  "network-require-validator-to-connect": true
}
EOF

mkdir -p /app/configs/chains/C

cat <<EOF > /app/configs/chains/C/config.json
{
  "state-sync-enabled": false,
  "snowman-api-enabled": false,
  "coreth-admin-api-enabled": false,
  "rpc-gas-cap": 2500000000,
  "rpc-tx-fee-cap": 100,
  "eth-apis": ["internal-public-eth","internal-public-blockchain","internal-public-transaction-pool","internal-public-tx-pool","internal-public-debug","internal-private-debug","debug-tracer","web3","public-eth","public-eth-filter","public-debug","private-debug","net"],
  "pruning-enabled": false
}
EOF

cat <<EOF > /app/rosetta-config.json
{
  "mode": "$METAL_MODE",
  "rpc_base_url": "$AVALANCHEGO_RPC_BASE_URL",
  "indexer_base_url": "$AVALANCHEGO_INDEXER_BASE_URL",
  "listen_addr": "0.0.0.0:8080",
  "network_id": 1,
  "network_name": "$METAL_NETWORK",
  "chain_id": $METAL_CHAIN,
  "genesis_block_hash": "$AVALANCHE_GENESIS_HASH"
}
EOF

# Execute a custom command instead of default on
if [ -n "$@" ]; then
  exec $@
fi

exec /app/rosetta-runner \
  -mode $METAL_MODE \
  -avalanche-bin /app/avalanchego \
  -avalanche-config /app/avalanchego-config.json \
  -rosetta-bin /app/rosetta-server \
  -rosetta-config rosetta-config.json
