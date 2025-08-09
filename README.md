<a id="readme-top"></a>

<!-- PROJECT SHIELDS -->
[![Build Status](https://img.shields.io/github/actions/workflow/status/android-sms-gateway/ca-backend/go.yml?branch=master&style=for-the-badge)](https://github.com/android-sms-gateway/ca-backend/actions)
[![Go Version](https://img.shields.io/github/go-mod/go-version/android-sms-gateway/ca-backend?style=for-the-badge)](go.mod)
[![License](https://img.shields.io/github/license/android-sms-gateway/ca-backend.svg?style=for-the-badge)](LICENSE)

<br />
<div align="center">
  <h3 align="center">🔒 Android SMS Gateway CA</h3>

  <p align="center">
    Private Certificate Authority for Secure Local Communications
    <br />
    <a href="https://ca.sms-gate.app/docs/index.html"><strong>Explore the API docs »</strong></a>
    <br />
    <br />
    <a href="https://github.com/android-sms-gateway/ca-backend/issues/new?labels=bug">Report Bug</a>
    ·
    <a href="https://github.com/android-sms-gateway/ca-backend/issues/new?labels=enhancement">Request Feature</a>
  </p>
</div>

<!-- TABLE OF CONTENTS -->
- [📖 About The Project](#-about-the-project)
  - [🛠️ Built With](#️-built-with)
- [🚀 Getting Started](#-getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [💻 Usage](#-usage)
  - [Method Comparison](#method-comparison)
  - [CLI Method](#cli-method)
  - [API Method](#api-method)
- [⚠️ Limitations](#️-limitations)
- [🚨 Migration Guide](#-migration-guide)
- [❓ FAQ](#-faq)
- [🤝 Contributing](#-contributing)
- [📄 License](#-license)


<!-- ABOUT THE PROJECT -->
## 📖 About The Project

This private Certificate Authority simplifies secure communications within local networks while maintaining security standards. By operating its own [Certificate Authority (CA)](https://en.wikipedia.org/wiki/Certificate_authority), the project eliminates common security pitfalls associated with self-signed certificates and manual certificate management.

> **Important** Security Value Proposition
> - **🌍 Solves private IP validation** - Public CAs cannot validate private IP addresses
> - **⚠️ Reduces security risks** - Eliminates manual certificate installation on client devices

The CA enforces strict security boundaries through multiple layers:

1. **Private IP Enforcement** - All issued certificates validated against RFC 1918 address ranges
2. **Key Management** - CA private key loaded securely (PEM/PKCS#8); certificates parsed using x509
3. **Request Validation** - CSRs validated to ensure SAN entries are private IPs (RFC 1918)

### 🛠️ Built With

- [![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
- [![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)
- [![Make](https://img.shields.io/badge/make-4C8A43?style=for-the-badge&logo=gnu-make&logoColor=white)](https://www.gnu.org/software/make/)

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- GETTING STARTED -->
## 🚀 Getting Started

### Prerequisites

- Go 1.24.1+ (for building from source)
- Docker (optional, for containerized deployment)
- OpenSSL (for manual key/CSR generation)
- curl and jq (for API examples)

### Installation

You don't need to install the CA locally to use it. You can use the [API](#api-method) or the [CLI](#cli-method) to issue a certificate with the project's CA at [ca.sms-gate.app](https://ca.sms-gate.app).

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE EXAMPLES -->
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

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LIMITATIONS -->
## ⚠️ Limitations

The Certificate Authority service has the following limitations:

- 🔐 Only issues certificates for private IP ranges:
    ```text
    10.0.0.0/8
    172.16.0.0/12
    192.168.0.0/16
    ```
- ⏳ Certificate validity: 1 year
- 📛 Maximum 1 `POST` request per minute

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- MIGRATION GUIDE -->
## 🚨 Migration Guide

Self-signed certificates will be deprecated after v2.0 release. It is recommended to use the project's CA instead.

Migration checklist:
- [ ] Replace self-signed certs before v2.0 release
- [ ] Update automation scripts to use CLI tool or API
- [ ] Rotate certificates every 1 year

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- FAQ -->
## ❓ FAQ

**Why don't I need to install CA on devices?**  
The root CA certificate is embedded in the SMSGate app (v1.31+).  
Note: other clients (browsers, third‑party services) that do not embed this CA will not trust these certificates unless you install the CA in their trust store.

**Certificate issuance failed?**  
Ensure your IP matches private ranges and hasn't exceeded quota

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- CONTRIBUTING -->
## 🤝 Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- LICENSE -->
## 📄 License

Distributed under the Apache-2.0 License. See [`LICENSE`](LICENSE) for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>
