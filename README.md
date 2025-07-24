# MQTT for golang

[![PkgGoDev](https://pkg.go.dev/badge/github.com/koron/go-mqtt)](https://pkg.go.dev/github.com/koron/go-mqtt)
[![Actions/Go](https://github.com/koron/go-mqtt/workflows/Go/badge.svg)](https://github.com/koron/go-mqtt/actions?query=workflow%3AGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/koron/go-mqtt)](https://goreportcard.com/report/github.com/koron/go-mqtt)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/koron/go-mqtt)

Yet another MQTT packages for golang.

This provides three MQTT related packages:

*   [packet](./packet) - MQTT packets encoder/decoder
*   [client](./client) - MQTT client library
*   [server](./server) - MQTT broker/server adapter

## Client

### How to connect with WebSocket

To connect MQTT server with WebSocket, use `ws://` scheme for `Addr`
field.

```go
clinet.Connect(client.Param{
    ID:   "wsclient-1234",
    Addr: "ws://localhost:8082/mqtt/over/websocket",
})
```

This will estimate `Origin` header to connect to WS server.
If you want to specify `Origin` set `Param.Options.WSOrigin` option field.

```go
clinet.Connect(client.Param{
    ID:   "wsclient-1234",
    Addr: "ws://localhost:8082/mqtt/over/websocket",
    Options: &client.Options{
        WSOrigin: "http://localhost:80/your/favorite/origin",
        // other fields are copied from client.DefaultOptions
        Version:      4,
        CleanSession: true,
        KeepAlive:    30,
    },
})
```

When you want to use secure WebSocket, try `wss://` scheme and
`Options.TLSConfig` field.

```go
clinet.Connect(client.Param{
    ID:   "wssclient-1234",
    Addr: "wss://localhost:8082/mqtt/over/websocket",
    Options: &client.Options{
        TLSConfig: &tls.Config{
            // your favorite TLS configurations.
        },
    },
}
```

## References

*   http://docs.oasis-open.org/mqtt/mqtt/v3.1.1/mqtt-v3.1.1.html
*   http://public.dhe.ibm.com/software/dw/webservices/ws-mqtt/mqtt-v3r1.html
