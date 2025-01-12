# llm-go-blueprint

## Overview
The codebase is organized into the following main directories:

```
.
├── cmd/
│   └── main.go           # Application entry point
├── pkg/
│   ├── app/             # Application setup and configuration
│   ├── handlers/        # HTTP handlers for API endpoints
│   ├── llm/             # LLM model adapters and interfaces
│   ├── middleware/      # HTTP middleware components
│   ├── routes/          # HTTP route definitions
│   └── service/         # Business logic and service layer
└── README.md
```

### cmd/
- **main.go**: Entry point of the application
  - Initializes and configures the HTTP server
  - Sets up graceful shutdown handling
  - Bootstraps application dependencies

### pkg/app/
- **app.go**: Application setup and HTTP server configuration
  - Configures routing
  - Applies middleware
  - Creates the main HTTP handler chain

### pkg/llm/
Core LLM (Large Language Model) integration functionality
- **adapter.go**: Interface definitions for LLM adapters
- **factory.go**: Factory pattern implementation for creating LLM instances

Key features:
- Adapter pattern for different LLM implementations
- Factory pattern for model instantiation
- Easy integration of new models
- Testable design with mock implementations

### pkg/middleware/
HTTP middleware components
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
