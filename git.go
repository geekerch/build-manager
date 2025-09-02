package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"gopkg.in/yaml.v2"
	"build-tool/config"
)

// =============================================================================
// Data Structures
// =============================================================================

// GitManager handles Git operations
type GitManager struct {
	currentConfig config.GitConfig
}

// Branch represents a Git branch with metadata
type Branch struct {
	Name        string `json:"name"`
	Date        string `json:"date"`
	Description string `json:"description"`
	CommitHash  string `json:"commit_hash"`
	IsRelease   bool   `json:"is_release"`
}

// BranchConfig represents configuration from config.yaml
type BranchConfig struct {
	Project      ProjectConfig     `yaml:"project"`
	Repositories RepositoryConfig  `yaml:"repositories"`
	Build        BuildSettings     `yaml:"build"`
	Deployment   DeploymentConfig  `yaml:"deployment"`
}

// ProjectConfig contains project-level settings
type ProjectConfig struct {
	Name           string `yaml:"name"`
	DockerRegistry string `yaml:"docker_registry"`
	Namespace      string `yaml:"namespace"`
}

// RepositoryConfig contains Git repository settings
type RepositoryConfig struct {
	GitlabBaseURL string   `yaml:"gitlab_base_url"`
	Modules       []Module `yaml:"modules"`
}

// Module represents a Git submodule
type Module struct {
	Name     string `yaml:"name"`
	RepoPath string `yaml:"repo_path"`
}

// BuildSettings contains build-related configuration
type BuildSettings struct {
	Platforms       []string `yaml:"platforms"`
	GenerateSwagger bool     `yaml:"generate_swagger"`
	ModulesDir      string   `yaml:"modules_dir"`
	SwaggerCommand  string   `yaml:"swagger_command"`
}

// DeploymentConfig contains deployment settings
type DeploymentConfig struct {
	Environments         []string     `yaml:"environments"`
	HealthCheckEndpoint  string       `yaml:"health_check_endpoint"`
	Docker               DockerConfig `yaml:"docker"`
}

// DockerConfig contains Docker-specific settings
type DockerConfig struct {
	BuildArgs      []string `yaml:"build_args"`
	TagFormat      string   `yaml:"tag_format"`
	RegistryFormat string   `yaml:"registry_format"`
}

// VersionInfo represents version information
type VersionInfo struct {
	VersionInfo struct {
		ReleaseDate string `json:"release_date"`
		ReleaseType string `json:"release_type"`
		Description string `json:"description"`
	} `json:"version_info"`
	Modules map[string]string `json:"modules"`
	Docker  struct {
		Tag string `json:"tag"`
	} `json:"docker"`
}

// =============================================================================
// Constructor and Configuration
// =============================================================================

// NewGitManager creates a new Git manager
func NewGitManager(defaultConfig config.GitConfig) *GitManager {
	return &GitManager{
		currentConfig: defaultConfig,
	}
}

// UpdateConfig updates the Git configuration
func (gm *GitManager) UpdateConfig(cfg config.GitConfig) {
	gm.currentConfig = cfg
}

// =============================================================================
// Git Operations
// =============================================================================

// GetAllBranches fetches all branches from Git repository
func (gm *GitManager) GetAllBranches() ([]Branch, error) {
	log.Printf("Fetching branches from repository: %s", gm.currentConfig.URL)

	// Create Git URL with token if provided
	gitURL := gm.createAuthenticatedURL()

	// Execute git ls-remote to get remote branches
	cmd := exec.Command("git", "ls-remote", "--heads", gitURL)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch remote branches: %v", err)
	}

	branches := []Branch{}
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Parse git ls-remote output: "commit_hash refs/heads/branch_name"
		parts := strings.Fields(line)
		if len(parts) != 2 {
			continue
		}

		commitHash := parts[0]
		refName := parts[1]
		
		// Extract branch name from refs/heads/branch_name
		if !strings.HasPrefix(refName, "refs/heads/") {
			continue
		}
		
		branchName := strings.TrimPrefix(refName, "refs/heads/")
		
		// Get detailed branch information
		branchInfo := gm.createBranchInfo(branchName, commitHash)
		branches = append(branches, branchInfo)
	}

	log.Printf("Found %d branches", len(branches))
	return branches, nil
}

// CloneOrPullBranch clones a branch or pulls latest if already exists
func (gm *GitManager) CloneOrPullBranch(branchName, targetDir string) error {
	if _, err := os.Stat(targetDir); err == nil {
		// Directory exists, pull latest changes
		log.Printf("Updating existing repository: %s", targetDir)
		cmd := exec.Command("git", "-C", targetDir, "pull", "origin", branchName)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to pull latest changes: %v\nOutput: %s", err, string(output))
		}
	} else {
		// Directory doesn't exist, clone the branch
		log.Printf("Cloning branch %s to %s", branchName, targetDir)
		
		// Create parent directory
		if err := os.MkdirAll(filepath.Dir(targetDir), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %v", err)
		}
		
		gitURL := gm.createAuthenticatedURL()
		cmd := exec.Command("git", "clone", "-b", branchName, "--single-branch", gitURL, targetDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to clone branch %s: %v\nOutput: %s", branchName, err, string(output))
		}
	}

	log.Printf("Successfully updated repository for branch %s", branchName)
	return nil
}

// =============================================================================
// Branch File Operations
// =============================================================================

// GetBranchConfig reads config.yaml from a specific branch
func (gm *GitManager) GetBranchConfig(branchName string) (*BranchConfig, error) {
	// First ensure we have the latest version of the branch
	targetDir := filepath.Join("repos", "temp", branchName)
	if err := gm.CloneOrPullBranch(branchName, targetDir); err != nil {
		return nil, fmt.Errorf("failed to get branch: %v", err)
	}

	// Read config.yaml
	configPath := filepath.Join(targetDir, "config.yaml")
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config.yaml: %v", err)
	}

	var config BranchConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config.yaml: %v", err)
	}

	return &config, nil
}

