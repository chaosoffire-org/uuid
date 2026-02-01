// scripts/fuzz_runner.go
// Cross-platform fuzzing test runner for UUID package
// A professional CLI tool built with Cobra

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// FuzzTest represents a discovered fuzz test
type FuzzTest struct {
	Package  string `json:"package"`
	Name     string `json:"name"`
	FilePath string `json:"file_path"`
}

// TestResult holds the result of a single fuzz test
type TestResult struct {
	Test     FuzzTest      `json:"test"`
	Status   string        `json:"status"` // "pass", "fail", "skip", "timeout"
	Duration time.Duration `json:"duration_ns"`
	Error    string        `json:"error,omitempty"`
}

// Config holds the runner configuration
type Config struct {
	FuzzTime   time.Duration
	Exclude    []string
	Verbose    bool
	Parallel   int
	CoverDir   string
	ReportFile string
	DryRun     bool
	SingleTest string
	RootDir    string
}

var (
	config          Config
	fuzzTestPattern = regexp.MustCompile(`func\s+(Fuzz\w+)\s*\(\s*\w+\s+\*testing\.F\s*\)`)
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "fuzz-runner",
	Short: "Cross-platform fuzzing test runner",
	Long: `A professional CLI tool to discover and run all fuzz tests in your Go project.

Supports multiple platforms (Windows, Linux, macOS) and provides detailed
reporting, coverage collection, and flexible test filtering.`,
	Example: `  # List all fuzz tests without running them
  fuzz-runner list

  # Run all fuzz tests with 10 second duration each
  fuzz-runner run --fuzztime 10s

  # Run a specific fuzz test
  fuzz-runner run --test FuzzParse --fuzztime 1m

  # Run with verbose output and save report
  fuzz-runner run -v --fuzztime 30s --report results.json`,
	Version: "1.0.0",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all discovered fuzz tests",
	Long:  "Scan the project and list all fuzz test functions without running them.",
	RunE:  runList,
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run fuzz tests",
	Long:  "Discover and execute fuzz tests one by one with the specified duration.",
	RunE:  runFuzz,
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Display system information",
	Long:  "Show system and Go runtime information useful for debugging.",
	Run:   runInfo,
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&config.RootDir, "root", "r", "", "Project root directory (auto-detected if not specified)")
	rootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v", false, "Enable verbose output")

	// Run command flags
	runCmd.Flags().DurationVarP(&config.FuzzTime, "fuzztime", "t", 10*time.Second, "Duration to run each fuzz test (e.g., 10s, 1m, 5m)")
	runCmd.Flags().StringSliceVarP(&config.Exclude, "exclude", "e", nil, "Test names to exclude (comma-separated)")
	runCmd.Flags().IntVarP(&config.Parallel, "parallel", "p", runtime.NumCPU(), "Number of parallel fuzzing processes")
	runCmd.Flags().StringVarP(&config.CoverDir, "dir", "d", "", "Directory to save coverage data")
	runCmd.Flags().StringVarP(&config.ReportFile, "output", "o", "", "Path to save the test report (JSON)")
	runCmd.Flags().StringVarP(&config.SingleTest, "filter", "f", "", "Run only the specified fuzz test")

	// List command flags
	listCmd.Flags().StringVarP(&config.SingleTest, "regex", "x", "", "Filter tests by name pattern")

	rootCmd.AddCommand(listCmd, runCmd, infoCmd)
}

func initConfig() error {
	if config.RootDir == "" {
		root, err := findProjectRoot()
		if err != nil {
			return fmt.Errorf("could not find project root: %w\nUse --root flag to specify the project root directory", err)
		}
		config.RootDir = root
	}
	return nil
}

func runInfo(cmd *cobra.Command, args []string) {
	printBanner()
	printSystemInfo()
}

