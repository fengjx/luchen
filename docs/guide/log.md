# 应用日志

## 默认情况

1. 本地开发环境日志会输出到 stdout
2. 非开发环境日志文件默认路径：./logs，即：程序当前运行路径的 logs 目录下
3. 环境变量`LUCHEN_LOG_DIR`可以设置日志输出目录

## 方法说明

```go
// RootLogger 返回默认 logger
func RootLogger() logger.Logger

// Logger 从 context 获得 logger
func Logger(ctx context.Context) logger.Logger

// WithLogger context 注入 logger
func WithLogger(ctx context.Context, logger logger.Logger)
```

内部是对zap的一层包装，更多使用细节参考：[zap](https://github.com/uber-go/zap)
