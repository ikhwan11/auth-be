package main

// @title Auth Service API
// @version 1.0
// @description Authentication and Authorization Service API
// @host localhost:8060
// @BasePath /

import (
	"fmt"
	"log"
	"os"

	_ "github.com/ikhwan11/auth-be/docs"

	"github.com/ikhwan11/auth-be/internal/config"
	"github.com/ikhwan11/auth-be/internal/container"
	"github.com/ikhwan11/auth-be/internal/database"
	"github.com/ikhwan11/auth-be/internal/router"
)

func main() {
	// CMD
	if len(os.Args) > 1 {
		if handled := handleCommand(os.Args[1:]); handled {
			return
		}
	}

	cfg := config.Load()

	// WIRING DATABASES
	fmt.Printf("%+v\n", cfg.AuthDB)

	employeeDB, err := database.NewConnection(cfg.EmployeeDB)
	if err != nil {
		log.Fatal("failed to connect employee db:", err)
	}
	defer employeeDB.Close()

	authDB, err := database.NewConnection(cfg.AuthDB)
	if err != nil {
		log.Fatal("failed to connect auth db:", err)
	}
	defer authDB.Close()

	// CONTAINER SECTION
	appContainer, err := container.New(
		employeeDB,
		authDB,
		cfg,
	)
	if err != nil {
		log.Fatal(err)
	}

	// ROUTES SECTION
	r := router.Setup(appContainer)

	log.Println("🚀 Server running on :" + cfg.AppPort)

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}

// handleCommand menangani perintah CLI custom (make:migration, migrate:up, dll).
// return true jika command dikenali & sudah dieksekusi (main tidak perlu lanjut ke server).
func handleCommand(args []string) bool {
	switch args[0] {
	case "make:migration":
		if len(args) < 2 {
			fmt.Println("Usage: go run ./cmd/api make:migration <migration_name>")
			os.Exit(1)
		}
		if err := makeMigration(args[1]); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		return true

	case "migrate:up":
		if err := runMigration("up"); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		return true

	case "migrate:down":
		if err := runMigration("down"); err != nil {
			fmt.Println("Error:", err)
			os.Exit(1)
		}
		return true

	default:
		fmt.Printf("Unknown command: %s\n", args[0])
		fmt.Println("Available commands: make:migration, migrate:up, migrate:down")
		os.Exit(1)
		return true
	}
}
