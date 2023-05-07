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
	C.bd_set_model_path(l.Container, C.CString(modelFNAME))

	result := bool(C.bd_load_model(l.Container))
	if !result {
		return fmt.Errorf("failed to load the model")
	}

	return nil
}

func (l *LLama) GetRemainCount() int {
	return int(C.bd_get_n_remain(l.Container))
}

func (l *LLama) PredictTokens() bool {
	return bool(C.bd_predict_tokens(l.Container))
}

func (l *LLama) SetIsInteracting(isInteracting bool) {
	C.bd_set_is_interacting(l.Container, C.bool(isInteracting))
}

func (l *LLama) GetEmbedSize() int {
	return int(C.bd_get_embd_size(l.Container))
}

func (l *LLama) GetEmbedString(idx int) string {
	id := C.bd_get_embed_id(l.Container, C.int(idx))
	embedCSTR := C.bd_get_embed_string(l.Container, id)
	embedSTR := C.GoString(embedCSTR)

	return embedSTR
}

func (l *LLama) CheckPromptOrContinue() bool {
	return bool(C.bd_check_prompt_or_continue(l.Container))
}

func (l *LLama) DropBackUserInput() {
	C.bd_dropback_user_input(l.Container)
}

func (l *LLama) Predict(conn *ws.Conn, handler ws.Codec) error {
	remainCOUNT := l.GetRemainCount()

	responseBufferBytes := []byte{}
	responseBuffer := ""

END:
	for remainCOUNT != 0 {
		ok := l.PredictTokens()
		if !ok {
			return fmt.Errorf("failed to predict the tokens")
		}

		// display text
		embdSIZE := l.GetEmbedSize()
		for i := 0; i < embdSIZE; i++ {
			select {
			case <-l.PredictStop:
				remainCOUNT = 0
				break END
			default:
				embedSTR := l.GetEmbedString(i)

				responseBufferBytes = append(responseBufferBytes, []byte(embedSTR)...)
				if !utf8.ValidString(embedSTR) {
					continue // Because connection is closed, don't send invalid UTF-8
				}

				if len(responseBufferBytes) > 0 {
					responseBuffer = string(responseBufferBytes)
					if !utf8.ValidString(responseBuffer) {
						continue
					}
				} else {
					responseBuffer += embedSTR
				}

				// fmt.Print(responseBuffer)

				err := handler.Send(conn, "$$__RESPONSE_PREDICT__$$\n$$__SEPARATOR__$$\n"+responseBuffer)
				if err != nil {
					fmt.Println("Send error:", err)
					remainCOUNT = 0
					break END
				}

				responseBufferBytes = []byte{}
				responseBuffer = ""
			}
		}

		ok = l.CheckPromptOrContinue()
		if !ok {
			break
		}
		l.DropBackUserInput()

		remainCOUNT = l.GetRemainCount()
	}

	err := handler.Send(conn, "$$__RESPONSE_PREDICT__$$\n$$__SEPARATOR__$$\n"+"\n$$__RESPONSE_DONE__$$\n")
	if err != nil {
		fmt.Println("Send error:", err)
		return err
	}

	return nil
}
