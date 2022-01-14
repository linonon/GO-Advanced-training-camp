# Function and Framework

## Funciton

- 變更功能：添加記錄，刪除記錄，清空歷史
- 讀取功能：按照timeline 返回 topN，點查獲取進度信息。
- 其他功能：暫停/回復記錄，首次觀察增加經驗。

## Framework 

### history-service
kafka 是為高吞吐設計，超高頻的寫入並不是最優的，所以內存聚合和分片算法比較重要，按照 uid 來 sharding 數據，寫放大仍然很大，這裏我們使用 `region sharding`，打包一組數據當作一個 kafka message（比如 `uid % 100` 進行數據打包）