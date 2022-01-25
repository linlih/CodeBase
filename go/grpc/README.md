# grpc

示例 grpc 代码 : 2022.01.25

grpc 重点的数据结构基于 protobuf 进行编码的，所以使用 grpc 之前需要先了解 protobuf 的使用

首先使用 protoc 生成 pd 的 go 代码，如果出现 grpc 生成程序无法找到的情况，就需要安装下 grpc 的生成工具：
```bash
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```
生成命令（参考的 just for fun Youtube视频中的命令已经失效了，修改为以下命令），其中数据序列化的 API 代码和 gRPC 的代码分开在两个文件了：
```bash
protoc todo.proto --go-grpc_out=./ --go_out=./
```

这个示例代码的逻辑是这样的，分为客户端和服务端，提供两个功能，一个是 add 将字符串序列化到 pd 文件中，一个是 list 将 pd 文件中的字符串反序列化出来返回给客户端并打印出来。

代码编写流程：

1. 编写 proto 文件，生成 *pd.go 和 *grpc.pd.go 文件
2. 编写服务端
   1. 创建 grpc server
   2. 注册 *grpc.pd.go 中的server，传入 grpc server 和相应的服务处理函数
   3. 启动 grpc server 服务
3. 编写客户端
   1. 使用 grpc.Dial 函数连接远程 grpc 服务
   2. 创建 *grpc.pd.go 中的client
   3. 调用 client 中的 rpc 函数，获取远程 rpc 服务结果

代码来自视频： https://www.youtube.com/watch?v=uolTUtioIrc
