#include <iostream>
#include <unordered_map>

enum Type {
    PROTOTYPE_1 = 0,
    PROTOTYPE_2
};


class Prototype {
protected:
    std::string prototype_name_;
    float prototype_field_;
public:
    Prototype() {}
    Prototype(std::string prototype_name) : prototype_name_(prototype_name) {}
    virtual ~Prototype() {}
    virtual Prototype* Clone() const = 0;
    virtual void Method(float prototype_field) {
        this->prototype_field_ = prototype_field;
        std::cout << "Call Method from " << prototype_name_ << " with field : " << prototype_field_ << std::endl;
    }
};


/*
 * ConcretePrototype1 是 Prototype 的子类，实现了 Clone 的方法。
 * 在这个示例中所有 Prototype 的数据成员都是在 stack 上
 * 如果类中含有指针类型，那么需要在 Copy-Constructor 上实现自己的深拷贝
 * */
class ConcretePrototype1 : public Prototype {
private:
    float concrete_prototype_field1_;
public:
    ConcretePrototype1(std::string prototype_name, float concrete_prototype_field) :
        Prototype(prototype_name), concrete_prototype_field1_(concrete_prototype_field) {}

    Prototype* Clone() const override {
        // 因为数据都在 stack 上，可以使用默认的拷贝构造函数，如果还有指针，需要自己实现拷贝构造函数进行深拷贝
        // 注意：这里返回的是一个新的对象，所以调用的 client 代码需要负责释放之前的内存
        // 或者这里这里可以改成智能指针 unique_pointer
        // 这里利用的是默认拷贝构造函数，只完成浅拷贝
        return new ConcretePrototype1(*this);
    }
};

class ConcretePrototype2 : public Prototype {
private:
    float concrete_prototype_field2_;
public:
    ConcretePrototype2(std::string prototype_name, float concrete_prototype_field) :
        Prototype(prototype_name), concrete_prototype_field2_(concrete_prototype_field) {}

    Prototype* Clone() const override {
        return new ConcretePrototype2(*this);
    }
};

class PrototypeFactory {
private:
    std::unordered_map<Type, Prototype*, std::hash<int>> prototype_;
public:
    PrototypeFactory() {
        // 创建后就初始化一个对象，之后都从这两个初始化的对象拷贝生成新对象
        prototype_[Type::PROTOTYPE_1] = new ConcretePrototype1("PROTOTYPE_1 ", 50.f);
        prototype_[Type::PROTOTYPE_2] = new ConcretePrototype1("PROTOTYPE_2 ", 60.f);
    }
    ~PrototypeFactory() {
        delete prototype_[Type::PROTOTYPE_1];
        delete prototype_[Type::PROTOTYPE_2];
    }
    
    Prototype *CreatePrototype(Type type) {
        return prototype_[type]->Clone();
    }
};

void Client(PrototypeFactory &prototype_factory) {
    std::cout << "Let's create a Prototype 1\n";
    Prototype *prototype = prototype_factory.CreatePrototype(Type::PROTOTYPE_1);
    prototype->Method(90);
    delete prototype;

    std::cout << "\n";

    std::cout << "Let's create a Prototype 2\n";
    prototype = prototype_factory.CreatePrototype(Type::PROTOTYPE_2);
    prototype->Method(10);
    delete prototype;
}

int main() {
    PrototypeFactory *prototype_factory = new PrototypeFactory();
    Client(*prototype_factory);
    delete prototype_factory;
    return 0;
}

