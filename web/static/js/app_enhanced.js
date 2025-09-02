// Enhanced Build Tool JavaScript - Based on Reference Project Patterns

// Global variables
let currentGitConfig = '';
let selectedBranch = '';
let buildInProgress = false;
let websocket = null;
let branches = [];

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
});

async function initializeApp() {
    try {
        addLog('🚀 初始化 Build Tool...', 'info');
        await loadGitConfigs();
        setupWebSocket();
        setupEventListeners();
        addLog('✅ 初始化完成', 'success');
    } catch (error) {
        addLog('❌ 初始化失败: ' + error.message, 'error');
    }
}

// Load available Git configurations
async function loadGitConfigs() {
    try {
        const response = await fetch('/api/git-configs');
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        
        const configs = await response.json();
        const select = document.getElementById('gitConfigSelect');
        
        // Clear existing options
        select.innerHTML = '<option value="">請選擇 Git 配置...</option>';
        
        // Add configurations to dropdown
        configs.forEach(config => {
            const option = document.createElement('option');
            option.value = config.key;
            option.textContent = `${config.description} (${config.key})`;
            select.appendChild(option);
        });
        
        addLog(`📂 載入了 ${configs.length} 個 Git 配置`, 'success');
    } catch (error) {
        addLog('❌ 載入 Git 配置失敗: ' + error.message, 'error');
        throw error;
    }
}

// Handle Git configuration change
async function onGitConfigChange() {
    const select = document.getElementById('gitConfigSelect');
    currentGitConfig = select.value;
    
    if (!currentGitConfig) {
        clearBranchData();
        return;
    }
    
    addLog(`🔄 切換到配置: ${currentGitConfig}`, 'info');
    
    try {
        select.classList.add('loading');
        await loadBranches(currentGitConfig);
        clearBranchSelection();
    } catch (error) {
        addLog('❌ 載入分支失敗: ' + error.message, 'error');
    } finally {
        select.classList.remove('loading');
    }
}

// Load branches for selected Git configuration
async function loadBranches(gitConfig) {
    try {
        const response = await fetch(`/api/branches/${gitConfig}`);
        if (!response.ok) throw new Error(`HTTP ${response.status}`);
        
        branches = await response.json();
        renderBranches();
        addLog(`📋 載入了 ${branches.length} 個分支`, 'success');
    } catch (error) {
        addLog('❌ 載入分支失敗: ' + error.message, 'error');
        throw error;
    }
}

// Render branch cards with enhanced styling
function renderBranches() {
    const grid = document.getElementById('branchGrid');
    
    if (!branches || branches.length === 0) {
        grid.innerHTML = '<div class="simple-placeholder">此配置下沒有可用分支</div>';
        return;
    }
    
    grid.innerHTML = branches.map(branch => `
        <div class="branch-card" data-branch="${branch}" onclick="selectBranch('${branch}')">
            <div class="branch-name">${branch}</div>
            <div class="branch-info">
                <div style="margin-top: 8px; color: #7f8c8d; font-size: 0.9rem;">
                    📅 分支類型: ${getBranchType(branch)}
                </div>
            </div>
        </div>
    `).join('');
}

// Get branch type for display
function getBranchType(branchName) {
    if (branchName === 'dev') return '開發分支';
    if (branchName === 'main' || branchName === 'master') return '主分支';
    if (/^\d{4}$/.test(branchName)) return '發布分支';
    return '其他分支';
}

