package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const repoAPI = "https://api.github.com/repos/kavinkumar999/kite-cli/releases/latest"

type githubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update Kite CLI to the latest version",
	Long: `Check for updates and install the latest version.

Examples:
  kite update          # Update to latest version
  kite update --check  # Check for updates without installing`,
	Run: runUpdate,
}

var checkOnly bool

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&checkOnly, "check", false, "Check for updates without installing")
}

func runUpdate(cmd *cobra.Command, args []string) {
	fmt.Println("Checking for updates...")

	latest, err := getLatestVersion()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
		os.Exit(1)
	}

	currentVersion := strings.TrimPrefix(Version, "v")
	latestVersion := strings.TrimPrefix(latest.TagName, "v")

	if currentVersion == latestVersion {
		fmt.Printf("✓ You're on the latest version (%s)\n", Version)
		return
	}

	if Version == "dev" {
		fmt.Println("You're running a development build.")
		fmt.Printf("Latest release: %s\n", latest.TagName)
		fmt.Printf("  → %s\n", latest.HTMLURL)
	} else {
		fmt.Printf("New version available: %s → %s\n", Version, latest.TagName)
	}

	if checkOnly {
		fmt.Printf("\nRun 'kite update' to install the latest version.\n")
		return
	}

	fmt.Println("\nUpdating...")

	if err := installLatest(); err != nil {
		fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
		fmt.Println("\nManual update:")
		fmt.Println("  curl -sSL https://raw.githubusercontent.com/kavinkumar999/kite-cli/main/install.sh | bash")
		os.Exit(1)
	}

	fmt.Println("✓ Update complete! Restart your terminal or run 'kite version' to verify.")
}

func getLatestVersion() (*githubRelease, error) {
	resp, err := http.Get(repoAPI)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return &githubRelease{TagName: "v0.0.0", HTMLURL: "https://github.com/kavinkumar999/kite-cli"}, nil
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func installLatest() error {
	// Check if Go is available for building from source
	if _, err := exec.LookPath("go"); err == nil {
		return buildFromSource()
	}

	// Fallback to install script
	return runInstallScript()
}

func buildFromSource() error {
	fmt.Println("Building from source...")

	tmpDir := os.TempDir() + "/kite-cli-update"
	os.RemoveAll(tmpDir)

	// Clone
	cmd := exec.Command("git", "clone", "--depth", "1", "https://github.com/kavinkumar999/kite-cli.git", tmpDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone failed: %w", err)
	}

	// Build
	buildCmd := exec.Command("go", "build", "-ldflags=-s -w", "-o", "kite", ".")
	buildCmd.Dir = tmpDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	// Install
	home, _ := os.UserHomeDir()
	installPath := home + "/bin/kite"

	src := tmpDir + "/kite"
	if runtime.GOOS == "windows" {
		src += ".exe"
		installPath += ".exe"
	}

	// Copy binary
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.WriteFile(installPath, input, 0755); err != nil {
		return err
	}

	// Cleanup
	os.RemoveAll(tmpDir)

	return nil
}

func runInstallScript() error {
	fmt.Println("Running install script...")

	cmd := exec.Command("bash", "-c", "curl -sSL https://raw.githubusercontent.com/kavinkumar999/kite-cli/main/install.sh | bash")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
