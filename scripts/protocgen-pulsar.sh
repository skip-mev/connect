# this script is for generating protobuf files for the new google.golang.org/protobuf API

set -eo pipefail

protoc_install_gopulsar() {
  go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
}

protoc_install_gopulsar

echo "Cleaning API directory"
(cd api; find ./ -type f \( -iname \*.pulsar.go -o -iname \*.pb.go -o -iname \*.cosmos_orm.go -o -iname \*.pb.gw.go \) -delete; find . -empty -type d -delete; cd ..)

echo "Generating API module"
(cd proto; buf generate --template buf.gen.pulsar.yaml --exclude-path connect/service)

echo "fixing types.pulsar.go"
sed -i.bak 's|cosmossdk.io/api/connect/types/v2|github.com/skip-mev/connect/v2/api/connect/types/v2|g' ./api/connect/types/v2/currency_pair.pulsar.go && rm ./api/connect/types/v2/currency_pair.pulsar.go.bak
sed -i.bak 's|cosmossdk.io/api/connect/oracle/v2|github.com/skip-mev/connect/v2/api/connect/oracle/v2|g' ./api/connect/types/v2/currency_pair.pulsar.go && rm ./api/connect/types/v2/currency_pair.pulsar.go.bak

echo "fixing oracle.pulsar.go"
sed -i.bak 's|cosmossdk.io/api/connect/types/v2|github.com/skip-mev/connect/v2/api/connect/types/v2|g' ./api/connect/oracle/v2/query.pulsar.go && rm ./api/connect/oracle/v2/query.pulsar.go.bak
sed -i.bak 's|cosmossdk.io/api/connect/types/v2|github.com/skip-mev/connect/v2/api/connect/types/v2|g' ./api/connect/oracle/v2/tx.pulsar.go && rm ./api/connect/oracle/v2/tx.pulsar.go.bak
sed -i.bak 's|cosmossdk.io/api/connect/types/v2|github.com/skip-mev/connect/v2/api/connect/types/v2|g' ./api/connect/oracle/v2/genesis.pulsar.go && rm ./api/connect/oracle/v2/genesis.pulsar.go.bak

echo "fixing market.pulsar.go"
sed -i.bak 's|cosmossdk.io/api/connect/types/v2|github.com/skip-mev/connect/v2/api/connect/types/v2|g' ./api/connect/marketmap/v2/market.pulsar.go && rm ./api/connect/marketmap/v2/market.pulsar.go.bak
sed -i.bak 's|cosmossdk.io/api/connect/types/v2|github.com/skip-mev/connect/v2/api/connect/types/v2|g' ./api/connect/marketmap/v2/query.pulsar.go && rm ./api/connect/marketmap/v2/query.pulsar.go.bak
sed -i.bak 's|cosmossdk.io/api/connect/oracle/v2|github.com/skip-mev/connect/v2/api/connect/oracle/v2|g' ./api/connect/marketmap/v2/market.pulsar.go && rm ./api/connect/marketmap/v2/market.pulsar.go.bak
