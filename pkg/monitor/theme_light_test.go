package monitor

import (
	"reflect"
	"testing"
)

func TestLightThemeIsCompleteAndValid(t *testing.T) {
	light := LightTheme()
	normalized, err := normalizedTheme(light)
	if err != nil {
		t.Fatalf("LightTheme has an invalid color: %v", err)
	}
	if normalized != light {
		t.Fatalf("LightTheme should not inherit dark defaults: got %+v", normalized)
	}

	// Every color slot must be set explicitly so a light monitor never falls
	// back to a dark-background value.
	v := reflect.ValueOf(light)
	for i := 0; i < v.NumField(); i++ {
		name := v.Type().Field(i).Name
		if name == "SyntaxTheme" {
			continue
		}
		if v.Field(i).String() == "" {
			t.Errorf("LightTheme.%s is empty", name)
		}
	}
}

func TestSetThemeLightUsesLightMarkdownPalette(t *testing.T) {
	m := Model{}
	if err := m.SetTheme(LightTheme()); err != nil {
		t.Fatal(err)
	}
	if got, want := m.MarkdownTheme.Colors.Text, themeColorHex(LightTheme().TextPrimary); got != want {
		t.Fatalf("markdown text color = %s, want %s", got, want)
	}
}
