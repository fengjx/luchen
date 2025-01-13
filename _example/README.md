# luchen 开发示例

## helloworld 简单示例

根据 proto 生成代码
```bash
lc pbgen -f helloworld/pb/greet.proto
```

运行服务
```bash
go run helloworld/main.go
```

## registrar 服务注册&发现

根据 proto 生成代码
```bash
lc pbgen -f registrar/pb/greet.proto
```

运行服务
```bash
go run registrar/main.go
```

