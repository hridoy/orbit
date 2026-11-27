package s3

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"orbit/internal/config"
	"orbit/internal/db"
)

func TestSSEKMSEnvelopeEncryption(t *testing.T) {
	eng := &S3Engine{
		cfg: &config.Config{
			MasterKey: bytes.Repeat([]byte{0x42}, 32),
		},
	}

	plain := []byte("Top-secret database backup payload encrypted with Orbit SSE-KMS!")
	packed, err := eng.encryptSSE(plain)
	if err != nil {
		t.Fatalf("encryptSSE: %v", err)
	}
	if bytes.Contains(packed, plain) {
		t.Fatal("plaintext leaked in SSE ciphertext")
	}
	if !isSSEEncrypted(packed) {
		t.Fatal("expected SSE magic header")
	}

	decrypted, err := eng.decryptSSE(packed)
	if err != nil {
		t.Fatalf("decryptSSE: %v", err)
	}
	if !bytes.Equal(decrypted, plain) {
		t.Fatalf("decrypted mismatch: got %q", string(decrypted))
	}

	// Tamper with ciphertext -> GCM auth must fail
	packed[len(packed)-1] ^= 0xff
	if _, err := eng.decryptSSE(packed); err == nil {
		t.Fatal("expected tamper detection failure")
	}
}

func TestJSONIAMPolicyEvaluation(t *testing.T) {
	policy := `{
		"Version": "2012-10-17",
		"Statement": [
			{"Effect": "Allow", "Action": ["s3:GetObject", "s3:PutObject"], "Resource": ["arn:aws:s3:::assets/prod/*"]},
			{"Effect": "Deny", "Action": ["s3:PutObject"], "Resource": ["arn:aws:s3:::assets/prod/locked/*"]}
		]
	}`

	if ok, err := evaluateIAMPolicyJSON(policy, http.MethodGet, "assets", "prod/logo.png"); err != nil || !ok {
		t.Fatalf("expected Allow for prod/logo.png GET, got ok=%v err=%v", ok, err)
	}
	if ok, err := evaluateIAMPolicyJSON(policy, http.MethodPut, "assets", "prod/locked/key.pem"); err != nil || ok {
		t.Fatalf("expected explicit Deny for prod/locked/* PUT, got ok=%v err=%v", ok, err)
	}
	if ok, err := evaluateIAMPolicyJSON(policy, http.MethodDelete, "assets", "prod/logo.png"); err != nil || ok {
		t.Fatalf("expected implicit Deny for DELETE, got ok=%v err=%v", ok, err)
	}
}

func TestMultiNodePeerReplicationAndFailover(t *testing.T) {
	tmp := t.TempDir()
	database, err := db.Open(filepath.Join(tmp, "orbit.db"))
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	defer database.Close()
	if err := database.Migrate(); err != nil {
		t.Fatalf("Migrate: %v", err)
	}
	store := db.NewStorageStore(database)
	eng := &S3Engine{
		cfg:   &config.Config{MasterKey: bytes.Repeat([]byte{0x24}, 32), StorageDir: tmp},
		store: store,
	}

	replicated := make(map[string][]byte)
	peerSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			buf := new(bytes.Buffer)
			_, _ = buf.ReadFrom(r.Body)
			replicated[r.URL.Path] = buf.Bytes()
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method == http.MethodGet {
			if b, ok := replicated[r.URL.Path]; ok {
				_, _ = w.Write(b)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer peerSrv.Close()

	encTok, _ := eng.sealTenantSecret("peer-secret")
	if err := store.CreatePeer(&db.StoragePeer{
		PeerID:       "peer_1",
		Name:         "node-2",
		Endpoint:     peerSrv.URL,
		AuthTokenEnc: encTok,
		Status:       "online",
	}); err != nil {
		t.Fatalf("CreatePeer: %v", err)
	}

	eng.replicatePutToPeers("assets", "backup.tar", "application/octet-stream", []byte("cluster-data"))
	if string(replicated["/s3/assets/backup.tar"]) != "cluster-data" {
		t.Fatalf("expected peer to receive replicated object, got %q", string(replicated["/s3/assets/backup.tar"]))
	}

	got, _, ok := eng.fetchFromPeers("assets", "backup.tar")
	if !ok || string(got) != "cluster-data" {
		t.Fatalf("expected failover read from peer, got ok=%v data=%q", ok, string(got))
	}
}
