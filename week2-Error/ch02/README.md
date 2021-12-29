# Error Type

## Sentinel Error
預定義錯誤。

- 不能依賴這個方法來判斷Error，Error方法存在於 `error`裡，主要方便工程師使用，但是不適合程序使用，這個通常用於記錄日誌，輸出到stdout中。
- Sentinel errors 成為 API 公共部分。
- Sentinel errors 在兩個包之間創建了依賴
  
結論： 盡量避免 Sentinel errors

## Error Types
實現了 error 的自定義類型。

```Go
type MyError struct {
    Msg string
    File string
    Line int
}

func (e *MyError) Error() string{
    return fmt.Sprintf("%s: %d: %s", e.File, e.Line, e.Msg)
}

func test() error {
    return &MyError{"Something happened", "server.go", 42}
}

func main() {
    err := test()
    switch err := err.(type) {
        case nil:
        // call succeeded, nothing to do
        case *MyError:
            fmt.Println("error occurred on line:", err.Line)
        default:
        // unknown error occurred
    }
}
```

好的例子: os.PathError

## Opaque errors