package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"io/fs"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"build-tool/config"
)

//go:embed web/static/*
var staticFiles embed.FS

//go:embed web/templates/* web/static/*
var templateFiles embed.FS

// BuildManager handles the build operations
type BuildManager struct {
	config     *config.Config
	gitManager *GitManager
	upgrader   websocket.Upgrader
}

// NewBuildManager creates a new build manager instance
func NewBuildManager(cfg *config.Config) *BuildManager {
	// Get first git config as default
	var defaultGitConfig config.GitConfig
	for _, gitConfig := range cfg.GitConfigs {
		defaultGitConfig = gitConfig
		break
	}
	
	return &BuildManager{
		config:     cfg,
		gitManager: NewGitManager(defaultGitConfig),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all origins for development
			},
		},
	}
}

func main() {
	// Load configuration
	cfg := config.LoadConfig("config.json")

	// Initialize build manager
	bm := NewBuildManager(cfg)

	// Setup routes
	router := bm.setupRoutes()

	// Create necessary directories
	createDirectories()

	// Print startup info
	printStartupInfo(cfg.Server.Port)

	// Start server
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, router))
}

// setupRoutes configures the HTTP routes
func (bm *BuildManager) setupRoutes() *mux.Router {
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/git-configs", bm.GetGitConfigs).Methods("GET")
	r.HandleFunc("/api/branches/{gitConfig}", bm.GetBranches).Methods("GET")
	r.HandleFunc("/api/config/{gitConfig}/{branch}", bm.GetConfig).Methods("GET")
	r.HandleFunc("/api/versions/{gitConfig}/{branch}", bm.GetVersions).Methods("GET")
	r.HandleFunc("/api/release-notes/{gitConfig}/{branch}", bm.GetReleaseNotes).Methods("GET")
	r.HandleFunc("/ws", bm.HandleWebSocket)

	// Serve static files from embedded FS
	staticSub, _ := fs.Sub(templateFiles, "web/static")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.FS(staticSub))))

	// UI route
	r.HandleFunc("/", bm.ServeUI).Methods("GET")

	return r
}

// createDirectories creates necessary directories
func createDirectories() {
	dirs := []string{"repos", "build-temp"}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Failed to create directory %s: %v", dir, err)
		}
	}
}

// printStartupInfo prints server startup information
func printStartupInfo(port string) {
	fmt.Printf("🚀 Build Tool 啟動中...\n")
	fmt.Printf("📱 Web UI: http://localhost:%s\n", port)
	fmt.Printf("🔌 WebSocket: ws://localhost:%s/ws\n", port)
	fmt.Printf("📁 Repos 目錄: %s\n", "repos")
	fmt.Printf("📁 Build 暫存目錄: %s\n", "build-temp")
}