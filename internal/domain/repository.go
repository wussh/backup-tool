package domain

// BackupRepository defines the interface for backup operations
type BackupRepository interface {
	// BackupPostgres performs a PostgreSQL backup
	BackupPostgres(config DatabaseConfig, method BackupMethod, backupPath, namespace string) error
	
	// BackupMySQL performs a MySQL backup
	BackupMySQL(config DatabaseConfig, method BackupMethod, backupPath, namespace string) error
	
	// BackupMariaDB performs a MariaDB backup
	BackupMariaDB(config DatabaseConfig, method BackupMethod, backupPath, namespace string) error
	
	// BackupMongoDB performs a MongoDB backup
	BackupMongoDB(config DatabaseConfig, method BackupMethod, backupPath, namespace, tempDir string) error
	
	// GetFileSize returns the size of a file or directory
	GetFileSize(path string, isDirectory bool) (string, error)
}
