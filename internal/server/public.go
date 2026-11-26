package server

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"orbit/internal/db"
	"orbit/internal/proxy"
)

// startPublicProxy serves apps on port 80. A request whose Host matches an app domain goes to
// that app. Any other request goes to the single live app when there is exactly one, which makes
// the server IP useful before DNS is set up. The dashboard (8787) is never reachable this way.
func (s *Server) startPublicProxy() {
	rp := proxy.NewReverseProxy(s.appStore, s.meter, singleLiveAppFallback(s.appStore))

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			rp.ReloadRoutes() // ports change on every deploy
		}
	}()

	// ACME HTTP-01 challenges are answered here before any app routing.
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if s.tlsEngine.ServeChallenge(w, r) {
			return
		}
		rp.ServeHTTP(w, r)
	})

	// HTTPS on 443: certificates are picked by SNI from the TLS domain engine.
	go func() {
		tlsSrv := &http.Server{
			Addr:              ":443",
			Handler:           handler,
			ReadHeaderTimeout: 10 * time.Second,
			TLSConfig: &tls.Config{
				MinVersion:     tls.VersionTLS12,
				GetCertificate: s.tlsEngine.GetCertificate,
			},
		}
		log.Printf("🔒 public HTTPS listening on :443")
		if err := tlsSrv.ListenAndServeTLS("", ""); err != nil {
			log.Printf("⚠️ public HTTPS not started: %v", err)
		}
	}()

	srv := &http.Server{
		Addr:              ":80",
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("🌐 public proxy listening on :80")
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("⚠️ public proxy not started: %v", err)
	}
}

// singleLiveAppFallback forwards to the only app that is live, and 404s otherwise.
func singleLiveAppFallback(appStore *db.AppStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apps, err := appStore.ListApps()
		if err != nil {
			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}
		var live []*db.App
		for _, a := range apps {
			if a.ActivePort > 0 && a.CurrentVersion != nil && *a.CurrentVersion != "" {
				live = append(live, a)
			}
		}
		if len(live) != 1 {
			http.Error(w, "no application is bound to this host", http.StatusNotFound)
			return
		}
		target := &url.URL{Scheme: "http", Host: fmt.Sprintf("127.0.0.1:%d", live[0].ActivePort)}
		httputil.NewSingleHostReverseProxy(target).ServeHTTP(w, r)
	})
}
