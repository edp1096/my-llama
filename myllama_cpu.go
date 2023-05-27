//go:build !clblast && !cuda
// +build !clblast,!cuda

package myllama

// #cgo LDFLAGS: -static -L. -lstdc++ -lllama -lbinding
import "C"
