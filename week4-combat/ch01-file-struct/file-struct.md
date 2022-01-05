# Standard Go Project Layout 

- `/cmd`: `/cmd/myapp`
- `/internal`: 分割共享、非共享的代碼。
- `/pkg`: 工具， 略

# Kit project layout

kit（基礎庫）必備特點：
- 統一
- Go標準庫方式佈局
- 高度抽象：為了可擴展
- 支持插件

# Service Application Project Layout

## v1 - 失血模型
api, cmd, configs, internal. README,CHANGELOG, OWNERS.

model -> dao -> service -> api, model struct 串接各層

還有可能把api的數據從頭傳到尾，這樣不好。

## v2 - 

`little copy better than little dependence`

- internal 
    - biz: 依賴倒置.
    - data: 將領域模型拿出來.
    - service: 
- PO(Persistent Object): 持久化對象，他跟持久層的数据结构形成一一對應的映射關係。

幾種工程化模型：
- 失血模型：Model 層僅有`getter`/`setter`方法， 業務邏輯和應用邏輯都放在`服務層`
- 貧血模型：包含**不依賴**`持久層`的業務邏輯，這裏的Domain Object 是不依賴於持久層的。
- 充血模型：包含所有業務邏輯
- 脹血模型：All in One，沒用

## Life cycle
資源初始化後，再啟動監聽服務

### **依賴注入**
非依賴注入：    內部初始化
依賴注入：      外部初始化後再傳入

- Example: `Google`'s `Wire`