// Enhanced branch selection with smooth animations
function selectBranch(branchName) {
    selectedBranch = branchName;
    
    // Update UI with animation
    document.querySelectorAll('.branch-card').forEach(card => {
        card.classList.remove('selected');
        card.style.transform = 'scale(1)';
    });
    
    const selectedCard = document.querySelector(`[data-branch="${branchName}"]`);
    if (selectedCard) {
        selectedCard.classList.add('selected');
        selectedCard.style.transform = 'scale(1.02)';
        
        // Smooth scroll to selected card if needed
        selectedCard.scrollIntoView({ behavior: 'smooth', block: 'nearest' });
    }
    
    // Update selected branch input
    document.getElementById('selectedBranch').value = branchName;
    
    // Enable build button with animation
    const buildBtn = document.getElementById('startBuild');
    buildBtn.disabled = false;
    buildBtn.classList.add('loading');
    
    // Load branch information
    loadBranchInfo(branchName).then(() => {
        buildBtn.classList.remove('loading');
    });
    
    addLog(`✅ 已選擇分支: ${branchName}`, 'success');
}

// Load branch information (release notes, versions, config)
async function loadBranchInfo(branchName) {
    if (!currentGitConfig || !branchName) return;
    
    try {
        // Show branch info section
        document.getElementById('branchInfoPlaceholder').style.display = 'none';
        document.getElementById('branchInfoData').style.display = 'block';
        
        // Show build section
        document.getElementById('buildPlaceholder').style.display = 'none';
        document.getElementById('buildPanelData').style.display = 'block';
        
        // Load all branch data in parallel
        const [releaseNotes, versions, config] = await Promise.all([
            loadReleaseNotes(branchName),
            loadVersions(branchName),
            loadConfig(branchName)
        ]);
        
        // Update UI
        updateBranchInfoDisplay(releaseNotes, versions, config);
        addLog(`📄 載入分支資訊完成: ${branchName}`, 'success');
        
    } catch (error) {
        addLog(`❌ 載入分支資訊失敗: ${error.message}`, 'error');
    }
}

// Load release notes
async function loadReleaseNotes(branchName) {
    try {
        const response = await fetch(`/api/release-notes/${currentGitConfig}/${branchName}`);
        if (!response.ok) return '暫無發布說明';
        return await response.text();
    } catch {
        return '載入失敗';
    }
}

// Load version information
async function loadVersions(branchName) {
    try {
        const response = await fetch(`/api/versions/${currentGitConfig}/${branchName}`);
        if (!response.ok) return {};
        return await response.json();
    } catch {
        return {};
    }
}

// Load configuration
async function loadConfig(branchName) {
    try {
        const response = await fetch(`/api/config/${currentGitConfig}/${branchName}`);
        if (!response.ok) return {};
        return await response.json();
    } catch {
        return {};
    }
}

// Update branch information display
function updateBranchInfoDisplay(releaseNotes, versions, config) {
    // Update release notes
    const releaseNotesEl = document.getElementById('releaseNotes');
    releaseNotesEl.textContent = releaseNotes.substring(0, 200) + (releaseNotes.length > 200 ? '...' : '');
    
    // Update version info
    const versionInfoEl = document.getElementById('versionInfo');
    const versionText = Object.keys(versions).length > 0 
        ? Object.entries(versions).slice(0, 3).map(([k, v]) => `${k}: ${v}`).join(', ')
        : '暫無版本資訊';
    versionInfoEl.textContent = versionText;
    
    // Update config info
    const configInfoEl = document.getElementById('configInfo');
    const configText = config.project?.name 
        ? `專案: ${config.project.name}`
        : '暫無配置資訊';
    configInfoEl.textContent = configText;
}

// Setup WebSocket connection
function setupWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
    websocket = new WebSocket(wsUrl);
    
    websocket.onopen = function() {
        addLog('🔌 WebSocket 連接已建立', 'success');
    };
    
    websocket.onmessage = function(event) {
        try {
            const data = JSON.parse(event.data);
            handleWebSocketMessage(data);
        } catch (error) {
            addLog('❌ WebSocket 消息解析失敗', 'error');
        }
    };
    
    websocket.onclose = function() {
        addLog('🔌 WebSocket 連接已斷開', 'warning');
        // Try to reconnect after 3 seconds
        setTimeout(setupWebSocket, 3000);
    };
    
    websocket.onerror = function(error) {
        addLog('❌ WebSocket 錯誤', 'error');
    };
}

