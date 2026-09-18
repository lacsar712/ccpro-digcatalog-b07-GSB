package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"digcatalog/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupWorkspaceTest(t *testing.T) *Handler {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.User{}, &models.Site{}, &models.Unit{}, &models.Material{}, &models.Find{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return &Handler{DB: db}
}

func callWorkspaceToday(t *testing.T, h *Handler) map[string]any {
	t.Helper()
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/api/workspace/today", nil)
	h.WorkspaceToday(c)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	var out map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return out
}

func TestWorkspaceToday(t *testing.T) {
	h := setupWorkspaceTest(t)

	site := models.Site{Name: "测试工地", Period: "商周"}
	h.DB.Create(&site)
	unit := models.Unit{SiteID: site.ID, Code: "T1"}
	h.DB.Create(&unit)

	// 东八区今日边界（与 handler 同一口径）；CreatedAt 用本地时间写入，与生产一致
	now := time.Now().In(cst8)
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, cst8).In(time.Local)

	seq := 0
	mk := func(createdAt time.Time) models.Find {
		seq++
		return models.Find{
			UnitID:       unit.ID,
			RegisterNo:   "REG-" + strconv.Itoa(seq),
			ArtifactType: "陶片",
			CreatedAt:    createdAt,
		}
	}

	// 今日 3 条（凌晨、现在、今晚 23 点）
	todayTimes := []time.Time{dayStart.Add(1 * time.Hour), time.Now(), dayStart.Add(23 * time.Hour)}
	for _, tt := range todayTimes {
		f := mk(tt)
		if err := h.DB.Create(&f).Error; err != nil {
			t.Fatalf("create find: %v", err)
		}
	}
	// 昨日 23 点（东八区）1 条 —— 不应计入今日
	yesterday := mk(dayStart.Add(-1 * time.Hour))
	h.DB.Create(&yesterday)
	// 软删除一条今日的 —— 不应计入
	deleted := mk(dayStart.Add(2 * time.Hour))
	h.DB.Create(&deleted)
	h.DB.Delete(&deleted)

	out := callWorkspaceToday(t, h)

	if got := out["todayFindCount"]; got != float64(3) {
		t.Errorf("todayFindCount = %v, want 3", got)
	}
	if got := out["siteCount"]; got != float64(1) {
		t.Errorf("siteCount = %v, want 1", got)
	}
	if got := out["unitCount"]; got != float64(1) {
		t.Errorf("unitCount = %v, want 1", got)
	}
	if got := out["findCount"]; got != float64(4) {
		t.Errorf("findCount = %v, want 4 (软删除不计)", got)
	}

	recent, ok := out["recentFinds"].([]any)
	if !ok {
		t.Fatalf("recentFinds not an array: %T", out["recentFinds"])
	}
	if len(recent) != 4 {
		t.Fatalf("len(recentFinds) = %d, want 4", len(recent))
	}
	first := recent[0].(map[string]any)
	if first["unitCode"] != "T1" || first["siteName"] != "测试工地" {
		t.Errorf("recent[0] unit/site = %v/%v, want T1/测试工地", first["unitCode"], first["siteName"])
	}
	// id 倒序：第一条应是最后创建且未软删的记录
	lastID := first["id"].(float64)
	for _, item := range recent[1:] {
		id := item.(map[string]any)["id"].(float64)
		if id >= lastID {
			t.Errorf("recentFinds 未按 id 倒序: %v >= %v", id, lastID)
		}
		lastID = id
	}
}

func TestWorkspaceTodayRecentLimit5(t *testing.T) {
	h := setupWorkspaceTest(t)
	site := models.Site{Name: "S", Period: "P"}
	h.DB.Create(&site)
	unit := models.Unit{SiteID: site.ID, Code: "T9"}
	h.DB.Create(&unit)
	for i := 0; i < 8; i++ {
		h.DB.Create(&models.Find{UnitID: unit.ID, RegisterNo: "R" + string(rune('A'+i)), ArtifactType: "骨器"})
	}
	out := callWorkspaceToday(t, h)
	recent := out["recentFinds"].([]any)
	if len(recent) != 5 {
		t.Errorf("len(recentFinds) = %d, want 5 (上限)", len(recent))
	}
	if got := out["todayFindCount"]; got != float64(8) {
		t.Errorf("todayFindCount = %v, want 8", got)
	}
}
