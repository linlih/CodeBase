#include <iostream>
#include <chrono>
#include <thread>

// 这个示例代码不能实际使用，因为线程不安全可能导致创建多个对象，而不是单例模式
// 编译方法：g++ NaviveSingleton.cc -lpthread
// 多执行几次可以看到偶尔会输出：FOO、FOO(创建了一个对象)，其他情况输出：FOO、BAR(创建了两个对象)

/*
 * Singleton 提供了一个GetInstance方法，作为构造函数的替代让Client去获取类对象
 * */
class Singleton {
    /*
     * 单例模式的构造函数应该设置为private，不应该暴露出来
     * */
protected:
    Singleton(const std::string value): value_(value) {}
    
    static Singleton* singleton_;

    std::string value_;
public:
    // 单例模式不能被克隆
    Singleton(Singleton &other) = delete;
    // 单例模式也不能被赋值
    void operator=(const Singleton &) = delete;
    /*
     * 设置为一个静态方法来获取对象
     * 第一次执行的时候创建对象并放到静态变量singleton_中
     * 后面的调用就会直接返回这个singleton_
     * */
    static Singleton *GetInstance(const std::string &value);
    void SomeBusinessLogic() {}
    std::string value() const {
        return value_;
    }
};

Singleton* Singleton::singleton_ = nullptr;

Singleton *Singleton::GetInstance(const std::string& value) {
    /*
     * 这里的实现存在一个问题：线程不安全
     * 如果两个线程同时判断了singleton_为空，那么就会创建两个对象，不符合单例模式了
     * */
    if(singleton_ == nullptr) {
        singleton_ = new Singleton(value);
    }
    return singleton_;
}


void ThreadFoo() {
    std::this_thread::sleep_for(std::chrono::milliseconds(1000));
    Singleton* singleton = Singleton::GetInstance("FOO");
    std::cout << singleton->value() << "\n";
}

void ThreadBar() {
    std::this_thread::sleep_for(std::chrono::milliseconds(1000));
    Singleton* singleton = Singleton::GetInstance("BAR");
    std::cout << singleton->value() << "\n";
}

int main() {
    std::cout << "If you see the same value, then singleton was reused!\n" <<
                 "If you see different values, then 2 singletons were created!\n\n" <<
                 "RESULT:\n";
    std::thread t1(ThreadFoo);
    std::thread t2(ThreadBar);
    t1.join();
    t2.join();
    return 0;
}
