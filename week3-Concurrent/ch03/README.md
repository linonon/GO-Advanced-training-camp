# Detecting Race Conditions With Go

`i++`: 雖然看上去是安全的，但是實際的組合語言是不安全的，應該用 `Mutex` or `Atomic` ( `channel` 偏重)

[Race Conditions Example](/week3-Concurrent/ch02/main.go)

# Data Race Conditions
1. `Atomic`
2. `Visable`

`Atomic.Value` 可以解決 map 指向的問題。應該看看文檔

# sync.atomic

`Copy-On-Write`：寫事複製，寫操作時複製全量老數據到一個新的對象中，攜帶上次新寫的數據，之後利用原子替換，更新調用者的變量。（變更指針）

- `BGSave` `Fork` `COW` 原理 
- 多數用在 `微服務降級` 或者 `local cache`

# Mutex
1. goroutine 1: 獲得Lock，休眠100ms後釋放
2. goroutine 2: 休眠100ms，再獲得lock然後釋放

結果： routine 1 獲得 Lock 遠超過 2

原因： 在 1 得到 Lock 的100ms 內， 2 會將獲得 Lock 的請求放到一個 FIFO 的隊列中，等 1 釋放之後，2 再喚醒，然後去嘗試獲得 Lock。 不過 1 釋放完以後又會馬上請求 Lock， 2 在喚醒的時候可能就得不到了。

![](/week3-Concurrent/pic/mutex-mode.png)
Figure Mutex Mode.

1. Barging: 提高吞吐量
2. HandsOff: 公平
3. Spinning: 快速，耗CPU

# errgroup

```go
func main(){
    g, ctx := errgroup.WithContext(context.Background())

    var a,b,c []int

    // 調用廣告
    g.Go(func() error{
        // a = xxx
        return errors.New("test")
    })

    // 調用ai
    g.Go(func() error{
        // b = xxx
    })

    // 調用運營平台
    g.Go(func() error{
        // c = xxx
    })

    err := g.Wait()
    fmt.Println(err)

    fmt.Println(ctx.Err())
}
```
 
# sync.Pool

ring buffer + double chain