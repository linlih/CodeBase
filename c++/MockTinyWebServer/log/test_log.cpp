#include <iostream>
#include <unistd.h>
#include "./log.h"

using namespace std;

int main() {
    Log::get_instance()->init("./test.log", 0, 8192, 10, 10);
    int m_close_log = 0;
    LOG_DEBUG("Test");
    LOG_INFO("info");

    sleep(1);

    return 0;
}