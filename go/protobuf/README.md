# protobuf

protobuf google提出的关于序列化结构体的一种数据格式，类似于xml，json这种格式。

```bash
go get google.golang.org/protobuf/cmd/protoc-gen-go
```

生成 proto 的 go 代码，使用以下命令：
```bash
protoc --go_out=. todo.proto
```

程序测试方法：
```bash
go run main.go add test
```

然后会生成 `mydb.pd`文件，里面存储的是序列化的结果。

执行以下命令可以列出 pd 里面的内容：
```bash
go run main.go list
```

使用流程上来总结：
1. 编写*.proto文件
2. 使用protoc命令生成对应的语言的代码文件，protoc支持的有go、c++、php、java、csharp、objc、python、ruby
3. 调用生成的代码
4. 使用proto.Marshal、proto.UnMarshal进行序列化和反序列化

db 文件也可以用 protoc 命令来进行解码(但是注意，这个代码对编码内容做了调整，加入了每次编码的长度信息，所以并不是一个纯 pd 的文件，所以用这个命令解码是会不成功的，这个命令用于解码proto.Marshal的原始内容)：
```bash
cat mydb.pd | protoc --decode_raw
```

代码来自 Youtube： https://www.youtube.com/watch?v=_jQ3i_fyqGA