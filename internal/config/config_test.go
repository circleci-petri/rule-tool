package config

import (
	"os"
	"testing"

	"github.com/circleci/llm-agent-rules/internal/git"
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

	// Editor should default to Cursor
	if cfg.Editor != EditorCursor {
		t.Errorf("Expected Editor to be Cursor, got %s", cfg.Editor)
	}

	// Git repository should be disabled by default
	if cfg.UseGitRepo {
		t.Errorf("Expected UseGitRepo to be false, got true")
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

func TestGetRulesDir(t *testing.T) {
	cfg := New()

	t.Run("No git repository", func(t *testing.T) {
		want := "/some/path/rules"
		cfg.UseGitRepo = false
		cfg.RulesRepoPath = "/some/path"

		got := cfg.GetRulesDir()

		// Should return the rules repo path
		if got != want {
			t.Errorf("Expected GetRulesDir to return %q, got %q", want, got)
		}
	})

	t.Run("With git repository", func(t *testing.T) {
		want := "/some/path/.cursor/rules"
		cfg.UseGitRepo = true
		cfg.GitRepo = &git.Repository{
			RulesPath: want,
			CloneDir:  "/some/path",
		}

		got := cfg.GetRulesDir()

		// Should return the git repository rules path
		if got != want {
			t.Errorf("Expected GetRulesDir to return %q, got %q", want, got)
		}
	})
}

func TestSetGitRepoURL(t *testing.T) {
	t.Run("Set a git repository URL", func(t *testing.T) {
		cfg := New()
		gitURL := "https://github.com/test/test-repo.git"

		cfg.SetGitRepoURL(gitURL)

		// Should set the Git repository URL
		if cfg.GitRepoURL != gitURL {
			t.Errorf("Expected GitRepoURL to be set to %q, got %q", gitURL, cfg.GitRepoURL)
		}

		// Should enable UseGitRepo
		if !cfg.UseGitRepo {
			t.Errorf("Expected UseGitRepo to be true, got false")
		}

		// Should initialize Git repository
		if cfg.GitRepo == nil {
			t.Errorf("Expected Git repository to be initialized, got nil")
		}

		// Should set the URL in the Git repository
		if cfg.GitRepo.URL != gitURL {
			t.Errorf("Expected Git repository URL to be set to %q, got %q", gitURL, cfg.GitRepo.URL)
		}
	})
	t.Run("Set an empty git repository URL", func(t *testing.T) {
		cfg := New()

		cfg.SetGitRepoURL("")

		// Should disable UseGitRepo
		if cfg.UseGitRepo {
			t.Errorf("Expected UseGitRepo to be false, got true")
		}

		// Should set GitRepo to nil
		if cfg.GitRepo != nil {
			t.Errorf("Expected Git repository to be nil, got %v", cfg.GitRepo)
		}
	})
}

func TestSetEditor(t *testing.T) {
	testCases := []struct {
		name     string
		editor   Editor
		expected Editor
	}{
		{
			name:     "Set Cursor editor",
			editor:   EditorCursor,
			expected: EditorCursor,
		},
		{
			name:     "Set Windsurf editor",
			editor:   EditorWindsurf,
			expected: EditorWindsurf,
		},
		{
			name:     "Set invalid editor",
			editor:   Editor("invalid"),
			expected: EditorCursor,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := New()

			cfg.SetEditor(string(tc.editor))

			// Should set the editor
			if cfg.Editor != tc.expected {
				t.Errorf("Expected Editor to be %s, got %s", tc.expected, cfg.Editor)
			}
		})
	}
}

func TestGetRulesDirectory(t *testing.T) {
	testCases := []struct {
		name   string
		editor Editor
		want   string
	}{
		{
			name:   "Cursor editor",
			editor: EditorCursor,
			want:   ".cursor/rules",
		},
		{
			name:   "Windsurf editor",
			editor: EditorWindsurf,
			want:   ".windsurf/rules",
		},
		{
			name:   "Invalid editor",
			editor: Editor("invalid"),
			want:   ".cursor/rules",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := New()

			cfg.Editor = tc.editor

			// Should return the expected rules directory
			if cfg.GetRulesDirectory() != tc.want {
				t.Errorf("Expected GetRulesDirectory to return %q, got %q", tc.want, cfg.GetRulesDirectory())
			}
		})
	}
}
