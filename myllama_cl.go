//go:build windows && clblast
// +build windows,clblast

package myllama

// #cgo LDFLAGS: -static -L. -lstdc++ -lllama_cl -lbinding_cl
import "C"
