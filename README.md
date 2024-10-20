# Clash Rule Rewrite
固定自己的代理规则，避免被机场订阅覆盖，同时保证节点更新。支持在线编辑备份、还原。

# 运行
```
go mod tidy
go run .
```

# 使用

## 访问 `http://127.0.0.1:8080/` 进行在线编辑

[![image.png](http://www.cdnjson.com/images/2024/10/21/image.png)](http://www.cdnjson.com/image/a89bg)

## 转换订阅 `http://127.0.0.1:8080/api/process?url=你的订阅地址` 
