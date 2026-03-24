// Copyright 2019 wongoo. All rights reserved.

package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/vogo/namer/internal/scoring"
	"github.com/vogo/namer/internal/web"
)

const version = "1.0.0"

func defaultConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".namer.conf"
	}
	return filepath.Join(home, ".namer.conf")
}

func printUsage() {
	fmt.Printf(`namer v%s - 中国姓名阴阳五行评分工具

用法:
  namer                       交互式批量评分（自动引导配置）
  namer -c <配置文件>          使用指定配置文件批量评分
  namer <姓> <名>             评估单个名字（使用已有配置中的生辰信息）
  namer <姓> <名> <生日>       评估单个名字（指定生日，格式: 2024-03-15）
  namer -web                  启动 Web 服务（随机端口，自动打开浏览器）
  namer -web -port 8080       启动 Web 服务（指定端口）
  namer -h                    显示帮助

示例:
  namer                       首次使用，交互式输入配置后批量评分
  namer 王 明轩               评估"王明轩"这个名字
  namer 王 明轩 2024-03-15    评估"王明轩"，指定出生日期
  namer -c my.conf            使用 my.conf 配置文件批量评分
  namer -web                  启动 Web 界面
  namer -web -port 3000       在 3000 端口启动 Web 界面

配置文件:
  默认路径: ~/.namer.conf
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

	// namer <姓> <名> [生日]
	if len(args) >= 2 && !strings.HasPrefix(args[0], "-") {
		runSingleScore(args)
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

	// 默认批量模式
	runBatchMode(configPath)
}

// runSingleScore 单个名字评分
func runSingleScore(args []string) {
	lastName := args[0]
	firstName := args[1]

	// 加载配置获取生辰信息
	cfg := &scoring.Config{}
	cfgPath := defaultConfigPath()
	if _, err := os.Stat(cfgPath); err == nil {
		_ = scoring.ReadConfigFile(cfgPath, cfg)
	}

	// 如果命令行提供了生日
	if len(args) >= 3 {
		parseBirthday(args[2], cfg)
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

	r := scoring.CalcScore(lastName, firstName, cfg.Year, cfg.Month, cfg.Day, cfg.Hour, cfg.Minute)
	scoring.PrintResult(lastName, firstName, r)
}

// parseBirthday 解析生日字符串 "2024-03-15" 或 "2024-03-15-10" 或 "2024-03-15-10-30"
func parseBirthday(s string, cfg *scoring.Config) {
	parts := strings.Split(s, "-")
	if len(parts) >= 3 {
		fmt.Sscanf(parts[0], "%d", &cfg.Year)
		fmt.Sscanf(parts[1], "%d", &cfg.Month)
		fmt.Sscanf(parts[2], "%d", &cfg.Day)
	}
	if len(parts) >= 4 {
		fmt.Sscanf(parts[3], "%d", &cfg.Hour)
	}
	if len(parts) >= 5 {
		fmt.Sscanf(parts[4], "%d", &cfg.Minute)
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
		if cfg.LastName == "" {
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
		cfg.LastName, cfg.Year, cfg.Month, cfg.Day, cfg.Hour, cfg.Minute,
		cfg.FirstNameKeyWords)
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
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

func flagValue(args []string, flag string) string {
	for i, a := range args {
		if a == flag && i+1 < len(args) {
			return args[i+1]
		}
	}
	return ""
}
