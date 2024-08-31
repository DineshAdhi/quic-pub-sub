
# Simple Pub/Sub in QUIC

### 1. Configure Certificates

Run `sh cert.sh`  to configure Self-Signed Certificates and add it to your Trust Store

### 2. Start Relay

`go run relay/relay.go`

By default, relay runs on Port `5000`, you can configure by passing the flag like this,

`go run relay/relay.go -port=<PORT>`

### 2. Start Publisher

`go run pub/pub.go -relay=<RELAY_ADDRESS> -pubkey=<PUBKEY>`

`RELAY_ADDRESS` is the Address where the Relay is hosted. By default its, `localhost:5000`

`PUBKEY` is your publishing key, by default this is configured as `testdata`

### 3. Start Subscriber

`go run sub/sub.go -relay=<RELAY_ADDRESS> -pubkey=<PUBKEY>`

`RELAYPORT` is optional unless it is configured differently in Relay.

`PUBKEY` is the publishing key configured by the Publisher.

# Notes
1. Initial handshake between the Relay and Pub / Sub is carried out in Bi-directional QUIC Stream. Here is where  the PubKey is shared.
2. The Actual Packets are send via QUIC Datagrams which are encrypted.
3. Publisher pushes the data (a sequentially increasing unsigned integer) to the relay at an interval of 1 second.
4. Subscribers can subscribe to the data using the PubKey.
