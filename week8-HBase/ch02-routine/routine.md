# 分布式事務

## 事務消息
通過憑證（消息），完成最終一致性。類似麥當勞買吃的。

如何可靠的保存消息憑證？需要解決 mysql 和 message 存儲的一致性問題

業務冪等性

## Transaction log tailing
上述保存消息的方式使得消息數據和業務數據緊耦合在一起，從架構上看不夠優雅，而容易誘發其他問題。

有一些業務場景，可以直接使用主表被 canal 訂閱使用，有一些業務場景自帶這類 message 表，比如訂單或者交易流水，可以直接使用這類流水表作為 message 表使用。

實是流式消費數據，在消費者 balance 或者 balance-job 必須努力送達到。

## Polling publisher
定時的輪訓 msg 表，把 status = 1 的消息統統拿出來消費，可以按照自贈 id 排序，保證順序消費。

Pull 模型總體來說不太好， Pull太猛對 Database 有一定壓力， Pull 頻次低了，延遲比較高。

## 冪等
解決消息重複投遞：
- 全局唯一ID + 去重表
    - 在執行之前，先去消息應用狀態表中查詢一邊，如果找到，說明是重複消息。沒找到才可以執行。
```sql
for each msg in queue
    Begin transaction
        select count(*) as cnt from message_applay where msg_id=msg.msg_id;
        if cnt==0 then
            update B set amount=amount+10000 where userID=1;
            insert into message_applay(msg_id) values(msg.msg_id);
    End transaction
commit;
```

## 2PC
兩階段提交協議（Two Phase Commitment Protocol）
- 事務協調者：負責協調多個參與者進行事務投票及提交
- 多個事務參與者：即本地事務執行者



## 2PC Message Queue
生產者集群
1. 發送準備消息
2. 執行本地事務
3. 發送確認消息

傳給 RocketMQ集群

消費者集群
1. 接受消息
2. 執行本地事務
3. 發送確認消息成功 ack

等於用一個中間件解耦，異步同步消息

## TCC
- Try: 預處理
- Confirm： 確認
- Cancel： 撤銷

amount + frozen amount