# Clash Rule Rewrite
固定自己的代理规则，避免被机场订阅覆盖，同时保证节点更新。支持在线编辑备份、还原。

# 运行
```
go mod tidy
go run .
```

# 使用

## 访问 `http://127.0.0.1:8080/` 进行在线编辑

![效果预览](https://cdn-fusion.imgcdn.store/i/2024/a0a8c220ad9322b6.png)

## 转换订阅 `http://127.0.0.1:8080/api/process?url=你的订阅地址` 
