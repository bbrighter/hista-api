package main

import (
	"fmt"
	"io"
	"os"

	gormSchema "ariga.io/atlas-provider-gorm/gormschema"
	"encore.app/shoppingList/entity"
)

// Define the models to generate migrations for.
var models = []any{
	&entity.Item{},
	&entity.Product{},
}

func main() {
	statements, err := gormSchema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, statements)
}
