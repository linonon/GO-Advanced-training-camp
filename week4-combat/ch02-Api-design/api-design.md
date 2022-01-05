# gRPC
protobuf： 好用，省電省流量
- 雙向流
- Header 壓縮
- 單 TCP 的多路復用：一個 Link 可以讓多個 goroutine 用（通過ID確認）

# API Project
為了統一搜索和規範API，建利一個bapis倉庫，整合所有對內對外API。
- API倉庫，方便跨部門協作
- 版本管理，基於Git控制
- 規範化檢查
- API design review
- 權限管理，目錄OWNERS

# API Compatibility
- 向後兼容（非破壞性）的修改
    - 給API服務定義添加API接口
    - 給請求消息添加字段
    - 給響應消息添加字段

- 向後不兼容（破壞性）的修改
    - 刪除或重命名服務，字段，方法和枚舉值
    - 修改字段的類型
    - 修改現有請求的可見行為
    - 給資源消息添加 Read/Load 字段

# API Naming Conventions
`Request URL`: `<package_name>.<version>.<service_name>/{method}`

`package <pakcage_name>.<version>`

# API Errors

大類錯誤(http.StatusXXX,400~500之類的) + 小類錯誤信息(myError.Enum)
## 小標準錯誤配合大量資源
- 狀態空間變小降低文檔的複雜性，在 Client 中提供了更好的慣用映射，降低客戶端的邏輯複雜性，同時不限制是否包含可操作信息

## 錯誤傳播
- 翻譯錯誤時應該注意兩點：
    - 隱藏實現詳細信息和機密信息（如id，密碼之類）
    - 調整負責該錯誤的一方。例如，從另一個服務接收 INVALID_ARGUMENT 錯誤的服務器應該將 INTERNAL 傳播給自己的調用者。

## 全局Error code
在每個服務傳播錯誤的時候，做一次翻譯，這樣保證每個服務 + 錯誤枚舉，應該是唯一的，而且在 proto 定義中是可以寫出來文檔的。

## API Design
FieldMask 部分更新的方案
```proto
service LibraryService {
    rpc UpdateBook(UpdateBookRequest) returns (Book);
}

message UpdateBookRequest{
    Book book = 1;
    google.protobuf.FieldMask mask = 2;
}
```