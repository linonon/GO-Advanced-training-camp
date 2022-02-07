# Storage Design

## DataBase Design
HBase

數據寫入： 
- rowkey：使用用戶 id md5 以後的頭兩位 + uid，避免 rowkey 熱點密集到一個 region 中，導致寫、讀熱點。
- PUT mid, values，只需要寫道 column_family 的 info 列簇。
    - obj_id + obj_type: 稿件業務1，稿件ID 100， 100_1作為列名
    - value 使用 protobuf序列化一個結構體接入，所以只需要單次更新 kv store

數據讀取：
- 列表獲取位 GET mid， 直接獲取1000條，在內存中排序和翻頁。

取捨：redis cache miss後，不回傳 HBase

## Cache Design

// TODO: