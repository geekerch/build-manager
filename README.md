# Build Tool

這是一個用於管理 CI/CD 構建流程的 Web 工具，支援從 Git 配置倉庫讀取配置並執行構建。

## 專案結構

```
build-tool/                    # Go 程式主資料夾
├── main.go                   # 主程式入口
├── go.mod                    # Go modules 設定
├── config/                   # 配置管理模組
│   └── config.go            # 配置結構定義和讀取函數
└── web/                     # 嵌入式 Web UI
    ├── templates/           # HTML 模板
    │   └── index.html
    └── static/              # 靜態資源
        ├── css/
        │   └── style.css
        └── js/
            └── app.js

build-config-repo/            # 配置 Git Repository (獨立)
├── README.md                # 配置倉庫說明
├── config.yaml              # 構建配置
├── versions.json            # 版本資訊
├── release-notes.md         # 發布說明
└── scripts/                 # 構建腳本
    ├── build.sh
    └── deploy.sh
```

## 功能特點

### 1. 分支管理
- 支援 dev 分支 (開發測試)
- 支援發布分支 (0901, 0902...)
- 動態切換分支並顯示對應配置

### 2. Web UI 功能
- **分支選擇**: 下拉選單選擇分支
- **Release Notes 顯示**: 自動載入選定分支的發布說明
- **版本資訊顯示**: 顯示子模組版本號
- **構建配置**: 選擇平台和構建步驟
- **即時狀態**: 顯示構建進度和狀態

### 3. API 端點
- `GET /api/branches` - 獲取可用分支列表
- `GET /api/config/:branch` - 獲取指定分支的配置
- `GET /api/versions/:branch` - 獲取指定分支的版本資訊
- `GET /api/release-notes/:branch` - 獲取指定分支的發布說明
- `POST /api/build` - 開始構建流程
- `GET /api/build/status/:id` - 獲取構建狀態

## 安裝和執行

### 前置需求
- Go 1.21 或更高版本
- Git
- Docker (用於構建，可選)

### 安裝依賴
```bash
cd build-tool
go mod tidy
```

### 編譯執行
```bash
cd build-tool
go build .
./build-tool
```

或直接運行：
```bash
cd build-tool  
go run .
```

程式將在 `http://localhost:8080` 啟動。

### 設定 Git Token (可選)
如果你的 Git 倉庫需要認證，請在 `config.json` 中設定 Token：

```json
{
  "git_configs": {
    "main": {
      "url": "https://gitlab.example.com/group/repo.git",
      "token": "your-pat-token-here",
      "description": "主要配置倉庫"
    }
  }
}
```

## 配置說明

### 環境變數
- `GITLAB_TOKEN`: GitLab 存取權杖 (必要)

### 配置檔案格式

#### config.yaml
包含專案設定、倉庫資訊、構建設定和部署配置。

#### versions.json
包含版本資訊和子模組版本號。

#### release-notes.md
包含該版本的發布說明和變更記錄。

## 工作流程

1. **分支管理**: dev 分支開發 → merge 到發布分支 (如 0901)
2. **配置讀取**: 程式從 Git 倉庫讀取選定分支的配置
3. **構建執行**: 根據配置執行構建、推送、部署流程
4. **狀態追蹤**: 即時顯示構建進度和結果

## 最近更新 (2025/09/02)

### 修復的UI問題
- [x] **修復JavaScript與HTML元素ID不匹配問題**
  - 統一了 `gitConfigSelect`, `branchGrid`, `branchInfoSection`, `buildSection` 等元素ID
  - 修復了分支選擇和資訊顯示功能

- [x] **修復API調用問題**
  - 修正了Git配置API的數據格式處理 (從物件改為陣列)
  - 修正了API路徑，支援 `/api/branches/{gitConfig}` 格式
  - 修正了WebSocket消息格式處理

- [x] **改善UI交互體驗**
  - 添加了分支卡片式選擇界面
  - 改善了構建狀態指示器和進度條
  - 添加了即時日誌顯示功能
  - 修復了構建請求通過WebSocket發送

