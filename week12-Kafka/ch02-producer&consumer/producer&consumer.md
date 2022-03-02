# Producer & Consumer

## Producer

為保證 Producer 發送的數據，能可靠的發送到指定的 Topic， Topic 的每個 Partition 收到 Producer 發送的數據後，都需要向 Producer 發送 ACK。如果 Producer 收到 ACK，就會進行下一輪的發送，否則重新發送數據。

- 選擇完分區後，生產者知道了消息所屬的主題和分區，它將這條記錄添加到相同主題和分區的批量消息中，另一個現成負責發送這些批量消息到對應的Kafka Broker
- 當 Broker 接收到消息後，如果成功寫入則返回一個包含消息的主題，分區和偏移的 RecordMetadata 對象，否則返回異常
- 生產者收到數據後，對於異常可能會進行重試。

TODO: Learning Flink

## Push VS PULL

Producer --push-> Broker --pull-> Consumer

Push 和 Pull 模式的優劣

- Push模式 很難適應消費速率不同的消費者，無腦Push容易出問題——拒絕服務以及網絡擁堵
- Pull模式 則可以根據 Consumer 的消費能力以適當的速率消費信息。

對 Kafka 而言， Pull模式更合適。 Pull 模式可以簡化 Broker 的設計， Consumer 可自主控制消費消息的速率，同時 Consumer 可以自己控制消費方式

同時 Kafka 也使用了長輪訓，以免傻等
