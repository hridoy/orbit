package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"orbit/internal/compute"
)

func RunPush(client *Client, srcDir, appID string) error {
	if appID == "" {
		return fmt.Errorf("app ID is required (pass --app=<app_id>)")
	}

	cleanDir, err := filepath.Abs(srcDir)
	if err != nil {
		return err
	}

	fmt.Printf("🚀 Orbit Push: Deploying '%s' to app '%s'...\n", cleanDir, appID)
	startTime := time.Now()

	// 1. Scan directory and compute hashes
	filesMap := make(map[string]string)
	hasher := sha256.New()

	err = filepath.Walk(cleanDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			name := info.Name()
			if name == ".git" || name == "node_modules" || name == ".orbit" {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(cleanDir, path)
		if err != nil {
			return err
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		h := sha256.New()
		if _, err := io.Copy(h, f); err != nil {
			return err
		}

		fileHash := hex.EncodeToString(h.Sum(nil))
		filesMap[relPath] = fileHash
		hasher.Write([]byte(fileHash))
		return nil
	})
	if err != nil {
		return fmt.Errorf("scan directory failed: %w", err)
	}

	versionHash := fmt.Sprintf("v_%x", hasher.Sum(nil)[:8])
	fmt.Printf("📦 Local manifest generated: %d files (Version: %s)\n", len(filesMap), versionHash)

	manifest := compute.FileManifest{
		VersionHash: versionHash,
		Files:       filesMap,
	}

	// 2. Diff with Agent
	var diff compute.DiffResult
	err = client.PostJSON(fmt.Sprintf("/api/v1/apps/%s/deploy/diff", appID), manifest, &diff)
	if err != nil {
		return fmt.Errorf("manifest diff failed: %w", err)
	}

	fmt.Printf("⚡ Delta Diff: %d files to upload (%d unchanged, reused on host)\n",
		len(diff.MissingFiles), len(filesMap)-len(diff.MissingFiles))

	// 3. Upload missing/modified delta files
	for i, relPath := range diff.MissingFiles {
		filePath := filepath.Join(cleanDir, relPath)
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("open file %s: %w", filePath, err)
		}

		uploadPath := fmt.Sprintf("/api/v1/apps/%s/deploy/upload?version=%s&path=%s",
			appID, url.QueryEscape(versionHash), url.QueryEscape(relPath))

		resp, err := client.Do("POST", uploadPath, file, "application/octet-stream")
		file.Close()
		if err != nil || resp.StatusCode >= 400 {
			return fmt.Errorf("upload %s failed: %v", relPath, err)
		}
		_ = resp.Body.Close()

		fmt.Printf("   [%d/%d] Uploaded %s\n", i+1, len(diff.MissingFiles), relPath)
	}

	// 4. Activate Zero-Downtime Deployment
	fmt.Println("🔄 Activating zero-downtime Blue/Green deployment...")
	var depResult map[string]any
	activatePayload := map[string]string{"version_hash": versionHash}
	err = client.PostJSON(fmt.Sprintf("/api/v1/apps/%s/deploy/activate", appID), activatePayload, &depResult)
	if err != nil {
		return fmt.Errorf("activation failed: %w", err)
	}

	duration := time.Since(startTime)
	fmt.Printf("✅ Deployment SUCCESSFUL in %s!\n", duration.Round(time.Millisecond))
	fmt.Printf("   App:     %s\n", appID)
	fmt.Printf("   Version: %s\n", versionHash)
	if port, ok := depResult["port_used"]; ok {
		fmt.Printf("   Port:    %v (Zero-downtime traffic cutover complete)\n", port)
	}
	return nil
}

var _ = strings.TrimSpace
