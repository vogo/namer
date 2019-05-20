// Copyright 2019 wongoo. All rights reserved.

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	dataFileSuffix          = ".data"
	dataCandidateFileSuffix = ".candidate.data"
)

var (
	baiduQimingTestUrlFormat = "https://sp0.baidu.com/5LMDcjW6BwF3otqbppnN2DJv/qiming.pae.baidu.com/data/namedetail?" +
		"year=%d&month=%02d&day=%02d&hour=%02d&min=%02d&timeType=0&gender=0" +
		"&flag=1&cb=jsonp1"
	baiduScorePreviousString = "\"key\":\"score\",\"value\":\""
	baiduTestUrlPrefix       string
	keyWords                 []rune

	config         = &Config{}
	scoreDB        = newScore(0)
	scoreStat      = make(map[int][]string, 100)
	candidates     = make(map[string]int, 100)
	lastNameLength = 0
)

type NameScore struct {
	Score int                 `json:"s"`
	More  map[rune]*NameScore `json:"m,omitempty"`
	Read  int                 `json:"r"` // 0: unread, 1: read
}

func newScore(score int) *NameScore {
	return &NameScore{Score: score, More: make(map[rune]*NameScore, 0)}
}

type Config struct {
	LastName          string `json:"last_name"`
	Year              int    `json:"year"`
	Month             int    `json:"month"`
	Day               int    `json:"day"`
	Hour              int    `json:"hour"`
	Minute            int    `json:"minute"`
	Gender            int    `json:"gender"`
	FirstNameKeyWords string `json:"first_name_key_words"`
	MinCandidateScore int    `json:"min_candidate_score"`
}

func buildNameTestUrlPrefix() {
	baiduTestUrlPrefix = fmt.Sprintf(baiduQimingTestUrlFormat, config.Year, config.Month, config.Day, config.Hour, config.Minute)
}

func buildNameTestUrl(firstName string) string {
	return fmt.Sprintf(baiduTestUrlPrefix+"&fName=%s&lName=%s&_=%d", url.QueryEscape(config.LastName), url.QueryEscape(firstName), int64(time.Now().UnixNano())/int64(time.Millisecond))
}

func getNameScore(firstName string) (int, error) {
	addr := buildNameTestUrl(firstName)
	bt, err := request(addr)
	if err != nil {
		return 0, err
	}
	content := string(bt)
	index := strings.Index(content, baiduScorePreviousString)
	if index <= 0 {
		return 0, errors.New("cant find score from response: " + content)
	}

	content = content[index+len(baiduScorePreviousString):]
	index = strings.Index(content, "\"")
	if index <= 0 {
		return 0, errors.New("cant find score from response: " + content)
	}

	scoreStr := content[:index]
	score, err := strconv.Atoi(scoreStr)
	log.Printf("name score: %s = %d", firstName, score)
	return score, err

}

