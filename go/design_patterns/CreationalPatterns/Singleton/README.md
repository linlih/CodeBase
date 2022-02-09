# Singleton 单例模式

单例模式比较好理解，就是程序中是生成一个对象，其他线程使用的时候都使用这个唯一的对象。

go 的示例中给了两个，一个是`Conceptual Example`，使用锁的方式来实现；一个是`Another Example`，使用的是 `sync.Once` 来实现的。

代码参考自：https://refactoring.guru/design-patterns/singleton/go/example