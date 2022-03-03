# Goroutine VS Thread

## Goroutine

### Goroutine 和 thread 的區別

#### 擴所容

Goroutine Stack 內存 2kb，不夠就會擴容

Thread的話，需要分配 1～8 mb 的內存，同時還有 `guard page`，所以Thread 大小不能隨意更改

#### 創建與銷毀

創建/銷毀，Thread 的話會有巨大的消耗。但 goroutine 是由 go runtime管理，創建和銷毀的消耗非常小。

#### 調度切換

Thread： 1000～1500 ns，指令要12000～18000

Goroutine：200ns， 指令要 2400～3600

所以 goroutine 的切換成本就比 thread小得多

#### 複雜性

Thread： 非常複雜，多個thread的 share memory容易出問題。不能大量創建線程，成本高。使用網絡多路復用，存在大量callback。對於應用服務線程門檻高，例如需要做第三方庫隔離，需要考慮創建線程池。

### M: N 模型

Go 創建 M個 Thread，之後再創建 N個 goroutine 放在這 M個Thread 上執行。

## GMP概念-Part1

- G： goroutine 的縮寫，無限制
- M： OS Thread，也被稱為 Machine，使用 struct runtime.m，所有 M 是有Thread Stack的。
- P

### GM調度器

Go 1.2前的調度器，限制 Go 併發程序的伸縮性

這個模型的問題：

- 單一全局互斥鎖和幾種狀態存儲
- Goroutine傳遞問題：M之間經常性傳遞可運行的 goroutine，這導致調度延遲增大以及額外的性能損耗。
- Per-M 持有內存緩存：M提前撈一批內存過來，不夠再撈

## GMP概念-Part2

引入來 local queue，因為 P 的存在，runtime 並不需要做一個集中式的 goroutine 調度，每一個 M 都會在 P‘s local queue， global queue 或者其他 P 隊列找 G 執行，減少全局鎖對性能的影響。

## Wokr-stealing

TODO: Lock-Free:<https://yizhi.ren/2017/09/19/reorder>

當一個 P 完成本地所有 G 之後，並且全局隊列為空的時候，會嘗試挑選一個受害者 P2，從 P2 的 G 隊列中竊取一半的 G。否則從全局隊列中獲取（當前個數/GOMAXPROCS）個 G。同時還要考慮不同的獲取順序（互為質數的步長）

## Syscall

TODO: 可以再看看

## Spinning thread

自旋情況:自旋的 M 最多只允許 GOMAXPROCS(Busy P)； 同時，有類型1就不會有類型2

- 類型1: M 不帶 P 的找 P 掛載（一有 P 釋放就結合）
- 類型2: M 帶 P 的找 G 運行（一有 runable 的 G 就執行）

在新 G 被創建，M 進入系統調用， M 從空閒被激活這三種狀態變化前，調度器會確保至少有一個自旋 M 存在（喚醒 或者 創建一個 M），除非沒有空閒的 P

## GMP 問題總結

### 單一全局互斥鎖 (Sched.Lock) 和 集中狀態存儲

G 被分成`全局隊列` 和 的`本地隊列`,