- [x] **修復版面跳動問題**
  - 使用佔位符訊息取代 `display: none` 隱藏方式
  - 設定區塊最小高度避免版面突然變化
  - 添加平滑過渡效果和友善的提示訊息
  - 分支資訊和構建配置區塊現在有固定空間，不會破壞版面

- [x] **重新設計佈局結構**
  - 改用兩欄式佈局：左側顯示分支資訊，右側顯示構建配置
  - 移除複雜的grid佈局，避免版面錯亂
  - 簡化資訊顯示：Release Notes、版本資訊、配置預覽垂直排列
  - 構建日誌區塊移到最下方，不再被其他內容影響
  - 添加響應式設計，小螢幕時自動切換為單欄佈局

- [x] **完善CSS樣式**
  - 添加了分支卡片樣式和懸停效果
  - 改善了按鈕狀態和日誌容器樣式
  - 添加了響應式設計支援
  - 新增佔位符訊息樣式，提升用戶體驗

### 🎨 UI強化更新 (基於參考專案模式)

- [x] **現代化視覺設計**
  - 採用參考專案 `/home/cch/eventcenter/manager` 的設計模式
  - 實現漸層背景和毛玻璃效果 (backdrop-filter)
  - 增強按鈕樣式：漸層背景、3D懸停效果、光暈動畫
  - 改進分支卡片：漸層背景、選中狀態動畫、頂部指示條

- [x] **增強交互體驗**
  - 添加豐富的動畫效果：懸停、選中、加載狀態
  - 實現狀態指示器動畫：脈動、波紋、震動效果
  - 改進進度條：漸層填充、閃爍動畫
  - 增強日誌容器：終端風格、掃描線動畫

- [x] **優化用戶界面**
  - 重新設計表單元素：圓角、陰影、焦點狀態
  - 增強日誌顯示：彩色圖標、時間戳格式、自動滾動
  - 改進信息展示卡片：懸停效果、漸層背景
  - 添加鍵盤快捷鍵：Ctrl+Enter 開始構建、Escape 停止

- [x] **響應式和可用性提升**
  - 完善移動端適配：單欄佈局、觸摸友好
  - 優化滾動條樣式：自定義顏色、圓角
  - 增加加載狀態指示：旋轉動畫、按鈕狀態
  - 實現自動重連機制：WebSocket 斷線重連

- [x] **技術架構改進**
  - 創建增強版 CSS (`style_enhanced.css`) 和 JS (`app_enhanced.js`)
  - 保持向後兼容的同時提供現代化體驗
  - 優化代碼結構：模塊化函數、錯誤處理、性能優化
  - 添加實用工具函數：文件大小格式化、時間格式化

## 下一步開發

### 已實作功能
- [x] Git 倉庫操作 (讀取分支、檔案)
- [x] 多 Git 配置支援 (含 PAT Token)
- [x] 實際構建流程執行
- [x] 構建狀態管理 (WebSocket 即時更新)
- [x] Release Notes 和版本資訊顯示
- [x] 配置檔案預覽
- [x] 腳本執行功能
- [x] **UI交互修復** (分支選擇、狀態顯示、即時日誌)

### 待實作功能
- [ ] 腳本輸出即時顯示
- [ ] 錯誤處理和回滾
- [ ] 部署後健康檢查
- [ ] 權限管理
- [ ] 構建歷史記錄

### 優化方向
- [ ] 添加配置驗證
- [ ] 改善錯誤提示
- [ ] 添加構建歷史
- [ ] 支援並行構建
- [ ] 添加通知功能

## 更新記錄

### 2025-09-02 UI 重新設計與優化

**改進內容：**

1. **全新的現代化布局設計**
   - 採用 CSS Grid 布局，實現響應式設計
   - 引入 Font Awesome 圖標，提升視覺效果
   - 使用漸層色彩和陰影效果，打造專業外觀

2. **可隱藏與調整的左側選單**
   - 左側選單包含 Git 配置選擇和分支列表
   - 支援隱藏/顯示功能（點擊頂部選單按鈕）
   - 支援拖拽調整寬度（200px - 500px）
   - 分支列表採用卡片式設計，顯示分支類型標籤

