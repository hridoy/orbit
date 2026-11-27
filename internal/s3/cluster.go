package s3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"orbit/internal/db"
)

func (s *S3Engine) registerClusterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/v1/storage/peers", s.handleListPeers)
	mux.HandleFunc("POST /api/v1/storage/peers", s.handleCreatePeer)
	mux.HandleFunc("DELETE /api/v1/storage/peers/{id}", s.handleDeletePeer)
	mux.HandleFunc("POST /api/v1/storage/peers/{id}/sync", s.handleSyncPeer)
}

func (s *S3Engine) handleListPeers(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	peers, err := s.store.ListPeers()
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	if peers == nil {
		peers = []*db.StoragePeer{}
	}
	// Quick live health probe with 1.5s timeout
	client := &http.Client{Timeout: 1500 * time.Millisecond}
	for _, p := range peers {
		token, _ := s.unsealTenantSecret(p.AuthTokenEnc)
		req, err := http.NewRequest(http.MethodGet, strings.TrimRight(p.Endpoint, "/")+"/s3/", nil)
		if err == nil {
			if token != "" {
				req.Header.Set("Authorization", "Bearer "+token)
			}
			resp, err := client.Do(req)
			if err == nil && resp.StatusCode < 500 {
				p.Status = "online"
				_ = resp.Body.Close()
			} else {
				p.Status = "offline"
				if resp != nil {
					_ = resp.Body.Close()
				}
			}
			_ = s.store.UpdatePeerStatus(p.PeerID, p.Status, false)
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"peers": peers})
}

func (s *S3Engine) handleCreatePeer(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	var req struct {
		Name      string `json:"name"`
		Endpoint  string `json:"endpoint"`
		AuthToken string `json:"auth_token"`
	}
	if err := json.NewDecoder(io.LimitReader(r.Body, 1<<14)).Decode(&req); err != nil {
		writeS3Error(w, http.StatusBadRequest, "InvalidRequest", err.Error())
		return
	}
	req.Name = strings.TrimSpace(req.Name)
	req.Endpoint = strings.TrimRight(strings.TrimSpace(req.Endpoint), "/")
	if req.Name == "" || req.Endpoint == "" {
		writeS3Error(w, http.StatusBadRequest, "InvalidRequest", "name and endpoint are required")
		return
	}
	idHex, err := randomHex(6)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	encTok, err := s.sealTenantSecret(req.AuthToken)
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	p := &db.StoragePeer{
		PeerID:       "peer_" + idHex,
		Name:         req.Name,
		Endpoint:     req.Endpoint,
		AuthTokenEnc: encTok,
		Status:       "unknown",
	}
	if err := s.store.CreatePeer(p); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"peer": p})
}

