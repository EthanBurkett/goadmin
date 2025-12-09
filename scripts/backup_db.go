//go:build ignore
// +build ignore

package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	// Parse command line flags
	backupDir := flag.String("dir", "backups", "Directory to store backups")
	dbPath := flag.String("db", "data.sqlite", "Path to SQLite database file")
	flag.Parse()

	// Ensure backup directory exists
	if err := os.MkdirAll(*backupDir, 0755); err != nil {
		log.Fatalf("Failed to create backup directory: %v", err)
	}

	// Generate backup filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupFile := filepath.Join(*backupDir, fmt.Sprintf("backup_%s.zip", timestamp))

	// Create backup
	if err := createBackup(*dbPath, backupFile); err != nil {
		log.Fatalf("Backup failed: %v", err)
	}

	fmt.Printf("✅ Backup created successfully: %s\n", backupFile)

	// Cleanup old backups (keep last 10)
	if err := cleanupOldBackups(*backupDir, 10); err != nil {
		log.Printf("Warning: Failed to cleanup old backups: %v", err)
	}
}

// createBackup creates a compressed backup of the database
func createBackup(dbPath, backupFile string) error {
	// Check if database file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		return fmt.Errorf("database file not found: %s", dbPath)
	}

	// Create zip file
	zipFile, err := os.Create(backupFile)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer zipFile.Close()

	// Create zip writer
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Add database file to zip
	if err := addFileToZip(zipWriter, dbPath, filepath.Base(dbPath)); err != nil {
		return fmt.Errorf("failed to add database to backup: %w", err)
	}

	// Also backup WAL file if it exists (SQLite write-ahead log)
	walPath := dbPath + "-wal"
	if _, err := os.Stat(walPath); err == nil {
		if err := addFileToZip(zipWriter, walPath, filepath.Base(walPath)); err != nil {
			log.Printf("Warning: Failed to backup WAL file: %v", err)
		}
	}

	// Also backup SHM file if it exists (SQLite shared memory)
	shmPath := dbPath + "-shm"
	if _, err := os.Stat(shmPath); err == nil {
		if err := addFileToZip(zipWriter, shmPath, filepath.Base(shmPath)); err != nil {
			log.Printf("Warning: Failed to backup SHM file: %v", err)
		}
	}

	return nil
}

// addFileToZip adds a file to the zip archive
func addFileToZip(zipWriter *zip.Writer, filePath, zipPath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Name = zipPath
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

// cleanupOldBackups removes old backup files, keeping only the most recent ones
func cleanupOldBackups(backupDir string, keepCount int) error {
	// Get all backup files
	files, err := filepath.Glob(filepath.Join(backupDir, "backup_*.zip"))
	if err != nil {
		return err
	}

	// If we have fewer files than keepCount, nothing to delete
	if len(files) <= keepCount {
		return nil
	}

	// Sort files by modification time (newest first)
	type fileInfo struct {
		path    string
		modTime time.Time
	}

	var fileInfos []fileInfo
	for _, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}
		fileInfos = append(fileInfos, fileInfo{
			path:    file,
			modTime: info.ModTime(),
		})
	}

	// Sort by modification time (newest first)
	for i := 0; i < len(fileInfos)-1; i++ {
		for j := i + 1; j < len(fileInfos); j++ {
			if fileInfos[i].modTime.Before(fileInfos[j].modTime) {
				fileInfos[i], fileInfos[j] = fileInfos[j], fileInfos[i]
			}
		}
	}

	// Delete old backups
	deleted := 0
	for i := keepCount; i < len(fileInfos); i++ {
		if err := os.Remove(fileInfos[i].path); err != nil {
			log.Printf("Warning: Failed to delete old backup %s: %v", fileInfos[i].path, err)
		} else {
			deleted++
		}
	}

	if deleted > 0 {
		fmt.Printf("🗑️  Cleaned up %d old backup(s)\n", deleted)
	}

	return nil
}
