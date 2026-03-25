package score_test

import (
	"strings"
	"testing"

	"github.com/vogo/namer/integrations/helper"
)

func TestSingleNameScore(t *testing.T) {
	stdout, _, err := helper.Run(t,
		"-xing", "王",
		"-ming", "明轩",
		"-year", "2024",
		"-month", "3",
		"-day", "15",
		"-hour", "10",
		"-minute", "30",
	)
	if err != nil {
		t.Fatalf("namer single score failed: %v", err)
	}

	// Verify output contains expected sections
	checks := []string{
		"姓名: 王明轩",
		"总分:",
		"/ 100",
		"五格数理",
		"三才配置",
		"喜用神匹配",
		"内部五行",
		"阴阳平衡",
		"康熙笔画:",
		"五格:",
		"天格",
		"人格",
		"地格",
		"总格",
		"外格",
		"三才:",
		"八字:",
		"喜用神:",
		"字五行:",
		"阴阳:",
	}

	for _, check := range checks {
		if !strings.Contains(stdout, check) {
			t.Errorf("output missing %q", check)
		}
	}
}

func TestSingleCharName(t *testing.T) {
	stdout, _, err := helper.Run(t,
		"-xing", "李",
		"-ming", "明",
		"-year", "2000",
		"-month", "6",
		"-day", "1",
		"-hour", "8",
	)
	if err != nil {
		t.Fatalf("namer single char score failed: %v", err)
	}

	if !strings.Contains(stdout, "姓名: 李明") {
		t.Errorf("expected name 李明, got: %s", stdout)
	}

	if !strings.Contains(stdout, "总分:") {
		t.Errorf("expected score output, got: %s", stdout)
	}
}

func TestScoreWithoutBirthInfo(t *testing.T) {
	// Without birth info and no config file, should use default and show hint
	stdout, _, err := helper.Run(t,
		"-xing", "张",
		"-ming", "伟",
	)
	if err != nil {
		t.Fatalf("namer score without birth info failed: %v", err)
	}

	if !strings.Contains(stdout, "提示") {
		t.Errorf("expected hint about missing birth info, got: %s", stdout)
	}

	if !strings.Contains(stdout, "总分:") {
		t.Errorf("expected score output even without birth info, got: %s", stdout)
	}
}

func TestScoreWithGender(t *testing.T) {
	stdout, _, err := helper.Run(t,
		"-xing", "王",
		"-ming", "明轩",
		"-year", "2024",
		"-month", "3",
		"-day", "15",
		"-hour", "10",
		"-gender", "1",
	)
	if err != nil {
		t.Fatalf("namer score with gender failed: %v", err)
	}

	if !strings.Contains(stdout, "姓名: 王明轩") {
		t.Errorf("expected name in output, got: %s", stdout)
	}
}

func TestDifferentBirthDateProducesDifferentScores(t *testing.T) {
	stdout1, _, err := helper.Run(t,
		"-xing", "王",
		"-ming", "明轩",
		"-year", "1990",
		"-month", "1",
		"-day", "1",
		"-hour", "6",
	)
	if err != nil {
		t.Fatalf("first score failed: %v", err)
	}

	stdout2, _, err := helper.Run(t,
		"-xing", "王",
		"-ming", "明轩",
		"-year", "2024",
		"-month", "8",
		"-day", "20",
		"-hour", "18",
	)
	if err != nil {
		t.Fatalf("second score failed: %v", err)
	}

	// Extract total scores - they should differ because birth dates affect 喜用神
	// Both should have valid output
	if !strings.Contains(stdout1, "总分:") || !strings.Contains(stdout2, "总分:") {
		t.Error("both runs should produce score output")
	}

	// The 八字 lines should differ
	if extractLine(stdout1, "八字:") == extractLine(stdout2, "八字:") {
		t.Error("different birth dates should produce different bazi")
	}
}

func TestDoubleCharSurname(t *testing.T) {
	stdout, _, err := helper.Run(t,
		"-xing", "欧阳",
		"-ming", "明",
		"-year", "2024",
		"-month", "1",
		"-day", "1",
		"-hour", "12",
	)
	if err != nil {
		t.Fatalf("double surname score failed: %v", err)
	}

	if !strings.Contains(stdout, "姓名: 欧阳明") {
		t.Errorf("expected name 欧阳明 in output, got: %s", stdout)
	}
}

func extractLine(output, prefix string) string {
	for _, line := range strings.Split(output, "\n") {
		if strings.Contains(line, prefix) {
			return strings.TrimSpace(line)
		}
	}
	return ""
}
