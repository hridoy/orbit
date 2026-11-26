package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"orbit/internal/db"
	"orbit/internal/security"
)

func newUsersTest(t *testing.T) (*db.DashboardUserStore, *UsersHandler) {
	t.Helper()
	conn, err := db.Open(filepath.Join(t.TempDir(), "orbit.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })
	if err := conn.Migrate(); err != nil {
		t.Fatal(err)
	}
	store := db.NewDashboardUserStore(conn)
	hash, _ := security.HashPassword("admin-pass-123")
	if err := store.Create(&db.DashboardUser{Username: "boss", PasswordHash: hash, Role: "admin"}); err != nil {
		t.Fatal(err)
	}
	return store, NewUsersHandler(store)
}

func TestInviteReturnsOneTimePasswordThatWorks(t *testing.T) {
	store, h := newUsersTest(t)
	raw, _ := json.Marshal(map[string]string{"username": "dev1", "role": "operator"})
	rec := httptest.NewRecorder()
	h.handleCreate(rec, httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(raw)))
	if rec.Code != http.StatusCreated {
		t.Fatalf("invite: %d %s", rec.Code, rec.Body)
	}
	var out map[string]string
	_ = json.Unmarshal(rec.Body.Bytes(), &out)
	u, _ := store.Get("dev1")
	if u == nil || !security.VerifyPassword(out["temporary_password"], u.PasswordHash) {
		t.Fatal("temporary password does not match the stored hash")
	}
}

func TestInviteRejectsBadRoleAndDuplicate(t *testing.T) {
	_, h := newUsersTest(t)
	raw, _ := json.Marshal(map[string]string{"username": "x1", "role": "root"})
	rec := httptest.NewRecorder()
	h.handleCreate(rec, httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(raw)))
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("bad role: %d", rec.Code)
	}
	raw, _ = json.Marshal(map[string]string{"username": "boss", "role": "viewer"})
	rec = httptest.NewRecorder()
	h.handleCreate(rec, httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewReader(raw)))
	if rec.Code != http.StatusConflict {
		t.Fatalf("duplicate: %d", rec.Code)
	}
}

func TestLastAdminCannotBeDemotedOrRemoved(t *testing.T) {
	_, h := newUsersTest(t)
	viewer := "viewer"
	raw, _ := json.Marshal(map[string]any{"role": viewer})
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/boss", bytes.NewReader(raw))
	req.SetPathValue("username", "boss")
	h.handleUpdate(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("demoting the last admin: %d", rec.Code)
	}
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodDelete, "/api/v1/users/boss", nil)
	req.SetPathValue("username", "boss")
	h.handleDelete(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("removing the last admin: %d", rec.Code)
	}
}

func TestMemberLookupReportsDisabled(t *testing.T) {
	store, _ := newUsersTest(t)
	hash, _ := security.HashPassword("viewer-pass-1")
	_ = store.Create(&db.DashboardUser{Username: "ann", PasswordHash: hash, Role: "viewer"})
	disabled := true
	_ = store.Update("ann", nil, &disabled, nil)
	role, active, found := memberLookup{users: store}.Lookup("ann")
	if !found || active || role != "viewer" {
		t.Fatalf("got role=%s active=%v found=%v", role, active, found)
	}
}
