package cgollama

// #cgo CXXFLAGS: -I../llama.cpp/examples -I../llama.cpp
// #cgo LDFLAGS: -L../ -static -lstdc++ -lbinding -lllama
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

func New() (*LLama, error) {
	container := C.bd_init_container()
	if container == nil {
		return nil, fmt.Errorf("failed to initialize the container")
	}

	return &LLama{Container: container}, nil
}

func (l *LLama) LoadModel(modelFNAME string) error {
	container := l.Container
	C.bd_set_model_path(container, C.CString(modelFNAME))

	result := bool(C.bd_load_model(container))
	if !result {
		return fmt.Errorf("failed to load the model")
	}

	return nil
}

func (l *LLama) GetRemainCount() int {
	container := l.Container
	return int(C.bd_get_n_remain(container))
}

func (l *LLama) SetIsInteracting(isInteracting bool) {
	container := l.Container
	C.bd_set_is_interacting(container, C.bool(isInteracting))
}

func (l *LLama) Predict(conn *ws.Conn, handler ws.Codec) error {
	container := l.Container
	remainCOUNT := int(C.bd_get_n_remain(container))

	responseBYTEs := []byte{}
	response := ""

END:
	for remainCOUNT != 0 {
		ok := bool(C.bd_predict_tokens(container))
		if !ok {
			return fmt.Errorf("failed to predict the tokens")
		}

		// display text
		embdSIZE := int(C.bd_get_embd_size(container))
		for i := 0; i < embdSIZE; i++ {
			select {
			case <-l.PredictStop:
				remainCOUNT = 0
				break END
			default:
				id := C.bd_get_embed_id(container, C.int(i))
				embedCSTR := C.bd_get_embed_string(container, id)
				embedSTR := C.GoString(embedCSTR)
				// fmt.Print(embedSTR)

				responseBYTEs = append(responseBYTEs, []byte(embedSTR)...)
				response = string(responseBYTEs)

				if !utf8.ValidString(embedSTR) {
					continue // Because connection is closed, don't send invalid UTF-8
				}

				err := handler.Send(conn, "$$__RESPONSE_PREDICT__$$\n$$__SEPARATOR__$$\n"+response)
				if err != nil {
					fmt.Println("Send error:", err)
					remainCOUNT = 0
					break END
				}
			}
		}

		ok = bool(C.bd_check_prompt_or_continue(container))
		if !ok {
			// remainCOUNT = 0
			break
		}
		C.bd_dropback_user_input(container)

		remainCOUNT = int(C.bd_get_n_remain(container))
	}

	err := handler.Send(conn, "$$__RESPONSE_PREDICT__$$\n$$__SEPARATOR__$$\n"+response+"\n$$__RESPONSE_DONE__$$\n")
	if err != nil {
		fmt.Println("Send error:", err)
		return err
	}

	return nil
}
