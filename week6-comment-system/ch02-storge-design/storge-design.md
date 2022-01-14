# Storge-design
- 主表
    - id(主鍵)
    - 盡可能自增
    - 犧牲`寫`性能，提高`讀`性能。把統計類數據（評論總數，回復總數，根評論總數）存在表中，寫的時候一起更新。

## 索引內容分離
Comment_index: 評論樓層的索引組織表，實際並不包含內容。
Comment_Content: 評論具體內容的表。


## 緩存設計
MySQL獨立存儲 + Redis 加速就足夠了。

Comment_subject_cache: 對應主題的緩存，value使用 protobuf 序列化的方式存入。

Comment_index_cache: 使用 redis `sortedset` 進行索引的緩存，索引即數據的組織順序，而非數據內容。

## Singleflight
對於熱門主題，如果存在緩存穿透的情況，會導致大量的統進程，跨進程的數據回源到存儲層，可能會引起存儲過載的情況，如何只交給同進程內，一個人去做加載存儲？

- 使用`歸併回源`的思路：
    - 同進程 只交給一個人去獲取 mysql 數據，然後批量返回。同時這個 lease owner 投遞一個 kafka 消息，做 index cache 和 recovery 操作。這樣可以大大減少 mysql 的壓力，以及大量穿透導致的密集寫 kafka 的問題。
    - 更進一步，後續的連續請求，仍然可能會短時 cache miss，我們可以在進程內設置一個 short-lived flag，標記最近有一個人投遞了 cache rebuild 的消息，直接 drop。

## 熱點

流量熱點是因為突然熱門的主題，被高頻次的訪問，因為底層的 cache 設計， 一般時按照主題 key 進行一致性 hash 來進行分片， 但是熱點 key 一定命中某一個節點，這時候remote cache 可能會變成瓶頸，因此做 cache 的升級 local cache 是必要的，我們一般使用單進程自適應發現熱點的思路，附加一個短時的 ttl local cache，可以在進程內吞掉大量的讀請求。

在內存中使用 hashmap 統計每個 key 的訪問頻次，這裡可以使用滑動窗口統計，即每個窗口中，維護一個 hashmap ，之後統計所有未過去的 bucket， 匯總所有 key 的數據。