func (s *S3Engine) handleDeletePeer(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	if err := s.store.DeletePeer(r.PathValue("id")); err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (s *S3Engine) handleSyncPeer(w http.ResponseWriter, r *http.Request) {
	if !requireAdmin(w, r) {
		return
	}
	peerID := r.PathValue("id")
	peers, err := s.store.ListPeers()
	if err != nil {
		writeS3Error(w, http.StatusInternalServerError, "InternalError", err.Error())
		return
	}
	var target *db.StoragePeer
	for _, p := range peers {
		if p.PeerID == peerID {
			target = p
			break
		}
	}
	if target == nil {
		writeS3Error(w, http.StatusNotFound, "NoSuchPeer", "peer not found")
		return
	}
	token, _ := s.unsealTenantSecret(target.AuthTokenEnc)
	client := &http.Client{Timeout: 10 * time.Second}

	buckets, _ := s.store.ListBuckets()
	synced := 0
	failed := 0
	for _, b := range buckets {
		// Ensure bucket exists on peer
		bReq, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/s3/%s", target.Endpoint, b.BucketID), nil)
		if token != "" {
			bReq.Header.Set("Authorization", "Bearer "+token)
		}
		bReq.Header.Set("X-Orbit-Replicated", "1")
		if bResp, err := client.Do(bReq); err == nil {
			_ = bResp.Body.Close()
		}

		objs, _ := s.store.ListObjects(b.BucketID, "", 1000)
		for _, o := range objs {
			raw, err := s.readRawObjectBytes(b.BucketID, o)
			if err != nil {
				failed++
				continue
			}
			oReq, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/s3/%s/%s", target.Endpoint, b.BucketID, o.ObjectKey), bytes.NewReader(raw))
			if token != "" {
				oReq.Header.Set("Authorization", "Bearer "+token)
			}
			oReq.Header.Set("Content-Type", o.ContentType)
			oReq.Header.Set("X-Orbit-Replicated", "1")
			oResp, err := client.Do(oReq)
			if err != nil || oResp.StatusCode >= 300 {
				failed++
			} else {
				synced++
			}
			if oResp != nil {
				_ = oResp.Body.Close()
			}
		}
	}
	status := "online"
	if failed > 0 && synced == 0 {
		status = "offline"
	}
	_ = s.store.UpdatePeerStatus(target.PeerID, status, true)
	writeJSON(w, http.StatusOK, map[string]any{
		"synced_objects": synced,
		"failed_objects": failed,
		"status":         status,
	})
}

func (s *S3Engine) readRawObjectBytes(bucketName string, obj *db.StorageObject) ([]byte, error) {
	var data []byte
	var err error
	if strings.HasPrefix(obj.DiskRelPath, "erasure://") {
		data, _, err = s.readErasureObject(bucketName, obj.ObjectID)
	} else {
		data, err = os.ReadFile(filepath.Join(s.BucketDir(bucketName), obj.DiskRelPath))
	}
	if err != nil {
		return nil, err
	}
	if isSSEEncrypted(data) {
		return s.decryptSSE(data)
	}
	return data, nil
}

func (s *S3Engine) replicatePutToPeers(bucketName, objectKey, contentType string, plain []byte) {
	peers, err := s.store.ListPeers()
	if err != nil || len(peers) == 0 {
		return
	}
	client := &http.Client{Timeout: 8 * time.Second}
	for _, p := range peers {
		token, _ := s.unsealTenantSecret(p.AuthTokenEnc)
		req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/s3/%s/%s", p.Endpoint, bucketName, objectKey), bytes.NewReader(plain))
		if err != nil {
			continue
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		req.Header.Set("Content-Type", contentType)
		req.Header.Set("X-Orbit-Replicated", "1")
		resp, err := client.Do(req)
		if err == nil {
			_ = resp.Body.Close()
			_ = s.store.UpdatePeerStatus(p.PeerID, "online", true)
		} else {
			_ = s.store.UpdatePeerStatus(p.PeerID, "offline", false)
		}
	}
}

func (s *S3Engine) replicateDeleteToPeers(bucketName, objectKey string) {
	peers, err := s.store.ListPeers()
	if err != nil || len(peers) == 0 {
		return
	}
	client := &http.Client{Timeout: 5 * time.Second}
	for _, p := range peers {
		token, _ := s.unsealTenantSecret(p.AuthTokenEnc)
		req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/s3/%s/%s", p.Endpoint, bucketName, objectKey), nil)
		if err != nil {
			continue
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		req.Header.Set("X-Orbit-Replicated", "1")
		if resp, err := client.Do(req); err == nil {
			_ = resp.Body.Close()
		}
	}
}

func (s *S3Engine) fetchFromPeers(bucketName, objectKey string) ([]byte, string, bool) {
	peers, err := s.store.ListPeers()
	if err != nil || len(peers) == 0 {
		return nil, "", false
	}
	client := &http.Client{Timeout: 5 * time.Second}
	for _, p := range peers {
		token, _ := s.unsealTenantSecret(p.AuthTokenEnc)
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/s3/%s/%s", p.Endpoint, bucketName, objectKey), nil)
		if err != nil {
			continue
		}
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		if resp.StatusCode == http.StatusOK {
			data, rErr := io.ReadAll(io.LimitReader(resp.Body, MaxObjectSize))
			ct := resp.Header.Get("Content-Type")
			_ = resp.Body.Close()
			if rErr == nil {
				return data, ct, true
			}
		} else {
			_ = resp.Body.Close()
		}
	}
	return nil, "", false
}
