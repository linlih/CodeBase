#include <iostream>

/*
 * 总结下：这个文件是AbstractFactory的示例代码
 * 定义了两种抽象类：
 * 一个是工厂抽象类，提供给客户代码调用这里面的接口去使用
 * 一个是产品抽象类，定义了工厂能够生产的产品抽象
 * 实现流程：
 * 1. 先定义好产品抽象类
 * 2. 实现产品抽象类中的方法去定义具体的产品
 * 3. 定义好工厂抽象类
 * 4. 实现工厂抽象类中的方法去生成具体的产品
 * 5. 创建相应的工厂，客户调用工厂抽象类的方法去生产产品
 * */


/*
 * 所有的产品都应该有一个基类，具体的产品都 应该继承这个基类
 * */
class AbstractProductA {
public:
    virtual ~AbstractProductA(){};
    virtual std::string UsefulFunctionA() const = 0;
};

/*
 * 具体的产品继承产品基类，实现基类中的方法
 * */
class ConcreteProductA1 : public AbstractProductA {
public:
    std::string UsefulFunctionA() const override {
        return "The result of product A1.";
    }
};

class ConcreteProductA2 : public AbstractProductA {
    std::string UsefulFunctionA() const override {
        return "The result of product A2.";
    }
};

/*
 * 这是另外一个产品的抽象 interface，每个产品之间是可以交互的，但是必须是相同的具体对象
 * 比如这里传入的协作产品是 AbstarctProduct 的具体对象
 * */
class AbstractProductB {
public:
    virtual ~AbstractProductB(){};
    virtual std::string UsefulFunctionB() const = 0;
    virtual std::string AnotherUsefulFunctionB(const AbstractProductA &collaborator) const = 0;
};

class ConcreteProductB1 : public AbstractProductB {
public:
    std::string UsefulFunctionB() const override {
        return "The result of the product of B1.";
    }
    std::string AnotherUsefulFunctionB(const AbstractProductA &collaborator) const override {
        const std::string result = collaborator.UsefulFunctionA();
        return "The result of the B1 collaborating with ( " + result + " ).";
    }
};

class ConcreteProductB2 : public AbstractProductB {
public:
    std::string UsefulFunctionB() const override {
        return "The result of the product of B2.";
    }
    std::string AnotherUsefulFunctionB(const AbstractProductA &collaborator) const override {
        const std::string result = collaborator.UsefulFunctionA();
        return "The result of the B2 collaborating with ( " + result + " ).";
    }
};

/*
 * 抽象工厂interface，定义一组方法返回不同的抽象类产品
 * 各个工厂之间可以协作，但是不同的产品是不兼容的
 * */
class AbstractFactory {
public:
    virtual AbstractProductA *CreateProductA() const = 0;
    virtual AbstractProductB *CreateProductB() const = 0;
};

/*
 * 实现具体的工厂去生产工厂能够生产的具体产品
 */
class ConcreteFactory1 : public AbstractFactory {
    AbstractProductA *CreateProductA() const override {
        return new ConcreteProductA1();
    }
    AbstractProductB *CreateProductB() const override {
        return new ConcreteProductB1();
    }
};

class ConcreteFactory2 : public AbstractFactory {
    AbstractProductA *CreateProductA() const override {
        return new ConcreteProductA2();
    }
    AbstractProductB *CreateProductB() const override {
        return new ConcreteProductB2();
    }
};


/*
 * 客户代码就可以通过抽象类的方法去生产相应的产品
 */
void clientCode(const AbstractFactory &factory) {
    const AbstractProductA *product_a = factory.CreateProductA();
    const AbstractProductB *product_b = factory.CreateProductB();
    std::cout << product_b->UsefulFunctionB() << std::endl;
    std::cout << product_b->AnotherUsefulFunctionB(*product_a) << std::endl;
    delete product_a;
    delete product_b;
}

int main() {
    std::cout << "Client: Testing client code with the first factory type:\n";
    ConcreteFactory1 *f1 = new ConcreteFactory1();
    clientCode(*f1);
    delete f1;
    std::cout << std::endl;
    std::cout << "Client: Testing client code with the first factory type:\n";
    ConcreteFactory2 *f2 = new ConcreteFactory2();
    clientCode(*f2);
    delete f2;
    return 0;
}



