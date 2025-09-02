package main

// uiHTML contains the embedded HTML UI - this will be read from the template file
// For now, we'll use a simple constant but in production this would be embedded
const uiHTML = `<!DOCTYPE html>
<html lang="zh-TW">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Build Tool</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1);
            overflow: hidden;
            display: flex;
            flex-direction: column;
            height: 100vh;
        }

        .header {
            background: linear-gradient(135deg, #2c3e50, #34495e);
            color: white;
            padding: 12px 30px;
            text-align: center;
            flex-shrink: 0;
        }

        .header h1 {
            font-size: 1.6em;
            margin-bottom: 2px;
            font-weight: 300;
        }

        .header p {
            font-size: 0.9em;
            margin: 0;
        }

        .main-layout {
            display: flex;
            flex: 1;
            overflow: hidden;
        }

        .sidebar {
            width: 350px;
            background: #f8f9fa;
            border-right: 2px solid #e9ecef;
            display: flex;
            flex-direction: column;
            overflow: hidden;
        }

        .main-content {
            flex: 1;
            padding: 15px;
            overflow-y: auto;
            display: flex;
            flex-direction: column;
        }

        .top-row {
            display: flex;
            gap: 15px;
            flex: 1;
            min-height: 0;
        }

        .info-panel {
            flex: 2;
            display: flex;
            flex-direction: column;
            min-height: 0;
        }

        .build-panel {
            flex: 1;
            background: #f8f9fa;
            border-radius: 10px;
            padding: 15px;
            border: 1px solid #e9ecef;
        }

        .bottom-row {
            margin-top: 15px;
            height: 300px;
        }

        .section {
            margin-bottom: 30px;
        }

        .section h2 {
            color: #2c3e50;
            margin-bottom: 15px;
            font-size: 1.4em;
            border-bottom: 2px solid #3498db;
            padding-bottom: 8px;
        }

        .sidebar .section {
            margin-bottom: 20px;
            padding: 15px;
        }

        .sidebar .section h2 {
            font-size: 1.2em;
            margin-bottom: 10px;
        }

        .branch-list {
            display: flex;
            flex-direction: column;
            gap: 10px;
            max-height: 400px;
            overflow-y: auto;
        }

        .branch-item {
            background: white;
            border: 2px solid #e9ecef;
            border-radius: 10px;
            padding: 15px;
            transition: all 0.3s ease;
            cursor: pointer;
        }

        .branch-item:hover {
            border-color: #3498db;
            transform: translateX(5px);
            box-shadow: 0 5px 15px rgba(52, 152, 219, 0.15);
        }

        .branch-item.selected {
            border-color: #27ae60;
            background: #d4edda;
        }

        .branch-name {
            font-size: 1.1em;
            font-weight: bold;
            color: #2c3e50;
            margin-bottom: 5px;
        }

        .branch-info {
            color: #6c757d;
            font-size: 0.8em;
            line-height: 1.4;
        }

        .git-config-section {
            border-bottom: 1px solid #dee2e6;
            padding-bottom: 15px;
            margin-bottom: 15px;
        }

        .info-tabs {
            display: flex;
            border-bottom: 2px solid #e9ecef;
            margin-bottom: 20px;
        }

        .tab-button {
            padding: 10px 20px;
            border: none;
            background: none;
            color: #6c757d;
            cursor: pointer;
            border-bottom: 2px solid transparent;
            transition: all 0.3s ease;
        }

        .tab-button.active {
            color: #3498db;
            border-bottom-color: #3498db;
        }

        .tab-content {
            display: none;
        }

        .tab-content.active {
            display: block;
        }

        .info-card {
            border: 1px solid #e0e0e0;
            border-radius: 10px;
            padding: 15px;
            background: #f8f9fa;
            margin-bottom: 20px;
        }

        .content-area {
            background: #1e1e1e;
            color: #00ff00;
            font-family: 'Courier New', monospace;
            padding: 12px;
            border-radius: 5px;
            flex: 1;
            white-space: pre-wrap;
            overflow-y: auto;
            margin-top: 8px;
            font-size: 0.9em;
        }

        .build-panel h2 {
            margin-top: 0;
            margin-bottom: 15px;
            font-size: 1.3em;
        }

        .form-group {
            margin-bottom: 15px;
        }

        .form-group label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: #2c3e50;
        }

        select, input {
            width: 100%;
            padding: 12px 16px;
            border: 2px solid #e9ecef;
            border-radius: 10px;
            font-size: 16px;
            transition: border-color 0.3s ease;
        }

        select:focus, input:focus {
            outline: none;
            border-color: #3498db;
        }

        .checkbox-group {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-top: 15px;
        }

        .checkbox-item {
            display: flex;
            align-items: center;
            gap: 8px;
        }

        .checkbox-item input[type="checkbox"] {
            width: auto;
            margin: 0;
        }

        .btn {
            padding: 15px 30px;
            border: none;
            border-radius: 10px;
            font-size: 16px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            margin-right: 15px;
            margin-bottom: 10px;
        }

        .btn-primary {
            background: linear-gradient(135deg, #3498db, #2980b9);
            color: white;
        }

        .btn-success {
            background: linear-gradient(135deg, #27ae60, #219a52);
            color: white;
        }

        .btn-warning {
            background: linear-gradient(135deg, #f39c12, #e67e22);
            color: white;
        }

        .btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 8px 20px rgba(52, 152, 219, 0.3);
        }

        .btn:disabled {
            opacity: 0.6;
            cursor: not-allowed;
            transform: none !important;
        }

        .status-indicator {
            display: inline-block;
            width: 12px;
            height: 12px;
            border-radius: 50%;
            margin-right: 8px;
        }

        .status-idle { background: #95a5a6; }
        .status-running { 
            background: #f39c12; 
            animation: pulse 1s infinite;
        }
        .status-success { background: #27ae60; }
        .status-error { background: #e74c3c; }

        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }

        .progress-bar {
            width: 100%;
            height: 8px;
            background: #e9ecef;
            border-radius: 4px;
            overflow: hidden;
            margin: 15px 0;
        }

        .progress-fill {
            height: 100%;
            background: linear-gradient(90deg, #3498db, #2980b9);
            width: 0%;
            transition: width 0.3s ease;
        }

        .log-container {
            background: #1e1e1e;
            color: #00ff00;
            font-family: 'Courier New', monospace;
            padding: 12px;
            border-radius: 8px;
            height: 240px;
            overflow-y: auto;
            border: 2px solid #333;
            font-size: 0.85em;
            line-height: 1.3;
        }

        @media (max-width: 768px) {
            .main-layout {
                flex-direction: column;
            }
            
            .sidebar {
                width: 100%;
                max-height: 300px;
            }
            
            .checkbox-group {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚀 Build Tool</h1>
            <p>選擇 Git 配置和分支，執行構建任務</p>
        </div>

        <div class="main-layout">
            <!-- 左側邊欄 -->
            <div class="sidebar">
                <!-- Git Configuration Selection -->
                <div class="section git-config-section">
                    <h2>📁 Git 配置</h2>
                    <select id="gitConfigSelect" onchange="onGitConfigChange()">
                        <option value="">請選擇 Git 配置...</option>
                    </select>
                </div>

                <!-- Branch Selection -->
                <div class="section">
                    <h2>📋 可用分支</h2>
                    <div class="branch-list" id="branchList">
                        <!-- 動態生成分支項目 -->
                    </div>
                    <div id="branchEmptyPlaceholder" style="color:#6c757d; font-size: 0.9em; padding:8px 2px;">請先選擇 Git 配置以載入分支</div>
                </div>
            </div>

            <!-- 主要內容區域 -->
            <div class="main-content">
                <!-- Top Row: Branch Information and Build Configuration -->
                <div class="top-row" id="topRow" style="display: none;">
                    <!-- Branch Information -->
                    <div class="info-panel" id="branchInfoSection" style="display: none;">
                        <h2>📄 分支資訊</h2>
                        
                        <!-- 標籤頁 -->
                        <div class="info-tabs">
                            <button class="tab-button active" onclick="switchTab('releaseNotes')">📝 Release Notes</button>
                            <button class="tab-button" onclick="switchTab('versionInfo')">🏷️ 版本資訊</button>
                            <button class="tab-button" onclick="switchTab('configInfo')">⚙️ 配置預覽</button>
                        </div>

                        <!-- 標籤內容 -->
                        <div id="releaseNotes" class="tab-content active">
                            <div class="content-area">請先選擇 Git 配置和分支以查看發布說明</div>
                        </div>
                        
                        <div id="versionInfo" class="tab-content">
                            <div class="content-area">請先選擇 Git 配置和分支以查看版本資訊</div>
                        </div>
                        
                        <div id="configInfo" class="tab-content">
                            <div class="content-area">請先選擇 Git 配置和分支以查看配置資訊</div>
                        </div>
                    </div>

                    <!-- Build Configuration -->
                    <div class="build-panel" id="buildSection" style="display: none;">
                        <h2>⚙️ 構建配置</h2>
                        <div class="form-group">
                            <label for="selectedBranch">選擇的分支:</label>
                            <input type="text" id="selectedBranch" readonly placeholder="請先選擇一個分支">
                        </div>

                        <div class="form-group">
                            <label>構建選項:</label>
                            <div class="checkbox-group">
                                <div class="checkbox-item">
                                    <input type="checkbox" id="pullRepos" checked>
                                    <label for="pullRepos">拉取配置倉庫</label>
                                </div>
                                <div class="checkbox-item">
                                    <input type="checkbox" id="buildImages" checked>
                                    <label for="buildImages">執行構建腳本</label>
                                </div>
                                <div class="checkbox-item">
                                    <input type="checkbox" id="pushHarbor" checked>
                                    <label for="pushHarbor">推送到 Harbor</label>
                                </div>
                                <div class="checkbox-item">
                                    <input type="checkbox" id="deploy">
                                    <label for="deploy">執行部署</label>
                                </div>
                            </div>
                        </div>

                        <div class="form-group">
                            <button class="btn btn-success" id="startBuild" disabled>
                                <span class="status-indicator status-idle"></span>
                                開始構建
                            </button>
                            <button class="btn btn-warning" id="stopBuild" disabled>停止構建</button>
                            <button class="btn btn-primary" id="refreshBranches">刷新分支</button>
                        </div>

                        <div class="progress-bar">
                            <div class="progress-fill" id="progressFill"></div>
                        </div>
                    </div>
                </div>

                <!-- Bottom Row: Build Logs -->
                <div class="bottom-row">
                    <div class="section">
                        <h2>📝 構建日誌</h2>
                        <div class="log-container" id="logContainer">
                            <div style="color: #00ff00;">等待選擇配置和分支...</div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <script>
        let gitConfigs = [];
        let branches = [];
        let selectedGitConfig = null;
        let selectedBranch = null;
        let buildRunning = false;
        let ws = null;

        // 初始化頁面
        async function init() {
            await loadGitConfigs();
            setupWebSocket();
            setupEventListeners();
        }

        // 載入 Git 配置
        async function loadGitConfigs() {
            try {
                const response = await fetch('/api/git-configs');
                gitConfigs = await response.json();
                renderGitConfigs();
                addLog('✅ Git 配置載入完成');
            } catch (error) {
                addLog('❌ 載入 Git 配置失敗: ' + error.message);
            }
        }

        // 渲染 Git 配置選項
        function renderGitConfigs() {
            const select = document.getElementById('gitConfigSelect');
            select.innerHTML = '<option value="">請選擇 Git 配置...</option>';
            
            gitConfigs.forEach(config => {
                const option = document.createElement('option');
                option.value = config.name;
                option.textContent = config.description + ' (' + config.name + ')';
                select.appendChild(option);
            });
        }

        // Git 配置改變時
        async function onGitConfigChange() {
            const select = document.getElementById('gitConfigSelect');
            const configName = select.value;
            
            if (!configName) {
                hideBranchInfo();
                return;
            }
            
            selectedGitConfig = configName;
            await loadBranches(configName);
            
            // 顯示上方資訊與構建區塊容器（但內部內容仍按需切換顯示）
            document.getElementById('topRow').style.display = 'flex';
            hideBranchInfo();
        }

        // 載入分支資料
        async function loadBranches(gitConfig) {
            try {
                const response = await fetch('/api/branches/' + gitConfig);
                branches = await response.json();
                renderBranches();
                addLog('✅ 分支資料載入完成: ' + branches.length + ' 個分支');
            } catch (error) {
                addLog('❌ 載入分支資料失敗: ' + error.message);
            }
        }

        // 渲染分支列表
        function renderBranches() {
            const list = document.getElementById('branchList');
            const placeholder = document.getElementById('branchEmptyPlaceholder');
            
            if (!branches || branches.length === 0) {
                list.innerHTML = '';
                placeholder.style.display = 'block';
                return;
            }
            
            placeholder.style.display = 'none';
            list.innerHTML = branches.map(branch => 
                '<div class="branch-item" data-branch="' + branch.name + '">' +
                    '<div class="branch-name">' + branch.name + '</div>' +
                    '<div class="branch-info">' +
                        '<div>📅 ' + branch.date + '</div>' +
                        '<div>🔗 ' + branch.commit_hash + '</div>' +
                    '</div>' +
                '</div>'
            ).join('');
        }

        // 建立 WebSocket 連接
        function setupWebSocket() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = protocol + '//' + window.location.host + '/ws';
            
            ws = new WebSocket(wsUrl);
            
            ws.onopen = function() {
                addLog('🔌 WebSocket 連接已建立');
            };
            
            ws.onmessage = function(event) {
                const data = JSON.parse(event.data);
                
                if (data.type === 'log') {
                    const logData = data.data;
                    addLog(logData.message, logData.type);
                } else if (data.type === 'progress') {
                    updateProgress(data.data.progress);
                }
            };
            
            ws.onclose = function() {
                addLog('❌ WebSocket 連接已斷開');
            };
            
            ws.onerror = function(error) {
                addLog('❌ WebSocket 錯誤: ' + error);
            };
        }

        // 設置事件監聽器
        function setupEventListeners() {
            // 分支選擇
            document.getElementById('branchList').addEventListener('click', function(e) {
                const item = e.target.closest('.branch-item');
                if (item) {
                    selectBranch(item.dataset.branch);
                }
            });

            // 構建按鈕
            document.getElementById('startBuild').addEventListener('click', startBuild);
            document.getElementById('stopBuild').addEventListener('click', stopBuild);
            document.getElementById('refreshBranches').addEventListener('click', refreshBranches);
        }

        // 選擇分支
        async function selectBranch(branchName) {
            selectedBranch = branchName;

            // 更新UI
            document.querySelectorAll('.branch-item').forEach(function(item) {
                item.classList.remove('selected');
            });
            document.querySelector('[data-branch="' + branchName + '"]').classList.add('selected');

            document.getElementById('selectedBranch').value = branchName;
            document.getElementById('startBuild').disabled = false;

            // 載入分支資訊
            await loadBranchInfo(selectedGitConfig, branchName);
            showBranchInfo();

            addLog('已選擇分支: ' + branchName);
        }

        // 標籤切換功能
        function switchTab(tabName) {
            // 移除所有 active 狀態
            document.querySelectorAll('.tab-button').forEach(button => {
                button.classList.remove('active');
            });
            document.querySelectorAll('.tab-content').forEach(content => {
                content.classList.remove('active');
            });
            
            // 設置當前 active 狀態
            event.target.classList.add('active');
            document.getElementById(tabName).classList.add('active');
        }

        // 載入分支資訊
        async function loadBranchInfo(gitConfig, branch) {
            try {
                addLog('📄 載入分支資訊: ' + branch);
                
                // 載入 Release Notes
                const notesResponse = await fetch('/api/release-notes/' + gitConfig + '/' + branch);
                if (notesResponse.ok) {
                    const notesData = await notesResponse.json();
                    document.querySelector('#releaseNotes .content-area').textContent = notesData.notes || '無發布說明';
                    addLog('✅ Release Notes 載入成功');
                } else {
                    document.querySelector('#releaseNotes .content-area').textContent = '載入發布說明失敗: ' + notesResponse.status;
                    addLog('❌ Release Notes 載入失敗: ' + notesResponse.status);
                }
                
                // 載入版本資訊
                const versionsResponse = await fetch('/api/versions/' + gitConfig + '/' + branch);
                if (versionsResponse.ok) {
                    const versionsData = await versionsResponse.json();
                    
                    let versionText = '';
                    if (versionsData.modules) {
                        for (const [key, value] of Object.entries(versionsData.modules)) {
                            versionText += key + ': ' + value + '\n';
                        }
                    }
                    if (versionsData.docker && versionsData.docker.tag) {
                        versionText += 'DOCKER_TAG: ' + versionsData.docker.tag + '\n';
                    }
                    if (versionsData.version_info) {
                        versionText += '\n發布資訊:\n';
                        versionText += 'Release Date: ' + (versionsData.version_info.release_date || 'N/A') + '\n';
                        versionText += 'Release Type: ' + (versionsData.version_info.release_type || 'N/A') + '\n';
                        versionText += 'Description: ' + (versionsData.version_info.description || 'N/A') + '\n';
                    }
                    
                    document.querySelector('#versionInfo .content-area').textContent = versionText || '無版本資訊';
                    addLog('✅ 版本資訊載入成功');
                } else {
                    document.querySelector('#versionInfo .content-area').textContent = '載入版本資訊失敗: ' + versionsResponse.status;
                    addLog('❌ 版本資訊載入失敗: ' + versionsResponse.status);
                }
                
                // 載入配置資訊
                const configResponse = await fetch('/api/config/' + gitConfig + '/' + branch);
                if (configResponse.ok) {
                    const configData = await configResponse.json();
                    document.querySelector('#configInfo .content-area').textContent = JSON.stringify(configData, null, 2);
                    addLog('✅ 配置資訊載入成功');
                } else {
                    document.querySelector('#configInfo .content-area').textContent = '載入配置資訊失敗: ' + configResponse.status;
                    addLog('❌ 配置資訊載入失敗: ' + configResponse.status);
                }
                
            } catch (error) {
                addLog('⚠️ 載入分支資訊時發生錯誤: ' + error.message);
                document.querySelector('#releaseNotes .content-area').textContent = '載入發布說明失敗: ' + error.message;
                document.querySelector('#versionInfo .content-area').textContent = '載入版本資訊失敗: ' + error.message;
                document.querySelector('#configInfo .content-area').textContent = '載入配置資訊失敗: ' + error.message;
            }
        }

        // 顯示分支資訊
        function showBranchInfo() {
            // 顯示資訊面板與構建面板
            document.getElementById('branchInfoSection').style.display = 'flex';
            document.getElementById('buildSection').style.display = 'block';
            
            // 滾動至日誌區域仍保持在視窗內，不被擠到下方
            const bottomRow = document.querySelector('.bottom-row');
            if (bottomRow) {
                bottomRow.style.height = '240px';
            }
        }

        // 隱藏分支資訊
        function hideBranchInfo() {
            document.getElementById('branchInfoSection').style.display = 'none';
            document.getElementById('buildSection').style.display = 'none';
            selectedBranch = null;
            document.getElementById('startBuild').disabled = true;
            
            // 確保底部日誌區域高度固定，不會因為內容上方隱藏而拉伸
            const bottomRow = document.querySelector('.bottom-row');
            if (bottomRow) {
                bottomRow.style.height = '240px';
            }
            
            // 清空內容
            document.querySelector('#releaseNotes .content-area').textContent = '請先選擇 Git 配置和分支以查看發布說明';
            document.querySelector('#versionInfo .content-area').textContent = '請先選擇 Git 配置和分支以查看版本資訊';
            document.querySelector('#configInfo .content-area').textContent = '請先選擇 Git 配置和分支以查看配置資訊';
        }

        // 開始構建
        function startBuild() {
            if (!selectedGitConfig || !selectedBranch || buildRunning || !ws || ws.readyState !== WebSocket.OPEN) {
                if (!ws || ws.readyState !== WebSocket.OPEN) {
                    addLog('❌ WebSocket 連接未建立，無法開始構建');
                }
                return;
            }

            buildRunning = true;
            updateBuildStatus('running');
            updateProgress(0);

            const buildRequest = {
                gitConfig: selectedGitConfig,
                branch: selectedBranch,
                pullRepos: document.getElementById('pullRepos').checked,
                buildImages: document.getElementById('buildImages').checked,
                pushHarbor: document.getElementById('pushHarbor').checked,
                deploy: document.getElementById('deploy').checked
            };

            // 發送構建請求到後端
            ws.send(JSON.stringify(buildRequest));
        }

        // 停止構建
        function stopBuild() {
            if (buildRunning) {
                buildRunning = false;
                updateBuildStatus('idle');
                addLog('⏹️ 構建已停止');
            }
        }

        // 刷新分支
        function refreshBranches() {
            if (selectedGitConfig) {
                loadBranches(selectedGitConfig);
            }
        }

        // 更新構建狀態
        function updateBuildStatus(status) {
            const indicator = document.querySelector('.status-indicator');
            const button = document.getElementById('startBuild');
            const stopBtn = document.getElementById('stopBuild');

            indicator.className = 'status-indicator status-' + status;

            switch (status) {
                case 'running':
                    button.disabled = true;
                    stopBtn.disabled = false;
                    break;
                case 'idle':
                case 'success':
                case 'error':
                    button.disabled = !selectedBranch;
                    stopBtn.disabled = true;
                    break;
            }
        }

        // 更新進度條
        function updateProgress(percent) {
            document.getElementById('progressFill').style.width = percent + '%';
        }

        // 添加日誌
        function addLog(message, type) {
            type = type || 'info';
            const container = document.getElementById('logContainer');
            const timestamp = new Date().toLocaleTimeString();
            
            let color = '#00ff00';
            if (type === 'success') color = '#27ae60';
            if (type === 'error') color = '#e74c3c';
            if (type === 'warning') color = '#f39c12';
            
            const entry = document.createElement('div');
            entry.style.color = color;
            entry.textContent = '[' + timestamp + '] ' + message;
            container.appendChild(entry);
            
            // 保持日誌容器固定高度，不會推擠其他區塊
            container.style.overflowY = 'auto';
            container.scrollTop = container.scrollHeight;
            
            // 構建結束狀態處理
            if (message.includes('🎉 構建完成') || message.includes('✅ 構建完成')) {
                buildRunning = false;
                updateBuildStatus('success');
            } else if (message.includes('❌') && buildRunning) {
                buildRunning = false;
                updateBuildStatus('error');
            }
        }

        // 頁面載入時初始化
        document.addEventListener('DOMContentLoaded', init);
    </script>
</body>
</html>`