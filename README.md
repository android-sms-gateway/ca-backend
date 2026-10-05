# 🔐 Certificate Authority

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stars][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]

A private Certificate Authority for the SMSGate ecosystem that issues short-lived certificates for private IP addresses, eliminating self-signed certificate management. See [Certificate Authority](https://docs.sms-gate.app/services/ca/) for the full guide.

## 📖 About

Public CAs cannot validate private IP addresses, so local SMSGate deployments need their own CA. This service accepts certificate signing requests via a REST API, validates that each request targets exactly one private IP address (RFC 1918), signs certificates valid for 365 days, and stores request state in Redis.

## 📚 Table of Contents

- [📖 About](#-about)
- [📚 Table of Contents](#-table-of-contents)
- [⭐ Features](#-features)
- [📦 Prerequisites](#-prerequisites)
- [🚀 Getting Started](#-getting-started)
- [⚙️ Configuration](#️-configuration)
- [📦 Deployment](#-deployment)
- [🔌 API Overview](#-api-overview)
- [📚 Documentation](#-documentation)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)

## ⭐ Features

- REST API for CSR submission and status tracking
- Strict validation: exactly one IP, CN matches the IP, private range only
- Asynchronous signing queue with Redis-backed state
- Certificates valid for 365 days
- Swagger/OpenAPI docs at `/docs`

## 📦 Prerequisites

- Go 1.25+ for building from source
- Redis (default `redis://localhost:6379/0`)
- CA certificate and private key (PEM, PKCS#8)

## 🚀 Getting Started

```bash
redis-server --daemonize yes
make build
CSR__CA_CERT_PATH=./ca/ca.crt \
CSR__CA_KEY_PATH=./ca/ca.key \
./bin/ca-backend
```

Submit a CSR for a private IP:

```bash
curl -sSf -X POST http://127.0.0.1:3000/api/v1/csr \
  -H 'Content-Type: application/json' \
  -d '{"type":"webhook","content":"<PEM CSR>"}'
```

## ⚙️ Configuration

Configuration is loaded from environment variables (double-underscore nesting) or an optional YAML file via `CONFIG_PATH`. See [.env.example](.env.example).

| Variable                  | Default                    | Description                      |
| ------------------------- | -------------------------- | -------------------------------- |
| `HTTP__ADDRESS`           | `127.0.0.1:3000`           | HTTP listen address              |
| `HTTP__PROXY_HEADER`      | *(empty)*                  | Trusted proxy header             |
| `HTTP__PROXIES`           | *(empty)*                  | Comma-separated trusted proxies  |
| `API__CORS_ALLOW_ORIGINS` | *(empty)*                  | CORS allowed origins (empty=off) |
| `STORAGE__URL`            | `redis://localhost:6379/0` | Redis connection URL             |
| `CSR__TTL`                | `24h`                      | CSR request lifetime             |
| `CSR__CA_CERT_PATH`       | *(required)*               | CA certificate file path         |
| `CSR__CA_KEY_PATH`        | *(required)*               | CA private key file path         |
| `CONFIG_PATH`             | *(empty)*                  | Optional YAML config file path   |

Full reference: [Certificate Authority](https://docs.sms-gate.app/services/ca/).

## 📦 Deployment

```bash
docker compose up -d
```

Or with Docker directly:

```bash
docker run -d --name ca-backend \
  -p 3000:3000 \
  -e STORAGE__URL=redis://redis:6379/0 \
  -e CSR__CA_CERT_PATH=/etc/ca/ca.crt \
  -e CSR__CA_KEY_PATH=/etc/ca/ca.key \
  -v ./ca:/etc/ca:ro \
  ghcr.io/android-sms-gateway/ca-backend:latest
```

See [compose.yaml](compose.yaml) for the full service definition (includes Redis).

## 🔌 API Overview

| Method | Path               | Description                      |
| ------ | ------------------ | -------------------------------- |
| POST   | `/api/v1/csr`      | Submit a CSR (202 Accepted)      |
| GET    | `/api/v1/csr/{id}` | Check CSR status and certificate |
| GET    | `/docs`            | Swagger/OpenAPI documentation    |
| GET    | `/health`          | Health probes                    |
| GET    | `/metrics`         | Prometheus metrics               |

## 📚 Documentation

- [Certificate Authority](https://docs.sms-gate.app/services/ca/)
- [Central docs](https://docs.sms-gate.app/)

## 🤝 Contributing

Contributions are welcome. Open an issue or pull request; follow the repo's `make fmt` / `make lint` / `make test` checks and the [README style guide](https://docs.sms-gate.app/).

## 📄 License

Apache-2.0. See [LICENSE](LICENSE).

<!-- Reference-style badge URLs: style=for-the-badge is mandatory -->
[contributors-shield]: https://img.shields.io/github/contributors/android-sms-gateway/ca-backend?style=for-the-badge
[contributors-url]: https://github.com/android-sms-gateway/ca-backend/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/android-sms-gateway/ca-backend?style=for-the-badge
[forks-url]: https://github.com/android-sms-gateway/ca-backend/network/members
[stars-shield]: https://img.shields.io/github/stars/android-sms-gateway/ca-backend?style=for-the-badge
[stars-url]: https://github.com/android-sms-gateway/ca-backend/stargazers
[issues-shield]: https://img.shields.io/github/issues/android-sms-gateway/ca-backend?style=for-the-badge
[issues-url]: https://github.com/android-sms-gateway/ca-backend/issues
[license-shield]: https://img.shields.io/github/license/android-sms-gateway/ca-backend?style=for-the-badge
[license-url]: https://github.com/android-sms-gateway/ca-backend/blob/main/LICENSE
