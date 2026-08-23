package testutils

import (
	OS "os"

	Logger "github.com/danyel/ecommerce/cmd/logger"
)

func PreInitTest() {
	err := OS.Setenv("DEBUG_ENABLED", "true")
	if err != nil {
		Logger.Log.Fatal(err)
		OS.Exit(0)
	}
}
