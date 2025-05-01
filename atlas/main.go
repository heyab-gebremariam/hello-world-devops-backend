package main

import (
	"fmt"
	"io"
	"os"

	"hello-world-devops-backend/models"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&models.Task{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	// Output the generated SQL statements
	io.WriteString(os.Stdout, stmts)
}
