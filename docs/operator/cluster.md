### cluster

Operations on a Temporal cluster

#### health

Check health of frontend service

**--address**="": host:port for Temporal frontend service

**--codec-auth**="": Authorization header to set for requests to Codec Server

**--codec-endpoint**="": Remote Codec Server Endpoint

**--color**="": when to use color: auto, always, never. (default: auto)

**--context-timeout**="": Optional timeout for context of RPC call in seconds (default: 5)

**--env**="": Env name to read the client environment variables from (default: default)

**--grpc-meta**="": gRPC metadata to send with requests. Format: key=value. Use valid JSON formats for value

**--namespace, -n**="": Temporal workflow namespace (default: default)

**--tls-ca-path**="": Path to server CA certificate

**--tls-cert-path**="": Path to x509 certificate

**--tls-disable-host-verification**: Disable tls host name verification (tls must be enabled)

**--tls-key-path**="": Path to private key

**--tls-server-name**="": Override for target server name

#### describe

Show information about the cluster

**--address**="": host:port for Temporal frontend service

**--codec-auth**="": Authorization header to set for requests to Codec Server

**--codec-endpoint**="": Remote Codec Server Endpoint

**--color**="": when to use color: auto, always, never. (default: auto)

**--context-timeout**="": Optional timeout for context of RPC call in seconds (default: 5)

**--env**="": Env name to read the client environment variables from (default: default)

**--fields**="": customize fields to print. Set to 'long' to automatically print more of main fields

**--grpc-meta**="": gRPC metadata to send with requests. Format: key=value. Use valid JSON formats for value

**--namespace, -n**="": Temporal workflow namespace (default: default)

**--output, -o**="": format output as: table, json, card. (default: table)

**--tls-ca-path**="": Path to server CA certificate

**--tls-cert-path**="": Path to x509 certificate

**--tls-disable-host-verification**: Disable tls host name verification (tls must be enabled)

**--tls-key-path**="": Path to private key

**--tls-server-name**="": Override for target server name

#### system

Show information about the system and capabilities

**--address**="": host:port for Temporal frontend service

**--codec-auth**="": Authorization header to set for requests to Codec Server

**--codec-endpoint**="": Remote Codec Server Endpoint

**--color**="": when to use color: auto, always, never. (default: auto)

**--context-timeout**="": Optional timeout for context of RPC call in seconds (default: 5)

**--env**="": Env name to read the client environment variables from (default: default)

**--fields**="": customize fields to print. Set to 'long' to automatically print more of main fields

**--grpc-meta**="": gRPC metadata to send with requests. Format: key=value. Use valid JSON formats for value

**--namespace, -n**="": Temporal workflow namespace (default: default)

**--output, -o**="": format output as: table, json, card. (default: table)

**--tls-ca-path**="": Path to server CA certificate

**--tls-cert-path**="": Path to x509 certificate

**--tls-disable-host-verification**: Disable tls host name verification (tls must be enabled)

**--tls-key-path**="": Path to private key

**--tls-server-name**="": Override for target server name

#### upsert

Add or update a remote cluster

**--address**="": host:port for Temporal frontend service

**--codec-auth**="": Authorization header to set for requests to Codec Server

**--codec-endpoint**="": Remote Codec Server Endpoint

**--color**="": when to use color: auto, always, never. (default: auto)

**--context-timeout**="": Optional timeout for context of RPC call in seconds (default: 5)

**--enable-connection**: Enable cross cluster connection

**--env**="": Env name to read the client environment variables from (default: default)

**--frontend-address**="": Frontend address of the remote cluster

**--grpc-meta**="": gRPC metadata to send with requests. Format: key=value. Use valid JSON formats for value

**--namespace, -n**="": Temporal workflow namespace (default: default)

**--tls-ca-path**="": Path to server CA certificate

**--tls-cert-path**="": Path to x509 certificate

**--tls-disable-host-verification**: Disable tls host name verification (tls must be enabled)

**--tls-key-path**="": Path to private key

**--tls-server-name**="": Override for target server name

#### list

List all remote clusters

**--address**="": host:port for Temporal frontend service

**--codec-auth**="": Authorization header to set for requests to Codec Server

**--codec-endpoint**="": Remote Codec Server Endpoint

**--color**="": when to use color: auto, always, never. (default: auto)

**--context-timeout**="": Optional timeout for context of RPC call in seconds (default: 5)

**--env**="": Env name to read the client environment variables from (default: default)

**--fields**="": customize fields to print. Set to 'long' to automatically print more of main fields

**--grpc-meta**="": gRPC metadata to send with requests. Format: key=value. Use valid JSON formats for value

**--limit**="": number of items to print (default: 0)

**--namespace, -n**="": Temporal workflow namespace (default: default)

**--no-pager, -P**: disable interactive pager

**--output, -o**="": format output as: table, json, card. (default: table)

**--pager**="": pager to use: less, more, favoritePager..

**--time-format**="": format time as: relative, iso, raw. (default: relative)

**--tls-ca-path**="": Path to server CA certificate

**--tls-cert-path**="": Path to x509 certificate

**--tls-disable-host-verification**: Disable tls host name verification (tls must be enabled)

**--tls-key-path**="": Path to private key

**--tls-server-name**="": Override for target server name

#### remove

Remove a remote cluster

**--address**="": host:port for Temporal frontend service

**--codec-auth**="": Authorization header to set for requests to Codec Server

**--codec-endpoint**="": Remote Codec Server Endpoint

**--color**="": when to use color: auto, always, never. (default: auto)

**--context-timeout**="": Optional timeout for context of RPC call in seconds (default: 5)

**--env**="": Env name to read the client environment variables from (default: default)

**--grpc-meta**="": gRPC metadata to send with requests. Format: key=value. Use valid JSON formats for value

**--name**="": Frontend address of the remote cluster

**--namespace, -n**="": Temporal workflow namespace (default: default)

**--tls-ca-path**="": Path to server CA certificate

**--tls-cert-path**="": Path to x509 certificate

**--tls-disable-host-verification**: Disable tls host name verification (tls must be enabled)

**--tls-key-path**="": Path to private key

**--tls-server-name**="": Override for target server name
