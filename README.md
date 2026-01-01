# kgym

A comprehensive gym management platform for managing members, subscriptions, and gym providers.

## 🚀 Quick Start

### Prerequisites

- **Go 1.25** or higher
- **Docker** (for running databases and services locally)
- **Make** (for running build commands)

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd kgym
   ```

2. **Install development tools**
   ```bash
   make tools-install
   ```
   This installs all required tools (protoc, mockgen, golangci-lint, etc.) in the `bin/` directory.

   **Add tools to your PATH:**
   ```bash
   export PATH="$PWD/bin:$PATH"
   ```
   To make this permanent, add the above line to your shell profile (e.g., `~/.bashrc`, `~/.zshrc`, or `~/.config/fish/config.fish`).

3. **Initialize contract dependencies**
   ```bash
   make contracts-deps
   ```

4. **Generate Protocol Buffer code**
   ```bash
   make contracts-protobuf-gen-go
   ```

5. **Install Go module dependencies**
   ```bash
   make gomod-all
   ```
   This runs `go mod tidy` for all services to ensure dependencies are up to date.

## 🛠️ Tech Stack

### Core
- **Language**: Go 1.25
- **Architecture**: Microservices with gRPC
- **API Gateway**: gRPC-Gateway (REST/HTTP)

### Databases & Storage
- **CockroachDB** (via pgx/v5) - Primary relational database
- **Redis** (go-redis/v9) - Caching and session storage
- **MinIO** - Object storage for files

### Communication & APIs
- **gRPC** - Inter-service communication
- **Protocol Buffers** - API contracts and serialization
- **gRPC-Gateway** - REST API gateway

### Observability
- **Prometheus** - Metrics collection
- **OpenTelemetry** - Distributed tracing

## 📊 Test Coverage

| Service | Coverage |
|---------|----------|
| **Gateway** | ![Gateway Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/kitanoyoru/kgym/master/.github/badges/gateway.json&label=gateway) |
| **SSO** | ![SSO Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/kitanoyoru/kgym/master/.github/badges/sso.json&label=sso) |
| **User** | ![User Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/kitanoyoru/kgym/master/.github/badges/user.json&label=user) |
| **File** | ![File Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/kitanoyoru/kgym/master/.github/badges/file.json&label=file) |

## 📧 Contact

For questions, suggestions, or collaboration opportunities, please reach out:

- **Name**: Alexandr Rutkowski
- **Email**: [kitanoyoru@icloud.com](mailto:kitanoyoru@icloud.com)

## 📄 License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.
