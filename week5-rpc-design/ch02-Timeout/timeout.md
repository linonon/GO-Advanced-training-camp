# Timeout
- 網路傳播具有不確定性
- 客戶端和服務端不一致的超時策略導致資源浪費
- “默認值”策略
- 高延遲服務導致 client 浪費資源等待，使用超時傳遞：進程間傳遞 + 跨進程傳遞

我們依賴的微服務的超時策略並不清楚，或者隨著業務迭代耗時超時發生了變化，意外的導致依賴者出現了超時。
- 服務提供者定義好latency SLO，更新到 gRPC Proto 定義
- kit基礎庫兜底默認超時，比如100ms，進行配置防禦保護，避免出現類似60s之類的超大超時策略
- 配置中心公共模板，對於未配置的服務使用公共配置。

可以在 proto 檔案中寫好接口的 SLO: `95th: 100ms, 99th 150ms`

## 超時傳遞
超時傳遞指的是把當前服務的剩餘 Quata 傳遞到下游服務中，繼承超時策略， 控制請求級別的全局超時控制。
- 進程內超時控制：一個請求在每個階段（網路請求）開始前，就要檢查是否還有足夠的剩餘來處理請求，以及繼承他的超時策略，使用 Go 標準庫的 `context.WithTimeout`

Example: Redis:
```go
func (c *asiiConn) Get(ctx context.Context, key string) (result *Item, err error){
    c.conn.SetWriteDeadline(shrinkDeadline(ctx,c.writeTimeout))
    if _, err = fmt.Fprintf(c.rw, "get %s\r\n", key); err != nil {
        ///
    }
}
```

## 跨進程傳遞
1. A gRPC 請求 B， 1s超時
2. B 使用了300ms 處理請求， 再轉發請求C。
3. C 配置了600ms 超時，但實際只用了500ms
4. 到其他的下游，發現餘量不足，取消傳遞。

- 雙峰分佈：95%的請求耗時在100ms內，5%的請求可能永遠不會完成（長超時）
- 對於監控不要只看 average，可以看看耗時分佈統計，比如 95th，99th，999th等
- 設置合理的超時，拒絕超長請求，或者當Server不可用要主動失敗。