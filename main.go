package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"

	// "go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

var (
	interval   = 60 * time.Second
	trialCount = 10
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <arg>", os.Args[0])
		return
	}
	target := os.Args[1]
	target_url, err := url.Parse(target)
	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.Background()
	exporter, err := otlpmetricgrpc.New(ctx)
	// exporter, err := stdoutmetric.New()
	if err != nil {
		fmt.Printf("Failed to create exporter: %v\n", err)
		return
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(target_url.String())))

	provider := metricsdk.NewMeterProvider(metricsdk.WithResource(res), metricsdk.WithReader(metricsdk.NewPeriodicReader(exporter, metricsdk.WithInterval(interval))))
	otel.SetMeterProvider(provider)
	defer provider.Shutdown(ctx)

	metricMeter := provider.Meter("synthetic", metric.WithInstrumentationAttributes(
		attribute.String("target", target_url.String()),
	))
	tcpHandshakeMeter, err := metricMeter.Float64Histogram("tcp.handshake", metric.WithUnit("s"))
	tlsHandshakeMeter, err := metricMeter.Float64Histogram("tls.handshake", metric.WithUnit("s"))
	if err != nil {
		fmt.Printf("Failed to create metric: %v\n", err)
		return
	}

	host := target_url.Host
	if host == "" {
		host = target_url.Path
	}
	if _, _, err := net.SplitHostPort(host); err != nil {
		// ポートが指定されていない場合は443を追加
		host = host + ":443"
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Printf("Measuring %s", host)
			for range trialCount {
				measure(ctx, host, target_url, tcpHandshakeMeter, tlsHandshakeMeter)
			}
		}
	}

}

func measure(ctx context.Context, host string, target_url *url.URL, tcpHandshakeMeter, tlsHandshakeMeter metric.Float64Histogram) {
	tcpConn, tcpHandshakeDuration, err := measureTCPHandshake(host)

	if err != nil {
		fmt.Printf("TCP handshake failed: %v\n", err)
		return
	}
	defer tcpConn.Close()

	conf := &tls.Config{
		ServerName:             target_url.Hostname(),
		SessionTicketsDisabled: true,
	}

	// TLSハンドシェイク計測
	client, tlsHandshakeDuration, err := measureTLSHandshake(tcpConn, conf)
	defer client.Close()

	if err != nil {
		fmt.Printf("TLS handshake failed: %v\n", err)
		return
	}

	localip, _, _ := net.SplitHostPort(tcpConn.LocalAddr().String())

	tcpHandshakeMeter.Record(ctx, tcpHandshakeDuration.Seconds(), metric.WithAttributes(attribute.String("sourceIp", localip)))
	tlsHandshakeMeter.Record(ctx, tlsHandshakeDuration.Seconds(), metric.WithAttributes(attribute.String("sourceIp", localip)))

}
