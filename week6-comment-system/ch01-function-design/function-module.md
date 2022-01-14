# Function module

在動手設計前，反覆思考，真正編碼的時間只有5%

# Framework modules

![](/week6-comment-system/ch01-function-design/pic/comment-framework.png)
<center>comment-framework</center>

## 先簡單劃清 module

- BFF(Backend for Frontend): `複雜評論業務的服務編排`，比如訪問帳號服務進行等級判定，同事需要在BFF面向移動端/WEB場景來設計API，這一層抽象吧評論本身的內容列表進行處理（`加載`、`分頁`、`排序`等）進行了`隔離`，關注在業務平台化邏輯上。
- Service: Comment-service: 服務層，去平台業務的邏輯，專注在評論功能的API實現上，比如發布、讀取、刪除等，關注在穩定性、可用性上。這樣讓上游可以靈活組織邏輯，把基礎能力和業務能力進行剝離。
- Job: comment-Job: 消息隊列最大的用途就是消峰處理。
- Admin: comment-Admin: `管理平台`，按照安全等級劃分服務，尤其劃分運營平台，他們會`共享服務層的存儲層`（MySQL、Redis）。運營體系的數據大量都是檢索，我們使用 canal 進行同步到 ES 中，整個數據的展示都是通過ES，在通過業務主鍵更新業務數據層，這樣運營端的查詢壓力就下方給了獨立的 fulltexvvt search 系統。
- Dependency: account-service、filter-service：整個評論服務還會依賴一些外部的gRPC服務，統一的平台業務邏輯在 comment BFF層收斂， 這裏 account-service 主要是帳號服務， filter-service 是敏感詞過濾服務。

## 具體模塊
- comment-service:
    - `讀`的核心邏輯
        - 先讀取緩存，在讀取存儲。 
        - 策略是多變的（先審後發，先發後審），但是讀，寫是不變的。
        - 讀多寫少：緩存設計很關鍵，其次是消息隊列的設計。存儲層很難進行`擴縮容`，可能需要擴縮容。
        - Cache rebuild：使用異步回源。
    - `寫`的核心邏輯：
        - 寫請求最終回穿透到存儲層：把存儲的直接衝擊下放到消息隊列
        - Kafka：全局並行，局部串行的生產消費方式。
        - 對於回源信息也是類似的思路。
- Comment：作為BFF，是面向端、平台、業務組合的服務，所以平台擴展的能力，我們都在Comment服務來實現，方便統一和准入平台，以統一的藉口形式提供平台化的能力。
    - 依賴其他 gRPC 服務，整合統一平台側的邏輯
    - 直接向端上提供接口，提供數據的讀寫接口，甚至可以整合端上，提供統一的端上SDK。
    - 需要對非核心依賴的 gRPC 服務進行降級，當這些服務不穩定時。