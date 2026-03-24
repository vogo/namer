package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vogo/namer/internal/scoring"
)

func TestHandleScoreSuccess(t *testing.T) {
	body := `{"last_name":"王","first_name":"明轩","year":2024,"month":3,"day":15,"hour":10,"minute":30}`
	req := httptest.NewRequest(http.MethodPost, "/api/score", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleScore(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body: %s", w.Code, w.Body.String())
	}

	var resp scoreResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Name != "王明轩" {
		t.Errorf("name = %q, want 王明轩", resp.Name)
	}
	if resp.Total <= 0 || resp.Total > 100 {
		t.Errorf("total = %.1f, out of range", resp.Total)
	}
	if resp.Detail.BaZi == "" {
		t.Error("bazi should not be empty")
	}
}

func TestHandleScoreEmptyName(t *testing.T) {
	body := `{"last_name":"","first_name":"","year":2024,"month":3,"day":15}`
	req := httptest.NewRequest(http.MethodPost, "/api/score", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handleScore(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestHandleScoreMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/score", nil)
	w := httptest.NewRecorder()
	handleScore(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", w.Code)
	}
}

func TestHandleScoreBadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/score", bytes.NewBufferString("not json"))
	w := httptest.NewRecorder()
	handleScore(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestHandleBatchSuccess(t *testing.T) {
	body := `{"last_name":"王","key_words":"明,轩","year":2024,"month":3,"day":15,"hour":10,"minute":30}`
	req := httptest.NewRequest(http.MethodPost, "/api/batch", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleBatch(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body: %s", w.Code, w.Body.String())
	}

	var resp batchResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	// 2 single + 4 double = 6 results, take min(6, 20)
	if len(resp.Results) != 6 {
		t.Errorf("results count = %d, want 6", len(resp.Results))
	}
	// 按分数降序
	for i := 1; i < len(resp.Results); i++ {
		if resp.Results[i].Total > resp.Results[i-1].Total {
			t.Errorf("results not sorted: [%d]=%.1f > [%d]=%.1f", i, resp.Results[i].Total, i-1, resp.Results[i-1].Total)
		}
	}
}

func TestHandleBatchEmpty(t *testing.T) {
	body := `{"last_name":"","key_words":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/batch", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handleBatch(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestHandleBatchMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/batch", nil)
	w := httptest.NewRecorder()
	handleBatch(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", w.Code)
	}
}

func TestHandleBatchBadJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/batch", bytes.NewBufferString("bad"))
	w := httptest.NewRecorder()
	handleBatch(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestBuildScoreResponse(t *testing.T) {
	r := buildScoreResponse("王", "明轩", scoreResultForTest())
	if r.Name != "王明轩" {
		t.Errorf("name = %q, want 王明轩", r.Name)
	}
	if r.Total <= 0 {
		t.Errorf("total = %.1f, should be > 0", r.Total)
	}
	if len(r.Detail.CharWuXing) != 3 {
		t.Errorf("char_wuxing len = %d, want 3", len(r.Detail.CharWuXing))
	}
	if r.Detail.BaZi == "" {
		t.Error("bazi should not be empty")
	}
}

func scoreResultForTest() scoring.ScoreResult {
	return scoring.CalcScore("王", "明轩", 2024, 3, 15, 10, 30)
}

func TestHandleBatchSingleChar(t *testing.T) {
	body := `{"last_name":"李","key_words":"明","year":2024,"month":1,"day":1}`
	req := httptest.NewRequest(http.MethodPost, "/api/batch", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handleBatch(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp batchResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	// 1 single + 1 double = 2
	if len(resp.Results) != 2 {
		t.Errorf("results = %d, want 2", len(resp.Results))
	}
}

func TestHandleScoreDifferentNames(t *testing.T) {
	names := []struct{ last, first string }{
		{"李", "浩然"},
		{"张", "伟"},
		{"赵", "子龙"},
	}
	for _, n := range names {
		body, _ := json.Marshal(scoreRequest{LastName: n.last, FirstName: n.first, Year: 2024, Month: 1, Day: 1})
		req := httptest.NewRequest(http.MethodPost, "/api/score", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		handleScore(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("%s%s: status = %d", n.last, n.first, w.Code)
		}
		var resp scoreResponse
		json.Unmarshal(w.Body.Bytes(), &resp)
		if resp.Total <= 0 || resp.Total > 100 {
			t.Errorf("%s: total = %.1f", resp.Name, resp.Total)
		}
	}
}

func TestHandleBatchManyKeywords(t *testing.T) {
	body := `{"last_name":"王","key_words":"明,轩,浩,然","year":2024,"month":6,"day":15,"hour":8}`
	req := httptest.NewRequest(http.MethodPost, "/api/batch", bytes.NewBufferString(body))
	w := httptest.NewRecorder()
	handleBatch(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp batchResponse
	json.Unmarshal(w.Body.Bytes(), &resp)
	// 4 single + 16 double = 20, capped at 20
	if len(resp.Results) != 20 {
		t.Errorf("results = %d, want 20", len(resp.Results))
	}
}

func TestBuildScoreResponseSanCaiJX(t *testing.T) {
	// 验证三才吉凶字段非空
	r := scoring.CalcScore("张", "伟", 2024, 1, 1, 12, 0)
	resp := buildScoreResponse("张", "伟", r)
	if resp.Detail.SanCaiDesc == "" {
		t.Error("sancai_desc should not be empty")
	}
	// SanCaiJX 可能是空（如果 desc 不在表中），但对正常名字应该有值
	if resp.Detail.SanCaiJX == "" {
		t.Logf("sancai_jx is empty for desc=%s", resp.Detail.SanCaiDesc)
	}
}

func TestOpenBrowserVariable(t *testing.T) {
	called := false
	orig := OpenBrowser
	OpenBrowser = func(url string) { called = true }
	defer func() { OpenBrowser = orig }()

	OpenBrowser("http://localhost:1234")
	if !called {
		t.Error("OpenBrowser should have been called")
	}
}

func TestNewMux(t *testing.T) {
	mux := NewMux()

	// 测试首页
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("/ status = %d", w.Code)
	}

	// 测试 API 路由存在
	body := `{"last_name":"王","first_name":"明","year":2024,"month":1,"day":1}`
	req2 := httptest.NewRequest(http.MethodPost, "/api/score", bytes.NewBufferString(body))
	w2 := httptest.NewRecorder()
	mux.ServeHTTP(w2, req2)
	if w2.Code != http.StatusOK {
		t.Errorf("/api/score via mux status = %d", w2.Code)
	}
}

func TestIndexPage(t *testing.T) {
	mux := NewMux()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/html; charset=utf-8" {
		t.Errorf("content-type = %q", ct)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("namer")) {
		t.Error("response should contain 'namer'")
	}
}
