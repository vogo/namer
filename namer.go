// Copyright 2019 wongoo. All rights reserved.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"slices"
	"strconv"
	"syscall"

	"github.com/vogo/namer/internal/scoring"
	"github.com/vogo/namer/internal/web"
)

const version = "1.0.0"

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".namer.json"
	}
	return filepath.Join(home, ".namer.json")
}

func printUsage() {
	fmt.Printf(`namer v%s - 中国姓名阴阳五行评分工具

用法:
  namer                       交互式批量评分（自动引导配置）
  namer -c <配置文件>          使用指定配置文件批量评分
  namer -xing 王 -keywords 明,轩,浩,然 -year 2024 -month 3 -day 15 -hour 10 -minute 30
                              通过命令行参数批量评分（无需配置文件）
  namer -xing 王 -ming 明轩 [-year 2024 -month 3 -day 15 -hour 10 -minute 30 -gender 1]
                              评估单个名字（通过参数指定所有信息）
  namer -web                  启动 Web 服务（随机端口，自动打开浏览器）
  namer -web -port 8080       启动 Web 服务（指定端口）
  namer -h                    显示帮助

示例:
  namer                       首次使用，交互式输入配置后批量评分
  namer -xing 王 -keywords 明,轩,浩,然 -year 2024 -month 3 -day 15 -hour 10 -minute 30
                              通过命令行参数批量评分
  namer -xing 王 -ming 明轩   评估"王明轩"（使用已有配置中的生辰信息）
  namer -xing 王 -ming 明轩 -year 2024 -month 3 -day 15 -hour 10 -minute 30
                              评估"王明轩"，指定完整出生信息
  namer -c my.conf            使用 my.conf 配置文件批量评分
  namer -web                  启动 Web 界面
  namer -web -port 3000       在 3000 端口启动 Web 界面

参数:
  -xing <姓>                  姓氏
  -ming <名>                  名字
  -keywords <备选字>          名字备选字（逗号分隔，如：明,轩,浩,然）
  -year <年>                  出生年份
  -month <月>                 出生月份（1-12）
  -day <日>                   出生日期（1-31）
  -hour <时>                  出生时辰（0-23）
  -minute <分>                出生分钟（0-59）
  -gender <性别>              性别（1=男, 2=女）
  -score <最低分>             最小候选分数（默认60）

配置文件:
  默认路径: ~/.namer.json
  首次运行时会交互式引导创建配置文件，配置项包括：
    - 姓氏、出生年月日时分、性别
    - 名字备选字（逗号分隔）

评分维度（总分100）:
  五格数理  30分    天格、人格、地格、总格、外格的数理吉凶
  三才配置  25分    天人地三才的五行生克关系
  喜用神    20分    名字五行是否补益八字喜用神
  内部五行  15分    姓名各字之间的五行生克
  阴阳平衡  10分    姓名各字笔画的阴阳搭配
`, version)
}

