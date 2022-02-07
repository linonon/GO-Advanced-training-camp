# Goim 長連接 TCP 編程

## 概覽

- Comet: 請求進來讓 Comet 託管，再起 RPC 連接給後面的 Logic 層
- Logic: 監控 Connect/Disconnect，可自定義auth，記錄 Session，業務也可以通過設備ID、用戶ID、RoomID、全局廣播進行消息推送。
- Job: 通過消息隊列進行消峰處理，並把消息推送到對應 Comet 節點

## 協議設計

現在的話，可以直接用 gRPC。

## 邊緣節點

Comet 長連接連續節點，通常部署在距離用戶比較近，通過TCP 或者Websocket 建立連接，並通過應用層 Heartbeat 進行保活檢測，保證連接可用性。

節點之間通過雲VPC專線通信，按地區部署分佈。

## 負載均衡

比較特殊。

## 心跳保活

- 長連接斷開的原因：
  - 被kill
  - NAT 超時
  - 網絡狀態發生變化
  - 網絡差、DHCP到期等
- 高效長連接方案
  - 進程保活（防殺）
  - 心跳保活（阻止NAT超時）
  - 斷線重連
- 自適應心跳
  - 心跳可選區間
  - 心跳增加步長（網絡順暢，窗口可以變大）
  - 心跳週期探測

## 用戶 Auth 和 Session 信息

- Connect
- Disconncet
- Session

## Comet

Bucket，連接通過 DevicedID進行管理，用於讀寫所拆散，並且實現房間消息推送，類似Nginx Worker

每個 Bucket 都有獨立的 Goroutine 和讀寫鎖優化：

```go
Bucket {
    channels map[string]*Channel
    rooms map[string]*Room
}
```

- Bucket
  - 維護當前信息通道和房間的信息，有獨立的 Goroutine 和 讀寫鎖優化，用戶可以自定義配置對應的 Buckets 數量，在大併發業務上有其明顯。
- Room
  - 維護房間的通道Channel，推送消息進行了合併寫，即 Batch Write，如果不合併寫，每來一個小的消息都通過長連結寫出去，系統 SysCall 調用的開銷會非常大
- Channel
  - 一個連結通道，Writer/Readre 就是對網絡 Conn 的封裝， CliProto 是一個 Ring buffer， 保存 Room 廣播或是直接發送過來的消息體。

## Logic

業務邏輯層，處理連接Auth、消息路由、用戶會話管理

- sdk：通過 TCP/Websocket 建立長連接，進行重連，心跳保活；
- goim：負責連結管理，提供消息長連接能力
- backend：處理業務邏輯，對消息過濾以及持久化等相關操作。

## Job

通過 Kafka 進行消息消峰，保證信息逐步推送成功。

- 支持的多種推送方式
  - Push(DevicesID, Message)
  - Push(UserID, Message)
  - Push(RoomID, Message)
  - Push(Message)