3. **分頁式分支資訊展示**
   - 將原本的分支資訊整合成 4 個分頁：
     - Release Notes（預設頁面）
     - 版本資訊
     - 配置資訊  
     - 構建配置
   - 每個分頁都有獨立的圖標和清晰的內容區域

4. **可隱藏的底部構建日誌面板**
   - 構建日誌移至底部面板
   - 支援隱藏/顯示功能
   - 支援拖拽調整高度（100px - 400px）
   - 日誌採用終端機風格設計，支援不同類型的日誌顏色

5. **改進的構建配置界面**
   - 構建選項採用網格布局的卡片式設計
   - 每個選項都有對應的圖標
   - 進度條顯示百分比數字
   - 按鈕狀態更加直觀

6. **響應式設計**
   - 支援手機和平板設備
   - 在小螢幕上左側選單自動變為覆蓋式
   - 構建選項在小螢幕上變為單列布局

**技術特色：**
- ✅ CSS Grid + Flexbox 現代布局
- ✅ 平滑動畫過渡效果
- ✅ 自定義滾動條樣式
- ✅ 拖拽調整面板大小
- ✅ WebSocket 自動重連機制
- ✅ 日誌條目數量限制（防止記憶體問題）

**使用體驗提升：**
- 🎨 更美觀的視覺設計
- 🖱️ 更直觀的操作方式
- 📱 更好的響應式支援
- ⚡ 更流暢的動畫效果
- 🔧 更靈活的界面調整

### 2025-09-02 Git 分支功能修正

**問題描述：**
Git 分支功能選擇不正常，前端無法正確載入和顯示分支列表。

**修正內容：**

1. **前端資源路徑錯誤修正**
   - 修正 `index.html` 中的 JavaScript 檔案路徑：`app_enhanced.js` → `app.js`
   - 修正 `index.html` 中的 CSS 檔案路徑：`style_enhanced.css` → `style.css`
   - 這些檔案路徑錯誤導致前端 JavaScript 無法載入，分支選擇功能完全失效

2. **配置檔案中無效 Git 倉庫修正**
   - 將 `config.json` 中的無效測試倉庫 `https://github.com/demo/build-config.git` 
   - 更新為有效的公開倉庫 `https://github.com/octocat/Hello-World.git`
   - 避免因為無效倉庫導致的 Git 操作失敗

**修正後功能：**
- ✅ Git 配置選擇下拉選單正常載入
- ✅ 分支列表能正確從 Git 倉庫獲取並顯示
- ✅ 分支選擇後能正確載入分支資訊（配置、版本、發布說明）
- ✅ WebSocket 連接正常，構建日誌即時顯示

**測試驗證：**
- 在 http://localhost:8082 可正常存取 Web UI
- API 端點 `/api/git-configs` 和 `/api/branches/{gitConfig}` 正常回應
- 前端 JavaScript 正確載入，分支選擇功能恢復正常

## 技術架構

- **後端**: Golang + Gorilla Mux + WebSocket
- **前端**: 原生 JavaScript + HTML/CSS (嵌入式)
- **部署**: 單一執行檔
- **配置**: YAML + JSON
- **版本控制**: Git (支援多倉庫)
- **認證**: PAT Token

## 使用說明

### 1. 設定配置檔案

編輯 `config.json`：

```json
{
  "server": {
    "port": "8080"
  },
  "git_configs": {
    "main": {
      "url": "https://gitlab.wise-paas.com/your-group/build-config.git",
      "token": "your-gitlab-token-here",
      "description": "主要構建配置倉庫"
    }
  }
}
```

### 2. 配置 Git 倉庫結構

在您的配置 Git 倉庫中，每個分支應包含：

```
├── config.yaml           # 構建配置
├── versions.json         # 版本資訊
├── release-notes.md      # 發布說明
└── scripts/              # 構建腳本
    ├── build.sh
    ├── push.sh
    └── deploy.sh
```

### 3. 執行流程

1. **選擇 Git 配置** - 選擇要使用的配置倉庫
2. **選擇分支** - 選擇 dev、0901、0902 等分支
3. **查看資訊** - 自動顯示 Release Notes、版本和配置
4. **配置構建** - 選擇要執行的步驟
5. **開始構建** - 即時查看執行過程和結果