# 網絡通訊協議

## Socket 抽象層
應用通過 套接字 向網絡發出請求或者應答網絡請求。
主要操作：
- 建立、接受連結
- 讀寫、關閉、超時
- 獲取地址、端口

## TCP 可靠連結、面向連結的協議

服務端流程：
- 監聽端口
- 接受客戶端請求連結
- 創建 goroutine 處理連結

客戶端流程：
- 建利與服務端的連結
- 進行數據收發
- 關閉連結

closewait 過多肯定不正常，可能是因為服務器負載過高，處理不好、代碼有bug等。

## UDP 不可靠連結，允許廣播或多播

- 不需要建利連結
- 不可靠的，沒有時序通信
- 數據包有長度
- 支持多播和廣播
- 低延遲，實時性好
- 應用於影片直播，遊戲同步

## HTTP 超文本傳輸協議
`\r\n` 進行分割
- 請求報文
    - ...
- 相應報文
    - ...

- Linux 常用 Net 指令
    - nload
    - tcpflow
    - ss
    - netstat
    - nmon
    - top

## gRPC 基於 HTTP2 協議擴展

- Request
    - Headers
        - method = POST
        - scheme = https
        - path = /api.echo.v1.Echo/SayHello..- content-type = application/grpc+proto
        - grpc-encoding = gzip
    - Data
        - 1 byte of zero (not compressed zip)
        - network order 4 bytes of proto message length.
        - serialized proto message

- Response
    - Headers
        - status = 200
    - Data
    - Trailers

### HTTP2 如何提升網絡速度
- 1.1 的優化
    - 增加了持久連結，每個請求進行了串行模式
    - 瀏覽器為每個域名維護`6個TCP`持久連結
    - CDN的域名分片
- 2 的多路復用
    - 請求數據二進制分幀層處理後，會轉換成請求ID編號的幀，通過協議棧發送給瀏覽器。
    - 服務器接收到所有幀之後，會將所有相同的ID合併為一條完整的請求信息。
    - 然後服務器處理該請求，並將處理的相應行、相應頭和相應體分別發送至二進制分幀層。
    - 同樣，二進制分幀層會將這些相應數據轉化成一個個帶有請求ID變好的幀，經過協議發送給瀏覽器。
    - 瀏覽器接收到相應幀之後，會根據ID編號將幀的數據提交給對應的請求。