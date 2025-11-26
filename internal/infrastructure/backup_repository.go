package infrastructure

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/wush/db-backup-tool/internal/domain"
)

// BackupRepositoryImpl implements domain.BackupRepository
type BackupRepositoryImpl struct{}

// NewBackupRepository creates a new backup repository
func NewBackupRepository() domain.BackupRepository {
	return &BackupRepositoryImpl{}
}

// BackupPostgres performs a PostgreSQL backup
func (r *BackupRepositoryImpl) BackupPostgres(config domain.DatabaseConfig, method domain.BackupMethod, backupPath, namespace string) error {
	cwd, _ := os.Getwd()
	
	switch method {
	case domain.BackupMethodDockerRun:
		cmd := exec.Command("docker", "run", "--rm",
			"-e", fmt.Sprintf("PGPASSWORD=%s", config.Password),
			"-v", fmt.Sprintf("%s/backup/postgres:/backup", cwd),
			fmt.Sprintf("postgres:%s", config.Version),
			"pg_dump", "-h", config.Host, "-U", config.User, config.Database)
		
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("docker run failed: %w", err)
		}
		return os.WriteFile(backupPath, output, 0644)
		
	case domain.BackupMethodDockerExec:
		cmd := exec.Command("docker", "exec", config.Container,
			"sh", "-c",
			fmt.Sprintf("PGPASSWORD='%s' pg_dump -h localhost -U %s %s",
				config.Password, config.User, config.Database))
		
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("docker exec failed: %w", err)
		}
		return os.WriteFile(backupPath, output, 0644)
		
	case domain.BackupMethodKubectlExec:
		cmd := exec.Command("kubectl", "exec", "-n", namespace, config.Pod, "--",
			"sh", "-c",
			fmt.Sprintf("PGPASSWORD='%s' pg_dump -h localhost -U %s %s",
				config.Password, config.User, config.Database))
		
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("kubectl exec failed: %w", err)
		}
		return os.WriteFile(backupPath, output, 0644)
	}
	
	return fmt.Errorf("unknown backup method: %s", method)
}

// BackupMySQL performs a MySQL backup
func (r *BackupRepositoryImpl) BackupMySQL(config domain.DatabaseConfig, method domain.BackupMethod, backupPath, namespace string) error {
	cwd, _ := os.Getwd()
	
	switch method {
	case domain.BackupMethodDockerRun:
		cmd := exec.Command("docker", "run", "--rm",
			"-v", fmt.Sprintf("%s/backup/mysql:/backup", cwd),
			fmt.Sprintf("mysql:%s", config.Version),
			"sh", "-c",
			fmt.Sprintf("mysqldump -h%s -u%s -p%s %s",
				config.Host, config.User, config.Password, config.Database))
		
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("docker run failed: %w", err)
		}
		return os.WriteFile(backupPath, output, 0644)
		
	case domain.BackupMethodDockerExec:
		cmd := exec.Command("docker", "exec", config.Container,
			"sh", "-c",
			fmt.Sprintf("mysqldump -h localhost -u%s -p%s %s",
				config.User, config.Password, config.Database))
		
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("docker exec failed: %w", err)
		}
		return os.WriteFile(backupPath, output, 0644)
		
	case domain.BackupMethodKubectlExec:
		cmd := exec.Command("kubectl", "exec", "-n", namespace, config.Pod, "--",
			"sh", "-c",
			fmt.Sprintf("mysqldump -h localhost -u%s -p%s %s",
				config.User, config.Password, config.Database))
		
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("kubectl exec failed: %w", err)
		}
		return os.WriteFile(backupPath, output, 0644)
	}
	
	return fmt.Errorf("unknown backup method: %s", method)
}

