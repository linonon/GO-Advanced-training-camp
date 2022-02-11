# Cache

## 緩存選型

### memcache

memcache 提供簡單的 kv cache 存儲，value 大小不超過1mb

memcache 使用了 slab（預先分配內存大小塊）做內存管理， 存在一定的浪費， 如果大量接近的item，建議調整 memcache 參數來優化每一個 slab 增長的 ratio、可以通過設置 slab_automove & slab_reassign 開啟 memcache 的動態/手動 move slab， 防止某些 slab 熱點導致內存足夠的情況下引發 LRU。（比如分割大小是20， 但是塞了一個 200 的數據進來，就會引發 LRU，之前的 slab 對象會被 LRU 算法踢掉）

大部分情況下，簡單 KV 推薦使用 memcache， 吞吐和相應都足夠好。

內存池設計：Nginx：ngx_pool_t, tcmalloc

### Redis

redis 有豐富的數據類型，支持增量方式的修改部分數據，比如排行榜，集合，數組等。

比較常用的方式是使用 redis 作為`數據索引`（比如評論列表的ID，播放歷史的列表ID集合，我們的關係列表ID）

redis 沒有使用內存池，一般使用 `jemalloc` 來優化內存分配，需要編譯時使用 `jemalloc` 代替 glib 的 malloc 使用。

redis 使用 大value 的話，還是會有一定的影響。

### Proxy

數據需要去中心化，每個微服務應該獨佔一個緩存節點。

- 從集中式訪問緩存到 Sidecar 訪問緩存：
  - 微服務強調去中心化
  - LVS 運帷困難，容易造成流量熱點，隨下游擴容而擴容，連結不均衡等問題。
  - Sidecar 伴生容器隨 App 容器啟動而啟動，配置簡化。

### 一致性Hash

一致性 hash 是將數據按照特徵值`映射`到一個首位相接的 hash 環上，同時也將節點（按照IP地址或者及其名hash）映射到這個環上。

對於數據，從數據在環上的位置開始，`順時針`找到的第一個節點即為數據的存儲節點。

- 餘數分布式算法：Node = HashCode(key) % (Number of servers)
  - 由於保存鍵的服務器會發生巨大變化而影響緩存的命中率
  - 因為服務器位置的偏移，會產生大量的 Cache Miss
- Consisten Hashing: 圓形 Hash 鏈表。
  - 只有在 圓 上的服務器的地點逆時針的第一台服務器上的 key 會收到影響。

![一致性Hash環](/week8-HBase/pic/一致性Hash環.png)

- Hash 算法的考量指標：
  - 平衡性（Balance）：盡可能分佈到所有的緩存中去
  - 單調性（Monotonicity）：如果已經有一些`內容1`通過哈希分派到相應的`緩衝`中，又有新的`緩衝區`加入到系統中。那麼 Hash 的結果應能夠保證原有已分配的`內容1`可以被映射到`新的緩衝區`中去，而不會被映射到舊的緩衝集合中的`其他緩衝區`
  - 分散性（Spread）：`相同內容`的存儲到不同緩衝中去，降低了系統存儲到效率，需要盡量降低分散性。
  - 負載（Load）：Hash 應該能夠盡量降低緩衝的負荷。
  - 平滑性（Smoothness）：緩存`服務器的數目平滑改變`和`緩存對象的平滑改變`是一致的。

    ![一致性Hash環-加點](/week8-HBase/pic/一致性Hash環-加點.png)

- 引入量虛擬節點
  - 通過在服務器IP或者主機名後面增加編號來實現
- 舉例：
  - 每台服務器計算三個虛擬節點。
  - 但是數據定位算法不變，只是多了一步`虛擬節點到實際節點的映射`，如過數據找到了NodeA_1/2/3，那麼數據就會定位到 NodeA實體節點上。

- 搶紅包：
  - 在網關層，使用一致性 Hash，對紅包 id 進行分片，命中到某一個邏輯服務器處理，在進程內做寫操作的合併，減少存儲層的單行鎖徵用。
  - 合併紅包計算邏輯，一次更新數據庫

- 有界負載一致性 Hahs：一致性Hahs升級版
  - 增加了負載信息從而進行決策。

### Slot

- Redis-Cluster
  - 按照16384槽按照節點數量進行平均分配，由節點進行管理。
  - 對每個 Key 按照 CRC16 規則進行 Hash 運算，把 hash 結果對 16383 進行取餘，把餘數發功給 Redis 節點。
  - Redis-Cluster 的節點之間會共享信息(Gossip)，每個節點都會知道哪個節點負責那個範圍對應數據槽。

### 數據一致性

Storege 和 Cache 同步更新容易出現數據不一致。

模擬 MySQL Slave 做數據複製，再把消息投遞到 Kafka， 保證至少一次消費：

- 同步操作DB
- 同步操作Cache
- 利用 Job 消費消息，重係補償一次緩存操作

