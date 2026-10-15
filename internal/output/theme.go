package output

import (
	"fmt"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/charmbracelet/x/term"
)

// ThemeEnv selects the terminal palette for CLI output and the standalone
// monitor: "dark" (default), "light", or "auto" (query the terminal).
const ThemeEnv = "TD_THEME"

// Palette holds the ANSI 256 colors used by plain CLI output. Each light value
// keeps the dark hue but is a deeper shade so it stays legible on white.
type Palette struct {
	Text       string
	Subtle     string
	Closed     string
	Success    string
	Error      string
	Warning    string
	Primary    string
	Secondary  string
	Accent     string
	IsDarkMode bool
}

// DarkPalette is td's original CLI palette.
var DarkPalette = Palette{
	Text: "255", Subtle: "241", Closed: "242", Success: "42", Error: "196",
	Warning: "214", Primary: "212", Secondary: "141", Accent: "45", IsDarkMode: true,
}

// LightPalette mirrors pkg/monitor.LightTheme for CLI output.
var LightPalette = Palette{
	Text: "235", Subtle: "244", Closed: "246", Success: "28", Error: "160",
	Warning: "130", Primary: "162", Secondary: "97", Accent: "24",
}

var currentPalette = DarkPalette

func init() {
	if name := os.Getenv(ThemeEnv); name != "" {
		resolved, err := ResolveThemeName(name)
		if err == nil {
			applyPalette(resolved)
		}
	}
}

// ResolveThemeName normalizes a theme setting to "dark" or "light". "auto"
// asks the terminal for its background color only when stdout is a terminal,
// so piped output never waits on a terminal query.
func ResolveThemeName(name string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", "dark":
		return "dark", nil
	case "light":
		return "light", nil
	case "auto":
		if !term.IsTerminal(os.Stdout.Fd()) || !term.IsTerminal(os.Stdin.Fd()) {
			return "dark", nil
		}
		if lipgloss.HasDarkBackground(os.Stdin, os.Stdout) {
			return "dark", nil
		}
		return "light", nil
	default:
		return "", fmt.Errorf("invalid theme %q (want dark, light, or auto)", name)
	}
}

// CurrentPalette returns the palette selected for this process.
func CurrentPalette() Palette {
	return currentPalette
}

func applyPalette(name string) {
	if name == "light" {
		currentPalette = LightPalette
	} else {
		currentPalette = DarkPalette
	}
	p := currentPalette
	subtleStyle = fg(p.Subtle)
	successStyle = fg(p.Success)
	errorStyle = fg(p.Error)
	warningStyle = fg(p.Warning)
	priorityStyle = fg(p.Primary)
	statusStyles = newStatusStyles(p)
}

func fg(c string) lipgloss.Style {
	return lipgloss.NewStyle().Foreground(lipgloss.Color(c))
}
