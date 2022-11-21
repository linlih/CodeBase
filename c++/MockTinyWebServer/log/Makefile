CXX ?= g++

test_log: test_log.cpp log.cpp log.h 
	$(CXX) -o test_log  $^ $(CXXFLAGS) -lpthread -Wformat-truncation
 
test_bq: test_block_queue.cpp ../lock/locker.h block_queue.h 
	$(CXX) -o test_bq  $^ $(CXXFLAGS) -lpthread

clean:
	rm test_bq