func main() {
	args := os.Args[1:]

	// 帮助
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help" || args[0] == "help") {
		printUsage()
		return
	}

	// 版本
	if len(args) == 1 && (args[0] == "-v" || args[0] == "--version") {
		fmt.Printf("namer v%s\n", version)
		return
	}

	// namer -web [-port 8080]
	if hasFlag(args, "-web") {
		port := 0
		if p := flagValue(args, "-port"); p != "" {
			port, _ = strconv.Atoi(p)
		}
		if err := web.Start(port); err != nil {
			fmt.Printf("web 服务启动失败: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// namer -xing 王 -ming 明轩 [其他参数] — 单个名字评分
	if flagValue(args, "-xing") != "" && flagValue(args, "-ming") != "" {
		runSingleScore(args)
		return
	}

	// namer -xing 王 -keywords 明,轩 [...] — 命令行参数批量评分
	if flagValue(args, "-xing") != "" && flagValue(args, "-keywords") != "" {
		runCLIBatchMode(args)
		return
	}

	// namer -c <file>
	configPath := ""
	for i, a := range args {
		if a == "-c" && i+1 < len(args) {
			configPath = args[i+1]
			break
		}
	}

	// 默认批量模式（配置文件 + 交互）
	runBatchMode(configPath)
}

// runSingleScore 单个名字评分
func runSingleScore(args []string) {
	xing := flagValue(args, "-xing")
	ming := flagValue(args, "-ming")

	// 加载配置获取生辰信息
	cfg := &scoring.Config{}
	cfgPath := defaultConfigPath()
	if _, err := os.Stat(cfgPath); err == nil {
		_ = scoring.ReadConfigFile(cfgPath, cfg)
	}

	// 命令行参数覆盖配置文件
	if v := flagValue(args, "-year"); v != "" {
		cfg.Year, _ = strconv.Atoi(v)
	}
	if v := flagValue(args, "-month"); v != "" {
		cfg.Month, _ = strconv.Atoi(v)
	}
	if v := flagValue(args, "-day"); v != "" {
		cfg.Day, _ = strconv.Atoi(v)
	}
	if v := flagValue(args, "-hour"); v != "" {
		cfg.Hour, _ = strconv.Atoi(v)
	}
	if v := flagValue(args, "-minute"); v != "" {
		cfg.Minute, _ = strconv.Atoi(v)
	}
	if v := flagValue(args, "-gender"); v != "" {
		cfg.Gender, _ = strconv.Atoi(v)
	}

	// 如果仍缺少生辰信息，使用默认值
	if cfg.Year <= 0 {
		fmt.Println("提示: 未找到生辰配置，使用默认生辰(2024-01-01 12:00)计算喜用神")
		fmt.Printf("      运行 namer 进入交互模式可配置生辰信息\n\n")
		cfg.Year = 2024
		cfg.Month = 1
		cfg.Day = 1
		cfg.Hour = 12
	}

	r := scoring.CalcScore(xing, ming, cfg.Year, cfg.Month, cfg.Day, cfg.Hour, cfg.Minute)
	scoring.PrintResult(xing, ming, r)
}

// runCLIBatchMode 通过命令行参数批量评分（无需配置文件）
func runCLIBatchMode(args []string) {
	cfg := &scoring.Config{}
	cfg.Xing = flagValue(args, "-xing")
	cfg.MingKeywords = flagValue(args, "-keywords")

	if v := flagValue(args, "-year"); v != "" {
		cfg.Year, _ = strconv.Atoi(v)
	}
	if v := flagValue(args, "-month"); v != "" {
		cfg.Month, _ = strconv.Atoi(v)
	}
	if v := flagValue(args, "-day"); v != "" {
		cfg.Day, _ = strconv.Atoi(v)
	}
	if v := flagValue(args, "-hour"); v != "" {
		cfg.Hour, _ = strconv.Atoi(v)
	}
	if v := flagValue(args, "-minute"); v != "" {
		cfg.Minute, _ = strconv.Atoi(v)
	}
	if v := flagValue(args, "-gender"); v != "" {
		cfg.Gender, _ = strconv.Atoi(v)
	}
	if v := flagValue(args, "-score"); v != "" {
		cfg.MinCandidateScore, _ = strconv.Atoi(v)
	}
	if cfg.MinCandidateScore <= 0 {
		cfg.MinCandidateScore = 60
	}

	if cfg.Year <= 0 || cfg.Month <= 0 || cfg.Day <= 0 {
		fmt.Println("错误: 批量评分需要指定出生信息 (-year -month -day)")
		fmt.Println("示例: namer -xing 王 -keywords 明,轩 -year 2024 -month 3 -day 15 -hour 10 -minute 30")
		os.Exit(1)
	}

	fmt.Printf("\n姓氏: %s | 生辰: %d-%02d-%02d %02d:%02d | 备选字: %s\n",
		cfg.Xing, cfg.Year, cfg.Month, cfg.Day, cfg.Hour, cfg.Minute,
		cfg.MingKeywords)
	fmt.Println("开始批量评分...")

	finishChan := make(chan int, 1)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("\n评分异常: %v\n", err)
			}
			finishChan <- 1
		}()
		scoring.NameScoring(cfg)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		fmt.Printf("\n收到中断信号: %s\n", sig)
	case <-finishChan:
		fmt.Println("\n评分完成!")
	}
}

// runBatchMode 批量评分模式
func runBatchMode(configPath string) {
	if configPath == "" {
		configPath = defaultConfigPath()
	}

	cfg := &scoring.Config{}

	// 尝试读取已有配置
	if _, err := os.Stat(configPath); err == nil {
		if err := scoring.ReadConfigFile(configPath, cfg); err != nil {
			fmt.Printf("配置文件读取失败: %v\n", err)
			fmt.Println("将重新引导配置...")
			cfg = &scoring.Config{}
		}
	}

	// 检查配置是否完整，不完整则交互补全
	if !cfg.IsComplete() {
		if cfg.Xing == "" {
			fmt.Println("=== namer 姓名评分工具 ===")
			fmt.Println("首次使用，请输入以下配置信息：")
		} else {
			fmt.Println("配置信息不完整，请补充以下内容：")
		}
		fmt.Println()
		scoring.PromptConfig(cfg)

		// 保存配置
		if err := scoring.WriteConfigFile(configPath, cfg); err != nil {
			fmt.Printf("配置保存失败: %v\n", err)
		} else {
			fmt.Printf("\n配置已保存到: %s\n", configPath)
		}
	}

	fmt.Printf("\n姓氏: %s | 生辰: %d-%02d-%02d %02d:%02d | 备选字: %s\n",
		cfg.Xing, cfg.Year, cfg.Month, cfg.Day, cfg.Hour, cfg.Minute,
		cfg.MingKeywords)
	fmt.Println("开始批量评分...")

	finishChan := make(chan int, 1)

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("\n评分异常: %v\n", err)
			}
			finishChan <- 1
		}()
		scoring.NameScoring(cfg)
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-signalChan:
		fmt.Printf("\n收到中断信号: %s\n", sig)
	case <-finishChan:
		fmt.Println("\n评分完成!")
	}
}

func hasFlag(args []string, flag string) bool {
	return slices.Contains(args, flag)
}

func flagValue(args []string, flag string) string {
	for i, a := range args {
		if a == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}
