# kgym

A comprehensive gym management platform for managing members, subscriptions, and gym providers.

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
| **Gateway** | ![Gateway Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/REPO_PLACEHOLDER/BRANCH_PLACEHOLDER/.github/badges/gateway.json&label=gateway) |
| **SSO** | ![SSO Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/REPO_PLACEHOLDER/BRANCH_PLACEHOLDER/.github/badges/sso.json&label=sso) |
| **User** | ![User Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/REPO_PLACEHOLDER/BRANCH_PLACEHOLDER/.github/badges/user.json&label=user) |
| **File** | ![File Coverage](https://img.shields.io/endpoint?url=https://raw.githubusercontent.com/REPO_PLACEHOLDER/BRANCH_PLACEHOLDER/.github/badges/file.json&label=file) |

## 📧 Contact

**Alexandr Rutkowski**
[kitanoyoru@icloud.com](mailto:kitanoyoru@icloud.com)
