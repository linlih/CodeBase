#include <iostream>
#include "block_queue.h"

using namespace std;

int main() {
    block_queue<int> bq(2);

    bq.push(1);
    bq.push(5);
    
    int a1, a2, a3;
    bq.pop(a1);
    bq.pop(a2);
    cout << a1 << " " << a2 << endl;
    bq.pop(a3);
    cout << a3 << endl;
    return 0;
}