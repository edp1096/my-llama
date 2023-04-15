package cgollama

// #cgo CXXFLAGS: -I../llama.cpp/examples -I../llama.cpp
// #cgo LDFLAGS: -L../ -lbinding -lm -lstdc++ -static
// #include "binding.h"
import "C"
import (
	"fmt"
	"unicode/utf8"
	"unsafe"

	ws "golang.org/x/net/websocket"
)

type LLama struct {
	Container   unsafe.Pointer
	State       unsafe.Pointer
	PredictStop chan bool
}

func New(modelFNAME string) (*LLama, error) {
	container := C.llama_init_container()
	if container == nil {
		return nil, fmt.Errorf("failed to initialize the container")
	}

	return &LLama{Container: container}, nil
}

func (l *LLama) LoadModel(modelFNAME string) error {
	container := l.Container
	C.llama_set_model_path(container, C.CString(modelFNAME))

	result := bool(C.llama_load_model(container))
	if !result {
		return fmt.Errorf("failed to load the model")
	}

	return nil
}

func (l *LLama) InitParams() error {
	container := l.Container
	result := bool(C.llama_init_params(container))

	if !result {
		return fmt.Errorf("failed to initialize the parameters")
	}

	return nil
}

func (l *LLama) SetupParams() {
	container := l.Container
	C.llama_setup_params(container)
}

func (l *LLama) GetRemainCount() int {
	container := l.Container
	return int(C.llama_get_n_remain(container))
}

func (l *LLama) Predict(conn *ws.Conn, handler ws.Codec, input string) error {
	container := l.Container
	remainCOUNT := int(C.llama_get_n_remain(container))

	responseBYTEs := []byte{}
	response := ""

	// promptLenth := len(input)
	// isFinish := false

END:
	for remainCOUNT > 0 {
		embdSIZE := int(C.llama_get_embd_size(container))

		for i := 0; i < embdSIZE; i++ {
			select {
			case <-l.PredictStop:
				break END
			default:
				id := C.llama_get_embed_id(container, C.int(i))
				embedCSTR := C.llama_get_embed_string(container, id)
				embedSTR := C.GoString(embedCSTR)
				fmt.Print(embedSTR)

				responseBYTEs = append(responseBYTEs, []byte(embedSTR)...)
				response = string(responseBYTEs)

				if !utf8.ValidString(embedSTR) {
					continue // Because connection is closed, don't send invalid UTF-8
				}

				err := handler.Send(conn, response)
				if err != nil {
					fmt.Println("Send error:", err)
					break END
				}
			}
		}

		remainCOUNT = int(C.llama_get_n_remain(container))
	}

	err := handler.Send(conn, response+"\n$$__RESPONSE_DONE__$$\n")
	if err != nil {
		fmt.Println("Send error:", err)
		return err
	}

	return nil
}
