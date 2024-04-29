# 应用日志

## 默认情况

1. 本地开发环境日志会输出到 stdout
2. 非开发环境日志文件默认路径：./logs，即：程序当前运行路径的 logs 目录下
3. 环境变量`LUCHEN_LOG_DIR`可以设置日志输出目录

## 方法说明

```go
// GetLogger 从上下文获取当前 logger
func GetLogger(ctx context.Context) logger.Logger

// WithLogger context 注入 logger
func WithLogger(ctx context.Context, fields ...zap.Field) context.Context

// GetLogDir 返回日志路径
func GetLogDir() string
```

## 日志打印

日志打印 api 

```go
func Debug(msg string, fields ...zap.Field)
func DebugCtx(ctx context.Context, msg string, fields ...zap.Field)
func Debugf(format string, args ...interface{})
func DebugfCtx(ctx context.Context, format string, args ...interface{})
func Error(msg string, fields ...zap.Field)
func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field)
func Errorf(format string, args ...interface{})
func ErrorfCtx(ctx context.Context, format string, args ...interface{})
func Fatal(msg string, fields ...zap.Field)
func FatalCtx(ctx context.Context, msg string, fields ...zap.Field)
func Fatalf(format string, args ...interface{})
func FatalfCtx(ctx context.Context, format string, args ...interface{})
func Info(msg string, fields ...zap.Field)
func InfoCtx(ctx context.Context, msg string, fields ...zap.Field)
func Infof(format string, args ...interface{})
func InfofCtx(ctx context.Context, format string, args ...interface{})
func NewKitLogger(name string, level logger.Level) kitlog.Logger
func Panic(msg string, fields ...zap.Field)
func PanicCtx(ctx context.Context, msg string, fields ...zap.Field)
func Panicf(format string, args ...interface{})
func PanicfCtx(ctx context.Context, format string, args ...interface{})
func Warn(msg string, fields ...zap.Field)
func WarnCtx(ctx context.Context, msg string, fields ...zap.Field)
func Warnf(format string, args ...interface{})
func WarnfCtx(ctx context.Context, format string, args ...interface{})
```
