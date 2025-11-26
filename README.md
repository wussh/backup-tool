# Database Backup Tools

A comprehensive suite of database backup tools supporting PostgreSQL, MySQL, MariaDB, and MongoDB with multiple backup methods.

## 📦 Available Tools

### 1. **Go Interactive Tool** (Clean Architecture) - Recommended ⭐
Professional-grade interactive backup tool with Clean Architecture principles.

### 2. **Bash Scripts** (Flexible Automation)
Lightweight, configuration-based backup scripts for automation.

---

## 🚀 Go Interactive Tool (Clean Architecture)

A professional-grade database backup tool built with Clean Architecture principles in Go.

### 🏗️ Architecture Overview

This project follows **Clean Architecture** (also known as Hexagonal Architecture or Ports and Adapters), ensuring:

- **Independence of Frameworks**: Business logic doesn't depend on external libraries
- **Testability**: Business rules can be tested without UI, database, or external elements
- **Independence of UI**: Easy to change UI without changing business logic
- **Independence of Database**: Business rules aren't bound to backup mechanisms
- **Independence of External Agents**: Business rules don't know about the outside world

### Layer Structure

```
cmd/backup/              # Application entry point
└── main.go             # Dependency injection & wiring

internal/
├── domain/             # Enterprise Business Rules (Entities)
│   ├── entity.go       # Domain entities and value objects
│   ├── repository.go   # Repository interfaces (ports)
│   └── service.go      # Service interfaces (ports)
│
├── usecase/            # Application Business Rules
│   └── backup_usecase.go  # Orchestrates backup workflow
│
├── infrastructure/     # Frameworks & Drivers (Adapters)
│   └── backup_repository.go  # Docker/kubectl implementation
│
└── delivery/           # Interface Adapters
    └── cli/
        ├── config_service.go   # User input handling
        └── output_service.go   # Output formatting
```

## 📦 Directory Structure

```
.
├── cmd/
│   └── backup/
│       └── main.go                    # Application entry point
│
├── internal/
│   ├── domain/                        # Domain Layer (innermost)
│   │   ├── entity.go                 # Core entities
│   │   ├── repository.go             # Repository interface
│   │   └── service.go                # Service interfaces
│   │
│   ├── usecase/                       # Use Case Layer
│   │   └── backup_usecase.go         # Business logic
│   │
│   ├── infrastructure/                # Infrastructure Layer (outermost)
│   │   └── backup_repository.go      # External tool implementation
│   │
│   └── delivery/                      # Delivery Layer (outermost)
│       └── cli/
│           ├── config_service.go     # CLI input handler
│           └── output_service.go     # CLI output handler
│
├── go.mod
└── README.md
```

## 🎯 Clean Architecture Principles Applied

### 1. **Domain Layer** (Enterprise Business Rules)

**Location**: `internal/domain/`

**Responsibility**: Core business entities and interfaces

**Dependencies**: None (pure business logic)

**Files**:
- `entity.go`: Defines core entities (DatabaseConfig, BackupConfig, BackupResult)
- `repository.go`: Defines BackupRepository interface (port)
- `service.go`: Defines ConfigService and OutputService interfaces (ports)

**Example**:
```go
// Domain entity - no dependencies
type DatabaseConfig struct {
    Type      DatabaseType
    Host      string
    Database  string
    // ...
}

// Repository interface (port) - defines what we need, not how
type BackupRepository interface {
    BackupPostgres(config DatabaseConfig, ...) error
    BackupMySQL(config DatabaseConfig, ...) error
    // ...
}
```

### 2. **Use Case Layer** (Application Business Rules)

**Location**: `internal/usecase/`

**Responsibility**: Orchestrates business workflows

**Dependencies**: Only domain layer

**Files**:
- `backup_usecase.go`: Implements the backup workflow logic

**Example**:
```go
// Use case depends only on interfaces (dependency inversion)
type BackupUsecase struct {
    backupRepo    domain.BackupRepository
    configService domain.ConfigService
    outputService domain.OutputService
}

// Business logic is clean and testable
func (uc *BackupUsecase) ExecuteInteractiveBackup() error {
    // 1. Get configuration
    // 2. Execute backups
    // 3. Report results
}
```

