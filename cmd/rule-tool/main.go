package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/circleci/llm-agent-rules/internal/config"
	"github.com/circleci/llm-agent-rules/internal/linker"
	"github.com/circleci/llm-agent-rules/internal/rules"
	"github.com/circleci/llm-agent-rules/internal/ui"
)

func main() {
	// Initialize configuration
	cfg := config.New()

	// Parse command-line flags
	repoPath := flag.String("repo-path", "", "Path to the rules repository (overrides RULE_TOOL_PATH environment variable if set)")
	targetPath := flag.String("target-path", "", "Path to the target project (overrides RULE_TARGET_PATH environment variable if set)")
	gitRepoURL := flag.String("git-repo", "", "URL of a Git repository containing rules (overrides RULE_GIT_REPO_URL environment variable if set)")
	editor := flag.String("editor", "", "Editor to use (cursor, vscode, neovim, vim) (overrides RULE_EDITOR environment variable if set)")
	nonInteractive := flag.Bool("non-interactive", false, "Run in non-interactive mode")
	dryRun := flag.Bool("dry-run", false, "Show what would be done without making changes")
	listRules := flag.Bool("list", false, "List available rules")
	linkRule := flag.String("link", "", "Link a specific rule or comma-separated list of rules")
	unlinkRule := flag.String("unlink", "", "Unlink a specific rule or comma-separated list of rules")
	verbose := flag.Bool("verbose", false, "Enable verbose output")
	flag.Parse()

	// Set paths from flags if provided (flags take precedence over environment variables)
	if *repoPath != "" {
		cfg.SetRulesRepoPath(*repoPath)
	}

	if *targetPath != "" {
		cfg.SetTargetProjectPath(*targetPath)
	}

	// Set Git repository URL if provided
	if *gitRepoURL != "" {
		cfg.SetGitRepoURL(*gitRepoURL)
	}

	// Set editor if provided via flag
	if *editor != "" {
		cfg.SetEditor(*editor)
	}

	// Display configuration source if verbose
	if *verbose {
		if *repoPath != "" {
			fmt.Println("Using rules path from command line flag")
		} else if os.Getenv(config.EnvRulesPath) != "" {
			fmt.Printf("Using rules path from %s environment variable\n", config.EnvRulesPath)
		} else {
			fmt.Println("Using current directory as rules path")
		}

		if *targetPath != "" {
			fmt.Println("Using target path from command line flag")
		} else if os.Getenv(config.EnvTargetPath) != "" {
			fmt.Printf("Using target path from %s environment variable\n", config.EnvTargetPath)
		} else {
			fmt.Println("Using current directory as target path")
		}

		if *gitRepoURL != "" {
			fmt.Println("Using Git repository URL from command line flag")
		} else if os.Getenv(config.EnvGitRepoURL) != "" {
			fmt.Printf("Using Git repository URL from %s environment variable\n", config.EnvGitRepoURL)
		}

		if *editor != "" {
			fmt.Println("Using editor from command line flag")
		} else if os.Getenv(config.EnvEditor) != "" {
			fmt.Printf("Using editor from %s environment variable\n", config.EnvEditor)
		} else {
			fmt.Println("Using default editor: cursor")
		}

		fmt.Printf("Rules will be linked to: %s\n", filepath.Join(cfg.TargetProjectPath, cfg.GetRulesDirectory()))
	}

	// Handle Git repository if URL is provided
	if cfg.UseGitRepo {
		if *verbose {
			fmt.Println("Cloning Git repository...")
		}

		err := cfg.GitRepo.Clone()
		if err != nil {
			fmt.Printf("Error cloning Git repository: %v\n", err)
			os.Exit(1)
		}

		// Register cleanup function to remove the cloned repository on exit
		defer func() {
			if *verbose {
				fmt.Println("Cleaning up cloned Git repository...")
			}
			err := cfg.GitRepo.Cleanup()
			if err != nil && *verbose {
				fmt.Printf("Error cleaning up Git repository: %v\n", err)
			}
		}()

		if *verbose {
			fmt.Printf("Git repository cloned to: %s\n", cfg.GitRepo.CloneDir)
		}
	} else {
		// Validate local rules repository path only if not using Git
		if !cfg.ValidateRulesRepoPath() {
			fmt.Printf("Invalid rules repository path: %s\n", cfg.RulesRepoPath)
			os.Exit(1)
		}
	}

	// Validate target project path
	if !cfg.ValidateTargetProjectPath() {
		fmt.Printf("Invalid target project path: %s\n", cfg.TargetProjectPath)
		os.Exit(1)
	}

	// If verbose, show the resolved absolute paths
	if *verbose {
		if !cfg.UseGitRepo {
			fmt.Printf("Resolved rules path: %s\n", cfg.RulesRepoPath)
		}
		fmt.Printf("Resolved target path: %s\n", cfg.TargetProjectPath)
		if cfg.UseGitRepo {
			fmt.Printf("Using Git repository: %s\n", cfg.GitRepoURL)
		}
	}

	// Initialize rules manager
	rulesDir := cfg.GetRulesDir()
	rulesManager := rules.NewManager(rulesDir)

	// Load rules
	err := rulesManager.LoadRules()
	if err != nil {
		fmt.Printf("Error loading rules: %v\n", err)
		os.Exit(1)
	}

	if len(rulesManager.Rules) == 0 {
		fmt.Println("No rules found in repository")
		os.Exit(1)
	}

	// Initialize linker
	linkerInstance := linker.NewLinker(cfg.TargetProjectPath)
	// Set the config for the linker
	linkerInstance.SetConfig(cfg)

	// If dry run is enabled, set it on the linker
	if *dryRun {
		linkerInstance.SetDryRun(true)
	}

	// If verbose is enabled, set it on the linker
	if *verbose {
		linkerInstance.SetVerbose(true)
	}

	// Check which rules are already installed and mark them as selected
	for _, rule := range rulesManager.Rules {
		rule.IsInstalled = linkerInstance.IsRuleLinked(rule)
		// Initialize Selected to match IsInstalled as a starting point
		rule.Selected = rule.IsInstalled
	}

	// Define styles based on mode
	var titleStyle, ruleNameStyle, descStyle lipgloss.Style

	if *nonInteractive {
		// Plain text for non-interactive mode
		titleStyle = lipgloss.NewStyle()
		ruleNameStyle = lipgloss.NewStyle()
		descStyle = lipgloss.NewStyle()
	} else {
		// Colorful styles for interactive mode
		titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FFFF"))
		ruleNameStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF69B4"))
		descStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEEB"))
	}

	// Common header
	fmt.Println(titleStyle.Render("Rule Tool CLI"))
	fmt.Println(titleStyle.Render("---------------"))
	// When using a Git repo, show the cloned directory path instead of the local rules repo path
	if cfg.UseGitRepo && cfg.GitRepo != nil && cfg.GitRepo.IsCloned() {
		fmt.Printf("Rules repository: %s\n", cfg.GitRepo.CloneDir)
	} else {
		fmt.Printf("Rules repository: %s\n", cfg.RulesRepoPath)
	}
	fmt.Printf("Target project: %s\n", cfg.TargetProjectPath)
	fmt.Printf("Editor: %s\n", cfg.Editor)
	fmt.Printf("Rules directory: %s\n", filepath.Join(cfg.TargetProjectPath, cfg.GetRulesDirectory()))
	fmt.Printf("Found %d rules\n", len(rulesManager.Rules))

	// If no editor is set and we're in interactive mode, prompt for editor selection
	if !*nonInteractive && *linkRule == "" && *unlinkRule == "" && !*listRules && *editor == "" && os.Getenv(config.EnvEditor) == "" {
		// Prompt for editor selection
		fmt.Println("\nPlease select your editor:")
		fmt.Println("1. Cursor (default)")
		fmt.Println("2. Windsurf")
		fmt.Print("Enter your choice (1-2): ")

		var choice int
		_, err := fmt.Scanf("%d", &choice)
		if err != nil || choice < 1 || choice > 2 {
			// Default to Cursor if invalid input
			choice = 1
			fmt.Println("Invalid choice, defaulting to Cursor")
		}

		switch choice {
		case 1:
			cfg.SetEditor(string(config.EditorCursor))
		case 2:
			cfg.SetEditor(string(config.EditorWindsurf))
		}

		// Update the linker with the new config
		linkerInstance.SetConfig(cfg)

		fmt.Printf("Selected editor: %s\n", cfg.Editor)
		fmt.Printf("Rules will be linked to: %s\n", filepath.Join(cfg.TargetProjectPath, cfg.GetRulesDirectory()))
	}

	// Handle non-interactive modes if requested
	if *nonInteractive || *listRules || *linkRule != "" || *unlinkRule != "" {
		// List all available rules
		if *listRules || *nonInteractive {
			fmt.Println("\nAvailable Rules:")
			for i, rule := range rulesManager.Rules {
				ruleName := rule.Name
				if rule.Topic != "" {
					ruleName = rule.Topic + "/" + rule.Name
				}
				fmt.Printf("%d. %s: %s\n",
					i+1,
					ruleNameStyle.Render(ruleName),
					descStyle.Render(rule.Description))
			}
		}

		// Link specific rules
		if *linkRule != "" {
			rulesToLink := strings.Split(*linkRule, ",")
			for _, ruleName := range rulesToLink {
				ruleName = strings.TrimSpace(ruleName)
				rule := rulesManager.GetRuleByName(ruleName)
				if rule != nil {
					if *verbose {
						fmt.Printf("Debug - Rule found: %s\n", ruleName)
						fmt.Printf("Debug - Rule topic: %s\n", rule.Topic)
						fmt.Printf("Debug - Rule path: %s\n", rule.Path)
					}

					if *dryRun {
						fmt.Printf("Would link rule: %s\n", ruleName)

						// Display subfolder structure if applicable
						if rule.Topic != "" && *verbose {
							fmt.Printf("Would maintain subfolder structure: %s\n", rule.Topic)
						}
					} else {
						err := linkerInstance.LinkRule(rule)
						if err != nil {
							fmt.Printf("Error linking rule %s: %v\n", ruleName, err)
						} else {
							fmt.Printf("Linked rule: %s\n", ruleName)
						}
					}
				} else {
					fmt.Printf("Rule not found: %s\n", ruleName)
				}
			}
		}

		// Unlink specific rules
		if *unlinkRule != "" {
			rulesToUnlink := strings.Split(*unlinkRule, ",")
			for _, ruleName := range rulesToUnlink {
				ruleName = strings.TrimSpace(ruleName)
				if *dryRun {
					fmt.Printf("Would unlink rule: %s\n", ruleName)
				} else {
					err := linkerInstance.UnlinkRule(ruleName)
					if err != nil {
						fmt.Printf("Error unlinking rule %s: %v\n", ruleName, err)
					} else {
						fmt.Printf("Unlinked rule: %s\n", ruleName)
					}
				}
			}
		}

		return
	}

	// Interactive mode - Initialize UI model
	model := ui.New(cfg, rulesManager, linkerInstance)

	// Log paths for debugging
	fmt.Printf("Starting UI with Rules repository: %s\n", cfg.RulesRepoPath)
	fmt.Printf("Starting UI with Target project: %s\n", cfg.TargetProjectPath)

	// Run the application with full screen and mouse support
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),       // Use alternate screen buffer
		tea.WithMouseCellMotion(), // Enable mouse support
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
