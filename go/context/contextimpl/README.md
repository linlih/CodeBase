# 实现 Context Package

代码来自 Youtube： https://www.youtube.com/watch?v=8M90t0KvEDY

# 实现过程

1. 首先定义 Context， 它的构成是一个包含四个方法：Deadline、Done、Err、Value 的 interface。
> Deadline 表示的是截止期限，到什么时间节点结束
> 
> Done 返回一个管道，这个context应该被关闭；这个方法用于 select 语句中。
> 
> 如果 Done 还没有被关闭，那么这个时候 Err 返回 nil。如果 Done 关闭了，则返回被关闭的原因。当 Err 已经返回一个非空的错误时，多次调用这个函数返回相同的错误。
> 
> Value 只用于请求级别的数据。比如说请求过程中用户的 uuid 等。要注意这个 key 需要这次相等比较的。实现中使用了反射来判断，reflect.TeOf(key).Comparable()。这个 key 的定义也最好是不允许导出的（小写字母开头），避免冲突

2. 首先定义一个 emptyCtx， 这里就有一个疑问了，为什么 emptyCtx 定义成一个 int，而不是 struct{}呢？ 解答：因为要保证 Background 和 Todo 不是同一个类型，如果使用 struct{}，返回空的struct，那么他们两个是相等的，如果返回&struct{}，他们也是一样的。所以将其定义为为int，那么返回 new(int) 的时候指向的就是不同的空间，这两个Backround 和 Todo 就不会相等了。

   emptyCtx不会被取消，没有 value，也没有 deadline，只是定义为了 Context 的四个接口
 
   基于 emptyCtx 定义了两个 Context， 一个是 background，另外一个是 todo；其中 background 通常是用在了 main 函数， init 函数以及测试中，作为输入请求的 top-level context；todo 的话是用在还不确定需要用到 Context 的地方。
 
4. 接下来定义的是 WithCancel，输入的话是 parent 的 Context，返回两个东西，一个是Context，一个是取消函数。同时这里的核心是需要启动一个 gorountine 去执行 select 轮询，如果父节点 Done 了，那么取消该 goroutine，或者是该 context 自己结束了退出。
 
    在 WithCancel 的定义中要注意，因为 err 可能出现并发读写导致 data race，所以 err 的读写需要加锁。

5. 然后定义的 WithDeadline，需要实现下 Deadline 的方法，WithDeadline 的时候使用 time.After 来实现超时控制，超时后执行 cancel 函数即可。同时为了回收资源，再返回取消函数的时候，需要该函数也能把 time.After 的定时器取消掉。

6. 定义 WithTimeout，这个比较简单，只要使用 WithDeadline 去定义就可以了。

7. 最后是定义 WithValue，需要实现下 Value 方法；在 WithValue 的方法中需要判断 key 是否满足可比较即可。

以上就是 Youtube 中参考的实现逻辑，相比与官方的实现版本而言，简化了一些细节上的考虑，但是功能逻辑都是遵循整个 Context 的设计思路的，值得学习。




