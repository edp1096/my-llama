//go:build cuda
// +build cuda

package myllama

// #cgo LDFLAGS: -static -L. -lstdc++ -lllama_cu -lbinding_cu
import "C"
