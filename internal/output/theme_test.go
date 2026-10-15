package output

import (
	"testing"

	"github.com/marcus/td/internal/models"
)

func TestResolveThemeName(t *testing.T) {
	for in, want := range map[string]string{"": "dark", "dark": "dark", "Light": "light", " light ": "light"} {
		got, err := ResolveThemeName(in)
		if err != nil || got != want {
			t.Errorf("ResolveThemeName(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	if _, err := ResolveThemeName("solarized"); err == nil {
		t.Error("expected an error for an unknown theme")
	}
}

func TestApplyPaletteSwitchesStyles(t *testing.T) {
	t.Cleanup(func() { applyPalette("dark") })

	applyPalette("light")
	if CurrentPalette() != LightPalette {
		t.Fatal("light palette not selected")
	}
	if got := statusStyles[models.StatusOpen].GetForeground(); got != fg(LightPalette.Accent).GetForeground() {
		t.Errorf("open status color not updated: %v", got)
	}

	applyPalette("dark")
	if CurrentPalette() != DarkPalette {
		t.Fatal("dark palette not restored")
	}
}
