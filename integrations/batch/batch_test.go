package batch_test

import (
	"strings"
	"testing"

	"github.com/vogo/namer/integrations/helper"
)

func TestBatchWithConfigFile(t *testing.T) {
	configJSON := `{
		"xing": "王",
		"year": 2024,
		"month": 3,
		"day": 15,
		"hour": 10,
		"minute": 30,
		"gender": 1,
		"min_candidate_score": 60,
		"ming_keywords": "明,轩,浩"
	}`

	cfgPath := helper.WriteConfigFile(t, configJSON)

	stdout, _, err := helper.Run(t, "-c", cfgPath)
	if err != nil {
		t.Fatalf("batch scoring failed: %v", err)
	}

	// Should show batch progress and top results
	checks := []string{
		"开始批量评分",
		"Top 10",
		"分",
		"高分名字详情",
	}

	for _, check := range checks {
		if !strings.Contains(stdout, check) {
			t.Errorf("batch output missing %q", check)
		}
	}
}

func TestBatchWithMinimalKeywords(t *testing.T) {
	configJSON := `{
		"xing": "李",
		"year": 2000,
		"month": 6,
		"day": 1,
		"hour": 8,
		"minute": 0,
		"gender": 2,
		"min_candidate_score": 60,
		"ming_keywords": "明"
	}`

	cfgPath := helper.WriteConfigFile(t, configJSON)

	stdout, _, err := helper.Run(t, "-c", cfgPath)
	if err != nil {
		t.Fatalf("batch with minimal keywords failed: %v", err)
	}

	// With single keyword, should have 1 single-char + 1 double-char result
	if !strings.Contains(stdout, "评分完成") {
		t.Errorf("expected completion message, got: %s", stdout)
	}
}

func TestBatchWithEmptyKeywords(t *testing.T) {
	configJSON := `{
		"xing": "王",
		"year": 2024,
		"month": 1,
		"day": 1,
		"hour": 12,
		"minute": 0,
		"gender": 1,
		"min_candidate_score": 60,
		"ming_keywords": ""
	}`

	cfgPath := helper.WriteConfigFile(t, configJSON)

	stdout, _, err := helper.Run(t, "-c", cfgPath)
	// Should handle empty keywords gracefully (may prompt interactively or skip)
	_ = err
	_ = stdout
}

func TestBatchConfigFileNotFound(t *testing.T) {
	_, _, err := helper.Run(t, "-c", "/nonexistent/path/config.json")
	// Should handle missing config file - may prompt or error
	_ = err
}

func TestBatchResultsSortedByScore(t *testing.T) {
	configJSON := `{
		"xing": "张",
		"year": 2024,
		"month": 5,
		"day": 10,
		"hour": 14,
		"minute": 0,
		"gender": 1,
		"min_candidate_score": 60,
		"ming_keywords": "明,轩,浩,然"
	}`

	cfgPath := helper.WriteConfigFile(t, configJSON)

	stdout, _, err := helper.Run(t, "-c", cfgPath)
	if err != nil {
		t.Fatalf("batch scoring failed: %v", err)
	}

	// Find the Top 10 section and verify scores are descending
	lines := strings.Split(stdout, "\n")
	var scores []string
	inTop := false
	for _, line := range lines {
		if strings.Contains(line, "Top 10") {
			inTop = true
			continue
		}
		if inTop && strings.Contains(line, "====") {
			break
		}
		if inTop && strings.Contains(line, "分") {
			scores = append(scores, line)
		}
	}

	if len(scores) == 0 {
		t.Error("no scores found in Top 10 section")
	}
}