// BackupMariaDB performs a MariaDB backup
func (r *BackupRepositoryImpl) BackupMariaDB(config domain.DatabaseConfig, method domain.BackupMethod, backupPath, namespace string) error {
	cwd, _ := os.Getwd()
	
	switch method {
	case domain.BackupMethodDockerRun:
		cmd := exec.Command("docker", "run", "--rm",
			"-v", fmt.Sprintf("%s/backup/mariadb:/backup", cwd),
			fmt.Sprintf("mariadb:%s", config.Version),
			"sh", "-c",
			fmt.Sprintf("mysqldump -h%s -u%s -p%s %s",
				config.Host, config.User, config.Password, config.Database))
		
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("docker run failed: %w", err)
		}
		return os.WriteFile(backupPath, output, 0644)
		
	case domain.BackupMethodDockerExec:
		cmd := exec.Command("docker", "exec", config.Container,
			"sh", "-c",
			fmt.Sprintf("mysqldump -h localhost -u%s -p%s %s",
				config.User, config.Password, config.Database))
		
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("docker exec failed: %w", err)
		}
		return os.WriteFile(backupPath, output, 0644)
		
	case domain.BackupMethodKubectlExec:
		cmd := exec.Command("kubectl", "exec", "-n", namespace, config.Pod, "--",
			"sh", "-c",
			fmt.Sprintf("mysqldump -h localhost -u%s -p%s %s",
				config.User, config.Password, config.Database))
		
		output, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("kubectl exec failed: %w", err)
		}
		return os.WriteFile(backupPath, output, 0644)
	}
	
	return fmt.Errorf("unknown backup method: %s", method)
}

// BackupMongoDB performs a MongoDB backup
func (r *BackupRepositoryImpl) BackupMongoDB(config domain.DatabaseConfig, method domain.BackupMethod, backupPath, namespace, tempDir string) error {
	cwd, _ := os.Getwd()
	
	switch method {
	case domain.BackupMethodDockerRun:
		cmd := exec.Command("docker", "run", "--rm",
			"-v", fmt.Sprintf("%s/backup/mongodb:/backup", cwd),
			fmt.Sprintf("mongo:%s", config.Version),
			"mongodump", "--host", config.Host, "--db", config.Database,
			"--out", fmt.Sprintf("/backup/%s", filepath.Base(backupPath)))
		
		return cmd.Run()
		
	case domain.BackupMethodDockerExec:
		timestamp := filepath.Base(backupPath)
		
		// Create backup inside container
		cmd := exec.Command("docker", "exec", config.Container,
			"mongodump", "--host", "localhost", "--db", config.Database,
			"--out", fmt.Sprintf("%s/%s", tempDir, timestamp))
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create backup in container: %w", err)
		}
		
		// Copy backup from container to host
		os.MkdirAll(backupPath, 0755)
		cmd = exec.Command("docker", "cp",
			fmt.Sprintf("%s:%s/%s/%s", config.Container, tempDir, timestamp, config.Database),
			backupPath+"/")
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to copy backup from container: %w", err)
		}
		
		// Cleanup inside container
		cmd = exec.Command("docker", "exec", config.Container,
			"rm", "-rf", fmt.Sprintf("%s/%s", tempDir, timestamp))
		cmd.Run()
		
		return nil
		
	case domain.BackupMethodKubectlExec:
		timestamp := filepath.Base(backupPath)
		
		// Create backup inside pod
		cmd := exec.Command("kubectl", "exec", "-n", namespace, config.Pod, "--",
			"mongodump", "--host", "localhost", "--db", config.Database,
			"--out", fmt.Sprintf("%s/%s", tempDir, timestamp))
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create backup in pod: %w", err)
		}
		
		// Copy backup from pod to host
		os.MkdirAll(backupPath, 0755)
		cmd = exec.Command("kubectl", "cp",
			fmt.Sprintf("%s/%s:%s/%s/%s", namespace, config.Pod, tempDir, timestamp, config.Database),
			backupPath+"/")
		
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to copy backup from pod: %w", err)
		}
		
		// Cleanup inside pod
		cmd = exec.Command("kubectl", "exec", "-n", namespace, config.Pod, "--",
			"rm", "-rf", fmt.Sprintf("%s/%s", tempDir, timestamp))
		cmd.Run()
		
		return nil
	}
	
	return fmt.Errorf("unknown backup method: %s", method)
}

// GetFileSize returns the size of a file or directory
func (r *BackupRepositoryImpl) GetFileSize(path string, isDirectory bool) (string, error) {
	var cmd *exec.Cmd
	if isDirectory {
		cmd = exec.Command("du", "-sh", path)
	} else {
		cmd = exec.Command("du", "-h", path)
	}
	
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get file size: %w", err)
	}
	
	fields := strings.Fields(string(output))
	if len(fields) > 0 {
		return fields[0], nil
	}
	
	return "unknown", nil
}
