package web_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/vogo/namer/integrations/helper"
)

func startWebServer(t *testing.T) (port int, cleanup func()) {
	t.Helper()

	bin := helper.NamerBinary(t)

	// Find a free port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	port = listener.Addr().(*net.TCPAddr).Port
	_ = listener.Close()

	cmd := exec.Command(bin, "-web", "-port", fmt.Sprintf("%d", port))
	cmd.Env = append(os.Environ(), "HOME="+t.TempDir())

	// Prevent browser from opening
	if runtime.GOOS == "darwin" {
		cmd.Env = append(cmd.Env, "BROWSER=false")
	}

	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start web server: %v", err)
	}

	// Wait for server to be ready
	deadline := time.Now().Add(5 * time.Second)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%d", port), 100*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			break
		}
		time.Sleep(50 * time.Millisecond)
	}

	return port, func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}
}

func TestWebHomePage(t *testing.T) {
	port, cleanup := startWebServer(t)
	defer cleanup()

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/", port))
	if err != nil {
		t.Fatalf("failed to GET /: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html, got %s", ct)
	}
}

func TestWebScoreAPI(t *testing.T) {
	port, cleanup := startWebServer(t)
	defer cleanup()

	body := `{
		"last_name": "王",
		"first_name": "明轩",
		"year": 2024,
		"month": 3,
		"day": 15,
		"hour": 10,
		"minute": 30
	}`

	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/api/score", port),
		"application/json",
		bytes.NewBufferString(body),
	)
	if err != nil {
		t.Fatalf("POST /api/score failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Verify response structure
	for _, key := range []string{"name", "total", "wuge", "sancai", "xiyong", "wuxing", "yinyang", "detail"} {
		if _, ok := result[key]; !ok {
			t.Errorf("response missing key %q", key)
		}
	}

	if result["name"] != "王明轩" {
		t.Errorf("expected name 王明轩, got %v", result["name"])
	}

	total, ok := result["total"].(float64)
	if !ok || total <= 0 || total > 100 {
		t.Errorf("expected total in (0, 100], got %v", result["total"])
	}
}

func TestWebScoreAPIMethodNotAllowed(t *testing.T) {
	port, cleanup := startWebServer(t)
	defer cleanup()

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/api/score", port))
	if err != nil {
		t.Fatalf("GET /api/score failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", resp.StatusCode)
	}
}

func TestWebScoreAPIEmptyName(t *testing.T) {
	port, cleanup := startWebServer(t)
	defer cleanup()

	body := `{"last_name": "", "first_name": ""}`
	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/api/score", port),
		"application/json",
		bytes.NewBufferString(body),
	)
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestWebBatchAPI(t *testing.T) {
	port, cleanup := startWebServer(t)
	defer cleanup()

	body := `{
		"last_name": "王",
		"year": 2024,
		"month": 3,
		"day": 15,
		"hour": 10,
		"minute": 30,
		"key_words": "明,轩"
	}`

	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/api/batch", port),
		"application/json",
		bytes.NewBufferString(body),
	)
	if err != nil {
		t.Fatalf("POST /api/batch failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	results, ok := result["results"].([]any)
	if !ok || len(results) == 0 {
		t.Error("expected non-empty results array")
	}

	// With 2 keywords: 2 single-char + 4 double-char = 6 results
	if len(results) != 6 {
		t.Errorf("expected 6 results, got %d", len(results))
	}
}

func TestWebBatchAPIEmptyKeywords(t *testing.T) {
	port, cleanup := startWebServer(t)
	defer cleanup()

	body := `{"last_name": "王", "key_words": ""}`
	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/api/batch", port),
		"application/json",
		bytes.NewBufferString(body),
	)
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestWebScoreAPIInvalidJSON(t *testing.T) {
	port, cleanup := startWebServer(t)
	defer cleanup()

	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/api/score", port),
		"application/json",
		bytes.NewBufferString("{invalid json}"),
	)
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestWebScoreResponseDetail(t *testing.T) {
	port, cleanup := startWebServer(t)
	defer cleanup()

	body := `{
		"last_name": "李",
		"first_name": "明",
		"year": 2000,
		"month": 6,
		"day": 1,
		"hour": 8,
		"minute": 0
	}`

	resp, err := http.Post(
		fmt.Sprintf("http://localhost:%d/api/score", port),
		"application/json",
		bytes.NewBufferString(body),
	)
	if err != nil {
		t.Fatalf("POST failed: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var result struct {
		Name   string  `json:"name"`
		Total  float64 `json:"total"`
		Detail struct {
			Strokes    []int    `json:"strokes"`
			TianGe     int      `json:"tian_ge"`
			RenGe      int      `json:"ren_ge"`
			DiGe       int      `json:"di_ge"`
			ZongGe     int      `json:"zong_ge"`
			WaiGe      int      `json:"wai_ge"`
			SanCaiDesc string   `json:"sancai_desc"`
			BaZi       string   `json:"bazi"`
			XiYongShen string   `json:"xiyong_shen"`
			CharWuXing []string `json:"char_wuxing"`
			YinYangPat string   `json:"yinyang_pat"`
		} `json:"detail"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}

	if result.Name != "李明" {
		t.Errorf("expected 李明, got %s", result.Name)
	}

	if len(result.Detail.Strokes) != 2 {
		t.Errorf("expected 2 strokes, got %d", len(result.Detail.Strokes))
	}

	if result.Detail.TianGe <= 0 {
		t.Error("expected positive TianGe")
	}

	if result.Detail.BaZi == "" {
		t.Error("expected non-empty BaZi")
	}

	if result.Detail.XiYongShen == "" {
		t.Error("expected non-empty XiYongShen")
	}

	if len(result.Detail.CharWuXing) != 2 {
		t.Errorf("expected 2 char wuxing, got %d", len(result.Detail.CharWuXing))
	}
}
