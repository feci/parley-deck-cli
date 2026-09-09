//go:build !darwin

package telemetry

import "os"

func syncRecord(file *os.File) error { return file.Sync() }
