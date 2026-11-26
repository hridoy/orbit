package server

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"regexp"
	"strings"

	"orbit/internal/db"
	"orbit/internal/security"
)

var usernameRe = regexp.MustCompile(`^[a-zA-Z0-9._-]{2,64}$`)

var validRoles = map[string]bool{"admin": true, "operator": true, "viewer": true}

// UsersHandler manages dashboard members. Only the admin role reaches these routes, enforced by
// the auth middleware.
type UsersHandler struct {
	users *db.DashboardUserStore
}

func NewUsersHandler(users *db.DashboardUserStore) *UsersHandler {
	return &UsersHandler{users: users}
}

func (h *UsersHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/users", h.handleList)
	mux.HandleFunc("POST /api/v1/users", h.handleCreate)
	mux.HandleFunc("PUT /api/v1/users/{username}", h.handleUpdate)
	mux.HandleFunc("DELETE /api/v1/users/{username}", h.handleDelete)
	mux.HandleFunc("POST /api/v1/users/me/password", h.handleOwnPassword)
}

type memberView struct {
	Username  string `json:"username"`
	Role      string `json:"role"`
	Disabled  bool   `json:"disabled"`
	CreatedAt string `json:"created_at"`
}

func (h *UsersHandler) handleList(w http.ResponseWriter, r *http.Request) {
	list, err := h.users.List()
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	out := []memberView{}
	for _, u := range list {
		out = append(out, memberView{Username: u.Username, Role: u.Role, Disabled: u.Disabled, CreatedAt: u.CreatedAt.Format("2006-01-02")})
	}
	writeAuthJSON(w, http.StatusOK, out)
}

// handleCreate adds a member with a generated one-time password. The password is returned once.
func (h *UsersHandler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Username string `json:"username"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<10)).Decode(&body); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	body.Username = strings.TrimSpace(body.Username)
	if !usernameRe.MatchString(body.Username) {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "username: 2 to 64 letters, digits, dots, dashes or underscores"})
		return
	}
	if !validRoles[body.Role] {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "role must be admin, operator, or viewer"})
		return
	}
	pass, err := randomPassword()
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	hash, err := security.HashPassword(pass)
	if err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if err := h.users.Create(&db.DashboardUser{Username: body.Username, PasswordHash: hash, Role: body.Role}); err != nil {
		code := http.StatusInternalServerError
		if errors.Is(err, db.ErrUserExists) {
			code = http.StatusConflict
		}
		writeAuthJSON(w, code, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusCreated, map[string]string{
		"username": body.Username, "role": body.Role, "temporary_password": pass,
	})
}

// handleUpdate changes a member's role or disabled flag. The last active admin cannot be demoted or disabled.
func (h *UsersHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("username")
	var body struct {
		Role     *string `json:"role"`
		Disabled *bool   `json:"disabled"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<10)).Decode(&body); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	if body.Role != nil && !validRoles[*body.Role] {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "role must be admin, operator, or viewer"})
		return
	}
	target, err := h.users.Get(name)
	if err != nil || target == nil {
		writeAuthJSON(w, http.StatusNotFound, map[string]string{"error": "member not found"})
		return
	}
	losesAdmin := target.Role == "admin" && !target.Disabled &&
		((body.Role != nil && *body.Role != "admin") || (body.Disabled != nil && *body.Disabled))
	if losesAdmin {
		if n, err := h.activeAdmins(); err != nil || n <= 1 {
			writeAuthJSON(w, http.StatusConflict, map[string]string{"error": "keep at least one active admin"})
			return
		}
	}
	if err := h.users.Update(name, body.Role, body.Disabled, nil); err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]string{"message": "updated"})
}

func (h *UsersHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("username")
	target, err := h.users.Get(name)
	if err != nil || target == nil {
		writeAuthJSON(w, http.StatusNotFound, map[string]string{"error": "member not found"})
		return
	}
	if target.Role == "admin" && !target.Disabled {
		if n, err := h.activeAdmins(); err != nil || n <= 1 {
			writeAuthJSON(w, http.StatusConflict, map[string]string{"error": "keep at least one active admin"})
			return
		}
	}
	if err := h.users.Delete(name); err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]string{"message": "removed"})
}

// handleOwnPassword lets any member change their own password. Current password is required.
func (h *UsersHandler) handleOwnPassword(w http.ResponseWriter, r *http.Request) {
	claims := security.GetClaims(r.Context())
	name := memberName(claims)
	if name == "" {
		writeAuthJSON(w, http.StatusForbidden, map[string]string{"error": "only dashboard members can change a password here"})
		return
	}
	var body struct {
		Current string `json:"current"`
		New     string `json:"new"`
	}
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4<<10)).Decode(&body); err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid request body"})
		return
	}
	u, err := h.users.Get(name)
	if err != nil || u == nil || !security.VerifyPassword(body.Current, u.PasswordHash) {
		writeAuthJSON(w, http.StatusUnauthorized, map[string]string{"error": "current password is wrong"})
		return
	}
	hash, err := security.HashPassword(body.New)
	if err != nil {
		writeAuthJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := h.users.Update(name, nil, nil, &hash); err != nil {
		writeAuthJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeAuthJSON(w, http.StatusOK, map[string]string{"message": "password changed"})
}

func (h *UsersHandler) activeAdmins() (int, error) {
	list, err := h.users.List()
	if err != nil {
		return 0, err
	}
	n := 0
	for _, u := range list {
		if u.Role == "admin" && !u.Disabled {
			n++
		}
	}
	return n, nil
}

// memberName extracts the username from a member token ("user:<name>:<id>").
func memberName(c *security.TokenClaims) string {
	if c == nil || !strings.HasPrefix(c.TokenID, security.MemberTokenPrefix) {
		return ""
	}
	rest := strings.TrimPrefix(c.TokenID, security.MemberTokenPrefix)
	if i := strings.LastIndex(rest, ":"); i >= 0 {
		return rest[:i]
	}
	return ""
}

// randomPassword makes a 16-character password from letters and digits.
func randomPassword() (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
	out := make([]byte, 16)
	for i := range out {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}
		out[i] = alphabet[n.Int64()]
	}
	return string(out), nil
}
