# carzone
> Car Management System  A comprehensive car management system from scratch using:

- GoLang for backend 🖥️
- PostgreSQL for database 🗄️
- Docker Compose for containerization 🐳
- Grafana for monitoring 📊
- JWT for authentication 🔒
- Telemetry for data collection 📡


# OpenTelemetery Go Packages 
Open Telemetery provides a set of APIs, libraries, agents, and instrumentation to enable observability in your application.

## Installation
To install the necessary packages, use the following commands:

```bash
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp
go get go.opentelemetry.io/otel/sdk
```

# Prometheus Go Packages
## Installation
To install the necessary packages, use the following commands:

```bash
go get github.com/prometheus/client_golang/prometheus/promhttp
go get github.com/prometheus/client_golang/prometheus/
```

### Grafana error solution: 
```
Error reading Prometheus: Post "http://localhost:9090/api/v1/query": dial tcp 127.0.0.1:9090: connect: connection refused
```
- Solution: https://github.com/grafana/grafana/issues/46434

- Update Trip Request
d1e2f3a4-b5c6-7d8e-9f0a-b1c2d3e4f5a6

{
  "car_id": "9b9437c4-3ed1-45a5-b240-0fe3e24e0e4e",
  "description": "Eldoret to Mombasa Route",
  "distance_km": 20,
  "driver_id": "b2c3d4e5-f6c3-4e6d-ac3f-6d3d3d3d3d3d",
  "end_location": "Mombasa",
  "end_time": "2025-01-27T15:04:05Z",
  "fuel_consumed_liters": 2,
  "start_location": "Eldoret",
  "start_time": "2025-01-27T15:06:00Z",
  "status": "Completed"
}