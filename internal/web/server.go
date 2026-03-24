package web

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os/exec"
	"runtime"
	"strings"

	"github.com/vogo/namer/internal/data"
	"github.com/vogo/namer/internal/scoring"
)

// scoreRequest 单个名字评分请求
type scoreRequest struct {
	LastName  string `json:"last_name"`
	FirstName string `json:"first_name"`
	Year      int    `json:"year"`
	Month     int    `json:"month"`
	Day       int    `json:"day"`
	Hour      int    `json:"hour"`
	Minute    int    `json:"minute"`
}

// batchRequest 批量评分请求
type batchRequest struct {
	LastName string `json:"last_name"`
	Year     int    `json:"year"`
	Month    int    `json:"month"`
	Day      int    `json:"day"`
	Hour     int    `json:"hour"`
	Minute   int    `json:"minute"`
	KeyWords string `json:"key_words"`
}

// scoreResponse 评分响应
type scoreResponse struct {
	Name    string      `json:"name"`
	Total   float64     `json:"total"`
	WuGe    float64     `json:"wuge"`
	SanCai  float64     `json:"sancai"`
	XiYong  float64     `json:"xiyong"`
	WuXing  float64     `json:"wuxing"`
	YinYang float64     `json:"yinyang"`
	Detail  scoreDetail `json:"detail"`
}

type scoreDetail struct {
	Strokes    []int    `json:"strokes"`
	TianGe     int      `json:"tian_ge"`
	RenGe      int      `json:"ren_ge"`
	DiGe       int      `json:"di_ge"`
	ZongGe     int      `json:"zong_ge"`
	WaiGe      int      `json:"wai_ge"`
	SanCaiDesc string   `json:"sancai_desc"`
	SanCaiJX   string   `json:"sancai_jx"`
	BaZi       string   `json:"bazi"`
	XiYongShen string   `json:"xiyong_shen"`
	CharWuXing []string `json:"char_wuxing"`
	YinYangPat string   `json:"yinyang_pat"`
}

type batchResponse struct {
	Results []scoreResponse `json:"results"`
}

func buildScoreResponse(lastName, firstName string, r scoring.ScoreResult) scoreResponse {
	allRunes := []rune(lastName + firstName)
	charWx := make([]string, len(r.CharWuXing))
	for i, wx := range r.CharWuXing {
		charWx[i] = wx.String()
	}

	strokeNames := make([]int, len(allRunes))
	for i, c := range allRunes {
		strokeNames[i] = data.CharStroke(c)
	}

	sanCaiJX := ""
	if jx, ok := data.SanCaiJiXiong[r.SanCaiDesc]; ok {
		sanCaiJX = jx.String()
	}

	return scoreResponse{
		Name:    lastName + firstName,
		Total:   r.Total,
		WuGe:    r.WuGeScore,
		SanCai:  r.SanCaiScore,
		XiYong:  r.XiYongScore,
		WuXing:  r.WuXingScore,
		YinYang: r.YinYangScore,
		Detail: scoreDetail{
			Strokes:    r.Strokes,
			TianGe:     r.WuGe.TianGe,
			RenGe:      r.WuGe.RenGe,
			DiGe:       r.WuGe.DiGe,
			ZongGe:     r.WuGe.ZongGe,
			WaiGe:      r.WuGe.WaiGe,
			SanCaiDesc: r.SanCaiDesc,
			SanCaiJX:   sanCaiJX,
			BaZi:       r.BaZi.String(),
			XiYongShen: r.XiYongShen.String(),
			CharWuXing: charWx,
			YinYangPat: r.YinYangPattern,
		},
	}
}

func handleScore(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req scoreRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.LastName == "" || req.FirstName == "" {
		http.Error(w, "姓和名不能为空", http.StatusBadRequest)
		return
	}

	result := scoring.CalcScore(req.LastName, req.FirstName, req.Year, req.Month, req.Day, req.Hour, req.Minute)
	resp := buildScoreResponse(req.LastName, req.FirstName, result)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func handleBatch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req batchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if req.LastName == "" || req.KeyWords == "" {
		http.Error(w, "姓氏和备选字不能为空", http.StatusBadRequest)
		return
	}

	words := strings.Split(req.KeyWords, ",")
	var keyWords []rune
	for _, w := range words {
		w = strings.TrimSpace(w)
		if w != "" {
			keyWords = append(keyWords, []rune(w)[0])
		}
	}

	type nameScore struct {
		lastName  string
		firstName string
		result    scoring.ScoreResult
	}

	var all []nameScore

	// 单字名
	for _, c := range keyWords {
		fn := string(c)
		res := scoring.CalcScore(req.LastName, fn, req.Year, req.Month, req.Day, req.Hour, req.Minute)
		all = append(all, nameScore{req.LastName, fn, res})
	}

	// 双字名
	for _, c1 := range keyWords {
		for _, c2 := range keyWords {
			fn := string([]rune{c1, c2})
			res := scoring.CalcScore(req.LastName, fn, req.Year, req.Month, req.Day, req.Hour, req.Minute)
			all = append(all, nameScore{req.LastName, fn, res})
		}
	}

	// 按分数排序取 Top 20
	for i := 0; i < len(all); i++ {
		for j := i + 1; j < len(all); j++ {
			if all[j].result.Total > all[i].result.Total {
				all[i], all[j] = all[j], all[i]
			}
		}
	}

	count := min(len(all), 20)
	resp := batchResponse{Results: make([]scoreResponse, count)}
	for i := 0; i < count; i++ {
		resp.Results[i] = buildScoreResponse(all[i].lastName, all[i].firstName, all[i].result)
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

// OpenBrowser 打开浏览器（可替换用于测试）
var OpenBrowser = func(url string) {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return
	}
	_ = cmd.Start()
}

// NewMux 创建 HTTP handler
func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write([]byte(indexHTML))
	})
	mux.HandleFunc("/api/score", handleScore)
	mux.HandleFunc("/api/batch", handleBatch)
	return mux
}

// Start 启动 web 服务
func Start(port int) error {
	mux := NewMux()

	var listener net.Listener
	var err error

	if port > 0 {
		listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	} else {
		listener, err = net.Listen("tcp", ":0")
	}
	if err != nil {
		return err
	}

	addr := listener.Addr().(*net.TCPAddr)
	url := fmt.Sprintf("http://localhost:%d", addr.Port)
	fmt.Printf("namer web 服务已启动: %s\n", url)

	OpenBrowser(url)

	return http.Serve(listener, mux)
}
