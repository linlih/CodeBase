# Go优先级队列的实现

参考自文章：Go学堂 https://mp.weixin.qq.com/s/cUyNrUmv7cjByCglr094WQ

这篇文章给出了一个优先级队列任务消费模型的实现，思想是没有问题的，代码上由于提供不完整，有很多错误。

同时缺少了很多细节，如果关闭WorkerManager会导致死锁，没有取消所有的Job导致的。

需要考虑的问题有：
1. 当worker manager停止的时候，需要关闭所有的worker协程，同时取消所有的Job
2. 如果有新加入的Job那么应该在worker已经停止的时候返回错误

这个实现还是有些不足，需要再思考思考。
