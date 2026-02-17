package providerconfig

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

// LoginConfig holds configuration for a provider's login flow.
type LoginConfig struct {
	Name        string // Display name (e.g., "Wrike", "Jira")
	EnvVar      string // Environment variable name (e.g., "WRIKE_TOKEN")
	TokenField  string // Field name in config (e.g., "token")
	HelpURL     string // URL for getting a token
	HelpSteps   string // Navigation steps to get a token
	Scopes      string // Required permissions/scopes
	TokenPrefix string // Optional token prefix for validation (e.g., "ghp_")
}

// LoginOptions configures the login behavior.
type LoginOptions struct {
	Force  bool      // Override existing token without prompting
	Stdin  io.Reader // Input reader (defaults to os.Stdin)
	Stdout io.Writer // Output writer (defaults to os.Stdout)
}

// LoginResult contains the result of a login attempt.
type LoginResult struct {
	Token      string // The entered token
	ConfigPath string // Path where config was saved
	Cancelled  bool   // True if user cancelled
}

// RunLogin executes an interactive login flow for a provider.
func RunLogin(ctx context.Context, mgr ConfigManager, lc LoginConfig, opts LoginOptions) (*LoginResult, error) {
	if opts.Stdin == nil {
		opts.Stdin = os.Stdin
	}
	if opts.Stdout == nil {
		opts.Stdout = os.Stdout
	}

	// Check for existing token
	if mgr.Exists() && !opts.Force {
		cfg, err := mgr.Read(ctx)
		if err == nil {
			token := cfg.GetString(lc.TokenField)
			if token != "" && !strings.HasPrefix(token, "${") {
				masked := MaskToken(token)
				fmt.Fprintf(opts.Stdout, "Token already configured: %s\n", masked)
				fmt.Fprint(opts.Stdout, "Override? [y/N]: ")

				reader := bufio.NewReader(opts.Stdin)
				response, _ := reader.ReadString('\n')
				response = strings.TrimSpace(strings.ToLower(response))

				if response != "y" && response != "yes" {
					return &LoginResult{Cancelled: true}, nil
				}
			}
		}
	}

	// Print help
	PrintTokenHelp(opts.Stdout, lc)

	// Prompt for token
	token, err := PromptToken(opts.Stdin, opts.Stdout, lc)
	if err != nil {
		return nil, err
	}
	if token == "" {
		return &LoginResult{Cancelled: true}, nil
	}

	// Load existing config or create new
	var cfg Config
	if mgr.Exists() {
		cfg, err = mgr.Read(ctx)
		if err != nil {
			cfg = NewConfig()
		}
	} else {
		cfg = NewConfig()
	}

	// Update token
	cfg = cfg.Set(lc.TokenField, token)

	// Write config
	if err := mgr.Write(ctx, cfg); err != nil {
		return nil, fmt.Errorf("save config: %w", err)
	}

	return &LoginResult{
		Token:      token,
		ConfigPath: mgr.Path(),
	}, nil
}

// PrintTokenHelp displays formatted guidance for getting a token.
func PrintTokenHelp(w io.Writer, lc LoginConfig) {
	fmt.Fprintln(w)
	fmt.Fprintf(w, "%s Token Setup\n", lc.Name)
	fmt.Fprintln(w, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintf(w, "📍 Get token: %s\n", lc.HelpURL)
	if lc.HelpSteps != "" {
		fmt.Fprintf(w, "📋 Steps:     %s\n", lc.HelpSteps)
	}
	if lc.Scopes != "" {
		fmt.Fprintf(w, "🔑 Required:  %s\n", lc.Scopes)
	}
	if lc.TokenPrefix != "" {
		fmt.Fprintf(w, "💡 Format:    Token starts with '%s'\n", lc.TokenPrefix)
	}
	fmt.Fprintln(w, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Fprintln(w)
}

// PromptToken interactively prompts the user for a token.
// Returns empty string if cancelled.
func PromptToken(r io.Reader, w io.Writer, lc LoginConfig) (string, error) {
	fmt.Fprintf(w, "Enter your %s API token (leave empty to cancel): ", lc.Name)

	// Try to read password securely if stdin is a terminal
	if f, ok := r.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		tokenBytes, err := term.ReadPassword(int(f.Fd()))
		fmt.Fprintln(w) // newline after hidden input
		if err != nil {
			return "", fmt.Errorf("read token: %w", err)
		}
		token := strings.TrimSpace(string(tokenBytes))
		if token != "" && lc.TokenPrefix != "" && !strings.HasPrefix(token, lc.TokenPrefix) {
			fmt.Fprintf(w, "Warning: Token doesn't start with expected prefix '%s'\n", lc.TokenPrefix)
		}
		return token, nil
	}

	// Fallback to regular input (for piped input)
	reader := bufio.NewReader(r)
	token, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("read token: %w", err)
	}
	token = strings.TrimSpace(token)
	if token != "" && lc.TokenPrefix != "" && !strings.HasPrefix(token, lc.TokenPrefix) {
		fmt.Fprintf(w, "Warning: Token doesn't start with expected prefix '%s'\n", lc.TokenPrefix)
	}
	return token, nil
}

// MaskToken returns a masked version of a token for display.
func MaskToken(token string) string {
	if len(token) <= 8 {
		return "*******"
	}
	return token[:4] + "..." + token[len(token)-4:]
}

// DetectExistingToken checks if a token is already available from environment or config.
// Returns the source description and masked value, or empty strings if not found.
func DetectExistingToken(ctx context.Context, mgr ConfigManager, lc LoginConfig) (source, masked string) {
	// Check environment variable
	if val := os.Getenv(lc.EnvVar); val != "" {
		return lc.EnvVar + " environment variable", MaskToken(val)
	}

	// Check config file
	if mgr.Exists() {
		cfg, err := mgr.Read(ctx)
		if err == nil {
			token := cfg.GetString(lc.TokenField)
			if token != "" && !strings.HasPrefix(token, "${") {
				return "config file", MaskToken(token)
			}
		}
	}

	return "", ""
}
