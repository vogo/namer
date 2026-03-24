package scoring

import (
	"os"
	"strings"
	"testing"
)

func TestReadConfigFile(t *testing.T) {
	content := `{
		"xing": "王",
		"year": 2024,
		"month": 3,
		"day": 15,
		"hour": 10,
		"minute": 30,
		"gender": 0,
		"min_candidate_score": 80,
		"ming_keywords": "明,轩"
	}`
	f, err := os.CreateTemp("", "namer_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(f.Name()) }()
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	cfg := &Config{}
	err = ReadConfigFile(f.Name(), cfg)
	if err != nil {
		t.Fatalf("ReadConfigFile error: %v", err)
	}
	if cfg.Xing != "王" {
		t.Errorf("Xing = %q, want 王", cfg.Xing)
	}
	if cfg.Year != 2024 {
		t.Errorf("Year = %d, want 2024", cfg.Year)
	}
	if cfg.Month != 3 {
		t.Errorf("Month = %d, want 3", cfg.Month)
	}
	if cfg.MingKeywords != "明,轩" {
		t.Errorf("MingKeywords = %q, want '明,轩'", cfg.MingKeywords)
	}
}

func TestReadConfigFileNotExist(t *testing.T) {
	cfg := &Config{}
	err := ReadConfigFile("/tmp/nonexistent_namer_test.json", cfg)
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestReadConfigFileInvalidJSON(t *testing.T) {
	f, err := os.CreateTemp("", "namer_test_bad_*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Remove(f.Name()) }()
	if _, err := f.WriteString("not json"); err != nil {
		t.Fatal(err)
	}
	_ = f.Close()

	cfg := &Config{}
	err = ReadConfigFile(f.Name(), cfg)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestNameScoring(t *testing.T) {
	cfg := &Config{
		Xing:              "王",
		Year:              2024,
		Month:             3,
		Day:               15,
		Hour:              10,
		Minute:            30,
		MinCandidateScore: 80,
		MingKeywords:      "明,轩",
	}
	// NameScoring 不应 panic
	NameScoring(cfg)
}

func TestNameScoringEmpty(t *testing.T) {
	cfg := &Config{
		Xing:         "王",
		Year:         2024,
		Month:        3,
		Day:          15,
		MingKeywords: "",
	}
	// 空备选字不应 panic
	NameScoring(cfg)
}

func TestNameScoringWithSpaces(t *testing.T) {
	cfg := &Config{
		Xing:         "王",
		Year:         2024,
		Month:        3,
		Day:          15,
		Hour:         10,
		Minute:       30,
		MingKeywords: " 明 , 轩 ",
	}
	NameScoring(cfg)
}

func TestWriteConfigFile(t *testing.T) {
	cfg := &Config{
		Xing:              "李",
		Year:              2000,
		Month:             6,
		Day:               15,
		Hour:              8,
		Minute:            0,
		Gender:            1,
		MinCandidateScore: 70,
		MingKeywords:      "浩,然",
	}

	f, err := os.CreateTemp("", "namer_write_test_*.json")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	defer func() { _ = os.Remove(f.Name()) }()

	err = WriteConfigFile(f.Name(), cfg)
	if err != nil {
		t.Fatalf("WriteConfigFile error: %v", err)
	}

	// 读回验证
	cfg2 := &Config{}
	err = ReadConfigFile(f.Name(), cfg2)
	if err != nil {
		t.Fatalf("ReadConfigFile error: %v", err)
	}
	if cfg2.Xing != "李" {
		t.Errorf("Xing = %q, want 李", cfg2.Xing)
	}
	if cfg2.Year != 2000 {
		t.Errorf("Year = %d, want 2000", cfg2.Year)
	}
	if cfg2.MingKeywords != "浩,然" {
		t.Errorf("MingKeywords = %q, want '浩,然'", cfg2.MingKeywords)
	}
}

func TestPromptConfigFrom(t *testing.T) {
	// Hour=0 和 Minute=0 是有效值(午夜零分)，不会被提示
	// Gender 需要 1 或 2，初始 0 会触发提示
	input := "王\n2024\n3\n15\n10\n30\n1\n明,轩,浩\n"
	cfg := &Config{Hour: -1, Minute: -1} // 标记为未设置
	PromptConfigFrom(cfg, strings.NewReader(input))

	if cfg.Xing != "王" {
		t.Errorf("Xing = %q, want 王", cfg.Xing)
	}
	if cfg.Year != 2024 {
		t.Errorf("Year = %d, want 2024", cfg.Year)
	}
	if cfg.Month != 3 {
		t.Errorf("Month = %d, want 3", cfg.Month)
	}
	if cfg.Day != 15 {
		t.Errorf("Day = %d, want 15", cfg.Day)
	}
	if cfg.Hour != 10 {
		t.Errorf("Hour = %d, want 10", cfg.Hour)
	}
	if cfg.Minute != 30 {
		t.Errorf("Minute = %d, want 30", cfg.Minute)
	}
	if cfg.Gender != 1 {
		t.Errorf("Gender = %d, want 1", cfg.Gender)
	}
	if cfg.MingKeywords != "明,轩,浩" {
		t.Errorf("MingKeywords = %q, want '明,轩,浩'", cfg.MingKeywords)
	}
}

func TestPromptConfigFromPartial(t *testing.T) {
	// 已有大部分配置，只补缺性别和备选字
	input := "1\n明,轩\n"
	cfg := &Config{
		Xing:   "李",
		Year:   2000,
		Month:  6,
		Day:    15,
		Hour:   10,
		Minute: 30,
	}
	PromptConfigFrom(cfg, strings.NewReader(input))

	if cfg.Xing != "李" {
		t.Errorf("Xing should remain 李, got %q", cfg.Xing)
	}
	if cfg.Year != 2000 {
		t.Errorf("Year should remain 2000, got %d", cfg.Year)
	}
	if cfg.Gender != 1 {
		t.Errorf("Gender = %d, want 1", cfg.Gender)
	}
	if cfg.MingKeywords != "明,轩" {
		t.Errorf("MingKeywords = %q, want '明,轩'", cfg.MingKeywords)
	}
}

func TestPromptConfigFromRetryInvalid(t *testing.T) {
	// 月份先输入无效值 13 再输入有效值 3
	input := "王\n2024\n13\n3\n15\n1\n明\n"
	cfg := &Config{}
	PromptConfigFrom(cfg, strings.NewReader(input))

	if cfg.Month != 3 {
		t.Errorf("Month = %d, want 3 (after retry)", cfg.Month)
	}
}

func TestPromptConfigFromEOF(t *testing.T) {
	// 输入不足时不应 panic
	input := "王\n"
	cfg := &Config{}
	PromptConfigFrom(cfg, strings.NewReader(input))
	if cfg.Xing != "王" {
		t.Errorf("Xing = %q, want 王", cfg.Xing)
	}
}

func TestConfigIsComplete(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		want bool
	}{
		{"完整", Config{Xing: "王", Year: 2024, Month: 3, Day: 15, MingKeywords: "明"}, true},
		{"缺姓", Config{Year: 2024, Month: 3, Day: 15, MingKeywords: "明"}, false},
		{"缺年", Config{Xing: "王", Month: 3, Day: 15, MingKeywords: "明"}, false},
		{"缺月", Config{Xing: "王", Year: 2024, Day: 15, MingKeywords: "明"}, false},
		{"缺日", Config{Xing: "王", Year: 2024, Month: 3, MingKeywords: "明"}, false},
		{"缺备选字", Config{Xing: "王", Year: 2024, Month: 3, Day: 15}, false},
		{"月份越界", Config{Xing: "王", Year: 2024, Month: 13, Day: 15, MingKeywords: "明"}, false},
		{"日期越界", Config{Xing: "王", Year: 2024, Month: 3, Day: 32, MingKeywords: "明"}, false},
	}
	for _, tt := range tests {
		if got := tt.cfg.IsComplete(); got != tt.want {
			t.Errorf("%s: IsComplete() = %v, want %v", tt.name, got, tt.want)
		}
	}
}
