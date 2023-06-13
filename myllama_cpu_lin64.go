//go:build linux && !clblast && !cuda
// +build linux,!clblast,!cuda

package myllama

// #cgo LDFLAGS: -static -L. -lstdc++ -lllama_lin64 -lbinding_lin64
import "C"