// Handle WebSocket messages
function handleWebSocketMessage(data) {
    switch (data.type) {
        case 'log':
            addLog(data.message, data.logType || 'info');
            break;
        case 'progress':
            updateProgress(data.progress);
            break;
        case 'status':
            updateBuildStatus(data.status);
            break;
        case 'build_complete':
            handleBuildComplete(data);
            break;
        default:
            console.log('Unknown WebSocket message type:', data.type);
    }
}

// Start build process
function startBuild() {
    if (!selectedBranch || buildInProgress || !websocket || websocket.readyState !== WebSocket.OPEN) {
        if (!websocket || websocket.readyState !== WebSocket.OPEN) {
            addLog('❌ WebSocket 連接未建立，無法開始構建', 'error');
        }
        return;
    }
    
    buildInProgress = true;
    updateBuildStatus('running');
    updateProgress(0);
    
    const buildRequest = {
        gitConfig: currentGitConfig,
        branch: selectedBranch,
        pullRepos: document.getElementById('pullRepos').checked,
        buildImages: document.getElementById('buildImages').checked,
        pushHarbor: document.getElementById('pushHarbor').checked,
        deploy: document.getElementById('deploy').checked
    };
    
    addLog('🚀 開始構建流程...', 'info');
    addLog(`📋 配置: ${currentGitConfig}, 分支: ${selectedBranch}`, 'info');
    
    // Send build request via WebSocket
    websocket.send(JSON.stringify(buildRequest));
}

// Stop build process
function stopBuild() {
    if (buildInProgress && websocket && websocket.readyState === WebSocket.OPEN) {
        websocket.send(JSON.stringify({ action: 'stop' }));
        buildInProgress = false;
        updateBuildStatus('idle');
        addLog('⏹️ 構建已停止', 'warning');
    }
}

// Refresh branches
async function refreshBranches() {
    if (!currentGitConfig) {
        addLog('❌ 請先選擇 Git 配置', 'warning');
        return;
    }
    
    const btn = document.getElementById('refreshBranches');
    btn.classList.add('loading');
    btn.disabled = true;
    
    try {
        await loadBranches(currentGitConfig);
        clearBranchSelection();
        addLog('🔄 分支列表已刷新', 'success');
    } catch (error) {
        addLog('❌ 刷新分支失敗: ' + error.message, 'error');
    } finally {
        btn.classList.remove('loading');
        btn.disabled = false;
    }
}

// Update build status
function updateBuildStatus(status) {
    const indicator = document.querySelector('.status-indicator');
    const buildBtn = document.getElementById('startBuild');
    const stopBtn = document.getElementById('stopBuild');
    
    // Remove all status classes
    indicator.className = 'status-indicator';
    indicator.classList.add(`status-${status}`);
    
    switch (status) {
        case 'running':
            buildBtn.disabled = true;
            stopBtn.disabled = false;
            buildBtn.innerHTML = '<span class="status-indicator status-running"></span>構建中...';
            break;
        case 'success':
            buildInProgress = false;
            buildBtn.disabled = !selectedBranch;
            stopBtn.disabled = true;
            buildBtn.innerHTML = '<span class="status-indicator status-success"></span>開始構建';
            break;
        case 'error':
            buildInProgress = false;
            buildBtn.disabled = !selectedBranch;
            stopBtn.disabled = true;
            buildBtn.innerHTML = '<span class="status-indicator status-error"></span>開始構建';
            break;
        case 'idle':
        default:
            buildInProgress = false;
            buildBtn.disabled = !selectedBranch;
            stopBtn.disabled = true;
            buildBtn.innerHTML = '<span class="status-indicator status-idle"></span>開始構建';
            break;
    }
}

// Update progress bar
function updateProgress(percent) {
    const progressFill = document.getElementById('progressFill');
    progressFill.style.width = `${Math.min(100, Math.max(0, percent))}%`;
}

