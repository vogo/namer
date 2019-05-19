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
	"strconv"
	"strings"
	"time"
)

var (
	baiduQimingTestUrlFormat = "https://sp0.baidu.com/5LMDcjW6BwF3otqbppnN2DJv/qiming.pae.baidu.com/data/namedetail?" +
		"year=%d&month=%02d&day=%02d&hour=%02d&min=%02d&timeType=0&gender=0" +
		"&flag=1&cb=jsonp1"
	baiduTestUrlPrefix string
	keyWords           []rune

	config    = &Config{}
	scoreData = newScore(0)
	scoreStat = make(map[int][]string, 100)

	scorePreviousString = "\"key\":\"score\",\"value\":\""
)

type NameScore struct {
	Score int                 `json:"score"`
	More  map[rune]*NameScore `json:"more,omitempty"`
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
	index := strings.Index(content, scorePreviousString)
	if index <= 0 {
		return 0, errors.New("cant find score from response: " + content)
	}

	content = content[index+len(scorePreviousString):]
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
	return json.Unmarshal(bt, config)
}

func readScoreData(file string) error {
	file += ".score"
	bt, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}
	return json.Unmarshal(bt, scoreData)
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

func statCalculate(score *NameScore, name string) {
	scoreStat[score.Score] = append(scoreStat[score.Score], name)
	for key, value := range score.More {
		statCalculate(value, name+string(key))
	}
}

func writeScoreData() {
	bt, err := json.Marshal(scoreData)
	if err != nil {
		log.Fatalf("marshal err: %v", err)
	}
	err = ioutil.WriteFile(*configFile+".score", bt, 0666)
	if err != nil {
		log.Fatalf("write score file error: %v", err)
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
	if scoreData == nil {
		scoreData = newScore(0)
	}

	for i := 0; i < len(keyWords); i++ {
		err := oneFirstNameScoring(scoreData, keyWords[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func oneFirstNameScoring(parent *NameScore, firstName rune) error {
	if _, ok := scoreData.More[firstName]; ok {
		return nil
	}

	score, err := getNameScore(string(firstName))
	if err != nil {
		return err
	}

	parent.More[firstName] = newScore(score)
	return nil
}

func twoFirstNameLoopScoring() error {
	for i := 0; i < len(keyWords)-1; i++ {
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
	parentNameScore, ok := scoreData.More[firstName1]
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

	if parentNameScore.More == nil {
		parentNameScore.More = make(map[rune]*NameScore, 2)
	}

	parentNameScore.More[firstName2] = newScore(score)
	return nil
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
	buildNameTestUrlPrefix()

	words := strings.Split(config.FirstNameKeyWords, ",")
	keyWords = make([]rune, len(words))
	for i := 0; i < len(words); i++ {
		keyWords[i] = []rune(words[i])[0]
	}
	return nil
}

func loopScoring() {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("namer error: %v", err)
		}
	}()

	err := nameScoring()
	if err != nil {
		panic(err)
	}
}

func main() {
	err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}
	loopScoring()
	writeScoreData()
	statCalculate(scoreData, config.LastName)
	printTop10()
}
