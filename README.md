# matteosalvatoreg

CV website
Simple Website for use as a CV link

## Overview

`matteosalvatoreg` is a modern, professional CV/portfolio website built with Go. It provides a clean and professional interface to showcase skills, experience, and education.

## Features

- **Professional CV Display**: Complete CV information including work experience, skills, and education
- **Responsive Design**: Works well on desktop and mobile devices
- **RESTful API**: JSON endpoints for CV data
- **Health Check**: Built-in health monitoring endpoint
- **Security Focused**: Implements security best practices

## API Endpoints

### CV Information
- **GET /api/cv** - Returns complete CV information

### System Status
- **GET /health** - Health check endpoint

## Technology Stack

- **Backend**: Go (Gin framework)
- **API**: RESTful JSON
- **Security**: HTTP security headers, CORS
- **Deployment**: Docker ready

## Running

### Prerequisites

```bash
go 1.21+
```

### Building and Running

```bash
# Build the application
go build

# Run the application
./matteosalvatoreg

# Or using Docker
docker build -t matteosalvatoreg .
docker run -p 8081:8081 matteosalvatoreg
```

### Testing

```bash
# Run tests
go test -v ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Usage

### Local Development

```bash
# Start the server
./matteosalvatoreg

# Access CV information
curl http://localhost:8081/api/cv

# Check system health
curl http://localhost:8081/health
```

### API Examples

#### Get CV Information
```bash
curl -X GET http://localhost:8081/api/cv
```

#### Health Check
```bash
curl -X GET http://localhost:8081/health
```

## Deployment

### Docker Deployment

```Dockerfile
FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go build -o matteosalvatoreg main.go

EXPOSE 8081

CMD ["./matteosalvatoreg"]
```

### Kubernetes Deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: matteosalvatoreg
  labels:
    app: matteosalvatoreg
spec:
  replicas: 2
  selector:
    matchLabels:
      app: matteosalvatoreg
  template:
    metadata:
      labels:
        app: matteosalvatoreg
    spec:
      containers:
      - name: matteosalvatoreg
        image: matteosalvatoreg:latest
        ports:
        - containerPort: 8081
        env:
        - name: NODE_ENV
          value: "production"
```

## Development

### Project Structure

```
.
├── main.go          # Main application
├── main_test.go     # Tests
├── go.mod          # Module configuration
└── README.md       # Documentation
```

### Features Implemented

- **Professional CV Showcase**: Complete CV information display
- **Clean Architecture**: Separation of concerns
- **Security First**: Security headers and validation
- **API-First**: RESTful API design
- **Developer Experience**: Comprehensive testing and documentation

## License

MIT License

## Future Enhancements

- [ ] Add social media links
- [ ] Implement CV download functionality
- [ ] Add admin dashboard
- [ ] Integrate with HR systems
- [ ] Add multi-language support
- [ ] Implement analytics integration

## Contributing

1. Fork the repository
2. Create a feature branch
3. Implement your changes
4. Add tests
5. Submit a pull request

## Contact

For questions or suggestions, please open an issue or contact the maintainer.