// Handle build completion
function handleBuildComplete(data) {
    buildInProgress = false;
    
    if (data.success) {
        updateBuildStatus('success');
        addLog('🎉 構建完成！', 'success');
        updateProgress(100);
    } else {
        updateBuildStatus('error');
        addLog('❌ 構建失敗', 'error');
        updateProgress(0);
    }
}

// Enhanced log message display with better formatting
function addLog(message, type = 'info') {
    const logContainer = document.getElementById('logContainer');
    const timestamp = new Date().toLocaleTimeString();
    
    let colorClass = '';
    let icon = '';
    switch(type) {
        case 'success':
            colorClass = 'color: #27ae60; font-weight: 600;';
            icon = '✅ ';
            break;
        case 'error':
            colorClass = 'color: #e74c3c; font-weight: 600;';
            icon = '❌ ';
            break;
        case 'warning':
            colorClass = 'color: #f39c12; font-weight: 600;';
            icon = '⚠️ ';
            break;
        case 'info':
            colorClass = 'color: #3498db;';
            icon = 'ℹ️ ';
            break;
        default:
            colorClass = 'color: #00ff00;';
            icon = '📝 ';
    }
    
    const logEntry = document.createElement('div');
    logEntry.style.cssText = colorClass + ' margin-bottom: 2px; padding: 2px 0; line-height: 1.4;';
    logEntry.innerHTML = `<span style="color: #666; font-size: 0.9em;">[${timestamp}]</span> ${icon}${message}`;
    
    logContainer.appendChild(logEntry);
    logContainer.scrollTop = logContainer.scrollHeight;
    
    // Auto-scroll animation
    logContainer.style.scrollBehavior = 'smooth';
    
    // Limit log entries to prevent memory issues
    const maxLogs = 1000;
    while (logContainer.children.length > maxLogs) {
        logContainer.removeChild(logContainer.firstChild);
    }
}

// Setup event listeners
function setupEventListeners() {
    // Git config change
    document.getElementById('gitConfigSelect').addEventListener('change', onGitConfigChange);
    
    // Build buttons
    document.getElementById('startBuild').addEventListener('click', startBuild);
    document.getElementById('stopBuild').addEventListener('click', stopBuild);
    document.getElementById('refreshBranches').addEventListener('click', refreshBranches);
    
    // Keyboard shortcuts
    document.addEventListener('keydown', function(event) {
        // Ctrl+Enter to start build
        if (event.ctrlKey && event.key === 'Enter') {
            event.preventDefault();
            if (!buildInProgress && selectedBranch) {
                startBuild();
            }
        }
        
        // Escape to stop build
        if (event.key === 'Escape' && buildInProgress) {
            event.preventDefault();
            stopBuild();
        }
    });
}

// Clear branch selection
function clearBranchSelection() {
    selectedBranch = '';
    document.getElementById('selectedBranch').value = '';
    document.getElementById('startBuild').disabled = true;
    
    // Hide info sections
    document.getElementById('branchInfoPlaceholder').style.display = 'block';
    document.getElementById('branchInfoData').style.display = 'none';
    document.getElementById('buildPlaceholder').style.display = 'block';
    document.getElementById('buildPanelData').style.display = 'none';
    
    // Clear selection styling
    document.querySelectorAll('.branch-card').forEach(card => {
        card.classList.remove('selected');
        card.style.transform = 'scale(1)';
    });
}

// Clear branch data
function clearBranchData() {
    branches = [];
    clearBranchSelection();
    document.getElementById('branchGrid').innerHTML = '<div class="simple-placeholder">請先選擇 Git 配置</div>';
}

// Utility function to format file sizes
function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

// Utility function to format duration
function formatDuration(seconds) {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    
    if (hours > 0) {
        return `${hours}h ${minutes}m ${secs}s`;
    } else if (minutes > 0) {
        return `${minutes}m ${secs}s`;
    } else {
        return `${secs}s`;
    }
}