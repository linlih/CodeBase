# Prototype 原型

这个例子是参考了操作系统中的文件系统，每个文件和文件夹用`inode`接口来表示，同时`inode`还拥有`clone`方法。

`file`和`folder`都实现了`inode`接口，在`clone`方法中复制的`name`中都加上了`clone`字样。




