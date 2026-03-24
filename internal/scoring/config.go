package scoring

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
)

// Config 评分配置
type Config struct {
	LastName          string `json:"last_name"`
	Year              int    `json:"year"`
	Month             int    `json:"month"`
	Day               int    `json:"day"`
	Hour              int    `json:"hour"`
	Minute            int    `json:"minute"`
	Gender            int    `json:"gender"`
	MinCandidateScore int    `json:"min_candidate_score"`
	FirstNameKeyWords string `json:"first_name_key_words"`
}

// IsComplete 检查配置是否完整
func (c *Config) IsComplete() bool {
	return c.LastName != "" &&
		c.Year > 0 &&
		c.Month > 0 && c.Month <= 12 &&
		c.Day > 0 && c.Day <= 31 &&
		c.FirstNameKeyWords != ""
}

// ReadConfigFile 读取配置文件
func ReadConfigFile(file string, cfg *Config) error {
	bt, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(bt, cfg)
}

// WriteConfigFile 写入配置文件
func WriteConfigFile(file string, cfg *Config) error {
	bt, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, bt, 0644)
}

// PromptConfig 交互式逐项提示用户输入缺失的配置
func PromptConfig(cfg *Config) {
	PromptConfigFrom(cfg, os.Stdin)
}

// PromptConfigFrom 从指定 reader 读取用户输入补全配置
func PromptConfigFrom(cfg *Config, r io.Reader) {
	scanner := bufio.NewScanner(r)
	if cfg.LastName == "" {
		cfg.LastName = promptString(scanner, "请输入姓氏（如：王）")
	}
	if cfg.Year <= 0 {
		cfg.Year = promptInt(scanner, "请输入出生年份（如：2024）")
	}
	if cfg.Month <= 0 || cfg.Month > 12 {
		cfg.Month = promptIntRange(scanner, "请输入出生月份（1-12）", 1, 12)
	}
	if cfg.Day <= 0 || cfg.Day > 31 {
		cfg.Day = promptIntRange(scanner, "请输入出生日期（1-31）", 1, 31)
	}
	if cfg.Hour < 0 || cfg.Hour > 23 {
		cfg.Hour = promptIntRange(scanner, "请输入出生时辰（0-23，如：10 表示上午10点）", 0, 23)
	}
	if cfg.Minute < 0 || cfg.Minute > 59 {
		cfg.Minute = promptIntRange(scanner, "请输入出生分钟（0-59）", 0, 59)
	}
	if cfg.Gender != 1 && cfg.Gender != 2 {
		cfg.Gender = promptIntRange(scanner, "请输入性别（1=男, 2=女）", 1, 2)
	}
	if cfg.MinCandidateScore <= 0 {
		cfg.MinCandidateScore = 60
	}
	if cfg.FirstNameKeyWords == "" {
		cfg.FirstNameKeyWords = promptString(scanner, "请输入名字备选字（逗号分隔，如：明,轩,浩,然）")
	}
}

func promptString(scanner *bufio.Scanner, prompt string) string {
	for {
		fmt.Printf("%s: ", prompt)
		if !scanner.Scan() {
			return ""
		}
		s := strings.TrimSpace(scanner.Text())
		if s != "" {
			return s
		}
		fmt.Println("  输入不能为空，请重新输入")
	}
}

func promptInt(scanner *bufio.Scanner, prompt string) int {
	for {
		fmt.Printf("%s: ", prompt)
		if !scanner.Scan() {
			return 0
		}
		var n int
		if _, err := fmt.Sscanf(scanner.Text(), "%d", &n); err == nil && n > 0 {
			return n
		}
		fmt.Println("  请输入有效的数字")
	}
}

func promptIntRange(scanner *bufio.Scanner, prompt string, min, max int) int {
	for {
		fmt.Printf("%s: ", prompt)
		if !scanner.Scan() {
			return min
		}
		var n int
		if _, err := fmt.Sscanf(scanner.Text(), "%d", &n); err == nil && n >= min && n <= max {
			return n
		}
		fmt.Printf("  请输入 %d 到 %d 之间的数字\n", min, max)
	}
}

// NameScoring 批量评分
func NameScoring(cfg *Config) {
	words := strings.Split(cfg.FirstNameKeyWords, ",")
	keyWords := make([]rune, 0, len(words))
	for _, w := range words {
		w = strings.TrimSpace(w)
		if w != "" {
			keyWords = append(keyWords, []rune(w)[0])
		}
	}

	if len(keyWords) == 0 {
		fmt.Println("没有备选字，请检查配置")
		return
	}

	type nameResult struct {
		name  string
		score float64
	}

	var results []nameResult

	total := len(keyWords) + len(keyWords)*len(keyWords)
	done := 0

	// 单字名
	for _, c := range keyWords {
		firstName := string(c)
		r := CalcScore(cfg.LastName, firstName, cfg.Year, cfg.Month, cfg.Day, cfg.Hour, cfg.Minute)
		fullName := cfg.LastName + firstName
		results = append(results, nameResult{name: fullName, score: r.Total})
		done++
		fmt.Printf("\r评分进度: %d/%d", done, total)
	}

	// 双字名
	for _, c1 := range keyWords {
		for _, c2 := range keyWords {
			firstName := string([]rune{c1, c2})
			r := CalcScore(cfg.LastName, firstName, cfg.Year, cfg.Month, cfg.Day, cfg.Hour, cfg.Minute)
			fullName := cfg.LastName + firstName
			results = append(results, nameResult{name: fullName, score: r.Total})
			done++
			fmt.Printf("\r评分进度: %d/%d", done, total)
		}
	}
	fmt.Println()

	// 按分数排序
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	// 打印 Top 10
	fmt.Println("\n========== Top 10 ==========")
	count := min(len(results), 10)
	for i := 0; i < count; i++ {
		r := results[i]
		fmt.Printf("%2d. %s  %.1f 分\n", i+1, r.name, r.score)
	}

	// 打印高分名字详情
	fmt.Println("\n========== 高分名字详情 ==========")
	lastRunes := []rune(cfg.LastName)
	for i := 0; i < count; i++ {
		nr := results[i]
		firstName := string([]rune(nr.name)[len(lastRunes):])
		r := CalcScore(cfg.LastName, firstName, cfg.Year, cfg.Month, cfg.Day, cfg.Hour, cfg.Minute)
		PrintResult(cfg.LastName, firstName, r)
	}
}
