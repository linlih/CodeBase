# Builder 模式

`iBuilder.go` 中实现了 `iBuilder` 的接口，和获取不同类别 `Builder` 的方法。

`normalBuilder.go` 和 `iglooBuilder.go` 中实现了具体的两种不同建造方式，实现了 `iBuilder` 接口中的方法。

`house.go` 定义的是建造的产品。

`director.go` 实现的是 `Director` 的角色，由它来调用 `iBuilder` 提供的创建步骤去建造产品。

`main.go` 中就是客户的代码了，创建不同的 `builder` 对象，然后交给 `director` 去建造。

# 如何应用
1. 首先需要确定所有产品是可以通过相同的建造步骤生产出来的，否则不能用这个模式。
2. 将这些相同的步骤定义成 base builder interface
3. 实现具体的 builder，实现 base builder interface 中的建造步骤。同时要注意，这个实现中需要包含一个返回建造结果的函数。比如 `normalBuilder.go` 中的 `gethouse` 函数。这个函数不应该定义在 `iBuilder` 中的原因是：各个构造器可能构造没有公共接口的产品。如果能确定这个产品是来自同一个类型，那可以写在`iBuilder`中。
4. 创建一个 `director` 类， 封装使用相同的 builder 对象来构建 product
5. 客户代码需要创建 builder 和 director，把 builder 对象传给 director，由 director 去建造。
6. 客户从 builder 获取建造的结果

# 和 Abstract Factory 有什么区别？
`Builder` 更关注的是一步一步构造一个复杂的对象，而 `Abstract Factory` 更关注创建一组相关的对象。`Abstract Factory`返回的就是产品了，而`Builder`可以让你在产品生产之前做一些额外的操作。

代码参考自：https://refactoring.guru/design-patterns/builder/go/example
