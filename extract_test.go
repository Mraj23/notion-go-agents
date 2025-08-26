package notion

import "testing"

func TestExtractNotionTitle(t *testing.T) {
	props := map[string]any{
		"Name": map[string]any{
			"title": []any{
				map[string]any{"plain_text": "Hello"},
			},
		},
	}
	if got := ExtractNotionTitle(props); got != "Hello" {
		t.Fatalf("expected Hello, got %q", got)
	}
}

func TestSelectPrintableProperties(t *testing.T) {
	props := map[string]any{
		"Name":    map[string]any{"title": []any{}},
		"Due":     map[string]any{"type": "date"},
		"Ignored": map[string]any{"type": "people"},
	}
	out := SelectPrintableProperties(props)
	if len(out) != 2 {
		t.Fatalf("expected 2 props, got %d", len(out))
	}
	if _, ok := out["Ignored"]; ok {
		t.Fatal("Ignored property should not be selected")
	}
}
