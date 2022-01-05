# Go Modules Proxy

GoPROXY 可以解決一些公司內部的使用問題：
- 訪問公司內網的 git server
- 防止公網倉庫變更或者小時，導致線上編譯失敗或者緊急回退失敗。
- 公司審計和安全需要
- 防止公司內部開發人員配置不當造成import path 洩露
- cache 熱點依賴， 降低公司公網出口帶寬

```sh
export GoPROXY=https://goproxy.io.direct
export GOPRIVATE=git.mycompany.com.github.com/private
```
通過 GOPROXY 加速公有庫的加載，通過 GOPRIVATE 實現繞過公司私有庫的權限認證。

# Unittest

- 測試法則：
    - 70% 小型測試
    - 20% 中型測試
    - 10% 大型測試

"自動化實現的，用於驗證一個單位函數或者獨立功能模塊的代碼是否按照預期工作，著重於典型功能性錯誤，數據損壞，錯誤條件和大小差一(off-by-one)"

- 單元測試的基礎要求：
    - 快速
    - 環境一致
    - 任意順序
    - 並行 

- 基於 docker-compose 實現跨平台跨域語言環境的容器依賴管理方案， 以解決運行unitest場景下的（mysql, redis, mc）容器依賴問題：
    - 本地安裝Docker
    - 無侵入式的環境初始化
    - 快速重置環境
    - 隨時隨地的運行
    - 語意式API聲明資源
    - 真實外不依賴，而非 in-process 模擬

利用go官方提供的：Subtests + Gomock 完成整個單元測試。