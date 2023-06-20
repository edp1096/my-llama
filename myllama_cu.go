//go:build windows && cuda
// +build windows,cuda

package myllama

// #cgo CXXFLAGS: -DGGML_USE_CUBLAS -Illama.cpp -Illama.cpp/examples
// #cgo LDFLAGS: -static -L. -lstdc++ -lllama_cu -lmyllama_cu
import "C"
