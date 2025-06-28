# go-http-metrics

This project measures the time taken for TCP and TLS handshakes to a specified HTTP(S) endpoint and exports the metrics using OpenTelemetry.  It's designed to provide insights into network latency and performance.

## Usage

1.  **Install dependencies:** Ensure you have Go installed.  Then, run `go mod tidy` in the project directory to download necessary packages.

2.  **Run the application:** Execute the program with the target URL as a command-line argument:

    ```bash
    go run main.go https://www.example.com
    ```

    Replace `https://www.example.com` with the URL you wish to measure.  The program will continuously monitor the handshake times.

3. **Configuration:**

*   The measurement interval is currently hardcoded to 60 seconds. This can be adjusted in the `main.go` file.
*   The number of trials per measurement is also hardcoded to 10 in the `main.go` file.  You can modify this value to control the number of handshakes measured before reporting metrics.
*   The exporter is currently configured to send metrics to an OTLP gRPC endpoint.  You can modify this to use a different exporter, such as the stdout exporter (commented out in `main.go`), if you wish to see the metrics in your console instead.


## Metrics

The program exports the following metrics:

*   `tcp.handshake`:  The duration of the TCP handshake in seconds.
*   `tls.handshake`: The duration of the TLS handshake in seconds.

## OpenTelemetry

The project uses OpenTelemetry for metrics exporting.  Make sure you have a compatible OpenTelemetry collector or backend configured to receive the metrics.  The exporter is currently configured for the OTLP gRPC protocol. You will need to configure a backend to receive these metrics.
