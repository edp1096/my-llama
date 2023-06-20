//go:build windows && clblast
// +build windows,clblast

package myllama

// #cgo CXXFLAGS: -DGGML_USE_CLBLAST -Illama.cpp -Illama.cpp/examples
// #cgo LDFLAGS: -static -L. -lstdc++ -lllama_cl -lmyllama_cl
import "C"
