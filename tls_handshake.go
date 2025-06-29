package main

import (
	"crypto/tls"
	"net"
	"time"
)

// measureTLSHandshakeは、既存のTCP接続tcpConn上でTLSハンドシェイクを行い、
// TLS接続（*tls.Conn）、ハンドシェイクにかかった時間、エラーを返します。
func measureTLSHandshake(tcpConn net.Conn, conf *tls.Config) (*tls.Conn, time.Duration, error) {
	client := tls.Client(tcpConn, conf)
	start := time.Now()
	err := client.Handshake()
	duration := time.Since(start)
	return client, duration, err
}
