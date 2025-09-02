package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"
	"io"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// =============================================================================
// Data Structures
// =============================================================================

// BuildRequest represents a build request from the frontend
type BuildRequest struct {
	GitConfig     string `json:"gitConfig"`
	Branch        string `json:"branch"`
	PullRepos     bool   `json:"pullRepos"`
	BuildImages   bool   `json:"buildImages"`
	PushHarbor    bool   `json:"pushHarbor"`
	Deploy        bool   `json:"deploy"`
}

// LogMessage represents a log message sent via WebSocket
type LogMessage struct {
	Timestamp string `json:"timestamp"`
	Message   string `json:"message"`
	Type      string `json:"type"` // info, success, error, warning
}

// =============================================================================
// HTTP Handlers
// =============================================================================

// GetGitConfigs returns all available Git configurations
func (bm *BuildManager) GetGitConfigs(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	gitConfigs := make([]map[string]string, 0)
	for name, config := range bm.config.GitConfigs {
		gitConfigs = append(gitConfigs, map[string]string{
			"name": name,
			"url":  config.URL,
			"description": config.Description,
		})
	}
	
	if err := json.NewEncoder(w).Encode(gitConfigs); err != nil {
		log.Printf("Error encoding git configs: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetBranches returns all available branches from specified Git repository
func (bm *BuildManager) GetBranches(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	gitConfigName := vars["gitConfig"]
	
	gitConfig, exists := bm.config.GitConfigs[gitConfigName]
	if !exists {
		http.Error(w, "Git configuration not found", http.StatusNotFound)
		return
	}
	
	// Update git manager with selected config
	bm.gitManager.UpdateConfig(gitConfig)
	
	// Fetch branches from Git repository
	branches, err := bm.gitManager.GetAllBranches()
	if err != nil {
		log.Printf("Error fetching branches from Git: %v", err)
		http.Error(w, "Failed to fetch branches from Git repository", http.StatusInternalServerError)
		return
	}
	
	if err := json.NewEncoder(w).Encode(branches); err != nil {
		log.Printf("Error encoding branches: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetConfig returns configuration for specified branch
func (bm *BuildManager) GetConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	gitConfigName := vars["gitConfig"]
	branchName := vars["branch"]
	
	gitConfig, exists := bm.config.GitConfigs[gitConfigName]
	if !exists {
		http.Error(w, "Git configuration not found", http.StatusNotFound)
		return
	}
	
	bm.gitManager.UpdateConfig(gitConfig)
	
	config, err := bm.gitManager.GetBranchConfig(branchName)
	if err != nil {
		log.Printf("Error fetching config for branch %s: %v", branchName, err)
		http.Error(w, "Failed to fetch branch configuration", http.StatusInternalServerError)
		return
	}
	
	if err := json.NewEncoder(w).Encode(config); err != nil {
		log.Printf("Error encoding config: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetVersions returns version information for specified branch
func (bm *BuildManager) GetVersions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	gitConfigName := vars["gitConfig"]
	branchName := vars["branch"]
	
	gitConfig, exists := bm.config.GitConfigs[gitConfigName]
	if !exists {
		http.Error(w, "Git configuration not found", http.StatusNotFound)
		return
	}
	
	bm.gitManager.UpdateConfig(gitConfig)
	
	versions, err := bm.gitManager.GetBranchVersions(branchName)
	if err != nil {
		log.Printf("Error fetching versions for branch %s: %v", branchName, err)
		http.Error(w, "Failed to fetch branch versions", http.StatusInternalServerError)
		return
	}
	
	if err := json.NewEncoder(w).Encode(versions); err != nil {
		log.Printf("Error encoding versions: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// GetReleaseNotes returns release notes for specified branch
func (bm *BuildManager) GetReleaseNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	
	vars := mux.Vars(r)
	gitConfigName := vars["gitConfig"]
	branchName := vars["branch"]
	
	gitConfig, exists := bm.config.GitConfigs[gitConfigName]
	if !exists {
		http.Error(w, "Git configuration not found", http.StatusNotFound)
		return
	}
	
	bm.gitManager.UpdateConfig(gitConfig)
	
	notes, err := bm.gitManager.GetBranchReleaseNotes(branchName)
	if err != nil {
		log.Printf("Error fetching release notes for branch %s: %v", branchName, err)
		http.Error(w, "Failed to fetch release notes", http.StatusInternalServerError)
		return
	}
	
	response := map[string]string{
		"branch": branchName,
		"notes":  notes,
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding release notes: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// ServeUI serves the UI template from embedded files
func (bm *BuildManager) ServeUI(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")

    f, err := templateFiles.Open("web/templates/index.html")
    if err != nil {
        log.Printf("Error opening template: %v", err)
        http.Error(w, "UI template not found", http.StatusInternalServerError)
        return
    }
    defer f.Close()

    data, err := io.ReadAll(f)
    if err != nil {
        log.Printf("Error reading template: %v", err)
        http.Error(w, "Failed to read UI template", http.StatusInternalServerError)
        return
    }

    if _, err := w.Write(data); err != nil {
        log.Printf("Error writing UI response: %v", err)
    }
}

// =============================================================================
// WebSocket Handler
// =============================================================================

// HandleWebSocket handles WebSocket connections for real-time logs
func (bm *BuildManager) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := bm.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Println("WebSocket client connected")

	// Send initial connection message
	bm.sendLogMessage(conn, "WebSocket 連接已建立", "info")

	// Keep connection alive and handle incoming messages
	for {
		var buildReq BuildRequest
		err := conn.ReadJSON(&buildReq)
		if err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Handle build request
		go bm.handleBuildRequest(conn, buildReq)
	}
}

// =============================================================================
// Build Process Handler
// =============================================================================

// handleBuildRequest processes a build request and sends real-time updates
func (bm *BuildManager) handleBuildRequest(conn *websocket.Conn, req BuildRequest) {
	bm.sendLogMessage(conn, fmt.Sprintf("🚀 開始構建分支 %s (Git: %s)", req.Branch, req.GitConfig), "info")

	// Update git manager with selected config
	gitConfig := bm.config.GitConfigs[req.GitConfig]
	bm.gitManager.UpdateConfig(gitConfig)

	progress := 0
	stepSize := 100 / bm.countSteps(req)

	// Execute build steps
	if req.PullRepos {
		if !bm.executePullRepos(conn, &progress, stepSize, req.GitConfig, req.Branch) {
			return
		}
	}

	if req.BuildImages {
		if !bm.executeBuildImages(conn, &progress, stepSize, req.GitConfig, req.Branch) {
			return
		}
	}

	if req.PushHarbor {
		if !bm.executePushHarbor(conn, &progress, stepSize, req.GitConfig, req.Branch) {
			return
		}
	}

	if req.Deploy {
		if !bm.executeDeploy(conn, &progress, req.GitConfig, req.Branch) {
			return
		}
	}

	bm.sendLogMessage(conn, "🎉 構建完成！", "success")
}

// =============================================================================
// Build Step Implementations
// =============================================================================

// executePullRepos executes the pull repositories step
func (bm *BuildManager) executePullRepos(conn *websocket.Conn, progress *int, stepSize int, gitConfig, branchName string) bool {
	bm.sendLogMessage(conn, "▶️ 拉取配置倉庫...", "info")
	bm.sendProgress(conn, *progress)

	// Create target directory
	targetDir := filepath.Join("repos", gitConfig, branchName)
	
	// Clone or pull the branch
	if err := bm.gitManager.CloneOrPullBranch(branchName, targetDir); err != nil {
		bm.sendLogMessage(conn, fmt.Sprintf("❌ 拉取失敗: %v", err), "error")
		return false
	}

	bm.sendLogMessage(conn, "✅ 拉取配置倉庫完成", "success")
	*progress += stepSize
	bm.sendProgress(conn, *progress)
	return true
}

// executeBuildImages executes the build images step
func (bm *BuildManager) executeBuildImages(conn *websocket.Conn, progress *int, stepSize int, gitConfig, branchName string) bool {
	bm.sendLogMessage(conn, "▶️ 執行構建腳本...", "info")
	bm.sendProgress(conn, *progress)

	// Execute build script from the cloned repository
	targetDir := filepath.Join("repos", gitConfig, branchName)
	if err := bm.gitManager.ExecuteBuildScript(targetDir, "scripts/build.sh", conn, bm.sendLogMessage); err != nil {
		bm.sendLogMessage(conn, fmt.Sprintf("❌ 構建失敗: %v", err), "error")
		return false
	}

	bm.sendLogMessage(conn, "✅ 構建腳本執行完成", "success")
	*progress += stepSize
	bm.sendProgress(conn, *progress)
	return true
}

// executePushHarbor executes the push to Harbor step
func (bm *BuildManager) executePushHarbor(conn *websocket.Conn, progress *int, stepSize int, gitConfig, branchName string) bool {
	bm.sendLogMessage(conn, "▶️ 推送到 Harbor...", "info")
	bm.sendProgress(conn, *progress)

	// Execute push script from the cloned repository (if exists)
	targetDir := filepath.Join("repos", gitConfig, branchName)
	if err := bm.gitManager.ExecuteBuildScript(targetDir, "scripts/push.sh", conn, bm.sendLogMessage); err != nil {
		bm.sendLogMessage(conn, fmt.Sprintf("⚠️ 推送腳本執行警告: %v", err), "warning")
		// Continue even if push script fails or doesn't exist
	}

	bm.sendLogMessage(conn, "✅ 推送步驟完成", "success")
	*progress += stepSize
	bm.sendProgress(conn, *progress)
	return true
}

// executeDeploy executes the deployment step
func (bm *BuildManager) executeDeploy(conn *websocket.Conn, progress *int, gitConfig, branchName string) bool {
	bm.sendLogMessage(conn, "▶️ 執行部署...", "info")
	bm.sendProgress(conn, *progress)

	// Execute deploy script from the cloned repository (if exists)
	targetDir := filepath.Join("repos", gitConfig, branchName)
	if err := bm.gitManager.ExecuteBuildScript(targetDir, "scripts/deploy.sh", conn, bm.sendLogMessage); err != nil {
		bm.sendLogMessage(conn, fmt.Sprintf("⚠️ 部署腳本執行警告: %v", err), "warning")
		// Continue even if deploy script fails or doesn't exist
	}

	bm.sendLogMessage(conn, "✅ 部署步驟完成", "success")
	*progress = 100
	bm.sendProgress(conn, *progress)
	return true
}

// =============================================================================
// WebSocket Communication Helpers
// =============================================================================

// sendLogMessage sends a log message via WebSocket
func (bm *BuildManager) sendLogMessage(conn *websocket.Conn, message, msgType string) {
	logMsg := LogMessage{
		Timestamp: time.Now().Format("15:04:05"),
		Message:   message,
		Type:      msgType,
	}

	if err := conn.WriteJSON(map[string]interface{}{
		"type": "log",
		"data": logMsg,
	}); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}
}

// sendProgress sends progress update via WebSocket
func (bm *BuildManager) sendProgress(conn *websocket.Conn, progress int) {
	if err := conn.WriteJSON(map[string]interface{}{
		"type": "progress",
		"data": map[string]int{"progress": progress},
	}); err != nil {
		log.Printf("WebSocket write error: %v", err)
	}
}

// =============================================================================
// Utility Functions
// =============================================================================

// countSteps counts the number of enabled build steps
func (bm *BuildManager) countSteps(req BuildRequest) int {
	count := 0
	if req.PullRepos {
		count++
	}
	if req.BuildImages {
		count++
	}
	if req.PushHarbor {
		count++
	}
	if req.Deploy {
		count++
	}
	if count == 0 {
		return 1 // Avoid division by zero
	}
	return count
}