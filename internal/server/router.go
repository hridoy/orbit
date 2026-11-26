package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"orbit/internal/buildinfo"
	"orbit/internal/compute"
	"orbit/internal/config"
	"orbit/internal/database"
	"orbit/internal/db"
	"orbit/internal/jail"
	"orbit/internal/proxy"
	"orbit/internal/s3"
	"orbit/internal/security"
	"orbit/internal/telemetry"
	"orbit/internal/update"
	"orbit/ui"
)

type Server struct {
	cfg          *config.Config
	db           *db.DB
	nodeStore    *db.NodeStore
	appStore     *db.AppStore
	dbStore      *db.DatabaseStore
	secStore     *db.SecretStore
	storage      *db.StorageStore
	tokenMgr     *security.TokenManager
	auth         *security.AuthMiddleware
	jail         *jail.Jail
	jailHdr      *jail.Handler
	computeHdr   *compute.ComputeHandler
	gitHdr       *compute.GitHandler
	platformHdr  *update.Handler
	databaseHdr  *database.DatabaseHandler
	s3Engine     *s3.S3Engine
	meter        *proxy.BandwidthMeter
	telemetryHdr *telemetry.TelemetryBroadcaster
	mux          *http.ServeMux
	httpServer   *http.Server
}

func NewServer(
	cfg *config.Config,
	dbConn *db.DB,
	vault *security.Vault,
	tokenMgr *security.TokenManager,
) *Server {
	nodeStore := db.NewNodeStore(dbConn)
	appStore := db.NewAppStore(dbConn)
	dbStore := db.NewDatabaseStore(dbConn)
	secStore := db.NewSecretStore(dbConn, vault)
	storageStore := db.NewStorageStore(dbConn)

	jailInst := jail.NewJail(cfg)
	jailHandler := jail.NewHandler(jailInst)
	authMiddleware := security.NewAuthMiddleware(tokenMgr)

	cgroups := compute.NewCgroupController()
	supervisor := compute.NewProcessSupervisor(cfg, cgroups)
	deployer := compute.NewDeployCoordinator(cfg, appStore, secStore, supervisor)
	computeHandler := compute.NewComputeHandler(appStore, secStore, cgroups, supervisor, deployer)
	gitHandler := compute.NewGitHandler(cfg, appStore, db.NewAppSourceStore(dbConn), supervisor, deployer)

	dbManager := database.NewDatabaseManager(cfg, dbStore, secStore, vault)
	databaseHandler := database.NewDatabaseHandler(dbStore, dbManager, vault)

	s3Engine := s3.NewS3Engine(cfg, storageStore, tokenMgr)

	meter := proxy.NewBandwidthMeter()
	collector := telemetry.NewTelemetryCollector(meter)
	telemetryHandler := telemetry.NewTelemetryBroadcaster(collector)

	mux := http.NewServeMux()

	s := &Server{
		cfg:          cfg,
		db:           dbConn,
		nodeStore:    nodeStore,
		appStore:     appStore,
		dbStore:      dbStore,
		secStore:     secStore,
		storage:      storageStore,
		tokenMgr:     tokenMgr,
		auth:         authMiddleware,
		jail:         jailInst,
		jailHdr:      jailHandler,
		computeHdr:   computeHandler,
		gitHdr:       gitHandler,
		platformHdr:  update.NewHandler(cfg.DataDir, buildinfo.Version),
		databaseHdr:  databaseHandler,
		s3Engine:     s3Engine,
		meter:        meter,
		telemetryHdr: telemetryHandler,
		mux:          mux,
	}

	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	// Public Healthz
	s.mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    "healthy",
			"node_id":   s.cfg.NodeID,
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		})
	})

	// Protected Node & Version Info
	s.mux.HandleFunc("GET /api/v1/node", func(w http.ResponseWriter, r *http.Request) {
		node, err := s.nodeStore.GetNode(s.cfg.NodeID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(node)
	})

	// Register Jailed File Browser routes
	s.jailHdr.RegisterRoutes(s.mux)

	// Register Compute & Deploy routes
	s.computeHdr.RegisterRoutes(s.mux)

	// Register GitHub install and update routes
	s.gitHdr.RegisterRoutes(s.mux)

	// Register agent self-update routes
	s.platformHdr.RegisterRoutes(s.mux)

	// Register Dedicated Native Database routes
	s.databaseHdr.RegisterRoutes(s.mux)

	// Register Sovereign S3 Engine routes
	s.s3Engine.RegisterRoutes(s.mux)

	// Register Telemetry & Live Metrics routes
	s.telemetryHdr.RegisterRoutes(s.mux)

	// Mount embedded Vue 3 SPA UI for root and web client navigation
	s.mux.Handle("/", ui.Handler())
}

// Handler returns the secured root HTTP handler wrapped with recovery, headers, and auth.
func (s *Server) Handler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"error": fmt.Sprintf("internal server panic: %v", rec),
				})
			}
		}()

		// Security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		// Auth wrapper
		s.auth.Wrap(s.mux).ServeHTTP(w, r)
	})
}

// ListenAndServe starts the agent HTTP server.
func (s *Server) ListenAndServe() error {
	// Bring back apps that were live before the last restart, without blocking the listener.
	go s.computeHdr.RestoreApps()

	// Serve apps publicly on port 80 (HTTP; HTTPS comes later).
	go s.startPublicProxy()

	s.httpServer = &http.Server{
		Addr:              s.cfg.ListenAddr,
		Handler:           s.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

// Close gracefully stops the server.
func (s *Server) Close() error {
	if s.httpServer != nil {
		return s.httpServer.Close()
	}
	return nil
}
