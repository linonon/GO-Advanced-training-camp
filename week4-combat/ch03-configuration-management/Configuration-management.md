# Configuration
- 環境變量
Region, Zone, Cluster, Environment, Color等信息，在線運行平台打入容器，提供kit庫讀取使用。
- 靜態配置
不推薦`on-the-fly`變更，應該走迭代發布的流程。
- 動態配置
 
- 全局配置

## Functional options
Self-referential functions and the design of options -- Rob Pike
Functional options for friendly APIs -- Dave cheney

```go
// type DialOption struct {
//     f func(*dialOptions)
// }

type DialOption func(*dialOptions)

func Dial(network, addr string, options ...DialOption) (Conn, error) {
    do := dialOptions{
        // 可放默認值
        dial: net.Dial,
    }

    for _, option := range options {
        // option.f(&do)
        option(&do)
    }
}
```
使用方法：
```go
func main() {
    c, _ := redis.Dial("tcp", "127.0.0.1:3389",
    redis.DialDatabase(0),
    redis.DialPassword("hello world"),
    redis.DialReadTimeout(10 * time.Second)
    )
}
```
可見，我們只可以通過新建 `DialXXX()` 來對配置進行修改，別人無法在服務運行時對變量進行修改

更改以後需要還原: 
```go
type option func(f *Foo) option

func Verbosity(v int) option {
    return func(f *Foo) option{
        prev := f.verbosity
        f.verbosity = v
        return Verbosity(prev)
    }
}

func DoSomethingVervosely(foo *Foo, verbosity int) {
    // Could combine the next two lines, 
    // with some loss of readability.
    prev := foo.Option(pkg.Verbosity(verbosity))
    defer foo.Option(prev)
    // 方便單元測試後切回去
    // ... do some stuff with foo under high verbosity.
}
```

## Hybrid APIs
json.Marshal -> 生成配置結構 -> 作為參數傳入（ X ）

## Configuration & APIs

```go
func Dial(network,address string, options ...DialOption) (Conn, error)
```
- 僅保留 options API
- config file 和 options struct 解耦

配置工具的實踐：
- 語意驗證（Semantic verification）
- Syntax
- Lint
- Format

如何解耦呢:
```go
func main() {
   c := /* yaml file or toml file etc. */
   r, _ := redis.Dial(c.Network(), c.Address(), c.Options()...)
}
```
Options():
```go
func (c *Config) Options() []redis.Options {
    return []redis.Options{
        redis.DialDatabase(c.Database),
        redis.DialPassword(c.Password),
        redis.DialReadTimeout(c.ReadTimeout),
        //...; 可以通過自定義 `DialXXX()` 去可控的新增配置需求
    }
}
```
package redis:
```go
// Option configures how we set up the connection.
type Option interface {
    apply(*options)
}
```
結合YAML實例：
1. 用`JSON`作為"膠水"，實現 `YAML` -> `JSON` -> `Protobuf`
```go
func ApplyYAML(s *redis.Config/* protobuf */, yml string) error {
    js, err := yaml.YAMLToJSON([]byte(yml))
    if err != nil {
        return err
    }
 
    return ApplyJSON(s, string(js))
}
```
2. main()
```go
func main() {
    c := new(redis.Config)
    _ = ApplyYAML(c, loadConfig())
    r, _ := redis.Dial(c.Network, c.Address, Options(c)...)
}
```

配置的目標：
- 避免複雜：全局配置化模板
- 多樣的配置：模板可以通過覆蓋其中的字段實現多樣化
- 簡單化努力：避免過多沒必要的默認配置
- 以基礎設施 -> 面向用戶進行改變：業務配置越少越好，基礎設施配置越多越好。
- 配置的必選項和可選項區分：可選：最佳實踐的值
- 配置的防禦：不合理的配置直接`Panic()`，避免上線出問題
- 權限和變更追蹤
- 配置的版本和應用對其
- 安全的配置變更：逐步部署，回滾更改，自動回滾