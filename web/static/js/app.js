// Global variables
let currentBranch = '';
let currentBuildId = '';
let currentGitConfig = '';

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    loadGitConfigs();
    setupEventListeners();
});

// Load available git configurations
async function loadGitConfigs() {
    try {
        const response = await fetch('/api/git-configs');
        const data = await response.json();
        
        const select = document.getElementById('gitConfigSelect');
        select.innerHTML = '<option value="">請選擇 Git 配置...</option>';
        
        // data is an array of git config objects
        data.forEach(config => {
            const option = document.createElement('option');
            option.value = config.name;
            option.textContent = `${config.name} - ${config.description || config.url}`;
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Failed to load git configs:', error);
        alert('無法載入 Git 配置列表');
    }
}

// Handle git config selection change
async function onGitConfigChange() {
    const select = document.getElementById('gitConfigSelect');
    const gitConfig = select.value;
    
    if (!gitConfig) {
        hideBranchSelection();
        return;
    }
    
    currentGitConfig = gitConfig;
    await loadBranches(gitConfig);
    
    // Reset branch info panels for new selection
    hideBranchInfo();
}

// Load available branches for selected git config
async function loadBranches(gitConfig) {
    try {
        const response = await fetch(`/api/branches/${gitConfig}`);
        const data = await response.json();
        
        const branchGrid = document.getElementById('branchGrid');
        branchGrid.innerHTML = '';
        
        const branches = Array.isArray(data) ? data : (data.branches || []);
        
        if (branches.length > 0) {
            branches.forEach(b => {
                const branchName = b.name || b;
                const date = b.date || '';
                const hash = b.commit_hash || b.commitHash || '';
                
                const card = document.createElement('div');
                card.className = 'branch-card';
                card.innerHTML = `
                    <h3>${branchName}</h3>
                    <p>${date ? `📅 ${date} ` : ''}${hash ? `🔗 ${hash}` : ''}</p>
                `;
                card.onclick = () => selectBranch(branchName);
                branchGrid.appendChild(card);
            });
        } else {
            branchGrid.innerHTML = '<p>沒有找到可用的分支</p>';
        }
    } catch (error) {
        console.error('Failed to load branches:', error);
        document.getElementById('branchGrid').innerHTML = '<p>載入分支失敗</p>';
    }
}

// Handle branch selection
async function selectBranch(branch) {
    if (!branch || !currentGitConfig) {
        hideBranchInfo();
        return;
    }
    
    currentBranch = branch;
    document.getElementById('selectedBranch').value = branch;
    await loadBranchInfo(currentGitConfig, branch);
    showBranchInfo();
}

// Load branch information
async function loadBranchInfo(gitConfig, branch) {
    try {
        // Load release notes
        const notesResponse = await fetch(`/api/release-notes/${gitConfig}/${branch}`);
        if (notesResponse.ok) {
            const notesData = await notesResponse.json();
            document.getElementById('releaseNotes').textContent = notesData.notes || '沒有發布說明';
        } else {
            document.getElementById('releaseNotes').textContent = `載入發布說明失敗 (${notesResponse.status})`;
        }
        
        // Load version information
        const versionsResponse = await fetch(`/api/versions/${gitConfig}/${branch}`);
        if (versionsResponse.ok) {
            const versionsData = await versionsResponse.json();
            let versionText = '';
            if (versionsData.modules) {
                for (const [key, value] of Object.entries(versionsData.modules)) {
                    versionText += `${key}: ${value}\n`;
                }
            }
            if (versionsData.docker && versionsData.docker.tag) {
                versionText += `DOCKER_TAG: ${versionsData.docker.tag}\n`;
            }
            if (versionsData.version_info) {
                versionText += `\n發布資訊:\n`;
                versionText += `Release Date: ${versionsData.version_info.release_date || 'N/A'}\n`;
                versionText += `Release Type: ${versionsData.version_info.release_type || 'N/A'}\n`;
                versionText += `Description: ${versionsData.version_info.description || 'N/A'}\n`;
            }
            document.getElementById('versionInfo').textContent = versionText || '沒有版本資訊';
        } else {
            document.getElementById('versionInfo').textContent = `載入版本資訊失敗 (${versionsResponse.status})`;
        }
        
        // Load configuration
        const configResponse = await fetch(`/api/config/${gitConfig}/${branch}`);
        if (configResponse.ok) {
            const configData = await configResponse.json();
            document.getElementById('configInfo').textContent = JSON.stringify(configData, null, 2);
        } else {
            document.getElementById('configInfo').textContent = `載入配置失敗 (${configResponse.status})`;
        }
        
    } catch (error) {
        console.error('Failed to load branch info:', error);
        document.getElementById('releaseNotes').textContent = '載入發布說明失敗';
        document.getElementById('versionInfo').textContent = '載入版本資訊失敗';
        document.getElementById('configInfo').textContent = '載入配置失敗';
    }
}

// Show branch information sections
function showBranchInfo() {
    // 隱藏佔位符，顯示實際內容
    const branchInfoPlaceholder = document.getElementById('branchInfoPlaceholder');
    const branchInfoData = document.getElementById('branchInfoData');
    
    const buildPlaceholder = document.getElementById('buildPlaceholder');
    const buildPanelData = document.getElementById('buildPanelData');
    
    if (branchInfoPlaceholder) branchInfoPlaceholder.style.display = 'none';
    if (branchInfoData) branchInfoData.style.display = 'block';
    
    if (buildPlaceholder) buildPlaceholder.style.display = 'none';
    if (buildPanelData) buildPanelData.style.display = 'block';
    
    document.getElementById('startBuild').disabled = false;
}

// Hide branch information sections
function hideBranchInfo() {
    // 顯示佔位符，隱藏實際內容
    const branchInfoPlaceholder = document.getElementById('branchInfoPlaceholder');
    const branchInfoData = document.getElementById('branchInfoData');
    
    const buildPlaceholder = document.getElementById('buildPlaceholder');
    const buildPanelData = document.getElementById('buildPanelData');
    
    if (branchInfoPlaceholder) branchInfoPlaceholder.style.display = 'block';
    if (branchInfoData) branchInfoData.style.display = 'none';
    
    if (buildPlaceholder) buildPlaceholder.style.display = 'block';
    if (buildPanelData) buildPanelData.style.display = 'none';
    
    document.getElementById('startBuild').disabled = true;
}

// Hide branch selection
function hideBranchSelection() {
    hideBranchInfo();
    document.getElementById('branchGrid').innerHTML = '';
}

// Setup event listeners
function setupEventListeners() {
    const startBuildBtn = document.getElementById('startBuild');
    const stopBuildBtn = document.getElementById('stopBuild');
    const refreshBtn = document.getElementById('refreshBranches');
    
    startBuildBtn.addEventListener('click', handleBuildStart);
    stopBuildBtn.addEventListener('click', handleBuildStop);
    refreshBtn.addEventListener('click', handleRefresh);
}

// Handle build start
async function handleBuildStart() {
    if (!currentBranch || !currentGitConfig) {
        alert('請先選擇 Git 配置和分支');
        return;
    }
    
    const steps = [];
    if (document.getElementById('pullRepos').checked) steps.push('pull');
    if (document.getElementById('buildImages').checked) steps.push('build');
    if (document.getElementById('pushHarbor').checked) steps.push('push');
    if (document.getElementById('deploy').checked) steps.push('deploy');
    
    if (steps.length === 0) {
        alert('請至少選擇一個構建步驟');
        return;
    }
    
    try {
        // Start WebSocket connection for real-time updates
        connectWebSocket();
        
        // Wait for WebSocket connection to be established
        if (ws && ws.readyState === WebSocket.OPEN) {
            // Send build request via WebSocket
            const buildRequest = {
                gitConfig: currentGitConfig,
                branch: currentBranch,
                pullRepos: steps.includes('pull'),
                buildImages: steps.includes('build'),
                pushHarbor: steps.includes('push'),
                deploy: steps.includes('deploy')
            };
            
            ws.send(JSON.stringify(buildRequest));
            updateBuildUI(true);
            addLogMessage('構建請求已發送...', 'info');
        } else {
            // Fallback: wait a bit for WebSocket connection
            setTimeout(() => {
                if (ws && ws.readyState === WebSocket.OPEN) {
                    const buildRequest = {
                        gitConfig: currentGitConfig,
                        branch: currentBranch,
                        pullRepos: steps.includes('pull'),
                        buildImages: steps.includes('build'),
                        pushHarbor: steps.includes('push'),
                        deploy: steps.includes('deploy')
                    };
                    
                    ws.send(JSON.stringify(buildRequest));
                    updateBuildUI(true);
                    addLogMessage('構建請求已發送...', 'info');
                } else {
                    alert('WebSocket 連接失敗，無法開始構建');
                }
            }, 1000);
        }
    } catch (error) {
        console.error('Build request failed:', error);
        alert('構建請求失敗');
    }
}

// Handle build stop
function handleBuildStop() {
    if (currentBuildId) {
        // TODO: Implement build stop API
        addLogMessage('停止構建功能尚未實作', 'warning');
    }
}

// Handle refresh
function handleRefresh() {
    if (currentGitConfig) {
        loadBranches(currentGitConfig);
        addLogMessage('已刷新分支列表', 'info');
    } else {
        loadGitConfigs();
        addLogMessage('已刷新 Git 配置列表', 'info');
    }
}

// WebSocket connection
let ws = null;

// Connect to WebSocket for real-time updates
function connectWebSocket() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        return;
    }
    
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
    ws = new WebSocket(wsUrl);
    
    ws.onopen = function() {
        addLogMessage('WebSocket 連接已建立', 'info');
    };
    
    ws.onmessage = function(event) {
        try {
            const data = JSON.parse(event.data);
            handleWebSocketMessage(data);
        } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
        }
    };
    
    ws.onclose = function() {
        addLogMessage('WebSocket 連接已關閉', 'warning');
        ws = null;
    };
    
    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
        addLogMessage('WebSocket 連接錯誤', 'error');
    };
}

