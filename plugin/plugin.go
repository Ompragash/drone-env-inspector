// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	// EnvName is a comma-separated list of environment variable names to inspect.
	EnvName string `envconfig:"PLUGIN_ENV_NAME"`

	// Secret determines whether to write to the secret output file.
	Secret bool `envconfig:"PLUGIN_SECRET"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) error {
	if err := validateArgs(args); err != nil {
		return err
	}

	envNames := strings.Split(args.EnvName, ",")

	for _, name := range envNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		value, exists := os.LookupEnv(name)
		if !exists {
			logrus.Warnf("environment variable %s does not exist\n", name)
		}

		logrus.Infof("inspecting environment variable: %s\n", name)

		var err error
		if args.Secret {
			err = writeSecretOutputToFile(name, value)
		} else {
			err = writeOutputToFile(name, value)
		}

		if err != nil {
			return err
		}

		if exists {
			logrus.Infof("successfully exported %s\n", name)
		} else {
			logrus.Infof("exported %s with empty value (variable does not exist)\n", name)
		}
	}

	return nil
}

// validateArgs validates the plugin arguments.
func validateArgs(args Args) error {
	if args.EnvName == "" {
		return fmt.Errorf("env_name is required")
	}
	return nil
}

// writeOutputToFile writes a key-value pair to the DRONE_OUTPUT file.
func writeOutputToFile(key, value string) error {
	outputFile, err := os.OpenFile(os.Getenv("DRONE_OUTPUT"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open output file: %w", err)
	}
	defer outputFile.Close()

	_, err = fmt.Fprintf(outputFile, "%s=%s\n", key, value)
	if err != nil {
		return fmt.Errorf("failed to write to output file: %w", err)
	}

	return nil
}

// writeSecretOutputToFile writes a key-value pair to the HARNESS_OUTPUT_SECRET_FILE.
func writeSecretOutputToFile(key, value string) error {
	outputFile, err := os.OpenFile(os.Getenv("HARNESS_OUTPUT_SECRET_FILE"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open secret output file: %w", err)
	}
	defer outputFile.Close()

	_, err = fmt.Fprintf(outputFile, "%s=%s\n", key, value)
	if err != nil {
		return fmt.Errorf("failed to write to secret output file: %w", err)
	}

	return nil
}

