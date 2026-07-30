package handlers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"coffee-site-server/internal/auth"
	"coffee-site-server/internal/httpx"
	"coffee-site-server/internal/middleware"
	"coffee-site-server/internal/models"
	"coffee-site-server/internal/store"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "無效的請求內容")
		return
	}
	email := normalizeEmail(req.Email)
	if email == "" || len(req.Password) < 6 {
		httpx.Error(w, http.StatusBadRequest, "請輸入有效的信箱與至少 6 碼密碼")
		return
	}

	if existing, err := h.Store.GetUserByEmail(email); err == nil && existing != nil {
		httpx.Error(w, http.StatusConflict, "此信箱已被註冊")
		return
	} else if err != nil && !errors.Is(err, store.ErrNotFound) {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}

	if _, err := h.Store.CreateUser(email, hash, models.PermissionMember, false); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "註冊失敗，請稍後再試")
		return
	}

	code, err := h.Store.CreateVerificationCode(email, models.PurposeRegister)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	log.Printf("[stub-email] register verification code for %s: %s", email, code)

	httpx.WriteJSON(w, http.StatusCreated, map[string]any{
		"message": "註冊成功，請輸入驗證碼完成帳號啟用",
		"devCode": code,
	})
}

type verifyRegisterRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (h *Handler) VerifyRegisterCode(w http.ResponseWriter, r *http.Request) {
	var req verifyRegisterRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "無效的請求內容")
		return
	}
	email := normalizeEmail(req.Email)

	if err := h.Store.ConsumeVerificationCode(email, req.Code, models.PurposeRegister); err != nil {
		writeCodeError(w, err)
		return
	}
	if err := h.Store.MarkUserVerified(email); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"message": "驗證成功，請登入"})
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "無效的請求內容")
		return
	}
	email := normalizeEmail(req.Email)

	user, err := h.Store.GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			httpx.Error(w, http.StatusUnauthorized, "帳號或密碼錯誤")
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	if !auth.CheckPassword(user.PasswordHash, req.Password) {
		httpx.Error(w, http.StatusUnauthorized, "帳號或密碼錯誤")
		return
	}
	if !user.Verified {
		httpx.Error(w, http.StatusForbidden, "帳號尚未完成驗證")
		return
	}

	token, err := auth.IssueToken(h.Cfg.JWTSecret, user.ID, user.Email, string(user.Permission))
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"token":      token,
		"email":      user.Email,
		"permission": user.Permission,
	})
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req forgotPasswordRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "無效的請求內容")
		return
	}
	email := normalizeEmail(req.Email)

	if _, err := h.Store.GetUserByEmail(email); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			// Do not leak whether the account exists.
			httpx.WriteJSON(w, http.StatusOK, map[string]any{"message": "若帳號存在，驗證碼已寄出"})
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}

	code, err := h.Store.CreateVerificationCode(email, models.PurposeReset)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	log.Printf("[stub-email] password reset code for %s: %s", email, code)

	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"message": "驗證碼已寄出",
		"devCode": code,
	})
}

type verifyResetCodeRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (h *Handler) VerifyResetCode(w http.ResponseWriter, r *http.Request) {
	var req verifyResetCodeRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "無效的請求內容")
		return
	}
	email := normalizeEmail(req.Email)

	if err := h.Store.ValidateVerificationCode(email, req.Code, models.PurposeReset); err != nil {
		writeCodeError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"valid": true})
}

type resetPasswordRequest struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"newPassword"`
}

func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetPasswordRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "無效的請求內容")
		return
	}
	email := normalizeEmail(req.Email)
	if len(req.NewPassword) < 6 {
		httpx.Error(w, http.StatusBadRequest, "密碼至少需要 6 碼")
		return
	}

	if err := h.Store.ConsumeVerificationCode(email, req.Code, models.PurposeReset); err != nil {
		writeCodeError(w, err)
		return
	}

	hash, err := auth.HashPassword(req.NewPassword)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	if err := h.Store.UpdateUserPassword(email, hash); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{"message": "密碼已重設，請重新登入"})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "請先登入")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"email":      claims.Email,
		"permission": claims.Permission,
	})
}

func (h *Handler) VerifyPermission(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.ClaimsFromContext(r.Context())
	httpx.WriteJSON(w, http.StatusOK, map[string]any{
		"authorized": ok && claims.Permission == "admin",
	})
}

func writeCodeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, store.ErrInvalidCode):
		httpx.Error(w, http.StatusBadRequest, "驗證碼錯誤")
	case errors.Is(err, store.ErrCodeExpired):
		httpx.Error(w, http.StatusBadRequest, "驗證碼已過期，請重新取得")
	default:
		httpx.Error(w, http.StatusInternalServerError, "伺服器錯誤")
	}
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
