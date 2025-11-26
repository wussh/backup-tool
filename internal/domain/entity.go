package domain

import "time"

// DatabaseType represents the type of database
type DatabaseType string

const (
	DatabaseTypePostgres DatabaseType = "postgres"
	DatabaseTypeMySQL    DatabaseType = "mysql"
	DatabaseTypeMariaDB  DatabaseType = "mariadb"
	DatabaseTypeMongoDB  DatabaseType = "mongodb"
)

// BackupMethod represents the method used for backup
type BackupMethod string

const (
	BackupMethodDockerRun   BackupMethod = "docker-run"
	BackupMethodDockerExec  BackupMethod = "docker-exec"
	BackupMethodKubectlExec BackupMethod = "kubectl-exec"
)

// DatabaseConfig holds configuration for a database
type DatabaseConfig struct {
	Type      DatabaseType
	Host      string
	Port      int
	User      string
	Password  string
	Database  string
	Version   string
	Container string // For docker-exec
	Pod       string // For kubectl-exec
}

// BackupConfig holds backup configuration
type BackupConfig struct {
	Method        BackupMethod
	Timestamp     time.Time
	BackupDir     string
	TempDir       string
	K8sNamespace  string
	Databases     []DatabaseConfig
}

// BackupResult represents the result of a backup operation
type BackupResult struct {
	DatabaseType DatabaseType
	Database     string
	Success      bool
	BackupPath   string
	Size         string
	Error        error
	Duration     time.Duration
}

// Validation methods
func (dt DatabaseType) IsValid() bool {
	switch dt {
	case DatabaseTypePostgres, DatabaseTypeMySQL, DatabaseTypeMariaDB, DatabaseTypeMongoDB:
		return true
	}
	return false
}

func (bm BackupMethod) IsValid() bool {
	switch bm {
	case BackupMethodDockerRun, BackupMethodDockerExec, BackupMethodKubectlExec:
		return true
	}
	return false
}

// String methods
func (dt DatabaseType) String() string {
	return string(dt)
}

func (bm BackupMethod) String() string {
	return string(bm)
}
