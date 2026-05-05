package main

import (
	"fmt"
	"io"
	"os"

	_ "ariga.io/atlas-go-sdk/recordriver"
	gormSchema "ariga.io/atlas-provider-gorm/gormschema"
	"encore.app/product_mgmt/instances"
)

var models = []any{
	&instances.Instance{},
}

func main() {
	statements, err := gormSchema.New("postgres").Load(models...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, statements)
}
