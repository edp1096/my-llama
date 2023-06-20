//go:build linux && cuda
// +build linux,cuda

package myllama

// #cgo CXXFLAGS: -DGGML_USE_CUBLAS -Illama.cpp -Illama.cpp/examples -I/usr/local/cuda/include
// #cgo LDFLAGS: -L. -L/usr/local/cuda/lib64 -lstdc++ -lllama_cu_lin64 -lmyllama_cu_lin64 -lcudart -lcublas -lcublasLt
import "C"
