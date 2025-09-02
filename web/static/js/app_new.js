// Global variables
let currentBranch = '';
let currentBuildId = '';
let currentGitConfig = '';
let currentTab = 'release-notes';
let ws = null;

// UI state
let sidebarCollapsed = false;
let bottomPanelCollapsed = false;
let sidebarWidth = 300;
let bottomPanelHeight = 200;

// Initialize the application
document.addEventListener('DOMContentLoaded', function() {
    loadGitConfigs();
    setupEventListeners();
    initializeUI();
});

// Initialize UI state
function initializeUI() {
    // Set initial tab
    switchTab('release-notes');
    
    // Initialize resize functionality
    setupResizeHandlers();
    
    // Connect WebSocket
    connectWebSocket();
}

// Load available git configurations
async function loadGitConfigs() {
    try {
        const response = await fetch('/api/git-configs');
        const data = await response.json();
        
        const select = document.getElementById('gitConfigSelect');
        select.innerHTML = '<option value="">請選擇 Git 配置...</option>';
        
        data.forEach(config => {
            const option = document.createElement('option');
            option.value = config.name;
            option.textContent = `${config.name} - ${config.description || config.url}`;
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Failed to load git configs:', error);
        addLogMessage('無法載入 Git 配置列表', 'error');
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
        addLogMessage(`正在載入 ${gitConfig} 的分支列表...`, 'info');
        
        const response = await fetch(`/api/branches/${gitConfig}`);
        const data = await response.json();
        
        const branchList = document.getElementById('branchList');
        branchList.innerHTML = '';
        
        const branches = Array.isArray(data) ? data : (data.branches || []);
        
        if (branches.length > 0) {
            branches.forEach(branch => {
                const branchName = branch.name || branch;
                const date = branch.date || '';
                const hash = branch.commit_hash || branch.commitHash || '';
                const isRelease = branch.is_release || false;
                
                const branchItem = document.createElement('div');
                branchItem.className = 'branch-item';
                branchItem.innerHTML = `
                    <div class="branch-name">${branchName}</div>
                    <div class="branch-meta">
                        ${date ? `<span><i class="fas fa-calendar"></i> ${date}</span>` : ''}
                        ${hash ? `<span><i class="fas fa-code-commit"></i> ${hash}</span>` : ''}
                        <span class="branch-tag ${isRelease ? 'release' : ''}">${isRelease ? 'Release' : 'Dev'}</span>
                    </div>
                `;
                branchItem.onclick = () => selectBranch(branchName);
                branchList.appendChild(branchItem);
            });
            
            addLogMessage(`成功載入 ${branches.length} 個分支`, 'success');
        } else {
            branchList.innerHTML = '<div class="branch-placeholder">沒有找到可用的分支</div>';
            addLogMessage('沒有找到可用的分支', 'warning');
        }
    } catch (error) {
        console.error('Failed to load branches:', error);
        document.getElementById('branchList').innerHTML = '<div class="branch-placeholder">載入分支失敗</div>';
        addLogMessage('載入分支失敗', 'error');
    }
}

// Handle branch selection
async function selectBranch(branch) {
    if (!branch || !currentGitConfig) {
        hideBranchInfo();
        return;
    }
    
    // Update UI selection
    document.querySelectorAll('.branch-item').forEach(item => {
        item.classList.remove('selected');
    });
    
    event.currentTarget.classList.add('selected');
    
    currentBranch = branch;
    document.getElementById('selectedBranch').value = branch;
    
    addLogMessage(`選擇分支: ${branch}`, 'info');
    
    await loadBranchInfo(currentGitConfig, branch);
    showBranchInfo();
}

// Load branch information
async function loadBranchInfo(gitConfig, branch) {
    try {
        addLogMessage(`正在載入分支 ${branch} 的詳細資訊...`, 'info');
        
        // Load release notes
        try {
            const notesResponse = await fetch(`/api/release-notes/${gitConfig}/${branch}`);
            if (notesResponse.ok) {
                const notesData = await notesResponse.json();
                const notesText = notesData.notes || '沒有發布說明';
                document.getElementById('releaseNotes').textContent = notesText;
                addLogMessage(`Release Notes 載入成功: ${notesText.length} 字元`, 'success');
            } else {
                const errorMsg = `載入發布說明失敗 (${notesResponse.status})`;
                document.getElementById('releaseNotes').textContent = errorMsg;
                addLogMessage(errorMsg, 'error');
            }
        } catch (error) {
            const errorMsg = '載入發布說明時發生錯誤: ' + error.message;
            document.getElementById('releaseNotes').textContent = errorMsg;
            addLogMessage(errorMsg, 'error');
        }
        
        // Load version information
        try {
            const versionsResponse = await fetch(`/api/versions/${gitConfig}/${branch}`);
            if (versionsResponse.ok) {
                const versionsData = await versionsResponse.json();
                let versionText = '';
                
                // 版本資訊
                if (versionsData.version_info) {
                    versionText += `📋 發布資訊:\n`;
                    versionText += `  Release Date: ${versionsData.version_info.release_date || 'N/A'}\n`;
                    versionText += `  Release Type: ${versionsData.version_info.release_type || 'N/A'}\n`;
                    versionText += `  Description: ${versionsData.version_info.description || 'N/A'}\n\n`;
                }
                
                // Docker 標籤
                if (versionsData.docker && versionsData.docker.tag) {
                    versionText += `🐳 Docker 資訊:\n`;
                    versionText += `  Tag: ${versionsData.docker.tag}\n\n`;
                }
                
                // 模組版本
                if (versionsData.modules) {
                    versionText += `📦 模組版本:\n`;
                    for (const [key, value] of Object.entries(versionsData.modules)) {
                        versionText += `  ${key}: ${value}\n`;
                    }
                }
                
                const finalText = versionText || '沒有版本資訊';
                document.getElementById('versionInfo').textContent = finalText;
                addLogMessage(`版本資訊載入成功`, 'success');
            } else {
                const errorMsg = `載入版本資訊失敗 (${versionsResponse.status})`;
                document.getElementById('versionInfo').textContent = errorMsg;
                addLogMessage(errorMsg, 'error');
            }
        } catch (error) {
            const errorMsg = '載入版本資訊時發生錯誤: ' + error.message;
            document.getElementById('versionInfo').textContent = errorMsg;
            addLogMessage(errorMsg, 'error');
        }
        
        // Load configuration
        try {
            const configResponse = await fetch(`/api/config/${gitConfig}/${branch}`);
            if (configResponse.ok) {
                const configData = await configResponse.json();
                document.getElementById('configInfo').textContent = JSON.stringify(configData, null, 2);
            } else {
                document.getElementById('configInfo').textContent = `載入配置失敗 (${configResponse.status})`;
            }
        } catch (error) {
            document.getElementById('configInfo').textContent = '載入配置時發生錯誤';
        }
        
        addLogMessage(`分支 ${branch} 資訊載入完成`, 'success');
        
    } catch (error) {
        console.error('Failed to load branch info:', error);
        addLogMessage('載入分支資訊失敗', 'error');
    }
}

// Show branch information sections
function showBranchInfo() {
    document.getElementById('contentPlaceholder').style.display = 'none';
    
    // Show tab content
    document.querySelectorAll('.tab-content').forEach(content => {
        if (content.id !== 'contentPlaceholder') {
            content.style.display = 'block';
        }
    });
    
    // Enable build button
    document.getElementById('startBuild').disabled = false;
}

// Hide branch information sections
function hideBranchInfo() {
    document.getElementById('contentPlaceholder').style.display = 'block';
    
    // Hide tab content but keep tabs visible
    document.querySelectorAll('.tab-content').forEach(content => {
        if (content.id !== 'contentPlaceholder') {
            content.style.display = 'none';
        }
    });
    
    // Disable build button
    document.getElementById('startBuild').disabled = true;
    
    // Clear branch selection
    document.querySelectorAll('.branch-item').forEach(item => {
        item.classList.remove('selected');
    });
}

// Hide branch selection
function hideBranchSelection() {
    hideBranchInfo();
    document.getElementById('branchList').innerHTML = '<div class="branch-placeholder">請先選擇 Git 配置</div>';
}

// Switch between tabs
function switchTab(tabName) {
    currentTab = tabName;
    
    // Update tab buttons
    document.querySelectorAll('.tab-btn').forEach(btn => {
        btn.classList.remove('active');
    });
    
    document.querySelector(`[onclick="switchTab('${tabName}')"]`).classList.add('active');
    
    // Update tab content - only show if branch is selected
    if (currentBranch) {
        document.querySelectorAll('.tab-content').forEach(content => {
            content.classList.remove('active');
        });
        
        document.getElementById(`${tabName}-content`).classList.add('active');
        document.getElementById('contentPlaceholder').style.display = 'none';
    }
}

// Setup event listeners
function setupEventListeners() {
    const startBuildBtn = document.getElementById('startBuild');
    const stopBuildBtn = document.getElementById('stopBuild');
    
    startBuildBtn.addEventListener('click', handleBuildStart);
    stopBuildBtn.addEventListener('click', handleBuildStop);
}

// Handle build start
async function handleBuildStart() {
    if (!currentBranch || !currentGitConfig) {
        addLogMessage('請先選擇 Git 配置和分支', 'error');
        return;
    }
    
    const steps = [];
    if (document.getElementById('pullRepos').checked) steps.push('pull');
    if (document.getElementById('buildImages').checked) steps.push('build');
    if (document.getElementById('pushHarbor').checked) steps.push('push');
    if (document.getElementById('deploy').checked) steps.push('deploy');
    
    if (steps.length === 0) {
        addLogMessage('請至少選擇一個構建步驟', 'error');
        return;
    }
    
    try {
        // Ensure WebSocket connection
        if (!ws || ws.readyState !== WebSocket.OPEN) {
            connectWebSocket();
            await new Promise(resolve => setTimeout(resolve, 1000));
        }
        
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
            addLogMessage(`開始構建 ${currentBranch} 分支...`, 'info');
            
            // Switch to build config tab to show progress
            switchTab('build-config');
        } else {
            addLogMessage('WebSocket 連接失敗，無法開始構建', 'error');
        }
    } catch (error) {
        console.error('Build request failed:', error);
        addLogMessage('構建請求失敗', 'error');
    }
}

// Handle build stop
function handleBuildStop() {
    if (currentBuildId) {
        addLogMessage('停止構建功能尚未實作', 'warning');
    }
    updateBuildUI(false);
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

// Toggle sidebar
function toggleSidebar() {
    sidebarCollapsed = !sidebarCollapsed;
    const layout = document.querySelector('.app-layout');
    
    if (sidebarCollapsed) {
        layout.classList.add('sidebar-collapsed');
    } else {
        layout.classList.remove('sidebar-collapsed');
    }
}

// Toggle bottom panel
function toggleBottomPanel() {
    bottomPanelCollapsed = !bottomPanelCollapsed;
    const layout = document.querySelector('.app-layout');
    const icon = document.getElementById('bottomPanelToggleIcon');
    
    if (bottomPanelCollapsed) {
        layout.classList.add('bottom-collapsed');
        icon.className = 'fas fa-chevron-up';
    } else {
        layout.classList.remove('bottom-collapsed');
        icon.className = 'fas fa-chevron-down';
    }
}

// Clear logs
function clearLogs() {
    document.getElementById('logContainer').innerHTML = '';
    addLogMessage('日誌已清空', 'info');
}

// Setup resize handlers
function setupResizeHandlers() {
    // Sidebar resize will be implemented here
    // Bottom panel resize will be implemented here
}

// Initialize sidebar resize
function initSidebarResize(e) {
    e.preventDefault();
    
    const startX = e.clientX;
    const startWidth = sidebarWidth;
    
    function doResize(e) {
        const newWidth = startWidth + (e.clientX - startX);
        if (newWidth >= 200 && newWidth <= 500) {
            sidebarWidth = newWidth;
            document.querySelector('.app-layout').style.gridTemplateColumns = `${newWidth}px 1fr`;
        }
    }
    
    function stopResize() {
        document.removeEventListener('mousemove', doResize);
        document.removeEventListener('mouseup', stopResize);
    }
    
    document.addEventListener('mousemove', doResize);
    document.addEventListener('mouseup', stopResize);
}

// Initialize bottom panel resize
function initBottomPanelResize(e) {
    e.preventDefault();
    
    const startY = e.clientY;
    const startHeight = bottomPanelHeight;
    
    function doResize(e) {
        const newHeight = startHeight - (e.clientY - startY);
        if (newHeight >= 100 && newHeight <= 400) {
            bottomPanelHeight = newHeight;
            document.querySelector('.app-layout').style.gridTemplateRows = `60px 1fr ${newHeight}px`;
        }
    }
    
    function stopResize() {
        document.removeEventListener('mousemove', doResize);
        document.removeEventListener('mouseup', stopResize);
    }
    
    document.addEventListener('mousemove', doResize);
    document.addEventListener('mouseup', stopResize);
}

// Connect to WebSocket for real-time updates
function connectWebSocket() {
    if (ws && ws.readyState === WebSocket.OPEN) {
        return;
    }
    
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;
    
    ws = new WebSocket(wsUrl);
    
    ws.onopen = function() {
        addLogMessage('WebSocket 連接已建立', 'success');
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
        
        // Try to reconnect after 3 seconds
        setTimeout(connectWebSocket, 3000);
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
    
    if (isBuilding) {
        startBtn.disabled = true;
        startBtn.innerHTML = '<i class="fas fa-spinner fa-spin"></i> 構建中...';
        stopBtn.disabled = false;
    } else {
        startBtn.disabled = !currentBranch;
        startBtn.innerHTML = '<i class="fas fa-play"></i> 開始構建';
        stopBtn.disabled = true;
    }
}

// Update progress bar
function updateProgress(progress) {
    const progressFill = document.getElementById('progressFill');
    const progressText = document.getElementById('progressText');
    
    progressFill.style.width = progress + '%';
    progressText.textContent = progress + '%';
}

// Update build status
function updateBuildStatus(statusData) {
    if (statusData.status === 'completed') {
        updateBuildUI(false);
        addLogMessage('🎉 構建完成！', 'success');
        updateProgress(100);
    } else if (statusData.status === 'failed') {
        updateBuildUI(false);
        addLogMessage('❌ 構建失敗！', 'error');
    }
}

// Add log message
function addLogMessage(message, type = 'info') {
    const logContainer = document.getElementById('logContainer');
    const timestamp = new Date().toLocaleTimeString();
    
    const logEntry = document.createElement('div');
    logEntry.className = `log-entry log-${type} fade-in`;
    logEntry.innerHTML = `<span class="log-time">[${timestamp}]</span> ${message}`;
    
    logContainer.appendChild(logEntry);
    logContainer.scrollTop = logContainer.scrollHeight;
    
    // Limit log entries to prevent memory issues
    const entries = logContainer.querySelectorAll('.log-entry');
    if (entries.length > 1000) {
        entries[0].remove();
    }
}