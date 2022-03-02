# Leader & Followers

## Leader

Leader 負責數據的讀寫， Follower 等待 Leader出錯，然後嘗試成為新的 Leader

Kafka 集群依賴 Zookeeper 集群，所以最簡單最直觀的方案是，所有 Follower 都在 ZooKeeper 上設置一個 Watch， 一旦 Leader 當機， Leader對應的 ephemeral znode 會自動刪除，此時所有 Follower 都嘗試創建該節點，只有一個會創建成功成為新的 Leader，其他的繼續當 Follower

缺點：

- 腦裂：每個 Follower 會 Watch，但是不能保證所有的 Follower Watch 的狀態是一致的。
- 羊群效應： 如果當機的 Broker 有多個 Partition， 會造成多個 Watch 被觸發，造成集群內大量的調整（重新訂閱）。
- 負載過重：過多 Partition的時候， Watch 也會過多

## Controller

Kafka 從所有的 Broker 中選出一個 Controller，所有 Partition 的 Leader 都由 Controller 決定。 Controller 會將 Leader 的改變直接通過 RPC 的方式通知需要為此做出相應的 Broker。

Kafka 集群 Controller 的選舉過程如下：

- 每個 Broker 都會在 Controller Path 上註冊一個 Watch
- 當 Controller 失敗時， 對應的 Path 會自動消失，此時該 Watch 被 fire，所有 ”活”的 Broker 都會去競選成為新的 Controller， 但是只會有一個競選成功。
- 競選失敗的重新在 Path上註冊 Watch

## High Watermark & Log End Offset

每個 Kafka 副本對象都由兩個重要的屬性： LEO 和 HW。 注意是所有的副本，而不只是 Leader 副本。

- LEO: 記錄該副本底層日誌(log)中下一條消息的位移值
- HW: 水位值，對於同一個副本對象而言，HW 不會大於 LEO。小於 HW 的所有消息都被認為是 “已備份” 的。 同理， leader 副本和 follower副本的 HW 更新是有區別的。

Leader 根據 Follower 發來的 Fetch 請求中的 fetch offset 來確定 remote LEO 的值

HW = min(LEO, remoteLEO)

只有 HW 以下的數據是 Consumer 可見的。
