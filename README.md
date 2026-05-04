# Go + OpenTelemetry Logging

Simple Go app demonstrating OpenTelemetry SDK logging with Splunk.

## Quick Test

### 1. Start Splunk + OTel Collector

```bash
podman-compose up -d
```

Wait ~30 seconds for Splunk to start.

### 2. Run the App

```bash
go run main.go
```

You should see logs being sent to the collector.

### 3. View in Splunk

Open http://localhost:8000
- Username: `admin`
- Password: `Admin123!`

Search query:
```
index=* source="my-app"
```

### 4. Stop

```bash
podman-compose down
```

## How it Works

```
Go App (SDK) → OTel Collector (localhost:4317) → Splunk HEC (port 8088)
```

The app uses OpenTelemetry SDK to send structured logs directly to the collector, which batches and forwards them to Splunk.

## Files

- `main.go` - Go application with OpenTelemetry logging
- `otel-config.yaml` - OTel Collector configuration
- `podman-compose.yml` - Local Splunk + Collector setup
- `go.mod` / `go.sum` - Go dependencies

## Requirements

- Go 1.21+
- Podman or Docker
