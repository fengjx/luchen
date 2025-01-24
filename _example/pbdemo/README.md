# pbdemo

基于 proto 协议开发接口示例

## 根据 proto 生成接口实现

```bash
lc pbgen -f pbdemo/pb/greet.proto
```

## 启动服务

```bash
go run main.go
```

## 测试接口

```bash
curl -i -X POST 'http://localhost:8080/say-hello' \
  -H 'Content-Type: application/json' \
  -d '{
    "name": "fengjx" 
  }'
```
