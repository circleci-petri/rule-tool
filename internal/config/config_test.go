package config

import (
	"os"
	"testing"
)

func TestEnvVariables(t *testing.T) {
	// Save current env vars to restore later
	oldRulesEnv := os.Getenv(EnvRulesPath)
	oldTargetEnv := os.Getenv(EnvTargetPath)
	defer func() {
		os.Setenv(EnvRulesPath, oldRulesEnv)
		os.Setenv(EnvTargetPath, oldTargetEnv)
	}()

	testCases := []struct {
		name            string
		rulesEnvValue   string
		targetEnvValue  string
		rulesFlagValue  string
		targetFlagValue string
		expectedRules   string
		expectedTarget  string
	}{
		{
			name:            "Flags take precedence over env vars",
			rulesEnvValue:   "/env/rules/path",
			targetEnvValue:  "/env/target/path",
			rulesFlagValue:  "/flag/rules/path",
			targetFlagValue: "/flag/target/path",
			expectedRules:   "/flag/rules/path",
			expectedTarget:  "/flag/target/path",
		},
		{
			name:            "Env vars used when flags not provided",
			rulesEnvValue:   "/env/rules/path",
			targetEnvValue:  "/env/target/path",
			rulesFlagValue:  "",
			targetFlagValue: "",
			expectedRules:   "/env/rules/path",
			expectedTarget:  "/env/target/path",
		},
		{
			name:            "Mix of flags and env vars",
			rulesEnvValue:   "/env/rules/path",
			targetEnvValue:  "/env/target/path",
			rulesFlagValue:  "/flag/rules/path",
			targetFlagValue: "",
			expectedRules:   "/flag/rules/path",
			expectedTarget:  "/env/target/path",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment variables
			os.Setenv(EnvRulesPath, tc.rulesEnvValue)
			os.Setenv(EnvTargetPath, tc.targetEnvValue)

			// Create new config which loads from env vars
			cfg := New()

			// Verify env vars were used
			if tc.rulesEnvValue != "" && tc.rulesFlagValue == "" {
				if cfg.RulesRepoPath != tc.rulesEnvValue {
					t.Errorf("Expected RulesRepoPath to be %q from env var, got %q", tc.rulesEnvValue, cfg.RulesRepoPath)
				}
			}

			if tc.targetEnvValue != "" && tc.targetFlagValue == "" {
				if cfg.TargetProjectPath != tc.targetEnvValue {
					t.Errorf("Expected TargetProjectPath to be %q from env var, got %q", tc.targetEnvValue, cfg.TargetProjectPath)
				}
			}

			// Apply flag values if provided (simulating command-line flags)
			if tc.rulesFlagValue != "" {
				cfg.SetRulesRepoPath(tc.rulesFlagValue)
			}

			if tc.targetFlagValue != "" {
				cfg.SetTargetProjectPath(tc.targetFlagValue)
			}

			// Verify final results
			if cfg.RulesRepoPath != tc.expectedRules {
				t.Errorf("Expected RulesRepoPath to be %q, got %q", tc.expectedRules, cfg.RulesRepoPath)
			}

			if cfg.TargetProjectPath != tc.expectedTarget {
				t.Errorf("Expected TargetProjectPath to be %q, got %q", tc.expectedTarget, cfg.TargetProjectPath)
			}
		})
	}
}

func TestDefaultValues(t *testing.T) {
	// Save current env vars to restore later
	oldRulesEnv := os.Getenv(EnvRulesPath)
	oldTargetEnv := os.Getenv(EnvTargetPath)
	defer func() {
		os.Setenv(EnvRulesPath, oldRulesEnv)
		os.Setenv(EnvTargetPath, oldTargetEnv)
	}()

	// Clear environment variables
	os.Unsetenv(EnvRulesPath)
	os.Unsetenv(EnvTargetPath)

	// Create new config
	cfg := New()

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Rules path should default to current directory
	if cfg.RulesRepoPath != cwd {
		t.Errorf("Expected RulesRepoPath to be current directory %q, got %q", cwd, cfg.RulesRepoPath)
	}

	// Target path should default to current directory
	if cfg.TargetProjectPath != cwd {
		t.Errorf("Expected TargetProjectPath to be current directory %q, got %q", cwd, cfg.TargetProjectPath)
	}
}

func TestDefaultPathValues(t *testing.T) {
	// Save current env vars to restore later
	oldRulesEnv := os.Getenv(EnvRulesPath)
	oldTargetEnv := os.Getenv(EnvTargetPath)
	defer func() {
		os.Setenv(EnvRulesPath, oldRulesEnv)
		os.Setenv(EnvTargetPath, oldTargetEnv)
	}()

	// Clear environment variables
	os.Unsetenv(EnvRulesPath)
	os.Unsetenv(EnvTargetPath)

	// Create new config
	cfg := New()

	// Get current directory
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Rules path should default to current directory when not specified
	if cfg.RulesRepoPath != cwd && cfg.RulesRepoPath != "" {
		t.Errorf("Rules path should either default to current directory %q or be empty, got %q",
			cwd, cfg.RulesRepoPath)
	}

	// Target path should default to current directory
	if cfg.TargetProjectPath != cwd {
		t.Errorf("Expected TargetProjectPath to be current directory %q, got %q",
			cwd, cfg.TargetProjectPath)
	}

	// Now validate that ValidateRulesRepoPath can work with default path
	if !cfg.ValidateRulesRepoPath() && cfg.RulesRepoPath == cwd {
		t.Errorf("ValidateRulesRepoPath should return true for current directory %q", cwd)
	}

	// Validate that ValidateTargetProjectPath works with default path
	if !cfg.ValidateTargetProjectPath() {
		t.Errorf("ValidateTargetProjectPath should return true for current directory %q", cwd)
	}
}
