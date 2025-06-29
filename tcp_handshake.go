package main

import (
	"net"
	"time"
)

// measureTCPHandshakeは指定したホストへのTCPハンドシェイクにかかった時間を計測し、
// 接続済みのnet.Conn、所要時間、エラーを返します。
func measureTCPHandshake(host string) (net.Conn, time.Duration, error) {
	start := time.Now()
	conn, err := net.Dial("tcp", host)
	duration := time.Since(start)
	return conn, duration, err
}
