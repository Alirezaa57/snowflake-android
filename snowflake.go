
package snowflakeandroid

import (
    "context"
    "io"
    "net"
    sflib "gitlab.torproject.org/tpo/anti-censorship/pluggable-transports/snowflake/v2/client/lib"
)

type Client struct {
    transport *sflib.Transport
    listener net.Listener
    cancel context.CancelFunc
}

func Start(port string) error {
    config := sflib.ClientConfig{
        BrokerURL: "https://snowflake-broker.torproject.net/",
        FrontDomains: []string{"www.google.com"},
        ICEAddresses: []string{"stun:stun.voipgate.com:3478"},
        Max: 3,
    }

    transport, err := sflib.NewSnowflakeClient(config)
    if err != nil { return err }

    ln, err := net.Listen("tcp", "127.0.0.1:"+port)
    if err != nil { return err }

    ctx,_ := context.WithCancel(context.Background())

    go func(){
        for {
            conn,err := ln.Accept()
            if err != nil { return }
            go handle(conn,transport)
        }
    }()

    <-ctx.Done()
    return nil
}

func handle(local net.Conn, t *sflib.Transport){
    defer local.Close()

    remote,err := t.Dial()
    if err != nil { return }
    defer remote.Close()

    go io.Copy(remote,local)
    io.Copy(local,remote)
}
