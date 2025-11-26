package cli

import (
	"fmt"
	"strings"

	"github.com/wush/db-backup-tool/internal/domain"
)

const (
	colorRed    = "\033[0;31m"
	colorGreen  = "\033[0;32m"
	colorYellow = "\033[1;33m"
	colorBlue   = "\033[0;34m"
	colorCyan   = "\033[0;36m"
	colorReset  = "\033[0m"
)

// OutputServiceImpl implements domain.OutputService
type OutputServiceImpl struct{}

// NewOutputService creates a new output service
func NewOutputService() domain.OutputService {
	return &OutputServiceImpl{}
}

// PrintHeader prints the application header
func (s *OutputServiceImpl) PrintHeader() {
	fmt.Println(colorBlue + "========================================")
	fmt.Println("  Interactive Database Backup Tool")
	fmt.Println("  Clean Architecture Edition")
	fmt.Println("  Supports: PostgreSQL, MySQL, MariaDB, MongoDB")
	fmt.Println("========================================" + colorReset)
	fmt.Println()
}

// PrintConfigSummary prints the backup configuration summary
func (s *OutputServiceImpl) PrintConfigSummary(config domain.BackupConfig) {
	fmt.Printf("\n%s=== Configuration Summary ===%s\n", colorCyan, colorReset)
	fmt.Printf("Backup Method: %s\n", config.Method)
	fmt.Printf("Timestamp: %s\n", config.Timestamp.Format("2006-01-02 15:04:05"))
	fmt.Printf("Backup Directory: %s\n", config.BackupDir)
	
	if config.Method == domain.BackupMethodKubectlExec {
		fmt.Printf("Kubernetes Namespace: %s\n", config.K8sNamespace)
	}
	
	fmt.Printf("\nDatabases to backup:\n")
	for i, db := range config.Databases {
		fmt.Printf("  %d. %s - %s (Host: %s)\n", i+1, db.Type, db.Database, db.Host)
	}
}

// PrintBackupStart prints backup start message
func (s *OutputServiceImpl) PrintBackupStart(dbType domain.DatabaseType, config domain.DatabaseConfig, method domain.BackupMethod) {
	fmt.Printf("%s[%s] Starting backup...%s\n", colorBlue, strings.ToUpper(dbType.String()), colorReset)
	fmt.Printf("  Method: %s\n", method)
	fmt.Printf("  Host: %s\n", config.Host)
	fmt.Printf("  Database: %s\n", config.Database)
	
	if method == domain.BackupMethodDockerExec {
		fmt.Printf("  Container: %s\n", config.Container)
	} else if method == domain.BackupMethodKubectlExec {
		fmt.Printf("  Pod: %s\n", config.Pod)
	}
}

// PrintBackupResult prints backup result
func (s *OutputServiceImpl) PrintBackupResult(result domain.BackupResult) {
	if result.Success {
		fmt.Printf("%s✓ Backup completed: %s (%s) [%s]%s\n\n",
			colorGreen, result.BackupPath, result.Size, result.Duration, colorReset)
	} else {
		fmt.Printf("%s✗ Backup failed: %v [%s]%s\n\n",
			colorRed, result.Error, result.Duration, colorReset)
	}
}

// PrintSummary prints final summary
func (s *OutputServiceImpl) PrintSummary(results []domain.BackupResult) {
	fmt.Printf("\n%s========================================%s\n", colorBlue, colorReset)
	fmt.Printf("%sBackup Process Completed!%s\n", colorGreen, colorReset)
	fmt.Printf("%s========================================%s\n", colorBlue, colorReset)
	
	successCount := 0
	failureCount := 0
	
	for _, result := range results {
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
	}
	
	fmt.Printf("\nResults:\n")
	fmt.Printf("  %sSuccessful: %d%s\n", colorGreen, successCount, colorReset)
	if failureCount > 0 {
		fmt.Printf("  %sFailed: %d%s\n", colorRed, failureCount, colorReset)
	}
	
	fmt.Println("\nBackup files:")
	for _, result := range results {
		if result.Success {
			fmt.Printf("  %s✓%s %s: %s (%s)\n",
				colorGreen, colorReset, result.DatabaseType, result.BackupPath, result.Size)
		} else {
			fmt.Printf("  %s✗%s %s: %v\n",
				colorRed, colorReset, result.DatabaseType, result.Error)
		}
	}
	fmt.Println()
}

// PrintError prints an error message
func (s *OutputServiceImpl) PrintError(message string) {
	fmt.Printf("%s✗ Error: %s%s\n", colorRed, message, colorReset)
}

// PrintSuccess prints a success message
func (s *OutputServiceImpl) PrintSuccess(message string) {
	fmt.Printf("%s✓ %s%s\n", colorGreen, message, colorReset)
}