func runList(cmd *cobra.Command, args []string) error {
	printBanner()

	if err := initConfig(); err != nil {
		return err
	}

	tests, err := discoverFuzzTests()
	if err != nil {
		return fmt.Errorf("error discovering fuzz tests: %w", err)
	}

	if len(tests) == 0 {
		fmt.Println("⚠️  No fuzz tests found!")
		return nil
	}

	fmt.Printf("📋 Found %d fuzz test(s):\n\n", len(tests))
	fmt.Println("┌────┬─────────────────────────────┬──────────────────┐")
	fmt.Println("│ #  │ Test Name                   │ Package          │")
	fmt.Println("├────┼─────────────────────────────┼──────────────────┤")
	for i, test := range tests {
		fmt.Printf("│ %-2d │ %-27s │ %-16s │\n", i+1, truncate(test.Name, 27), truncate(test.Package, 16))
	}
	fmt.Println("└────┴─────────────────────────────┴──────────────────┘")

	return nil
}

func runFuzz(cmd *cobra.Command, args []string) error {
	printBanner()

	if err := initConfig(); err != nil {
		return err
	}

	printSystemInfo()

	tests, err := discoverFuzzTests()
	if err != nil {
		return fmt.Errorf("error discovering fuzz tests: %w", err)
	}

	if len(tests) == 0 {
		fmt.Println("⚠️  No fuzz tests found!")
		return nil
	}

	fmt.Printf("📋 Discovered %d fuzz test(s)\n", len(tests))
	fmt.Printf("⏱️  Fuzz time per test: %s\n\n", config.FuzzTime)

	results := executeFuzzTests(tests)

	printSummary(results)

	if config.ReportFile != "" {
		if err := writeReport(results); err != nil {
			return fmt.Errorf("error writing report: %w", err)
		}
	}

	// Exit with error if any test failed
	for _, r := range results {
		if r.Status == "fail" {
			return fmt.Errorf("some fuzz tests failed")
		}
	}

	return nil
}

func printBanner() {
	fmt.Println()
	fmt.Println("  ╔═══════════════════════════════════════════════════════╗")
	fmt.Println("  ║          🧪 UUID Fuzz Test Runner v1.0.0              ║")
	fmt.Println("  ╚═══════════════════════════════════════════════════════╝")
	fmt.Println()
}

func printSystemInfo() {
	fmt.Println("  📊 System Information")
	fmt.Println("  ─────────────────────")
	fmt.Printf("     OS:        %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("     CPUs:      %d\n", runtime.NumCPU())
	fmt.Printf("     Go:        %s\n", runtime.Version())
	if config.RootDir != "" {
		fmt.Printf("     Project:   %s\n", config.RootDir)
	}
	fmt.Println()
}

func findProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	dir := wd
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.work")); err == nil {
			return dir, nil
		}
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			parent := filepath.Dir(dir)
			if parent != dir {
				if _, err := os.Stat(filepath.Join(parent, "go.work")); err == nil {
					return parent, nil
				}
			}
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("could not find go.mod or go.work in any parent directory")
}

func discoverFuzzTests() ([]FuzzTest, error) {
	var tests []FuzzTest

	fmt.Printf("  📂 Scanning: %s\n\n", config.RootDir)

	err := filepath.Walk(config.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			name := info.Name()
			if name == "vendor" || name == ".git" || name == "testdata" || name == "scripts" {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		matches := fuzzTestPattern.FindAllStringSubmatch(string(content), -1)
		for _, match := range matches {
			testName := match[1]

			if isExcluded(testName) {
				continue
			}

			if config.SingleTest != "" && !strings.Contains(testName, config.SingleTest) {
				continue
			}

			relPath, _ := filepath.Rel(config.RootDir, filepath.Dir(path))
			pkgPath := "./" + filepath.ToSlash(relPath)
			if relPath == "." {
				pkgPath = "."
			}

			tests = append(tests, FuzzTest{
				Package:  pkgPath,
				Name:     testName,
				FilePath: path,
			})
		}

		return nil
	})

	return tests, err
}

func isExcluded(name string) bool {
	for _, exc := range config.Exclude {
		if name == exc || strings.Contains(name, exc) {
			return true
		}
	}
	return false
}

