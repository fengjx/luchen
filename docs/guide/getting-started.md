# 快速开始

编写一个http协议的helloworld程序 

<<< @/snippets/hello.go

启动服务
```bash
$ go run main.go
```

测试
```bash
curl http://localhost:8080/say-hello\?name\=foo 
hello: foo
```

你可能会认为代码过于复杂，但是，根据过往大型项目的实践来看，必要的代码分层对于多人协作开发的项目至关重要，可以保持代码的可维护和可扩展性，这对于长期维护的项目收益巨大。

[参考源码](https://github.com/fengjx/luchen/blob/master/_example/helloworld/main.go)
