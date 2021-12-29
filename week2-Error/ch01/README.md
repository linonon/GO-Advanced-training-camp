# Error

推薦的 New() 規則
`errors.New("packagename: error string")`

## Panic() != (Java:Exception)
Go panic 代表程序掛了，後面不可以繼續操作。

- 簡單
- 考慮失敗，不是成功
- 沒有隱藏的控制流
- 調用者控制error
- error are valuse