這樣做保證了時效性和一致性

有時候因為CacheMiss的回填會導致 v2 被複寫成 v1，這時候用 SETNX（SET if Not eXits），避免回填操作複寫 。

### 多級緩存

微服務拆分細粒度原子業務下的整合服務（聚合服務），用於提供粗粒度的接口，以及二級緩存加速，減少扇出的 rpc 網絡請求，減少延遲。

最重要的是保證多級緩存的一致性：

- 清理的優先級是由要求的，從下游往上游清理。
- 下游的緩存 expire TTL 要大於上游，避免穿透回源。

通過 DDD 思路區分不同的 Usercase。

### 熱點緩存

- 小表廣播，從 RemoteCache 提升為 LocalCache，App 定時更新，甚至可以讓運營平台支持廣播刷新 LocalCache
- 主動監控防禦預熱，比如直播房間頁在高在線情況下直接外掛服務主動防禦。
- 基礎庫框架支持熱點發現，自動短時的 short-live cache；
- 多 Cluster 支持
  - 多 Key 設計： 使用多副本，減小節點熱點的問題
- 副本策略總是會有`一致性`的問題
- 使用 `delete` 把 `高熱的key` 刪除了以後，別的查找就會因為 `Cache Miss` 而把`請求`透傳到數據庫。
  - 或許可以通過設置 `假delete`（先移到一個臨時區域並做上標記），請求進來以後會知道這個數據已經被刪除了，然後再通過業務決定要不要使用這個數據。

### 穿透緩存

- singlefly：對關鍵字進行一致性 Hash，使其某一個維度的 key 一定命中某個節點，然後在節點內使用 Mutex Lock，保證歸併回源，但是對批量查詢無解。
- 分布式鎖：不推薦，因為無法確保是否真正鎖住。
- 隊列：如果 Cache Miss，交由隊列聚合一個 Key，來 load 數據回寫緩存，對於 miss 當前請求可以使用 singlefly 保證回源，如評論架構實現。適合回源加載`數據重`的任務。比如評論 miss 只返回第一頁，但是需要構建完成評論數據索引。
- lease：這是一個 int64 的 token，與客戶端請求的 key 綁定，對於過時設置，在寫入時驗證 lease，可以解決這個問題。當 Client 在沒有獲取到 lease 時，可以稍等一下再訪問 cache， 這時 cache 中往往已有數據。

解決穿透緩存這個問題的核心：

1. 只讓一個人去 DB 取數據
2. 只讓一個人去構建 Cache
    1. 同時可以使用 Kafka 消息隊列的方式去構建 Cache，可以減少大量的內存壓力。

TODO: 了解 CRDT 解決緩存穿透

## 緩存技巧

### Incast Congestion

簡單來說就是：我們一次發送大量的包，這會導致 Router/Switch 分配不過來。

我們可以通過`自適應調整消息聚合力度的大小`的方法來解決這個問題，但是通常不會在業務代碼裡面實現，一般在 proxy 層面實現。

### 小技巧

- key 盡可能設置小，減少資源佔用。
- redis value 可以用 int 就不要用 string。
  - 對於小於 N 的 Value，redis 內部有 shared_object 緩存。
- 拆分 Key。 同樣的 Hashes key 回落到同一個 redis node，所以會產生熱點。
  - 考慮對 Hash 進行拆分（sharding 等方法）位小的 hash。
- 空緩存設置，避免每次請求 miss 都直接打到 DB。
  - 可以對 Key 進行內部加密，加大攻擊者構建出可用的 Key 的難度
- 讀失敗後的寫緩存策略（降級後一半讀失敗不出發回寫緩存）
- 序列化使用 protobuf，盡可能減少 size。
- 工具化膠水代碼（go generator）
  
### memcache 小技巧

- flag 使用： 標示 compress、encoding、large value等
- memcache 支持 gets， 盡量讀取， 盡可能的 pipeline，減少網絡往返
- 使用二進制協議，支持 pipeline delete， UDP 讀取，TCP 更新。

### redis 小技巧

- 增量更新一致性：
  - 如果使用 Exist，可能發現存在後，新增的時候 key 又沒了效果。
  - EXPIRE、ZADD/HSET 等，保證索引結構體務必存在的情況下去操作新增數據
- BITSET：存儲每日登陸用戶，單個標記位置，為了避免單個 BITSET 過大或者熱點，需要使用 region sharding
- List存儲：用類似 Stack PUSH/POP 操作。
  - 抽獎的獎池:
    - 先根據中獎概率篩選用戶，再讓中獎了的去獎池爭搶
    - 不同策略的獎池各自打散成一個 List，然後 POP 獎品
  - 核心：`打散`和`熱點`問題
- Sortedset: 用於 翻頁、排序、有序的集合，杜絕 zrange 或者 zrevrange
- 盡可能 PIPELINE 指令，避免集合過大。避免 for 循環。
