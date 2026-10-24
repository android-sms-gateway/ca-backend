# 🔐 Certificate Authority

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stars][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![License][license-shield]][license-url]

A private Certificate Authority for the SMSGate ecosystem that issues certificates valid for one year for private IP addresses, eliminating self-signed certificate management. See [Certificate Authority](https://docs.sms-gate.app/services/ca/) for the full guide.

## 📖 About

Public CAs cannot validate private IP addresses, so local SMSGate deployments need their own CA. This service accepts certificate signing requests via a REST API, validates that each request targets exactly one private IP address (RFC 1918 for IPv4, RFC 4193 `fc00::/7` for IPv6), signs certificates valid for 365 days, and stores request state in Redis.

You don't need to install the CA locally to use it. You can use the [API](#api-method) or the [CLI](#cli-method) to issue a certificate with the project's CA at [ca.sms-gate.app](https://ca.sms-gate.app).

## 📚 Table of Contents

- [🔐 Certificate Authority](#-certificate-authority)
  - [📖 About](#-about)
  - [📚 Table of Contents](#-table-of-contents)
  - [💻 Usage](#-usage)
    - [Method Comparison](#method-comparison)
    - [CLI Method](#cli-method)
    - [API Method](#api-method)
  - [⭐ Features](#-features)
  - [📦 Prerequisites](#-prerequisites)
  - [🚀 Getting Started](#-getting-started)
  - [⚙️ Configuration](#️-configuration)
  - [📦 Deployment](#-deployment)
  - [🔌 API Overview](#-api-overview)
  - [📚 Documentation](#-documentation)
  - [🤝 Contributing](#-contributing)
  - [📄 License](#-license)

## 💻 Usage

### Method Comparison

| Feature         | CLI Method 🖥️ | API Method 🌐      |
| --------------- | ------------ | ----------------- |
| Difficulty      | ⭐ Easy       | ⭐⭐ Medium         |
| Customization   | ❌ No         | ✅ Available       |
| Automation      | ✅ Full       | ❌ Manual          |
| Recommended For | Most users ✅ | CI/CD pipelines 🤖 |

### CLI Method

You can use the [SMSGate CLI](https://github.com/android-sms-gateway/cli/releases/latest) to issue a certificate.

1. 📥 **Generate Certificate**
    ```bash
    # Generate webhook certificate
    ./smsgate-ca webhooks --out=server.crt --keyout=server.key 192.168.1.10
    ```

2. 🔐 **Install Certificates**
    ```bash
    # Nginx example
    ssl_certificate /path/to/server.crt;
    ssl_certificate_key /path/to/server.key;
    ```

### API Method

1. 🔑 **Generate Key Pair**
    ```bash
    openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out server.key
    ```

2. 📝 **Create Config**
    ```ini
    # server.cnf
    [req]
    distinguished_name = req_distinguished_name
    x509_extensions = v3_req
    prompt = no
    
    [req_distinguished_name]
    CN = 192.168.1.10  # replace with your private IP
    
    [v3_req]
    keyUsage = nonRepudiation, digitalSignature, keyEncipherment
    extendedKeyUsage = serverAuth
    subjectAltName = @alt_names
    
    [alt_names]
    IP.0 = 192.168.1.10
    ```

3. 📋 **Generate CSR**
    ```bash
    openssl req -new -key server.key -out server.csr -extensions v3_req \
      -config ./server.cnf
    ```

4. 📨 **Submit CSR**
    ```sh
    jq -Rs '{content: .}' < server.csr | \
    curl -sSf -X POST \
      -H "Content-Type: application/json" \
      -d @- \
      https://ca.sms-gate.app/api/v1/csr
    ```

    You will receive a Request ID in the response.

5. 🕒 **Check Status**
    ```bash
    curl https://ca.sms-gate.app/api/v1/csr/REQ_12345 # replace with your Request ID
    ```

6. 📥 **Save Certificate**  
    When the request is approved, the certificate content will be provided in the `certificate` field of the response. Save the certificate content to the file `server.crt`.

7. 🔐 **Install Certificate**  
    Install the `server.crt` and `server.key` (from step 1) files to your server.

Full API documentation is available [here](https://ca.sms-gate.app/docs/index.html).

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

> **Security:** the published port exposes an unauthenticated certificate-signing endpoint. Restrict it to trusted clients or place the service behind an authenticated gateway.

```bash
docker compose up -d
```

Or with Docker directly — the service and Redis must share a user-defined network so the `redis` hostname resolves:

```bash
docker network create ca-net

docker run -d --name ca-redis --network ca-net \
  -v ca-redis-data:/data \
  --network-alias redis redis:7

docker run -d --name ca-backend --network ca-net \
  -p 3000:3000 \
  -e HTTP__ADDRESS=:3000 \
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
| GET    | `/metrics`         | Prometheus metrics               |

## 📚 Documentation

- [Certificate Authority](https://docs.sms-gate.app/services/ca/)
- [Central docs](https://docs.sms-gate.app/)

## 🤝 Contributing

Contributions are welcome. Open an issue or pull request; follow the repo's `make fmt` / `make lint` / `make test` checks and the [contributing guide](https://docs.sms-gate.app/contributing/).

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
[license-url]: https://github.com/android-sms-gateway/ca-backend/blob/master/LICENSE
