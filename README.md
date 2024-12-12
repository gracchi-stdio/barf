# Barf App

## Design Patterns & Architecture

### Patterns Used:

#### Repository Pattern

- Abstracts data persistence
- Separate interfaces from implementations
- Makes testing easier with mock repositories


#### Service Layer Pattern

- Contains business logic
- Coordinates between multiple repositories
- Handles transactions and complex operations


#### Dependency Injection

- Loose coupling between components
- Better testability
- Uses constructor injection


#### Factory Pattern

- For database connections
- For creating service instances
- For repository instantiation


### Project Structure
bookmanager/
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Configuration management
│   ├── domain/
│   │   ├── book.go          # Domain models
│   │   └── inventory.go
│   ├── repository/
│   │   ├── interfaces.go     # Repository interfaces
│   │   └── postgres/
│   │       ├── book.go
│   │       └── inventory.go
│   ├── service/
│   │   ├── book.go
│   │   └── inventory.go
│   └── handler/
│       └── http/
│           ├── book.go
│           └── inventory.go
├── pkg/
│   ├── validator/
│   │   └── validator.go      # Common validation utilities
│   └── logger/
│       └── logger.go         # Logging utilities
└── web/
└── src/
├── components/
│   ├── BookList.tsx
│   └── InventoryGrid.tsx
└── pages/
├── Books.tsx
└── Inventory.tsx