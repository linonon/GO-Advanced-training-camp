# Memory alloca

## Stack & Heap

Go 有兩個地方可以分配內存： 一個全局 Heap 空間來動態分配內存，另一個是每個 goroutine 都有的自身 Stack 空間。

### Stack

Stack 區的內存一般由編譯器自動進行分配和釋放，其中存儲著函數的入参以及局部變量，這些參數會隨著函數的`創建而創建`，函數的`返回而銷毀`

### Heap

Heap 區的內存一般由編譯器和工程師自己共同進行管理分配，交給 Runtime GC 來釋放。Heap 上分配必須找到一塊足夠大的內存來存放新的變量數據。後續釋放時，垃圾回收器`掃描堆空間尋找不再被使用的對象`。

Stack 分配廉價， Heap 分配昂貴。

## Var in Stack or Heap?

### Go FAQ

如果 GO 編譯器在函數返回後無法證明變量未被引用，則編譯器必須在會被垃圾回收的堆上分配變量以避免懸空指針錯誤。此外，如果局部變量非常大，將他存儲在堆上而不是棧上可能更有意義。

## 逃逸分析

`go build -gcflags -m`: 可以標記出逃逸的變量。

逃逸分析在大多數語言裡屬於靜態分析：在編譯器由靜態代碼分析來決定一個值是否能夠被分配在Stack Frame 上，還是需要逃逸到 Heap 上。

- 減少 GC 壓力，Stack 上到變量，隨著函數退出後系統直接回收，不需要標記後再清除。
- 減少內存碎片的產生
- 減輕分配堆內存的開銷，提高程序的運行速度。

## 超過 Stack Frame

就會被分配到 Heap，避免閉包外部調用失敗。

## Contiguous stacks

採用複製 Stack 的實現方式，在 Hot Split 場景中不會頻繁釋放內存，即不像分配一個新的內存塊並鏈接到老的 Stack 內存塊，而是會分配一個兩倍大的內存塊，並把老的內存塊內容複製到新的內存裡，當 Stack 縮減回之前大小時，我們不需要做任何事情。

- runtime.newstack 分配更大的 Stack 內存空間
- runtime.copystack 將 舊Stack 中的內容複製到 新Stack 中
- 將指向 舊Stack 對應的變量指針重新指向 新Stack
- runtime.stackfree 銷毀並回收 舊Stack 的內存空間

`如果 Stack區 的空間使用率不超過 1/4，那麼在垃圾回收的時候使用 runtime.shrinkstack 進行 Stack縮容，同樣適用 copystack 完成後續操作`

## Stack 擴容

Go 運行式判斷 Stack空間 是否足夠，所以在 Call function中會插入 runtime.morestack。

## 內存管理

需要解決的問題：

- 內存碎片：
- 大鎖：

### 重要概念

- page： 內存頁，一塊 8K 大小的內存空間。 Go 與操作系統之間的內存申請和釋放，都以 page 為單位。
- span：內存塊，一個或多個連續的 page 組成一個 span
- sizeclass： 空間規格，每個 span 都帶一個 sizeclass，標記著該 span 中的 page 應該如何使用。
- object： 對象，用來存儲一個變量數據內存空間，一個 span 在初始化時，會被切割成一堆等大的 object（ memcache 思想與其類似），假設 object 的大小事 16B， span 大小是8K，那麼span 中的 page 就會被初始化 8K/16B = 512 個 object。
