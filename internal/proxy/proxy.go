package proxy

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"

	"orbit/internal/db"
)

type ReverseProxy struct {
	mu       sync.RWMutex
	routes   map[string]int // domain -> active_port
	appIDs   map[string]string // domain -> appID
	meter    *BandwidthMeter
	appStore *db.AppStore
	fallback http.Handler
}

func NewReverseProxy(appStore *db.AppStore, meter *BandwidthMeter, fallback http.Handler) *ReverseProxy {
	rp := &ReverseProxy{
		routes:   make(map[string]int),
		appIDs:   make(map[string]string),
		meter:    meter,
		appStore: appStore,
		fallback: fallback,
	}
	rp.ReloadRoutes()
	return rp
}

// ReloadRoutes synchronizes the domain mapping from the SQLite database.
func (p *ReverseProxy) ReloadRoutes() {
	p.mu.Lock()
	defer p.mu.Unlock()

	apps, err := p.appStore.ListApps()
	if err != nil {
		return
	}

	p.routes = make(map[string]int)
	p.appIDs = make(map[string]string)

	for _, app := range apps {
		if app.Domain != nil && *app.Domain != "" {
			domain := strings.ToLower(strings.TrimSpace(*app.Domain))
			p.routes[domain] = app.ActivePort
			p.appIDs[domain] = app.AppID
		}
	}
}

// SetRoute directly maps a domain to a port and app ID.
func (p *ReverseProxy) SetRoute(domain, appID string, port int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	cleanDomain := strings.ToLower(strings.TrimSpace(domain))
	p.routes[cleanDomain] = port
	p.appIDs[cleanDomain] = appID
}

// ServeHTTP inspects r.Host and routes to the target app or fallback.
func (p *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	host := r.Host
	if colonIdx := strings.Index(host, ":"); colonIdx != -1 {
		host = host[:colonIdx]
	}
	host = strings.ToLower(strings.TrimSpace(host))

	p.mu.RLock()
	port, hasRoute := p.routes[host]
	appID := p.appIDs[host]
	p.mu.RUnlock()

	// If no domain route, check if fallback handler exists
	if !hasRoute || port <= 0 {
		if p.fallback != nil {
			p.fallback.ServeHTTP(w, r)
			return
		}
		http.Error(w, fmt.Sprintf("No application registered for host %s", host), http.StatusNotFound)
		return
	}

	// Metering wrappers
	p.meter.RecordRequestStart(appID)
	defer p.meter.RecordRequestEnd(appID)

	countingW := NewCountingResponseWriter(w, appID, p.meter)
	if r.Body != nil {
		r.Body = NewCountingReader(r.Body, appID, p.meter)
	}

	targetURL, _ := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", port))

	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.Transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// Set standard reverse proxy headers
	origDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		origDirector(req)
		req.Header.Set("X-Forwarded-Host", r.Host)
		if clientIP, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
			if prior := req.Header.Get("X-Forwarded-For"); prior != "" {
				clientIP = prior + ", " + clientIP
			}
			req.Header.Set("X-Forwarded-For", clientIP)
			req.Header.Set("X-Real-IP", clientIP)
		}
		if r.TLS != nil {
			req.Header.Set("X-Forwarded-Proto", "https")
		} else {
			req.Header.Set("X-Forwarded-Proto", "http")
		}
	}

	proxy.ErrorHandler = func(rw http.ResponseWriter, req *http.Request, err error) {
		countingW.WriteHeader(http.StatusBadGateway)
		_, _ = countingW.Write([]byte(fmt.Sprintf("Orbit Reverse Proxy: Bad Gateway (Target port %d unreachable: %v)", port, err)))
	}

	proxy.ServeHTTP(countingW, r)
}

var _ = context.Background