### 3. **Infrastructure Layer** (Frameworks & Drivers)

**Location**: `internal/infrastructure/`

**Responsibility**: Implements external integrations

**Dependencies**: Domain layer (implements interfaces)

**Files**:
- `backup_repository.go`: Implements BackupRepository using Docker/kubectl

**Example**:
```go
// Adapter implementing the port
type BackupRepositoryImpl struct{}

// Implements domain.BackupRepository interface
func (r *BackupRepositoryImpl) BackupPostgres(...) error {
    // Docker/kubectl specific implementation
}
```

### 4. **Delivery Layer** (Interface Adapters)

**Location**: `internal/delivery/cli/`

**Responsibility**: Handles user interaction

**Dependencies**: Domain layer (implements interfaces)

**Files**:
- `config_service.go`: CLI-based configuration input
- `output_service.go`: CLI-based output formatting

**Example**:
```go
// Adapter implementing the port
type ConfigServiceImpl struct {
    reader *bufio.Reader
}

// Implements domain.ConfigService interface
func (s *ConfigServiceImpl) SelectBackupMethod() (domain.BackupMethod, error) {
    // CLI-specific input handling
}
```

### 5. **Main** (Composition Root)

**Location**: `cmd/backup/main.go`

**Responsibility**: Dependency injection and wiring

**Example**:
```go
func main() {
    // Dependency Injection (all dependencies resolved here)
    backupRepo := infrastructure.NewBackupRepository()
    configService := cli.NewConfigService()
    outputService := cli.NewOutputService()
    
    // Wire up use case
    backupUsecase := usecase.NewBackupUsecase(
        backupRepo,
        configService,
        outputService,
    )
    
    // Execute
    backupUsecase.ExecuteInteractiveBackup()
}
```

## 🚀 Building and Running the Go Tool

### Build the application
```bash
go build -o bin/backup ./cmd/backup
```

### Run directly
```bash
go run ./cmd/backup/main.go
```

### Run the built binary
```bash
./bin/backup
```

### Install globally
```bash
go install ./cmd/backup
```

### Interactive Flow Example

```
========================================
  Interactive Database Backup Tool
  Clean Architecture Edition
  Supports: PostgreSQL, MySQL, MariaDB, MongoDB
========================================

Select backup method:
  1. docker-run    (Use temporary container)
  2. docker-exec   (Exec into existing Docker container)
  3. kubectl-exec  (Exec into Kubernetes pod)

Enter choice [1-3]: 3

Kubernetes Namespace [default]: production

Select databases to backup:
  1. PostgreSQL
  2. MySQL
  3. MariaDB
  4. MongoDB
  5. All databases

Enter choices (comma-separated, e.g., 1,2,4): 1

=== Configuring POSTGRES ===
PostgreSQL Host [postgres]: prod-postgres
PostgreSQL User [postgres]: admin
Database Name [mydb]: production_db
PostgreSQL Password: ********
PostgreSQL Version [15]: 15
Pod Name [postgres-0]: postgres-primary-0

=== Configuration Summary ===
Backup Method: kubectl-exec
Timestamp: 2025-11-26 10:21:59
Backup Directory: backup
Kubernetes Namespace: production

Databases to backup:
  1. postgres - production_db (Host: prod-postgres)

Proceed with backup? (y/n): y

[POSTGRES] Starting backup...
  Method: kubectl-exec
  Host: prod-postgres
  Database: production_db
  Pod: postgres-primary-0
✓ Backup completed: backup/postgres/production_db_2025-11-26_10-22-01.sql (145M) [2.3s]

========================================
Backup Process Completed!
========================================

Results:
  Successful: 1

Backup files:
  ✓ postgres: backup/postgres/production_db_2025-11-26_10-22-01.sql (145M)
```

## 🧪 Testing Strategy

Clean Architecture makes testing much easier:

### Unit Tests (Domain Layer)
```go
// Test entities and value objects
func TestDatabaseType_IsValid(t *testing.T) {
    // Pure business logic testing
}
```

