package main

import (
	"fmt"
	"net"
	"time"
)

// measureTCPHandshakeは指定したホストへのTCPハンドシェイクにかかった時間を計測し、
// 接続済みのnet.Conn、所要時間、エラーを返します。
func measureTCPHandshake(host string) (net.Conn, time.Duration, error) {
	start := time.Now()
	conn, err := net.Dial("tcp", host)
	duration := time.Since(start)
	fmt.Println(conn.RemoteAddr().String())
	fmt.Printf("Trace start: %v, end: %v\n", start.String(), time.Now().String())
	return conn, duration, err
}
