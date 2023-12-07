package handlers

import (
	"8th_pract_go/internal/config"
	"encoding/json"
	"fmt"
	"github.com/gorilla/securecookie"
	"log"
	"net/http"
	"time"
)

type ApiCfg struct {
	CookieName   string
	SecureCookie *securecookie.SecureCookie
}

func NewApi(cookiename, encrKey string) *ApiCfg {
	return &ApiCfg{
		CookieName:   cookiename,
		SecureCookie: securecookie.New([]byte(encrKey), nil),
	}
}

// SetApiCookie Метод структуры ApiCfg
// Кодирование и запись куки
func (cfg *ApiCfg) SetApiCookie(w http.ResponseWriter, r config.Request) {
	encoded, err := cfg.SecureCookie.Encode(cfg.CookieName, r)
	if err != nil {
		log.Println("Error while encoding cookie:", err)
		http.Error(w, "Error while encoding cookie:", 400)
		return
	}

	// Записываем куки
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.CookieName,
		Value:    encoded,
		Path:     "/",
		MaxAge:   360,
		Secure:   true,
		HttpOnly: true,
	})
}

// LinearHandler Линейная ручка
func (cfg *ApiCfg) LinearHandler(w http.ResponseWriter, r *http.Request) {
	var req config.Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", 400)
	}

	var res config.Response
	res.Data = fmt.Sprintf("Creating cookie with Email:%s, Password:%s)", req.Email, req.Password)
	cfg.SetApiCookie(w, req)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Can't encode JSON", 400)
	}
}

// ConcurrentHandler Конкурентная ручка
func (cfg *ApiCfg) ConcurrentHandler(w http.ResponseWriter, r *http.Request) {
	var req config.Request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON", 400)
	}

	var res config.Response

	go func(reqCopy *config.Request, resCopy *config.Response) {
		resCopy.Data = fmt.Sprintf("Creating cookie with Email:%s, Password:%s", reqCopy.Email, reqCopy.Password)
		cfg.SetApiCookie(w, *reqCopy)
	}(&req, &res)

	time.Sleep(5 * time.Second)

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		http.Error(w, "Can't encode JSON", 500)
	}
}

// GetCookie чтение куки
func (cfg *ApiCfg) GetCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cfg.CookieName)
	if err != nil {
		http.Error(w, "Cookie not found", 404)
		return
	}

	var req config.Request

	err = cfg.SecureCookie.Decode(cfg.CookieName, cookie.Value, &req)
	if err != nil {
		http.Error(w, "Error while decoding cookie", 500)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(req)
	if err != nil {
		http.Error(w, "Can't encode JSON", 500)
	}
}
