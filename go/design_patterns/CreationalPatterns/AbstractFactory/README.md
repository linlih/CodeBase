# 抽象工厂方法

抽象出一个创建工厂的方法，客户也就是这里的 main.go 都是与抽象类定义的接口交互，也就是 iSportsFacotry.go 中定义的接口.

抽象工厂类通过参数来创建不同的工厂，比如在 getSportsFactory 中传入的 adidas 和 nike。工厂创建完毕后，调用工厂提供的方法去生产产品。

同时对产品也进行了抽象，iShoe 和 iShirt，在 go 中没有继承的概念，使用的是结构体组合，每个品牌自己去定义自己产品 shoe 和 shirt。

代码参考自：https://refactoring.guru/design-patterns/abstract-factory/go/example
