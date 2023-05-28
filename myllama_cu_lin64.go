//go:build linux && cuda
// +build linux,cuda

package myllama

// #cgo CXXFLAGS: -Ivendors/llama.cpp -Ivendors/llama.cpp/examples -I/usr/local/cuda/include
// #cgo LDFLAGS: -L. -L/usr/local/cuda/lib64 -lstdc++ -lllama_cu_lin64 -lbinding_cu_lin64 -lcudart -lcublas -lcublasLt
import "C"