func request(addr string) ([]byte, error) {
	time.Sleep(time.Millisecond * 300)
	resp, err := http.Get(addr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func readConfigFile(file string) error {
	bt, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bt, config)
	if err != nil {
		return err
	}

	lastNameLength = len(config.LastName)
	return nil
}

func readScoreData(file string) error {
	file += dataFileSuffix
	bt, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	return json.Unmarshal(bt, scoreDB)
}

func readCandidateData(file string) error {
	file += dataCandidateFileSuffix
	bt, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	return json.Unmarshal(bt, &candidates)
}

func statCalculate(score *NameScore, name string) {
	scoreStat[score.Score] = append(scoreStat[score.Score], name)
	for key, value := range score.More {
		statCalculate(value, name+string(key))
	}
}

func writeScoreData() {
	bt, err := json.Marshal(scoreDB)
	if err != nil {
		log.Fatalf("marshal err: %v", err)
	}
	err = ioutil.WriteFile(*configFile+dataFileSuffix, bt, 0666)
	if err != nil {
		log.Fatalf("write score file error: %v", err)
	}
}

func printCandidates() {
	fmt.Println("----------候选名单-------------")
	for key, value := range candidates {
		fmt.Println(key, ":", value)
	}
}
func writeCandidateData() {
	bt, err := json.Marshal(candidates)
	if err != nil {
		log.Fatalf("marshal err: %v", err)
	}
	err = ioutil.WriteFile(*configFile+dataCandidateFileSuffix, bt, 0666)
	if err != nil {
		log.Fatalf("write candidate file error: %v", err)
	}
}

func nameScoring() error {
	err := oneFirstNameLoopScoring()
	if err != nil {
		return err
	}
	err = twoFirstNameLoopScoring()
	if err != nil {
		return err
	}
	return nil
}

func oneFirstNameLoopScoring() error {
	if scoreDB == nil {
		scoreDB = newScore(0)
	}

	for i := 0; i < len(keyWords); i++ {
		err := oneFirstNameScoring(scoreDB, keyWords[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func oneFirstNameScoring(parent *NameScore, firstName rune) error {
	if _, ok := scoreDB.More[firstName]; ok {
		return nil
	}

	score, err := getNameScore(string(firstName))
	if err != nil {
		return err
	}

	addToTree(scoreDB, score, firstName)

	return nil
}

func twoFirstNameLoopScoring() error {
	for i := 0; i < len(keyWords); i++ {
		for j := 1; j < len(keyWords); j++ {
			firstName1 := keyWords[i]
			firstName2 := keyWords[j]

			err := twoFirstNameScoring(firstName1, firstName2)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func twoFirstNameScoring(firstName1, firstName2 rune) error {
	parentNameScore, ok := scoreDB.More[firstName1]
	if !ok {
		return errors.New(fmt.Sprintf("no name score for: %v", firstName1))
	}
	if _, ok := parentNameScore.More[firstName2]; ok {
		return nil
	}

	score, err := getNameScore(string([]rune{firstName1, firstName2}))
	if err != nil {
		return err
	}

	addToTree(scoreDB, score, firstName1, firstName2)
	return nil
}

func addToTree(nameScore *NameScore, score int, firstName ...rune) {
	parent := nameScore
	size := len(firstName)
	for i := 0; i < size; i++ {
		if parent.More == nil {
			parent.More = make(map[rune]*NameScore, 2)
		}
		name := firstName[i]
		if _, ok := parent.More[name]; ok {
			continue
		}
		if i < size-1 {
			parent.More[name] = newScore(0)
			continue
		}
		parent.More[name] = newScore(score)
	}
}

var (
	configFile = flag.String("c", "", "config file")
)

func parseConfig() error {
	flag.Parse()
	err := readConfigFile(*configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %v", *configFile, err)
	}

	err = readScoreData(*configFile)
	if err != nil {
		return fmt.Errorf("failed to read score file: %v", err)
	}

	err = readCandidateData(*configFile)
	if err != nil {
		return fmt.Errorf("failed to read candidate file: %v", err)
	}

	buildNameTestUrlPrefix()

	words := strings.Split(config.FirstNameKeyWords, ",")
	keyWords = make([]rune, len(words))
	for i := 0; i < len(words); i++ {
		keyWords[i] = []rune(words[i])[0]
	}

	return nil
}

var (
	finishChan = make(chan int, 1)
)

func loopScoring() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("namer error: %v", err)
		}
		finishChan <- 1
	}()

	err := nameScoring()
	if err != nil {
		panic(err)
	}
}

func printTop10() {
	found := 0
	for i := 100; i >= 0 && found < 10; i-- {
		if names, ok := scoreStat[i]; ok {
			found++
			fmt.Printf("score: %d, names: %v\n\n", i, names)
		}
	}
}

func startCandidate() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("candidate error: %v", err)
		}
		finishChan <- 1
	}()
	for i := 100; i >= config.MinCandidateScore; i-- {
		if names, ok := scoreStat[i]; ok {
			for _, name := range names {
				if _, ok := candidates[name]; ok {
					continue
				}
				runes := []rune(name[lastNameLength:])
				nameScore := scoreDB
				for _, word := range runes {
					nameScore = nameScore.More[word]
				}
				if nameScore.Read == 1 {
					continue
				}

				fmt.Printf("是否加入候选: %s, 分数: %d  --> (y/n): ", name, nameScore.Score)

				for {
					yn := ""
					_, err := fmt.Scanln(&yn)
					if err != nil {
						fmt.Println(err)
						break
					}
					if yn == "y" {
						nameScore.Read = 1
						candidates[name] = nameScore.Score
						break
					}
					if yn == "n" {
						nameScore.Read = 1
						break
					}
				}
			}
		}
	}
}

func main() {
	err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	go loopScoring()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	select {
	case sig := <-signalChan:
		log.Printf("receive stop signal: %s", sig)
	case <-finishChan:
		log.Println("scoring finish")
	}

	writeScoreData()
	statCalculate(scoreDB, config.LastName)
	printTop10()

	go startCandidate()

	select {
	case sig := <-signalChan:
		log.Printf("receive stop signal: %s", sig)
	case <-finishChan:
		log.Println("candidate finish")
	}

	writeScoreData()
	writeCandidateData()
	printCandidates()
}
