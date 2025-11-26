package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/wush/db-backup-tool/internal/domain"
)

// ConfigServiceImpl implements domain.ConfigService
type ConfigServiceImpl struct {
	reader *bufio.Reader
}

// NewConfigService creates a new config service
func NewConfigService() domain.ConfigService {
	return &ConfigServiceImpl{
		reader: bufio.NewReader(os.Stdin),
	}
}

// SelectBackupMethod prompts user to select backup method
func (s *ConfigServiceImpl) SelectBackupMethod() (domain.BackupMethod, error) {
	fmt.Println("Select backup method:")
	fmt.Println("  1. docker-run    (Use temporary container)")
	fmt.Println("  2. docker-exec   (Exec into existing Docker container)")
	fmt.Println("  3. kubectl-exec  (Exec into Kubernetes pod)")
	
	for {
		fmt.Print("\nEnter choice [1-3]: ")
		input, _ := s.reader.ReadString('\n')
		input = strings.TrimSpace(input)
		
		switch input {
		case "1":
			return domain.BackupMethodDockerRun, nil
		case "2":
			return domain.BackupMethodDockerExec, nil
		case "3":
			return domain.BackupMethodKubectlExec, nil
		default:
			fmt.Println(colorRed + "Invalid choice. Please enter 1, 2, or 3." + colorReset)
		}
	}
}

// SelectDatabases prompts user to select databases to backup
func (s *ConfigServiceImpl) SelectDatabases() ([]domain.DatabaseType, error) {
	fmt.Println("\nSelect databases to backup:")
	fmt.Println("  1. PostgreSQL")
	fmt.Println("  2. MySQL")
	fmt.Println("  3. MariaDB")
	fmt.Println("  4. MongoDB")
	fmt.Println("  5. All databases")
	
	fmt.Print("\nEnter choices (comma-separated, e.g., 1,2,4): ")
	input, _ := s.reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "5" {
		return []domain.DatabaseType{
			domain.DatabaseTypePostgres,
			domain.DatabaseTypeMySQL,
			domain.DatabaseTypeMariaDB,
			domain.DatabaseTypeMongoDB,
		}, nil
	}
	
	choices := strings.Split(input, ",")
	var selected []domain.DatabaseType
	
	for _, choice := range choices {
		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			selected = append(selected, domain.DatabaseTypePostgres)
		case "2":
			selected = append(selected, domain.DatabaseTypeMySQL)
		case "3":
			selected = append(selected, domain.DatabaseTypeMariaDB)
		case "4":
			selected = append(selected, domain.DatabaseTypeMongoDB)
		}
	}
	
	if len(selected) == 0 {
		return nil, fmt.Errorf("no databases selected")
	}
	
	return selected, nil
}

// GetKubernetesNamespace prompts user for Kubernetes namespace
func (s *ConfigServiceImpl) GetKubernetesNamespace() (string, error) {
	fmt.Println()
	namespace := s.promptInput("Kubernetes Namespace", "default")
	return namespace, nil
}

// ConfigureDatabase prompts user to configure a specific database
func (s *ConfigServiceImpl) ConfigureDatabase(dbType domain.DatabaseType, method domain.BackupMethod) (domain.DatabaseConfig, error) {
	config := domain.DatabaseConfig{
		Type: dbType,
	}
	
	fmt.Printf("\n%s=== Configuring %s ===%s\n", colorBlue, strings.ToUpper(dbType.String()), colorReset)
	
	switch dbType {
	case domain.DatabaseTypePostgres:
		config.Host = s.promptInput("PostgreSQL Host", "postgres")
		config.User = s.promptInput("PostgreSQL User", "postgres")
		config.Database = s.promptInput("Database Name", "mydb")
		config.Password = s.promptPassword("PostgreSQL Password")
		config.Version = s.promptInput("PostgreSQL Version", "15")
		
		if method == domain.BackupMethodDockerExec {
			config.Container = s.promptInput("Container Name", "test-postgres")
		} else if method == domain.BackupMethodKubectlExec {
			config.Pod = s.promptInput("Pod Name", "postgres-0")
		}
		
	case domain.DatabaseTypeMySQL:
		config.Host = s.promptInput("MySQL Host", "mysql")
		config.User = s.promptInput("MySQL User", "root")
		config.Database = s.promptInput("Database Name", "mydb")
		config.Password = s.promptPassword("MySQL Password")
		config.Version = s.promptInput("MySQL Version", "8")
		
		if method == domain.BackupMethodDockerExec {
			config.Container = s.promptInput("Container Name", "test-mysql")
		} else if method == domain.BackupMethodKubectlExec {
			config.Pod = s.promptInput("Pod Name", "mysql-0")
		}
		
	case domain.DatabaseTypeMariaDB:
		config.Host = s.promptInput("MariaDB Host", "mariadb")
		config.User = s.promptInput("MariaDB User", "root")
		config.Database = s.promptInput("Database Name", "mydb")
		config.Password = s.promptPassword("MariaDB Password")
		config.Version = s.promptInput("MariaDB Version", "11")
		
		if method == domain.BackupMethodDockerExec {
			config.Container = s.promptInput("Container Name", "test-mariadb")
		} else if method == domain.BackupMethodKubectlExec {
			config.Pod = s.promptInput("Pod Name", "mariadb-0")
		}
		
	case domain.DatabaseTypeMongoDB:
		config.Host = s.promptInput("MongoDB Host", "mongodb")
		config.Database = s.promptInput("Database Name", "mydb")
		config.Version = s.promptInput("MongoDB Version", "7")
		
		if method == domain.BackupMethodDockerExec {
			config.Container = s.promptInput("Container Name", "test-mongodb")
		} else if method == domain.BackupMethodKubectlExec {
			config.Pod = s.promptInput("Pod Name", "mongodb-0")
		}
	}
	
	return config, nil
}

// ConfirmBackup asks user to confirm backup operation
func (s *ConfigServiceImpl) ConfirmBackup(config domain.BackupConfig) (bool, error) {
	fmt.Print("\nProceed with backup? (y/n): ")
	input, _ := s.reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	return input == "y" || input == "yes", nil
}

// Helper methods
func (s *ConfigServiceImpl) promptInput(prompt, defaultValue string) string {
	fmt.Printf("%s [%s]: ", prompt, defaultValue)
	input, _ := s.reader.ReadString('\n')
	input = strings.TrimSpace(input)
	
	if input == "" {
		return defaultValue
	}
	return input
}

func (s *ConfigServiceImpl) promptPassword(prompt string) string {
	fmt.Printf("%s: ", prompt)
	input, _ := s.reader.ReadString('\n')
	return strings.TrimSpace(input)
}
