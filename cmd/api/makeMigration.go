package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

const migrationDir = "migrations"

func makeMigration(name string) error {
	if name == "" {
		return fmt.Errorf("migration name tidak boleh kosong")
	}

	version, err := nextMigrationVersion(migrationDir)
	if err != nil {
		return err
	}

	upFile := filepath.Join(migrationDir, fmt.Sprintf("%s_%s.up.sql", version, name))
	downFile := filepath.Join(migrationDir, fmt.Sprintf("%s_%s.down.sql", version, name))

	upContent := fmt.Sprintf("CREATE TABLE %s (\n\n);\n", name)
	downContent := fmt.Sprintf("DROP TABLE IF EXISTS %s;\n", name)

	if err := os.MkdirAll(migrationDir, 0o755); err != nil {
		return err
	}

	if err := os.WriteFile(upFile, []byte(upContent), 0o644); err != nil {
		return fmt.Errorf("gagal membuat file up: %w", err)
	}
	if err := os.WriteFile(downFile, []byte(downContent), 0o644); err != nil {
		return fmt.Errorf("gagal membuat file down: %w", err)
	}

	fmt.Println("Migration berhasil dibuat:")
	fmt.Println(" -", upFile)
	fmt.Println(" -", downFile)
	return nil
}

// nextMigrationVersion mencari nomor urut terakhir di folder migrations
// dan mengembalikan nomor berikutnya dengan format 6 digit (000001, 000002, dst)
func nextMigrationVersion(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return "000001", nil
		}
		return "", err
	}

	maxVersion := 0
	re := regexp.MustCompile(`^(\d+)_`)

	for _, entry := range entries {
		matches := re.FindStringSubmatch(entry.Name())
		if len(matches) == 2 {
			v, err := strconv.Atoi(matches[1])
			if err == nil && v > maxVersion {
				maxVersion = v
			}
		}
	}

	return fmt.Sprintf("%06d", maxVersion+1), nil
}
