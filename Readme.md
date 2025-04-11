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

