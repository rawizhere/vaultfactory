# VaultFactory

Secure data storage and synchronization system.

## Features

- Secure storage for various data types (passwords, text, files, bank cards)
- JWT authentication with refresh tokens
- AES-256-GCM data encryption
- Cross-device synchronization
- CLI interface
- HTTP REST API

## Architecture

- **Server** - HTTP server for data storage and synchronization
- **Client** - CLI application for server interaction

## Project Structure

```
vaultfactory/
├── cmd/                 # Main applications
│   ├── server/          # Server application
│   └── client/          # Client application
├── configs/             # Configuration files
├── internal/            # Private application code
│   ├── server/          # Server logic
│   ├── client/          # Client logic
│   └── shared/          # Shared components
└── scripts/             # SQL migrations and scripts
```

## Tech Stack

- Go 1.25
- PostgreSQL
- JWT authentication
- AES-256-GCM encryption
- Argon2 password hashing