// Handle WebSocket messages
function handleWebSocketMessage(data) {
    if (data.type === 'log') {
        addLogMessage(data.data.message, data.data.type);
    } else if (data.type === 'progress') {
        updateProgress(data.data.progress);
    } else if (data.type === 'status') {
        updateBuildStatus(data.data);
    }
}

// Update build UI state
function updateBuildUI(isBuilding) {
    const startBtn = document.getElementById('startBuild');
    const stopBtn = document.getElementById('stopBuild');
    const statusIndicator = startBtn.querySelector('.status-indicator');
    
    if (isBuilding) {
        startBtn.disabled = true;
        stopBtn.disabled = false;
        if (statusIndicator) statusIndicator.className = 'status-indicator status-running';
    } else {
        // Keep disabled if no branch selected
        startBtn.disabled = !currentBranch;
        stopBtn.disabled = true;
        if (statusIndicator) statusIndicator.className = 'status-indicator status-idle';
    }
}

// Update progress bar
function updateProgress(progress) {
    const progressFill = document.getElementById('progressFill');
    progressFill.style.width = progress + '%';
}

// Update build status
function updateBuildStatus(statusData) {
    if (statusData.status === 'completed') {
        updateBuildUI(false);
        addLogMessage('構建完成！', 'success');
        updateProgress(100);
    } else if (statusData.status === 'failed') {
        updateBuildUI(false);
        addLogMessage('構建失敗！', 'error');
    }
}

// Add log message
function addLogMessage(message, type = 'log') {
    const logContainer = document.getElementById('logContainer');
    const timestamp = new Date().toLocaleTimeString();
    
    const logEntry = document.createElement('div');
    logEntry.className = `log-entry log-${type}`;
    logEntry.innerHTML = `<span class="log-time">[${timestamp}]</span> ${message}`;
    
    logContainer.appendChild(logEntry);
    logContainer.scrollTop = logContainer.scrollHeight;
}