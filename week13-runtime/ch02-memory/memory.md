# Memory allocat

## Stack & Heap

Go 有兩個地方可以分配內存： 一個全局 Heap 空間來動態分配內存，另一個是每個 goroutine 都有的自身 Stack 空間。

### Stack

Stack 區的內存一般由編譯器自動進行分配和釋放，其中存儲著函數的入参以及局部變量，這些參數會隨著函數的`創建而創建`，函數的`返回而銷毀`
