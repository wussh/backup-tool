package usecase

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/wush/db-backup-tool/internal/domain"
)

// BackupUsecase implements backup business logic
type BackupUsecase struct {
	backupRepo    domain.BackupRepository
	configService domain.ConfigService
	outputService domain.OutputService
}

// NewBackupUsecase creates a new backup usecase
func NewBackupUsecase(
	backupRepo domain.BackupRepository,
	configService domain.ConfigService,
	outputService domain.OutputService,
) *BackupUsecase {
	return &BackupUsecase{
		backupRepo:    backupRepo,
		configService: configService,
		outputService: outputService,
	}
}

// ExecuteInteractiveBackup runs the interactive backup process
func (uc *BackupUsecase) ExecuteInteractiveBackup() error {
	uc.outputService.PrintHeader()
	
	// Step 1: Select backup method
	method, err := uc.configService.SelectBackupMethod()
	if err != nil {
		return fmt.Errorf("failed to select backup method: %w", err)
	}
	
	// Step 2: Select databases
	dbTypes, err := uc.configService.SelectDatabases()
	if err != nil {
		return fmt.Errorf("failed to select databases: %w", err)
	}
	
	// Step 3: Get Kubernetes namespace if using kubectl-exec
	k8sNamespace := "default"
	if method == domain.BackupMethodKubectlExec {
		ns, err := uc.configService.GetKubernetesNamespace()
		if err != nil {
			return fmt.Errorf("failed to get kubernetes namespace: %w", err)
		}
		k8sNamespace = ns
	}
	
	// Step 4: Configure each database
	var dbConfigs []domain.DatabaseConfig
	for _, dbType := range dbTypes {
		config, err := uc.configService.ConfigureDatabase(dbType, method)
		if err != nil {
			return fmt.Errorf("failed to configure %s: %w", dbType, err)
		}
		dbConfigs = append(dbConfigs, config)
	}
	
	// Step 5: Build backup config
	backupConfig := domain.BackupConfig{
		Method:       method,
		Timestamp:    time.Now(),
		BackupDir:    "backup",
		TempDir:      "/tmp/db-backups",
		K8sNamespace: k8sNamespace,
		Databases:    dbConfigs,
	}
	
	// Step 6: Print summary and confirm
	uc.outputService.PrintConfigSummary(backupConfig)
	confirmed, err := uc.configService.ConfirmBackup(backupConfig)
	if err != nil {
		return fmt.Errorf("failed to get confirmation: %w", err)
	}
	if !confirmed {
		uc.outputService.PrintError("Backup cancelled by user")
		return nil
	}
	
	// Step 7: Execute backups
	results := uc.executeBackups(backupConfig)
	
	// Step 8: Print summary
	uc.outputService.PrintSummary(results)
	
	return nil
}

// executeBackups performs the actual backup operations
func (uc *BackupUsecase) executeBackups(config domain.BackupConfig) []domain.BackupResult {
	var results []domain.BackupResult
	
	timestamp := config.Timestamp.Format("2006-01-02_15-04-05")
	
	for _, dbConfig := range config.Databases {
		result := uc.backupDatabase(dbConfig, config.Method, timestamp, config.K8sNamespace, config.TempDir)
		results = append(results, result)
		uc.outputService.PrintBackupResult(result)
	}
	
	return results
}

// backupDatabase performs backup for a single database
func (uc *BackupUsecase) backupDatabase(
	dbConfig domain.DatabaseConfig,
	method domain.BackupMethod,
	timestamp string,
	namespace string,
	tempDir string,
) domain.BackupResult {
	startTime := time.Now()
	
	result := domain.BackupResult{
		DatabaseType: dbConfig.Type,
		Database:     dbConfig.Database,
		Success:      false,
	}
	
	// Print backup start message
	uc.outputService.PrintBackupStart(dbConfig.Type, dbConfig, method)
	
	// Create backup directory
	backupDir := filepath.Join("backup", dbConfig.Type.String())
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		result.Error = fmt.Errorf("failed to create backup directory: %w", err)
		result.Duration = time.Since(startTime)
		return result
	}
	
	var backupPath string
	var err error
	
	// Execute backup based on database type
	switch dbConfig.Type {
	case domain.DatabaseTypePostgres:
		backupPath = filepath.Join(backupDir, fmt.Sprintf("%s_%s.sql", dbConfig.Database, timestamp))
		err = uc.backupRepo.BackupPostgres(dbConfig, method, backupPath, namespace)
		
	case domain.DatabaseTypeMySQL:
		backupPath = filepath.Join(backupDir, fmt.Sprintf("%s_%s.sql", dbConfig.Database, timestamp))
		err = uc.backupRepo.BackupMySQL(dbConfig, method, backupPath, namespace)
		
	case domain.DatabaseTypeMariaDB:
		backupPath = filepath.Join(backupDir, fmt.Sprintf("%s_%s.sql", dbConfig.Database, timestamp))
		err = uc.backupRepo.BackupMariaDB(dbConfig, method, backupPath, namespace)
		
	case domain.DatabaseTypeMongoDB:
		backupPath = filepath.Join(backupDir, timestamp)
		err = uc.backupRepo.BackupMongoDB(dbConfig, method, backupPath, namespace, tempDir)
	}
	
	result.Duration = time.Since(startTime)
	result.BackupPath = backupPath
	
	if err != nil {
		result.Error = err
		return result
	}
	
	// Get backup size
	isDirectory := dbConfig.Type == domain.DatabaseTypeMongoDB
	size, err := uc.backupRepo.GetFileSize(backupPath, isDirectory)
	if err != nil {
		result.Error = fmt.Errorf("backup created but failed to get size: %w", err)
		return result
	}
	
	result.Size = size
	result.Success = true
	
	return result
}
