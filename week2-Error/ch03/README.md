# Handling Errors 

## Indented flow is for errors
無錯誤的正常流程代碼，將成為一條直線，而不是縮進的代碼。

推薦：
```go
f, err := os.Open(path)
if err != nil {
    // handle error
}
// do stuff
```
不推薦：
```go
f, err := os.Open(path)
if err == nil {
    // do stuff
}
// hanle error
```

Error的暫存返回
```go
func CountLines(r io.Reader) (int, error) {
    sc := bufio.NewScanner(r)
    // type Scanner struct {
    //     //..
    //     Err error
    // }
    // func (s *Scanner) Err(){}
    lines := 0

    for sc.Scan() {
        // Scan() 成功， 返回 true
        // 如果掃描出錯或掃到最後了，sc.Err = errors.New("Eof")
        line++
    }

    return lines, sc.Err()
}
```

errWriter
```go
type errWriter struct {
    io.Writer
    err error
}

func (e *errWriter) Write(buf []byte) (int, error) {
    if e.err != nil {
        return 0, e.err
    }

    var n int
    n, e.err = e.Writer.Write(buf)
    return n, nil
}
```

## Wrap errors
- `log 噪音`: 與 error 無關的且對調試沒有幫助的信息
- `記錄 log 的原因`: 因為有程序失敗了，而 log 包含了答案