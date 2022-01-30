#include <iostream>
#include <vector>

/*
 * 一般来说只有当这个产品是很负责的，并且需要扩展功能，才去使用 Builder 模式
 * 不同于其他的创建模式，不同的具体产品生产者可以可以生成不相关的产品。也就是说不同的建造者可能不一定使用相同的接口。
 */

class Product1 {
public:
    std::vector<std::string> parts_;
    void ListParts() const {
        std::cout << "Product parts:";
        for (size_t i=0; i < parts_.size(); i++) {
            if (parts_[i] == parts_.back()) {
                std::cout << parts_[i];
            } else {
                std::cout << parts_[i] << ", ";
            }
        }
        std::cout << "\n\n";
    }
};

class Builder{
public:
    virtual ~Builder(){}
    virtual void ProducePartA() const = 0;
    virtual void ProducePartB() const = 0;
    virtual void ProducePartC() const = 0;
};

class ConcreteBuilder1 : public Builder {
private:
    Product1* product;
public:
    ConcreteBuilder1() {
        this->Reset();
    }

    ~ConcreteBuilder1() {
        delete product;
    }

    void Reset() {
        this->product = new Product1();
    }

    void ProducePartA() const override {
        this->product->parts_.push_back("PartA1");
    }

    void ProducePartB() const override {
        this->product->parts_.push_back("PartB1");
    }

    void ProducePartC() const override {
        this->product->parts_.push_back("PartC1");
    }

    /*
     * 具体的buidler需要提供自己的方法来获取建造结果，因为不同建造者可能返回不同的建造结果，所以不遵循相同的接口。
     * 所以这个返回产品的结果，不能定义在Base的Builder里面。
     * 通常返回了一个产品之后，建造者应该开始建造新的产品，所以在这里调用了Reset函数。非固定写法，也可以要求client显式调用Reset函数。
     *
     * 同时要注意这里的内存，因为在类中定义的产品是一个指针，如果返回了产品之后，就把生产的产品内存交给client，
     * 那么需要client对产品的内存进行管理。这里也可以使用Smart Pointer来避免内存泄露。
     * */
    Product1* GetProduct() {
        Product1* result = this->product;
        this->Reset();
        return result;
    }
};

class Director {
private:
    Builder* builder;
public:
    void set_builder(Builder *builder) {
        this->builder = builder;
    }

    void BuildMinimalViableProduct() {
        this->builder->ProducePartA();
    }

    void BuildFullFeatureProduct() {
        this->builder->ProducePartA();
        this->builder->ProducePartB();
        this->builder->ProducePartC();
    }
};

void ClientCode(Director& director) {
    ConcreteBuilder1* builder = new ConcreteBuilder1();
    director.set_builder(builder);
    std::cout << "Standard basic product:\n";
    director.BuildMinimalViableProduct();

    Product1* p = builder->GetProduct();
    p->ListParts();
    delete p;

    std::cout << "Standard full feature product:\n";
    director.BuildFullFeatureProduct();

    p = builder->GetProduct();
    p->ListParts();
    delete p;

    // Builder 模式不一定需要一个 Director 类
    std::cout << "Custom product:\n";
    builder->ProducePartA();
    builder->ProducePartC();
    p = builder->GetProduct();
    p->ListParts();
    delete p;

    delete builder;
}



int main() {
    Director* director = new Director();
    ClientCode(*director);
    delete director;
    return 0;
}
