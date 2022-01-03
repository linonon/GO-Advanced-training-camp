# Goroutin

## Never start a goroutine without knowning when it will stop

當你創建 Goroutine 時，一定要問自己兩個問題
1. 什麼時候會關閉
2. 有什麼手段可以結束它

例子：
```go
func serveApp(){
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(resp http.ResponseWriter, req *http.Request) {
        fmt.Fprintln(resp, "Hello linonon")
    })
    http.ListenAndServe("0.0.0.0:8080", mux)
}

func serveDebug(){
    http.ListenAndServe("127.0.0.1:8001", http.DefaultServeMux)
}

func main(){
    go serveDebug()
    go serveApp()
    select {}
}
```
確保 serveApp 和 serveDebug 將它們到併發性留給調用者

`log.Fatal` 調用到是 `os.Exit`, `Defers`不會被調用

## Application Lifecycle