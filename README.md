# Glasnik

A web scraping tool with proxy server and registry capabilities.

## Components

- **Proxy Server**: Handles HTTP requests with authentication
- **Registry**: Manages and tracks available proxy servers
- **CLI Tool**: Command-line interface for all operations

## Building and Running

### Using Docker Compose (Recommended)

Start all services with a single command:

```bash
make docker-run
```

This will start:

- Registry server on port 8080
- Proxy server on port 8081
- MongoDB on port 27017

Stop all services:

```bash
make docker-stop
```

### Running Individual Components

#### Registry Server

Build and run the registry:

```bash
make docker-build-registry
make docker-run-registry
```

The registry will be available at `http://localhost:8080`

#### Proxy Server

Build and run the proxy:

```bash
make docker-build-proxy
make docker-run-proxy
```

The proxy will be available at `http://localhost:8081`

### Manual Setup (Without Docker)

1. Build the CLI tool:

```bash
make build
```

2. Start the registry:

```bash
./bin/glasnik registry --port 8080
```

3. Start a proxy server:

```bash
./bin/glasnik proxy \
  --host localhost \
  --port 8081 \
  --username myuser \
  --password mypass \
  --registry http://localhost:8080
```

## API Endpoints

### Registry Server

- `POST /register`: Register a new proxy server
  - Request body: ProxyInfo JSON object

### Proxy Server

- `GET /details`: Get proxy server details
  - Requires Basic Authentication
- All other paths: Proxy the request (requires Basic Authentication)

## Environment Variables

### Proxy Server

- `PROXY_USERNAME`: Username for proxy authentication
- `PROXY_PASSWORD`: Password for proxy authentication

## Example Usage

Make a request through the proxy:

```bash
curl -x http://localhost:8081 \
  -U myuser:mypass \
  http://example.com
```

Check proxy details:

```bash
curl -u myuser:mypass \
  http://localhost:8081/details
```

Original API endpoint:

```bash
POST http://localhost:3000/api/request-njuska
{
  "filter": "https://www.njuskalo.hr/iznajmljivanje-stanova?geo%5BlocationIds%5D=2691%2C2698",
  "token": "12345"
}
```
