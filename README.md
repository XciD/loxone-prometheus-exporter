# Loxone Prometheus exporter

[![version](https://img.shields.io/badge/status-beta-orange.svg)](https://github.com/XciD/loxone-prometheus-exporter)
[![Build Status](https://travis-ci.org/XciD/loxone-prometheus-exporter.svg?branch=master)](https://travis-ci.org/XciD/loxone-prometheus-exporter)
[![Go Report Card](https://goreportcard.com/badge/github.com/XciD/loxone-prometheus-exporter)](https://goreportcard.com/report/github.com/XciD/loxone-prometheus-exporter)
[![codecov](https://codecov.io/gh/XciD/loxone-prometheus-exporter/branch/master/graph/badge.svg)](https://codecov.io/gh/XciD/loxone-prometheus-exporter)
[![Pulls](https://img.shields.io/docker/pulls/xcid/loxone-prometheus-exporter.svg)](https://hub.docker.com/r/xcid/loxone-prometheus-exporter)
[![Layers](https://shields.beevelop.com/docker/image/layers/xcid/loxone-prometheus-exporter/latest.svg)](https://hub.docker.com/r/xcid/loxone-prometheus-exporter)
[![Size](https://shields.beevelop.com/docker/image/image-size/xcid/loxone-prometheus-exporter/latest.svg)](https://hub.docker.com/r/xcid/loxone-prometheus-exporter)


Usage:

```
go build -o exporter
./exporter --host loxone:8000 --user xcid --password test
```

Docker
```
docker run -it --name loxone-prometheus-exporter -p 8080:8080 xcid/loxone-prometheus-exporter --host loxone:8000 --user xcid --password test
```