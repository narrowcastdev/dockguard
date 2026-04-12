package report

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/narrowcastdev/dockguard/internal/rule"
)

// isTTY returns true if the given writer is an interactive terminal.
func isTTY(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	fi, err := f.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// ANSI color codes.
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
)

// PrintTerminal writes a human-readable report to w.
func PrintTerminal(w io.Writer, findings []rule.Finding, serviceCount int, version string) {
	color := isTTY(w)

	// Header.
	if color {
		fmt.Fprintf(w, "\n🛡️  %sdockguard %s%s — Docker Compose Security Scanner\n\n", colorBold, version, colorReset)
	} else {
		fmt.Fprintf(w, "\ndockguard %s — Docker Compose Security Scanner\n\n", version)
	}

	if len(findings) == 0 {
		if color {
			fmt.Fprintf(w, "✅ %d services scanned — no issues found!\n", serviceCount)
		} else {
			fmt.Fprintf(w, "%d services scanned — no issues found!\n", serviceCount)
		}
		return
	}

	// Sort findings by severity (critical first).
	sorted := make([]rule.Finding, len(findings))
	copy(sorted, findings)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Severity > sorted[j].Severity
	})

	// Group by severity.
	groups := map[rule.Severity][]rule.Finding{}
	for _, f := range sorted {
		groups[f.Severity] = append(groups[f.Severity], f)
	}

	// Print each group.
	for _, sev := range []rule.Severity{rule.Critical, rule.Warning, rule.Info} {
		items := groups[sev]
		if len(items) == 0 {
			continue
		}
		printSeverityHeader(w, sev, color)
		for _, f := range items {
			printFinding(w, f, color)
		}
		fmt.Fprintln(w)
	}

	// Summary.
	critCount := len(groups[rule.Critical])
	warnCount := len(groups[rule.Warning])
	infoCount := len(groups[rule.Info])

	if color {
		fmt.Fprintf(w, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
		fmt.Fprintf(w, "📊 %d services scanned — 🔴 %d critical · 🟡 %d warnings · 🔵 %d info\n",
			serviceCount, critCount, warnCount, infoCount)
		fmt.Fprintf(w, "💡 Run with --fix to generate a hardened compose file.\n\n")
	} else {
		fmt.Fprintf(w, "----------------------------------------------------\n")
		fmt.Fprintf(w, "%d services scanned — %d critical, %d warnings, %d info\n",
			serviceCount, critCount, warnCount, infoCount)
		fmt.Fprintf(w, "Run with --fix to generate a hardened compose file.\n\n")
	}
}

func printSeverityHeader(w io.Writer, sev rule.Severity, color bool) {
	if color {
		switch sev {
		case rule.Critical:
			fmt.Fprintf(w, "🔴 %s%sCRITICAL%s\n", colorBold, colorRed, colorReset)
		case rule.Warning:
			fmt.Fprintf(w, "🟡 %s%sWARNING%s\n", colorBold, colorYellow, colorReset)
		case rule.Info:
			fmt.Fprintf(w, "🔵 %s%sINFO%s\n", colorBold, colorBlue, colorReset)
		}
		return
	}
	fmt.Fprintf(w, "%s\n", sev.String())
}

func printFinding(w io.Writer, f rule.Finding, color bool) {
	icon := "  "
	switch f.Severity {
	case rule.Critical:
		icon = "  ✗"
	case rule.Warning:
		icon = "  ⚠"
	case rule.Info:
		icon = "  ℹ"
	}

	if color {
		var c string
		switch f.Severity {
		case rule.Critical:
			c = colorRed
		case rule.Warning:
			c = colorYellow
		case rule.Info:
			c = colorBlue
		}
		fmt.Fprintf(w, "%s%s %s: %s%s  %s[%s]%s\n",
			c, icon, f.Service, f.Message, colorReset, colorDim, f.Rule, colorReset)
		return
	}
	fmt.Fprintf(w, "%s %s: %s  [%s]\n", icon, f.Service, f.Message, f.Rule)
}

// PrintJSON writes findings as a JSON array to w.
func PrintJSON(w io.Writer, findings []rule.Finding) error {
	type jsonFinding struct {
		Rule       string `json:"rule"`
		Service    string `json:"service"`
		Severity   string `json:"severity"`
		Message    string `json:"message"`
		Suggestion string `json:"suggestion"`
	}

	out := make([]jsonFinding, len(findings))
	for i, f := range findings {
		out[i] = jsonFinding{
			Rule:       f.Rule,
			Service:    f.Service,
			Severity:   f.Severity.String(),
			Message:    f.Message,
			Suggestion: f.Suggestion,
		}
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
