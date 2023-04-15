package cgollama

// #cgo CXXFLAGS: -I../llama.cpp/examples -I../llama.cpp
// #cgo LDFLAGS: -L../ -lbinding -lm -lstdc++ -static
// #include "binding.h"
import "C"
import (
	"fmt"
	"unicode"
	"unsafe"

	ws "golang.org/x/net/websocket"
)

type LLama struct {
	Container   unsafe.Pointer
	State       unsafe.Pointer
	PredictStop chan bool
}

func New(modelFNAME string, opts ...ModelOption) (*LLama, error) {
	// mo := NewModelOptions(opts...)
	// modelPath := C.CString(modelFNAME)

	return &LLama{}, nil
}

func isASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}

	return true
}

func (l *LLama) GetInitialParams(text string, prompt string, antiprompt string, po PredictOptions) (unsafe.Pointer, int) {

	// input := C.CString(text)
	// reversePrompt := C.CString(antiprompt)

	siba := C.llama_init_container()

	remainCOUNT := 0

	return siba, remainCOUNT
}

func (l *LLama) GetContinueParams(text string, antiprompt string, params unsafe.Pointer, predVARs unsafe.Pointer, po PredictOptions) (unsafe.Pointer, unsafe.Pointer, int) {
	// input := C.CString(text)
	// reversePrompt := C.CString(antiprompt)

	remainCOUNT := 0

	return params, predVARs, remainCOUNT
}

func (l *LLama) Predict(conn *ws.Conn, handler ws.Codec, input string, siba unsafe.Pointer, remainCOUNT int) error {
	// responseBYTEs := []byte{}
	response := ""

	// promptLenth := len(input)
	// isFinish := false

END:
	for remainCOUNT > 0 {
		break END
	}

	err := handler.Send(conn, response+"\n$$__RESPONSE_DONE__$$\n")
	if err != nil {
		fmt.Println("Send error:", err)
		return err
	}

	return nil
}
