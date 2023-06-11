//go:build windows && clblast
// +build windows,clblast

package myllama

// #cgo CXXFLAGS: -DGGML_USE_CLBLAST -Ivendors/llama.cpp -Ivendors/llama.cpp/examples
// #cgo LDFLAGS: -static -L. -lstdc++ -lllama_cl -lbinding_cl
import "C"