// GetBranchVersions reads versions.json from a specific branch
func (gm *GitManager) GetBranchVersions(branchName string) (*VersionInfo, error) {
	// First ensure we have the latest version of the branch
	targetDir := filepath.Join("repos", "temp", branchName)
	if err := gm.CloneOrPullBranch(branchName, targetDir); err != nil {
		return nil, fmt.Errorf("failed to get branch: %v", err)
	}

	// Read versions.json
	versionsPath := filepath.Join(targetDir, "versions.json")
	data, err := ioutil.ReadFile(versionsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read versions.json: %v", err)
	}

	var versions VersionInfo
	if err := json.Unmarshal(data, &versions); err != nil {
		return nil, fmt.Errorf("failed to parse versions.json: %v", err)
	}

	return &versions, nil
}

// GetBranchReleaseNotes reads release-notes.md from a specific branch
func (gm *GitManager) GetBranchReleaseNotes(branchName string) (string, error) {
	// First ensure we have the latest version of the branch
	targetDir := filepath.Join("repos", "temp", branchName)
	if err := gm.CloneOrPullBranch(branchName, targetDir); err != nil {
		return "", fmt.Errorf("failed to get branch: %v", err)
	}

	// Read release-notes.md
	notesPath := filepath.Join(targetDir, "release-notes.md")
	data, err := ioutil.ReadFile(notesPath)
	if err != nil {
		return "", fmt.Errorf("failed to read release-notes.md: %v", err)
	}

	return string(data), nil
}

// =============================================================================
// Script Execution
// =============================================================================

// ExecuteBuildScript executes a script from the cloned repository
func (gm *GitManager) ExecuteBuildScript(repoDir, scriptPath string, conn *websocket.Conn, logFunc func(*websocket.Conn, string, string)) error {
	fullScriptPath := filepath.Join(repoDir, scriptPath)
	
	// Check if script exists
	if _, err := os.Stat(fullScriptPath); os.IsNotExist(err) {
		return fmt.Errorf("script not found: %s", scriptPath)
	}

	logFunc(conn, fmt.Sprintf("🔧 執行腳本: %s", scriptPath), "info")

	// Make script executable
	if err := os.Chmod(fullScriptPath, 0755); err != nil {
		return fmt.Errorf("failed to make script executable: %v", err)
	}

	// Set up environment variables
	env := os.Environ()
	if gm.currentConfig.Token != "" {
		env = append(env, "GITLAB_TOKEN="+gm.currentConfig.Token)
		env = append(env, "GIT_TOKEN="+gm.currentConfig.Token)
	}

	// Execute script
	cmd := exec.Command("bash", fullScriptPath)
	cmd.Dir = repoDir
	cmd.Env = env

	// Create pipes for real-time output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %v", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start script: %v", err)
	}

	// Read output in goroutines
	go gm.readOutput(stdout, conn, logFunc, "info")
	go gm.readOutput(stderr, conn, logFunc, "error")

	// Wait for command to complete
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("script execution failed: %v", err)
	}

	logFunc(conn, fmt.Sprintf("✅ 腳本執行完成: %s", scriptPath), "success")
	return nil
}

// =============================================================================
// Helper Functions
// =============================================================================

// createAuthenticatedURL creates a Git URL with authentication
func (gm *GitManager) createAuthenticatedURL() string {
	if gm.currentConfig.Token == "" {
		return gm.currentConfig.URL
	}

	// Handle different Git URL formats
	url := gm.currentConfig.URL
	if strings.HasPrefix(url, "https://") {
		// For HTTPS URLs, inject token
		// https://github.com/user/repo.git -> https://token@github.com/user/repo.git
		url = strings.Replace(url, "https://", fmt.Sprintf("https://%s@", gm.currentConfig.Token), 1)
	}
	
	return url
}

// createBranchInfo creates branch information
func (gm *GitManager) createBranchInfo(branchName, commitHash string) Branch {
	return Branch{
		Name:        branchName,
		Date:        time.Now().Format("2006-01-02"), // Would be better to get actual commit date
		Description: gm.generateDescription(branchName),
		CommitHash:  commitHash[:8], // Short hash
		IsRelease:   gm.isReleaseBranch(branchName),
	}
}

// generateDescription generates a description for the branch
func (gm *GitManager) generateDescription(branchName string) string {
	if gm.isReleaseBranch(branchName) {
		return fmt.Sprintf("Release branch %s", branchName)
	}
	return fmt.Sprintf("Development branch %s", branchName)
}

// isReleaseBranch checks if a branch is a release branch
func (gm *GitManager) isReleaseBranch(branchName string) bool {
	// Match patterns like: release/0804, 0901, 0902, etc.
	releasePatterns := []string{
		"release/", "rel/", "v",
	}

	for _, pattern := range releasePatterns {
		if strings.HasPrefix(branchName, pattern) {
			return true
		}
	}

	// Check if it's a 4-digit pattern like 0901, 0902
	if len(branchName) == 4 {
		for _, char := range branchName {
			if char < '0' || char > '9' {
				return false
			}
		}
		return true
	}

	return false
}

// readOutput reads command output and sends to WebSocket
func (gm *GitManager) readOutput(pipe interface{}, conn *websocket.Conn, logFunc func(*websocket.Conn, string, string), msgType string) {
	// Implementation depends on the pipe type, for now just a placeholder
	// In real implementation, you'd read from the pipe and send each line via logFunc
}