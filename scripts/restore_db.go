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
)

func main() {
	// Parse command line flags
	backupFile := flag.String("backup", "", "Path to backup zip file (required)")
	dbPath := flag.String("db", "data.sqlite", "Path where database should be restored")
	force := flag.Bool("force", false, "Force restore even if database exists")
	flag.Parse()

	if *backupFile == "" {
		flag.Usage()
		log.Fatal("Error: --backup flag is required")
	}

	// Check if backup file exists
	if _, err := os.Stat(*backupFile); os.IsNotExist(err) {
		log.Fatalf("Backup file not found: %s", *backupFile)
	}

	// Check if database already exists
	if _, err := os.Stat(*dbPath); err == nil && !*force {
		log.Fatalf("Database already exists at %s. Use --force to overwrite.", *dbPath)
	}

	// Restore backup
	if err := restoreBackup(*backupFile, *dbPath); err != nil {
		log.Fatalf("Restore failed: %v", err)
	}

	fmt.Printf("✅ Database restored successfully to: %s\n", *dbPath)
}

// restoreBackup extracts the database from a backup zip file
func restoreBackup(backupFile, dbPath string) error {
	// Open the zip file
	reader, err := zip.OpenReader(backupFile)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer reader.Close()

	// Extract all files from the zip
	for _, file := range reader.File {
		if err := extractFile(file, filepath.Dir(dbPath)); err != nil {
			return fmt.Errorf("failed to extract %s: %w", file.Name, err)
		}
		fmt.Printf("Extracted: %s\n", file.Name)
	}

	return nil
}

// extractFile extracts a single file from the zip archive
func extractFile(file *zip.File, destDir string) error {
	// Create destination path
	destPath := filepath.Join(destDir, file.Name)

	// Create parent directories if needed
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	// Open the file in the zip
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy contents
	_, err = io.Copy(destFile, srcFile)
	return err
}

// listBackups lists all available backups
func listBackups(backupDir string) error {
	files, err := filepath.Glob(filepath.Join(backupDir, "backup_*.zip"))
	if err != nil {
		return err
	}

	if len(files) == 0 {
		fmt.Println("No backups found")
		return nil
	}

	fmt.Println("Available backups:")
	for i, file := range files {
		info, err := os.Stat(file)
		if err != nil {
			continue
		}

		base := filepath.Base(file)

		fmt.Printf("%d. %s (%.2f MB, %s)\n",
			i+1,
			base,
			float64(info.Size())/(1024*1024),
			info.ModTime().Format("2006-01-02 15:04:05"),
		)
	}

	return nil
}
