package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"chuan/internal/services"
)

type Handler struct {
	webrtcService *services.WebRTCService
}

func NewHandler() *Handler {
	return &Handler{
		webrtcService: services.NewWebRTCService(),
	}
}

// HandleWebRTCWebSocket 处理WebRTC信令WebSocket连接
func (h *Handler) HandleWebRTCWebSocket(w http.ResponseWriter, r *http.Request) {
	h.webrtcService.HandleWebSocket(w, r)
}

// CreateRoomHandler 创建房间API - 简化版本，不处理无用参数
func (h *Handler) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应为JSON格式
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "方法不允许",
		})
		return
	}

	// 创建新房间（忽略请求体中的无用参数）
	code := h.webrtcService.CreateNewRoom()
	log.Printf("创建房间成功: %s", code)

	// 构建响应
	response := map[string]interface{}{
		"success": true,
		"code":    code,
		"message": "房间创建成功",
	}

	json.NewEncoder(w).Encode(response)
}

// WebRTCRoomStatusHandler WebRTC房间状态API
func (h *Handler) WebRTCRoomStatusHandler(w http.ResponseWriter, r *http.Request) {
	// 设置响应为JSON格式
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "方法不允许",
		})
		return
	}

	// 从查询参数获取房间代码
	code := r.URL.Query().Get("code")
	if code == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "缺少房间代码",
		})
		return
	}

	// 获取房间状态
	status := h.webrtcService.GetRoomStatus(code)

	json.NewEncoder(w).Encode(status)
}

// GetRoomStatusHandler 获取房间状态API
func (h *Handler) GetRoomStatusHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "方法不允许",
		})
		return
	}

	// 获取房间码
	code := r.URL.Query().Get("code")
	if code == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "房间码不能为空",
		})
		return
	}

	// 获取房间状态
	status := h.webrtcService.GetRoomStatus(code)
	json.NewEncoder(w).Encode(status)
}

const AuthCookieName = "auth_token"

// ConfigHandler 返回认证配置
func (h *Handler) ConfigHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allowCode := os.Getenv("ALLOW_CODE")
	response := map[string]interface{}{
		"auth_enabled": allowCode != "",
	}
	json.NewEncoder(w).Encode(response)
}

// LoginHandler 处理登录请求
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "方法不允许", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "无效的请求", http.StatusBadRequest)
		return
	}

	allowCode := os.Getenv("ALLOW_CODE")
	if allowCode != "" && req.Code == allowCode {
		// 登录成功，设置 cookie
		expiration := time.Now().Add(24 * time.Hour)
		cookie := http.Cookie{
			Name:     AuthCookieName,
			Value:    "authenticated",
			Expires:  expiration,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		}
		http.SetCookie(w, &cookie)
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true})
		return
	}

	// 登录失败
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "message": "无效的访问码"})
}

// CheckAuthHandler 检查认证状态
func (h *Handler) CheckAuthHandler(w http.ResponseWriter, r *http.Request) {
	allowCode := os.Getenv("ALLOW_CODE")
	if allowCode == "" {
		json.NewEncoder(w).Encode(map[string]interface{}{"authenticated": true})
		return
	}

	cookie, err := r.Cookie(AuthCookieName)
	if err != nil || cookie.Value != "authenticated" {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]interface{}{"authenticated": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"authenticated": true})
}
