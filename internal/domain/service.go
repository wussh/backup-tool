package domain

// ConfigService defines the interface for configuration operations
type ConfigService interface {
	// SelectBackupMethod prompts user to select backup method
	SelectBackupMethod() (BackupMethod, error)
	
	// SelectDatabases prompts user to select databases to backup
	SelectDatabases() ([]DatabaseType, error)
	
	// GetKubernetesNamespace prompts user for Kubernetes namespace
	GetKubernetesNamespace() (string, error)
	
	// ConfigureDatabase prompts user to configure a specific database
	ConfigureDatabase(dbType DatabaseType, method BackupMethod) (DatabaseConfig, error)
	
	// ConfirmBackup asks user to confirm backup operation
	ConfirmBackup(config BackupConfig) (bool, error)
}

// OutputService defines the interface for output operations
type OutputService interface {
	// PrintHeader prints the application header
	PrintHeader()
	
	// PrintConfigSummary prints the backup configuration summary
	PrintConfigSummary(config BackupConfig)
	
	// PrintBackupStart prints backup start message
	PrintBackupStart(dbType DatabaseType, config DatabaseConfig, method BackupMethod)
	
	// PrintBackupResult prints backup result
	PrintBackupResult(result BackupResult)
	
	// PrintSummary prints final summary
	PrintSummary(results []BackupResult)
	
	// PrintError prints an error message
	PrintError(message string)
	
	// PrintSuccess prints a success message
	PrintSuccess(message string)
}
