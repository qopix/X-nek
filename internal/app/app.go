package app

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/qopix/nek/internal/model"
	"github.com/qopix/nek/internal/protocol"
	"github.com/qopix/nek/internal/storage"
	"net/http"
	"strings"
	"time"
)

type Server struct {
	Store                            *storage.Store
	Token, GeneratedToken, PublicURL string
}

func New(path, token, publicURL string) (*Server, error) {
	st, e := storage.Open(path)
	if e != nil {
		return nil, e
	}
	s := &Server{Store: st, Token: token, PublicURL: strings.TrimRight(publicURL, "/")}
	if s.Token == "" {
		s.Token, e = randomToken(32)
		if e != nil {
			return nil, e
		}
		s.GeneratedToken = s.Token
	}
	return s, nil
}
func (s *Server) Handler() http.Handler {
	m := http.NewServeMux()
	m.HandleFunc("/", s.static)
	m.HandleFunc("/api/health", s.health)
	m.HandleFunc("/api/overview", s.auth(s.overview))
	m.HandleFunc("/api/protocols", s.auth(s.protocols))
	m.HandleFunc("/api/nodes", s.auth(s.nodes))
	m.HandleFunc("/api/users", s.auth(s.users))
	m.HandleFunc("/api/users/", s.auth(s.userAction))
	m.HandleFunc("/api/subscriptions", s.auth(s.subscriptions))
	m.HandleFunc("/sub/", s.subscriptionDocument)
	return secure(m)
}
func (s *Server) auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		got := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if got == "" || !equal(got, s.Token) {
			jsonOut(w, 401, map[string]string{"error": "unauthorized"})
			return
		}
		next(w, r)
	}
}
func (s *Server) static(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "web/index.html")
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	jsonOut(w, 200, map[string]any{"ok": true, "service": "x-nek", "time": time.Now().UTC()})
}
func (s *Server) overview(w http.ResponseWriter, r *http.Request) {
	st := s.Store.View()
	jsonOut(w, 200, map[string]any{"nodes": len(st.Nodes), "users": len(st.Users), "subscriptions": len(st.Subscriptions), "protocols": len(protocol.Catalog()), "product": "X-NEK"})
}
func (s *Server) protocols(w http.ResponseWriter, r *http.Request) {
	jsonOut(w, 200, protocol.Catalog())
}
func (s *Server) nodes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		jsonOut(w, 200, s.Store.View().Nodes)
		return
	}
	if r.Method != http.MethodPost {
		method(w)
		return
	}
	var in struct {
		Name     string `json:"name"`
		Address  string `json:"address"`
		AgentURL string `json:"agent_url"`
	}
	if !decode(w, r, &in) {
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Address = strings.TrimSpace(in.Address)
	if in.Name == "" || in.Address == "" {
		jsonOut(w, 400, map[string]string{"error": "name and address required"})
		return
	}
	n := model.Node{ID: id(), Name: in.Name, Address: in.Address, AgentURL: in.AgentURL, Status: "unknown", CreatedAt: time.Now().UTC()}
	if e := s.Store.Update(func(st *model.State) error { st.Nodes = append(st.Nodes, n); return nil }); e != nil {
		fail(w, e)
		return
	}
	jsonOut(w, 201, n)
}
func (s *Server) users(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		jsonOut(w, 200, s.Store.View().Users)
		return
	}
	if r.Method != http.MethodPost {
		method(w)
		return
	}
	var in struct {
		Username string `json:"username"`
	}
	if !decode(w, r, &in) {
		return
	}
	in.Username = strings.TrimSpace(in.Username)
	if in.Username == "" || len(in.Username) > 80 {
		jsonOut(w, 400, map[string]string{"error": "username required (max 80 chars)"})
		return
	}
	u := model.User{ID: id(), Username: in.Username, Enabled: true, CreatedAt: time.Now().UTC()}
	if e := s.Store.Update(func(st *model.State) error {
		for _, x := range st.Users {
			if strings.EqualFold(x.Username, u.Username) {
				return errors.New("username already exists")
			}
		}
		st.Users = append(st.Users, u)
		return nil
	}); e != nil {
		fail(w, e)
		return
	}
	jsonOut(w, 201, u)
}
func (s *Server) userAction(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(p) != 4 || p[0] != "api" || p[1] != "users" || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	enabled := p[3] == "enable"
	if p[3] != "enable" && p[3] != "disable" {
		http.NotFound(w, r)
		return
	}
	found := false
	e := s.Store.Update(func(st *model.State) error {
		for i := range st.Users {
			if st.Users[i].ID == p[2] {
				st.Users[i].Enabled = enabled
				found = true
				return nil
			}
		}
		return errors.New("user not found")
	})
	if e != nil {
		fail(w, e)
		return
	}
	if !found {
		jsonOut(w, 404, map[string]string{"error": "user not found"})
		return
	}
	jsonOut(w, 200, map[string]any{"ok": true, "id": p[2], "enabled": enabled})
}
func (s *Server) subscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		st := s.Store.View()
		out := make([]map[string]any, 0, len(st.Subscriptions))
		for _, x := range st.Subscriptions {
			out = append(out, map[string]any{"id": x.ID, "user_id": x.UserID, "protocol": x.Protocol, "node_id": x.NodeID, "expires_at": x.ExpiresAt, "enabled": x.Enabled, "url": s.PublicURL + "/sub/" + x.Token})
		}
		jsonOut(w, 200, out)
		return
	}
	if r.Method != http.MethodPost {
		method(w)
		return
	}
	var in struct {
		UserID   string `json:"user_id"`
		Protocol string `json:"protocol"`
		NodeID   string `json:"node_id"`
		Days     int    `json:"days"`
	}
	if !decode(w, r, &in) {
		return
	}
	if in.Days != 7 && in.Days != 30 && in.Days != 90 {
		jsonOut(w, 400, map[string]string{"error": "days must be 7, 30, or 90"})
		return
	}
	known := false
	for _, p := range protocol.Catalog() {
		if p.ID == in.Protocol {
			known = true
		}
	}
	if !known {
		jsonOut(w, 400, map[string]string{"error": "unknown protocol id"})
		return
	}
	st := s.Store.View()
	userOK := false
	for _, u := range st.Users {
		if u.ID == in.UserID && u.Enabled {
			userOK = true
		}
	}
	if !userOK {
		jsonOut(w, 400, map[string]string{"error": "enabled user not found"})
		return
	}
	if in.NodeID != "" {
		ok := false
		for _, n := range st.Nodes {
			if n.ID == in.NodeID {
				ok = true
			}
		}
		if !ok {
			jsonOut(w, 400, map[string]string{"error": "node not found"})
			return
		}
	}
	tok, e := randomToken(24)
	if e != nil {
		fail(w, e)
		return
	}
	x := model.Subscription{ID: id(), UserID: in.UserID, Protocol: in.Protocol, NodeID: in.NodeID, ExpiresAt: time.Now().UTC().Add(time.Duration(in.Days) * 24 * time.Hour), Token: tok, Enabled: true, CreatedAt: time.Now().UTC()}
	if e = s.Store.Update(func(st *model.State) error { st.Subscriptions = append(st.Subscriptions, x); return nil }); e != nil {
		fail(w, e)
		return
	}
	jsonOut(w, 201, map[string]any{"id": x.ID, "user_id": x.UserID, "protocol": x.Protocol, "expires_at": x.ExpiresAt, "enabled": x.Enabled, "url": s.PublicURL + "/sub/" + x.Token})
}
func (s *Server) subscriptionDocument(w http.ResponseWriter, r *http.Request) {
	token := strings.TrimPrefix(r.URL.Path, "/sub/")
	st := s.Store.View()
	for _, x := range st.Subscriptions {
		if equal(x.Token, token) && x.Enabled && time.Now().Before(x.ExpiresAt) {
			for _, u := range st.Users {
				if u.ID == x.UserID && u.Enabled {
					jsonOut(w, 200, map[string]any{"format": "nek-subscription-v1", "status": "provisioning-required", "protocol": x.Protocol, "message": "Placeholder only: a transport adapter must provision a real client configuration."})
					return
				}
			}
		}
	}
	http.NotFound(w, r)
}
func decode(w http.ResponseWriter, r *http.Request, v any) bool {
	defer r.Body.Close()
	if e := json.NewDecoder(r.Body).Decode(v); e != nil {
		jsonOut(w, 400, map[string]string{"error": "invalid JSON"})
		return false
	}
	return true
}
func jsonOut(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func method(w http.ResponseWriter) { jsonOut(w, 405, map[string]string{"error": "method not allowed"}) }
func fail(w http.ResponseWriter, e error) {
	status := 500
	if strings.Contains(e.Error(), "already exists") {
		status = 409
	}
	if strings.Contains(e.Error(), "not found") {
		status = 404
	}
	jsonOut(w, status, map[string]string{"error": e.Error()})
}
func id() string {
	b := make([]byte, 12)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	return hex.EncodeToString(b)
}
func randomToken(n int) (string, error) {
	b := make([]byte, n)
	if _, e := rand.Read(b); e != nil {
		return "", e
	}
	return hex.EncodeToString(b), nil
}
func equal(a, b string) bool {
	ha := sha256.Sum256([]byte(a))
	hb := sha256.Sum256([]byte(b))
	var x byte
	for i := range ha {
		x |= ha[i] ^ hb[i]
	}
	return x == 0
}
func secure(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		next.ServeHTTP(w, r)
	})
}
