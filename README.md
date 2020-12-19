# Transporter go

[![Release](https://img.shields.io/badge/release-0.1.0-blue)](https://github.com/ednailson/transporter-go/releases)
![Coverage](https://img.shields.io/badge/coverage-100%25-success)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)
[![GoDoc](https://img.shields.io/badge/godoc-reference-9cf)](https://godoc.org/github.com/ednailson/transporter-go)


A library that creates server abstracting the transport layer protocols

Supported protocols

* [TCP](https://en.wikipedia.org/wiki/Transmission_Control_Protocol)
* [UDP](https://en.wikipedia.org/wiki/User_Datagram_Protocol)

Check out the **[CHANGELOG](CHANGELOG.md)**.

## Downloading

    go get github.com/ednailson/transporter-go

## Creating a new server

```go
server, _ := transporter.New("udp", "127.0.0.1", 3399)
```

It will create a **UDP** server on `127.0.0.1` at the `3399` port

The first parameter is the protocol, the library only supports `tcp`, `tcp4`, `tcp6`, `udp`, `udp4` or `udp6`.

## Reading messages

you need to create the handler function to handle the received requests.

After creating the server and the handler function, you just need to start listening to receive the messages.

```go
//Printing all the received messages
fn := func(conn connection.Connection) {
	log.Println(string(conn.Message()))
}

//Starting to listen
_ = server.Listen(fn)
```

It will print all the received messages on the server

## Closing the server

You just need to call the `Close()` function on the `server`

```go
_ = server.Close()
```

# Developer

[Júnior Vilas Boas](http://ednailson.github.io)
