#include <iostream>
#include <algorithm>

// 这个和AdapterConceptualExample.cc没有太大区别，这里是用继承来实现
// AdapterConceptualExample 是需要将 adaptee 传给 adapter 而已

class Target {
public:
    virtual ~Target() = default;
    virtual std::string Request() const {
        return "Target: the default target's behavior.";
    }
};

class Adaptee {
public:
    std::string SpecificRequest() const {
        return ".eetpadA eht roivaheb laicepS";
    }
};

class Adapter : public Target, public Adaptee {
public:
    Adapter() {}
    std::string Request() const override {
        std::string to_reverse = SpecificRequest();
        std::reverse(to_reverse.begin(), to_reverse.end());
        return "Adapter: (TRANSLATED) " + to_reverse;
    }
};

void ClientCode(const Target *target) {
    std::cout << target->Request();
}

int main() {

    std::cout << "Client: I can work with just fine with the Target object:\n";
    Target *target = new Target;
    ClientCode(target);
    std::cout << "\n\n";

    Adaptee *adaptee = new Adaptee;
    std::cout << "Client: The Adaptee class has a weird interface. See, I don't understand it:\n";
    std::cout << "Adaptee: " << adaptee->SpecificRequest();
    std::cout << "\n\n";
    std::cout << "Client: But I can work with it via the Adapter:\n";
    Adapter *adapter = new Adapter;
    ClientCode(adapter);

    std::cout << "\n";
    delete target;
    delete adaptee;
    delete adapter;
 
    return 0;
}
