package main

import (
	"errors"
	"os"
	"strings"
)

var lockfilesToPackageManagers = map[string]string{
	"pnpm-lock.yaml":    "pnpm",
	"pnpm-lock.yml":     "pnpm",
	"package-lock.json": "npm",
	"bun.lockb":         "bun",
	"bun.lock":          "bun",
	"yarn.lock":         "yarn",
	"deno.lock":         "deno",
}

func detectPackageManager() string {
	cwd, err := os.Getwd()
	if err != nil {
		printErrorFatal("Unable to get current directory", err)
	}

	dirEntry, err := os.ReadDir(cwd)
	if err != nil {
		printErrorFatal("Unable to read contents of current directory", err)
	}

	var lockfiles []string
	for _, entry := range dirEntry {
		if entry.IsDir() {
			continue
		}

		if _, exists := lockfilesToPackageManagers[entry.Name()]; exists {
			lockfiles = append(lockfiles, entry.Name())
		}
	}

	if len(lockfiles) > 1 {
		multipeLockfilesErr := errors.New("- " + strings.Join(lockfiles, "\n- "))
		printErrorFatal("Found multiple lockfiles", multipeLockfilesErr)
	}

	if len(lockfiles) == 0 {
		return "npm"
	}
	return lockfilesToPackageManagers[lockfiles[0]]
}
