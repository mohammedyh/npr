package main

import (
	"errors"
	"os"
	"strings"
)

type DetectionError struct {
	Title string
	Err   error
}

func (e *DetectionError) Error() string { return e.Err.Error() }

var lockfilesToPackageManagers = map[string]string{
	"pnpm-lock.yaml":    "pnpm",
	"pnpm-lock.yml":     "pnpm",
	"package-lock.json": "npm",
	"bun.lockb":         "bun",
	"bun.lock":          "bun",
	"yarn.lock":         "yarn",
	"deno.lock":         "deno",
}

func detectPackageManager() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", &DetectionError{"Unable to get current directory", err}
	}

	entries, err := os.ReadDir(cwd)
	if err != nil {
		return "", &DetectionError{"Unable to read contents of current directory", err}
	}

	var lockfiles []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if _, exists := lockfilesToPackageManagers[entry.Name()]; exists {
			lockfiles = append(lockfiles, entry.Name())
		}
	}

	if len(lockfiles) > 1 {
		multipleLockfilesErr := errors.New("- " + strings.Join(lockfiles, "\n- "))
		return "", &DetectionError{Title: "Found multiple lockfiles", Err: multipleLockfilesErr}
	}

	if len(lockfiles) == 0 {
		return "npm", nil
	}
	return lockfilesToPackageManagers[lockfiles[0]], nil
}
