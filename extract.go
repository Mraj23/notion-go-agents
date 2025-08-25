package notion

import "strings"

func ExtractNotionTitle(props map[string]any) string {
	if props == nil {
		return ""
	}
	for _, v := range props {
		m, ok := v.(map[string]any)
		if !ok {
			continue
		}
		if t, _ := m["type"].(string); t == "title" {
			if arr, ok := m["title"].([]any); ok {
				var b strings.Builder
				for _, it := range arr {
					im, _ := it.(map[string]any)
					if pt, _ := im["plain_text"].(string); pt != "" {
						b.WriteString(pt)
					}
				}
				return strings.TrimSpace(b.String())
			}
		}
	}
	for _, key := range []string{"Name", "Title", "name", "title"} {
		if v, ok := props[key]; ok {
			m, _ := v.(map[string]any)
			if arr, ok := m["title"].([]any); ok {
				var b strings.Builder
				for _, it := range arr {
					im, _ := it.(map[string]any)
					if pt, _ := im["plain_text"].(string); pt != "" {
						b.WriteString(pt)
					}
				}
				return strings.TrimSpace(b.String())
			}
			if arr, ok := m["rich_text"].([]any); ok {
				var b strings.Builder
				for _, it := range arr {
					im, _ := it.(map[string]any)
					if pt, _ := im["plain_text"].(string); pt != "" {
						b.WriteString(pt)
					}
				}
				return strings.TrimSpace(b.String())
			}
		}
	}
	return ""
}

func SelectPrintableProperties(props map[string]any) map[string]any {
	out := make(map[string]any)
	for k, v := range props {
		if k == "title" || k == "Name" || k == "Title" {
			out[k] = v
			continue
		}
		switch vv := v.(type) {
		case map[string]any:
			if t, _ := vv["type"].(string); t == "status" || t == "date" || t == "number" || t == "url" || t == "select" || t == "multi_select" {
				out[k] = vv
			}
		}
	}
	return out
}
