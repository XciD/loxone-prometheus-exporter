# Loxone Prometheus exporter

Usage:

```
go build -o exporter
./exporter -host loxone:8000 -user xcid -password test
```

Docker
```
docker run -it --name loxone-prometheus-exporter -p 8080:8080 xcid/loxone-prometheus-exporter -host loxone:8000 -user xcid -password test
```