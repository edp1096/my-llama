//go:build linux && clblast
// +build linux,clblast

package myllama

// #cgo CXXFLAGS: -DGGML_USE_CLBLAST -Illama.cpp -Illama.cpp/examples
// #cgo LDFLAGS: -L. -lstdc++ -lllama_cl_lin64 -lmyllama_cl_lin64 -lclblast_lin64 -lOpenCL_lin64
import "C"
