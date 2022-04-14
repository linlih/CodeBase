# NULL Object

空对象模式，因为有时候我们创建失败的时候，需要返回相应的空内容，但是由于空的指针或者其他空的形式，不好统一判断。这个时候就可以创建一个NULL Object，这个Object是继承Real Object相同的接口，同时实现了相应的接口，那么在对于实际抽象出来的接口实现的对象就可以进行统一的操作，而无需担心是否出现了空对象指针的问题。

参考代码：https://www.geeksforgeeks.org/null-object-design-pattern/





