# llm-go-blueprint

## Overview
The codebase is organized into the following main directories:

```
.
├── cmd/
│   └── main.go           # Application entry point
├── files/
│   └── config.yaml       # Default configuration file
├── pkg/
│   ├── app/              # Application setup and configuration
│   ├── handlers/         # HTTP handlers for API endpoints
│   ├── llm/             
│   │   └── model/        # LLM model adapters and interfaces
│   │   └── validation/   # Validation functionality for LLM responses and requests
│   │       └── schemas/  # CUE schema definitions
│   ├── middleware/       # HTTP middleware components
│   ├── routes/           # HTTP route definitions
│   └── service/          # Business logic and service layer
│   └── run/              # Application bootstrapping and configuration
└── README.md
```

### cmd/
- **main.go**: Entry point of the application
  - Initializes and configures the HTTP server
  - Sets up graceful shutdown handling
  - Bootstraps application dependencies

### files/
- **config.yaml**: Default configuration file

### pkg/app/
- **app.go**: Application setup and HTTP server configuration
  - Configures routing
  - Applies middleware
  - Creates the main HTTP handler chain

### pkg/llm/model/
Core LLM (Large Language Model) integration functionality
- **adapter.go**: Interface definitions for LLM adapters
- **factory.go**: Factory pattern implementation for creating LLM instances

Key features:
- Adapter pattern for different LLM implementations
- Factory pattern for model instantiation
- Easy integration of new models
- Testable design with mock implementations

### pkg/llm/validation/
Core validation functionality for LLM responses and requests.
- **validation.go**: Interface definitions for validators
- **response.go**: Response schema validator implementation
- **schemas/**: CUE schema definitions
  - Defines response formats
  - Enforces type safety

Key features:
- Schema-based validation
- CUE to OpenAPI conversion
- Strict validation rules
- Multiple schema support
- Multiple validator support
- Extensible validator interface

### pkg/llm/validation/schemas/
Example CUE schema definitions for response validation
- **animalResponse.cue**: Animal response schema
- **personResponse.cue**: Person response schema

### pkg/middleware/
HTTP middleware components.
- **logging.go**: Request logging middleware
  - Logs HTTP method, path, duration, and status code
  - Provides request tracing and monitoring

### pkg/routes/
Listing all routes and maps the entire API surface.
- **routes.go**: HTTP route definitions
  - Defines API endpoints
  - Maps URLs to handlers
  - Groups related endpoints

### pkg/handlers/
Handlers: HTTP concerns (request parsing, validation, response writing)
- **handlers.go**: HTTP handlers for API endpoints
  - Handles request processing
  - Returns responses
  - Maps request data to service methods

### pkg/service/
Business logic and service layer. The service layer remains consistent regardless of the underlying model.
- **query.go**: Query service implementation
<!-- 
TODO: add additional service functionality
- **embedding.go**: Embedding service implementation
- **health.go**: Health check service -->

Features:
- Clear separation from transport layer
- Model-agnostic business logic
- Reusable service components

### pkg/run/
Application bootstrapping and configuration management
- **config.go**: Configuration structure and loading logic
- **run.go**: Main application setup and coordination

Features:
- Configuration management (YAML, environment variables)
- Graceful shutdown handling
- Application lifecycle coordination
- Command-line flag parsing

## Configuration

The application uses a YAML configuration file located in `files/config.yaml`.

The configuration file can be overridden using the `--config` flag.

```
# Use default config path
go run cmd/main.go

# Use custom config file
go run cmd/main.go --config ./custom-path/config.yaml

# Override with environment variable
PORT="9090" go run cmd/main.go
```

> **_NOTE:_**  Config loads with precedence: env vars > config file > defaults.

