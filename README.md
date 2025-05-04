# 🔀 Go Load Balancer

A lightweight load balancer built in Go with support for:
- Round-robin algorithm
- Dockerized backends
- Health checks
- JSON-based backend configuration

## 🔧 Features

- Round-robin routing
- Dockerized setup using Docker Compose
- Configurable backends from `config/backends.json`
- Easily extensible (TLS, sticky sessions, weighted routing)

## 🚀 Run

```bash
make docker-up