### Use Case Tests
```go
// Mock the dependencies
type MockBackupRepository struct {
    mock.Mock
}

func TestBackupUsecase_ExecuteInteractiveBackup(t *testing.T) {
    // Test business logic with mocks
    mockRepo := new(MockBackupRepository)
    mockConfig := new(MockConfigService)
    mockOutput := new(MockOutputService)
    
    usecase := NewBackupUsecase(mockRepo, mockConfig, mockOutput)
    // Test the workflow
}
```

### Integration Tests
```go
// Test with real implementations
func TestBackupRepository_BackupPostgres(t *testing.T) {
    repo := NewBackupRepository()
    // Test actual Docker commands
}
```

## 🔄 Dependency Flow

```
main.go (Composition Root)
    ↓
    ├─→ infrastructure.BackupRepository (adapter)
    ├─→ cli.ConfigService (adapter)
    ├─→ cli.OutputService (adapter)
    ↓
usecase.BackupUsecase (business logic)
    ↓ (depends on interfaces only)
    ├─→ domain.BackupRepository (interface)
    ├─→ domain.ConfigService (interface)
    └─→ domain.OutputService (interface)
```

**Key Principle**: Dependencies point inward. Domain has zero dependencies.

## 🎨 Design Patterns Used

1. **Repository Pattern**: Abstracts data access (BackupRepository)
2. **Service Pattern**: Encapsulates operations (ConfigService, OutputService)
3. **Dependency Injection**: All dependencies injected in main.go
4. **Interface Segregation**: Small, focused interfaces
5. **Single Responsibility**: Each layer has one reason to change

## ✨ Benefits of This Architecture

### 1. **Testability**
- Business logic can be tested without Docker/kubectl
- Mock implementations for all interfaces
- Fast unit tests without external dependencies

### 2. **Maintainability**
- Clear separation of concerns
- Easy to understand structure
- Changes isolated to specific layers

### 3. **Flexibility**
- Swap Docker for direct database connections
- Change from CLI to Web UI without touching business logic
- Add new backup methods without changing use cases

### 4. **Scalability**
- Add new databases by extending interfaces
- Parallel execution can be added in use case layer
- Easy to add features like scheduling, notifications

## 🔮 Future Enhancements

### Easy to Add:
1. **Web UI**: Add `internal/delivery/http/` without touching business logic
2. **REST API**: Add `internal/delivery/api/` alongside CLI
3. **Different Storage**: Add S3/GCS implementation of repository
4. **Scheduling**: Add scheduler in use case layer
5. **Monitoring**: Add observability in infrastructure layer
6. **Configuration Files**: Add config file parser in delivery layer

### Example: Adding Web UI
```go
// internal/delivery/http/handler.go
type BackupHandler struct {
    backupUsecase *usecase.BackupUsecase
}

func (h *BackupHandler) HandleBackup(w http.ResponseWriter, r *http.Request) {
    // Same use case, different delivery mechanism
    h.backupUsecase.ExecuteInteractiveBackup()
}
```

## 📚 References

