//go:build windows && !clblast && !cuda
// +build windows,!clblast,!cuda

package myllama

// #cgo LDFLAGS: -static -L. -lstdc++ -lllama -lmyllama
import "C"
