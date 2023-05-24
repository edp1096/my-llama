//go:build clblast
// +build clblast

package myllama

// #cgo LDFLAGS: -static -L. -lstdc++ -lllama_cl -lbinding_cl
import "C"
