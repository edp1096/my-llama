CC = gcc

#
# Compile flags
#

# keep standard at C11 and C++11
CFLAGS   = -I./llama.cpp -I. -O3 -DNDEBUG -std=c11 -fPIC
CXXFLAGS = -I./llama.cpp -I. -I./llama.cpp/examples -I./examples -O3 -DNDEBUG -std=c++11 -fPIC
LDFLAGS  =

# warnings
CFLAGS   += -Wall -Wextra -Wpedantic -Wcast-qual -Wdouble-promotion -Wshadow -Wstrict-prototypes -Wpointer-arith -Wno-unused-function
CXXFLAGS += -Wall -Wextra -Wpedantic -Wcast-qual -Wno-unused-function

# Architecture specific
# TODO: probably these flags need to be tweaked on some architectures
#       feel free to update the Makefile for your architecture and send a pull request or issue
ifeq ($(UNAME_M),$(filter $(UNAME_M),x86_64 i686))
# Use all CPU extensions that are available:
	CFLAGS += -march=native -mtune=native
endif

#
# Print build information
#

$(info I llama.cpp build info: )
$(info I CFLAGS:   $(CFLAGS))
$(info I CXXFLAGS: $(CXXFLAGS))
$(info I LDFLAGS:  $(LDFLAGS))
$(info )

build:
	$(MAKE) libbinding.a libllama.a
#	go env -w CGO_LDFLAGS="-O2 -g $(CGO_LDFLAGS)"
	go build -a -trimpath -ldflags="-w -s" -o bin/
#	go env -w CGO_LDFLAGS="-O2 -g"
# ifdef USE_CLBLAST
# 	cp openclblast/lib/clblast.dll bin/
# endif

llama.cpp/ggml.o:
#	$(MAKE) CC=$(CC) CFLAGS+='$(CFLAGS_ADD)' -C llama.cpp ggml.o
	$(MAKE) CC=$(CC) -C llama.cpp ggml.o

llama.cpp/llama.o:
	$(MAKE) CC=$(CC) -C llama.cpp llama.o

llama.cpp/common.o:
	$(MAKE) CC=$(CC) -C llama.cpp common.o

binding.o: llama.cpp/ggml.o llama.cpp/llama.o llama.cpp/common.o $(OBJS)
#	$(CXX) $(CXXFLAGS) -static $(LDFLAGS_ADD) $(CFLAGS_ADD) -I./llama.cpp -I./llama.cpp/examples cgollama/binding.cpp -o cgollama/binding.o -c $(LDFLAGS)
	$(CXX) $(CXXFLAGS) -static -I./llama.cpp -I./llama.cpp/examples cgollama/binding.cpp -o cgollama/binding.o -c $(LDFLAGS)

libllama.a: llama.cpp/ggml.o llama.cpp/common.o llama.cpp/llama.o $(OBJS)
	ar src libllama.a llama.cpp/ggml.o llama.cpp/common.o llama.cpp/llama.o $(OBJS)

libbinding.a: binding.o
	ar src libbinding.a cgollama/binding.o


build_for_cuda:
	$(MAKE) libbinding.a_for_cuda
	go build -trimpath -ldflags="-w -s" -o bin/

libbinding.a_for_cuda:
	$(CXX) $(CXXFLAGS) -I./llama.cpp -I./llama.cpp/examples cgollama/binding.cpp -o cgollama/binding.o -c $(LDFLAGS)
	ar src libbinding.a cgollama/binding.o


clean:
	rm -rf *.o
	rm -rf *.a
	rm -rf cgollama/*.o
	rm -rf cgollama/*.a
	$(MAKE) -C llama.cpp clean