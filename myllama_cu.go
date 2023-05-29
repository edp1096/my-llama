//go:build windows && cuda
// +build windows,cuda

package myllama

// #cgo LDFLAGS: -static -L. -lstdc++ -lllama_cu -lbinding_cu
import "C"