- [The Clean Architecture by Uncle Bob](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/)
- [Go Project Layout](https://github.com/golang-standards/project-layout)

---

## 📜 Bash Scripts (Flexible Automation)

Lightweight, environment-based backup scripts for automated backups.

### 📁 Available Scripts

Located in `scripts/` directory:

1. **`backup-all-flexible.sh`** - Backup all configured databases
2. **`backup-postgres-flexible.sh`** - PostgreSQL only
3. **`backup-mysql-flexible.sh`** - MySQL only
4. **`backup-mariadb-flexible.sh`** - MariaDB only
5. **`backup-mongodb-flexible.sh`** - MongoDB only

### 🔧 Setup

1. **Copy the example environment file**:
```bash
cp .env.example .env
```

2. **Edit `.env` with your configuration**:
```bash
# Choose backup method
BACKUP_METHOD=kubectl-exec  # or docker-run, docker-exec

# For kubectl-exec
K8S_NAMESPACE=production
PG_POD=postgres-primary-0
MYSQL_POD=mysql-primary-0
MARIADB_POD=mariadb-primary-0
MONGO_POD=mongodb-primary-0

# Database credentials
PG_HOST=postgres
PG_USER=postgres
PG_DB=production_db
PG_PASS=secure_password

# ... (configure other databases)
```

3. **Make scripts executable**:
```bash
chmod +x scripts/*.sh
```

### 🚀 Usage

#### Backup All Databases
```bash
./scripts/backup-all-flexible.sh
```

#### Backup Specific Database
```bash
./scripts/backup-postgres-flexible.sh
./scripts/backup-mysql-flexible.sh
./scripts/backup-mariadb-flexible.sh
./scripts/backup-mongodb-flexible.sh
```

### 📋 Backup Methods

All scripts support three backup methods:

#### 1. **docker-run** (Default)
Uses temporary containers - no running container needed.
```bash
BACKUP_METHOD=docker-run
```

Best for:
- Remote databases
- No local database containers
- Clean, isolated backups

#### 2. **docker-exec**
Executes commands in existing Docker containers.
```bash
BACKUP_METHOD=docker-exec
PG_CONTAINER=test-postgres
MYSQL_CONTAINER=test-mysql
```

Best for:
- Local Docker databases
- Docker Compose setups
- Development environments

#### 3. **kubectl-exec**
Executes commands in Kubernetes pods.
```bash
BACKUP_METHOD=kubectl-exec
K8S_NAMESPACE=production
PG_POD=postgres-0
MYSQL_POD=mysql-0
```

Best for:
- Kubernetes deployments
- Production environments
- Cloud-native setups

### 📂 Backup Output Structure

```
backup/
├── postgres/
│   ├── mydb_2025-11-26_10-30-00.sql
│   └── mydb_2025-11-26_14-30-00.sql
├── mysql/
│   ├── mydb_2025-11-26_10-30-00.sql
│   └── mydb_2025-11-26_14-30-00.sql
├── mariadb/
│   └── mydb_2025-11-26_10-30-00.sql
└── mongodb/
    └── 2025-11-26_10-30-00/
        └── mydb/
            ├── collection1.bson
            └── collection1.metadata.json
```

### ⚙️ Environment Variables

#### Backup Method Configuration
```bash
BACKUP_METHOD=docker-run|docker-exec|kubectl-exec
BACKUP_TEMP_DIR=/tmp/db-backups  # Temp dir for docker-exec/kubectl-exec
```

#### Docker Container Names (for docker-exec)
```bash
PG_CONTAINER=test-postgres
MYSQL_CONTAINER=test-mysql
MARIADB_CONTAINER=test-mariadb
MONGO_CONTAINER=test-mongodb
```

#### Kubernetes Pod Names (for kubectl-exec)
```bash
K8S_NAMESPACE=default
PG_POD=postgres-0
MYSQL_POD=mysql-0
MARIADB_POD=mariadb-0
MONGO_POD=mongodb-0
```

#### Database Credentials
```bash
# PostgreSQL
PG_HOST=postgres
PG_USER=postgres
PG_DB=mydb
PG_PASS=password
PG_VERSION=15

# MySQL
MYSQL_HOST=mysql
MYSQL_USER=root
MYSQL_PASS=password
MYSQL_DB=mydb
MYSQL_VERSION=8

# MariaDB
MARIADB_HOST=mariadb
MARIADB_USER=root
MARIADB_PASS=password
MARIADB_DB=mydb
MARIADB_VERSION=11

# MongoDB
MONGO_HOST=mongodb
MONGO_DB=mydb
MONGO_VERSION=7
```

### 🔄 Automation with Cron

Add to crontab for scheduled backups:

```bash
# Daily backup at 2 AM
0 2 * * * cd /path/to/backup-tool && ./scripts/backup-all-flexible.sh

# Every 6 hours
0 */6 * * * cd /path/to/backup-tool && ./scripts/backup-all-flexible.sh

# Weekly backup (Sunday at 3 AM)
0 3 * * 0 cd /path/to/backup-tool && ./scripts/backup-all-flexible.sh
```

### 📊 Script Output Example

```bash
$ ./scripts/backup-all-flexible.sh

=========================================
  Flexible Database Backup Automation
  Method: kubectl-exec
  Timestamp: 2025-11-26_10-30-00
=========================================

Creating backup directories...
✓ Directories created

[PostgreSQL] Starting backup...
  Method: kubectl-exec
  Host: postgres
  Database: mydb
  User: postgres
  Pod: postgres-0
  Namespace: production
✓ PostgreSQL backup completed: backup/postgres/mydb_2025-11-26_10-30-00.sql (145M)

[MySQL] Starting backup...
  Method: kubectl-exec
  Host: mysql
  Database: mydb
  User: root
  Pod: mysql-0
  Namespace: production
✓ MySQL backup completed: backup/mysql/mydb_2025-11-26_10-30-00.sql (87M)

========================================
Backup process completed!
========================================

Backup method: kubectl-exec
Backup location: /home/user/backup-tool/backup/

Recent backups:
backup/postgres/mydb_2025-11-26_10-30-00.sql
backup/mysql/mydb_2025-11-26_10-30-00.sql
```

---

## 🐳 Testing with Docker Compose

A `docker-compose.yml` is provided for testing all database types locally.

### Start Test Databases
```bash
docker-compose up -d
```

### Check Status
```bash
docker-compose ps
```

### Test Backups
```bash
# Configure for docker-exec method
cat > .env << EOF
BACKUP_METHOD=docker-exec
PG_CONTAINER=test-postgres
MYSQL_CONTAINER=test-mysql
MARIADB_CONTAINER=test-mariadb
MONGO_CONTAINER=test-mongodb

PG_HOST=localhost
PG_USER=postgres
PG_DB=mydb
PG_PASS=password

MYSQL_HOST=localhost
MYSQL_USER=root
MYSQL_PASS=password
MYSQL_DB=mydb

MARIADB_HOST=localhost
MARIADB_USER=root
MARIADB_PASS=password
MARIADB_DB=mydb

MONGO_HOST=localhost
MONGO_DB=mydb
EOF

# Run backup
./scripts/backup-all-flexible.sh
```

### Stop Test Databases
```bash
docker-compose down
```

### Clean Up (including volumes)
```bash
docker-compose down -v
```

---

## 📊 Comparison: Go vs Bash

| Feature | Go Tool | Bash Scripts |
|---------|---------|--------------|
| **Setup** | Compile once | Edit .env file |
| **User Experience** | ✅ Interactive prompts | 📝 Pre-configured |
| **Configuration** | ✅ Step-by-step | 📝 Environment variables |
| **Namespace Support** | ✅ Prompted | 📝 K8S_NAMESPACE var |
| **Error Handling** | ✅ Comprehensive | ⚠️ Basic |
| **Flexibility** | ✅ High | ⚠️ Medium |
| **Automation** | ⚠️ Manual run | ✅ Cron-friendly |
| **Dependencies** | Go binary only | Bash + tools |
| **Cross-platform** | ✅ Yes | ⚠️ Unix-like only |
| **Code Quality** | ✅ Clean Architecture | 📝 Functional |
| **Testability** | ✅ Easy to mock | ⚠️ Limited |
| **Best for** | Interactive use | Automation/CI/CD |

### When to Use Which?

**Use Go Tool When:**
- You need interactive configuration
- Running manual backups
- Want guided setup process
- Need cross-platform support
- Prefer compiled binaries

**Use Bash Scripts When:**
- Setting up automated/scheduled backups
- Integrating with CI/CD pipelines
- Need minimal dependencies
- Have existing .env configuration
- Running on cron jobs

---

## 🤝 Contributing

### For Go Tool
When adding features, follow these principles:
1. Start with domain entities and interfaces
2. Implement business logic in use cases
3. Create adapters in infrastructure/delivery
4. Wire everything in main.go

### For Bash Scripts
1. Maintain compatibility with all three backup methods
2. Add error handling and validation
3. Keep output formatting consistent
4. Update .env.example with new variables
5. Keep dependencies pointing inward

---