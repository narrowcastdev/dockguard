package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/narrowcastdev/dockguard/internal/compose"
	"github.com/narrowcastdev/dockguard/internal/hardener"
	"github.com/narrowcastdev/dockguard/internal/report"
	"github.com/narrowcastdev/dockguard/internal/rule"

	// Register all rules via init().
	_ "github.com/narrowcastdev/dockguard/internal/rule"
)

var version = "dev"

func main() {
	exitCode := run()
	os.Exit(exitCode)
}

func run() int {
	var (
		fix        bool
		outputPath string
		jsonOutput bool
		severity   string
	)

	flag.BoolVar(&fix, "fix", false, "generate a hardened compose file")
	flag.StringVar(&outputPath, "o", "docker-compose.hardened.yml", "output path for hardened file")
	flag.BoolVar(&jsonOutput, "json", false, "output findings as JSON")
	flag.StringVar(&severity, "severity", "info", "minimum severity to display (info, warning, critical)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: dockguard [flags] <docker-compose.yml>\n\nFlags:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		fmt.Fprintf(os.Stderr, "\ndockguard: missing required argument: path to docker-compose.yml\n")
		return 1
	}

	inputPath := flag.Arg(0)

	minSeverity, err := parseSeverity(severity)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dockguard: %v\n", err)
		return 1
	}

	f, err := compose.ParseFile(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dockguard: %v\n", err)
		return 1
	}

	findings := rule.Run(f)

	// Filter by minimum severity.
	filtered := filterBySeverity(findings, minSeverity)

	// Output report.
	if jsonOutput {
		if err := report.PrintJSON(os.Stdout, filtered); err != nil {
			fmt.Fprintf(os.Stderr, "dockguard: writing JSON: %v\n", err)
			return 1
		}
	} else {
		report.PrintTerminal(os.Stdout, filtered, len(f.Services), version)
	}

	// Generate hardened file if requested.
	if fix {
		fixes := rule.Fixes(f)
		if err := hardener.HardenFile(inputPath, outputPath, version, fixes); err != nil {
			fmt.Fprintf(os.Stderr, "dockguard: %v\n", err)
			return 1
		}
		if !jsonOutput {
			fmt.Fprintf(os.Stderr, "✅ Hardened file written to %s\n", outputPath)
		}
	}

	// Exit code based on max severity of ALL findings (not filtered).
	switch rule.MaxSeverity(findings) {
	case rule.Critical:
		return 2
	case rule.Warning:
		return 1
	}

	return 0
}

func parseSeverity(s string) (rule.Severity, error) {
	switch s {
	case "info":
		return rule.Info, nil
	case "warning":
		return rule.Warning, nil
	case "critical":
		return rule.Critical, nil
	}
	return rule.Info, fmt.Errorf("invalid severity: %q (must be info, warning, or critical)", s)
}

func filterBySeverity(findings []rule.Finding, min rule.Severity) []rule.Finding {
	var result []rule.Finding
	for _, f := range findings {
		if f.Severity >= min {
			result = append(result, f)
		}
	}
	return result
}