func executeFuzzTests(tests []FuzzTest) []TestResult {
	results := make([]TestResult, 0, len(tests))
	totalTests := len(tests)

	for i, test := range tests {
		fmt.Println("  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
		fmt.Printf("  🧪 [%d/%d] %s\n", i+1, totalTests, test.Name)
		fmt.Printf("     📦 Package: %s\n", test.Package)
		fmt.Printf("     ⏱️  Duration: %s\n", config.FuzzTime)
		fmt.Println("  ─────────────────────────────────────────────────────────")

		result := runSingleFuzzTest(test)
		results = append(results, result)

		switch result.Status {
		case "pass":
			fmt.Printf("  ✅ PASS (%s)\n\n", result.Duration.Round(time.Millisecond))
		case "fail":
			fmt.Printf("  ❌ FAIL (%s)\n", result.Duration.Round(time.Millisecond))
			if result.Error != "" && config.Verbose {
				fmt.Printf("     Error: %s\n", result.Error)
			}
			fmt.Println()
		case "skip":
			fmt.Printf("  ⏭️  SKIP: %s\n\n", result.Error)
		case "timeout":
			fmt.Printf("  ⏰ TIMEOUT (%s)\n\n", result.Duration.Round(time.Millisecond))
		}
	}

	return results
}

func runSingleFuzzTest(test FuzzTest) TestResult {
	start := time.Now()

	args := []string{
		"test",
		"-fuzz", fmt.Sprintf("^%s$", test.Name),
		"-fuzztime", config.FuzzTime.String(),
		"-run", fmt.Sprintf("^%s$", test.Name),
	}

	if config.Parallel > 0 {
		args = append(args, "-parallel", fmt.Sprintf("%d", config.Parallel))
	}

	if config.CoverDir != "" {
		args = append(args, "-test.gocoverdir", config.CoverDir)
	}

	if config.Verbose {
		args = append(args, "-v")
	}

	args = append(args, test.Package)

	cmd := exec.Command("go", args...)
	cmd.Dir = config.RootDir

	var outputBuilder strings.Builder
	if config.Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = &outputBuilder
		cmd.Stderr = &outputBuilder
	}

	err := cmd.Run()
	duration := time.Since(start)

	result := TestResult{
		Test:     test,
		Duration: duration,
	}

	if err != nil {
		result.Status = "fail"
		if exitError, ok := err.(*exec.ExitError); ok {
			result.Error = fmt.Sprintf("exit code %d", exitError.ExitCode())
		} else {
			result.Error = err.Error()
		}
		if !config.Verbose {
			result.Error += "\n" + outputBuilder.String()
		}
	} else {
		result.Status = "pass"
	}

	return result
}

func printSummary(results []TestResult) {
	if len(results) == 0 {
		return
	}

	fmt.Println("  ═══════════════════════════════════════════════════════════")
	fmt.Println("                          📊 SUMMARY")
	fmt.Println("  ═══════════════════════════════════════════════════════════")

	var passed, failed, skipped, timeout int
	var totalDuration time.Duration

	for _, r := range results {
		totalDuration += r.Duration
		switch r.Status {
		case "pass":
			passed++
		case "fail":
			failed++
		case "skip":
			skipped++
		case "timeout":
			timeout++
		}
	}

	fmt.Println()
	fmt.Printf("     ✅ Passed:   %d\n", passed)
	fmt.Printf("     ❌ Failed:   %d\n", failed)
	fmt.Printf("     ⏭️  Skipped:  %d\n", skipped)
	fmt.Printf("     ⏰ Timeout:  %d\n", timeout)
	fmt.Println("     ───────────────────")
	fmt.Printf("     📋 Total:    %d\n", len(results))
	fmt.Printf("     ⏱️  Duration: %s\n", totalDuration.Round(time.Second))

	if failed > 0 {
		fmt.Println("\n     ❌ Failed Tests:")
		for _, r := range results {
			if r.Status == "fail" {
				fmt.Printf("        • %s (%s)\n", r.Test.Name, r.Test.Package)
			}
		}
	}

	fmt.Println()
	fmt.Println("  ═══════════════════════════════════════════════════════════")
}

func writeReport(results []TestResult) error {
	report := struct {
		Generated string       `json:"generated"`
		Platform  string       `json:"platform"`
		GoVersion string       `json:"go_version"`
		Results   []TestResult `json:"results"`
	}{
		Generated: time.Now().Format(time.RFC3339),
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
		GoVersion: runtime.Version(),
		Results:   results,
	}

	file, err := os.Create(config.ReportFile)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return err
	}
	writer.Flush()

	fmt.Printf("\n  📄 Report saved to: %s\n", config.ReportFile)
	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
