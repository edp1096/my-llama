package cgollama

// #cgo CXXFLAGS: -I../llama.cpp/examples -I../llama.cpp
// #cgo LDFLAGS: -L../ -lbinding -lm -lstdc++ -static
// #include "binding.h"
import "C"
import (
	"fmt"
	"unsafe"
)

type LLama struct {
	state unsafe.Pointer
}

func New(model string, opts ...ModelOption) (*LLama, error) {
	mo := NewModelOptions(opts...)
	modelPath := C.CString(model)
	result := C.load_model(modelPath, C.int(mo.ContextSize), C.int(mo.Parts), C.int(mo.Seed), C.bool(mo.F16Memory), C.bool(mo.MLock))
	if result == nil {
		return nil, fmt.Errorf("failed loading model")
	}

	return &LLama{state: result}, nil
}

func (l *LLama) Predict(text string, po PredictOptions) error {
	input := C.CString(text)

	params := C.llama_allocate_params(input, C.int(po.Seed), C.int(po.Threads), C.int(po.Tokens), C.int(po.TopK),
		C.float(po.TopP), C.float(po.Temperature), C.float(po.Penalty), C.int(po.Repeat), C.bool(po.IgnoreEOS), C.bool(po.F16KV))

	predVARs := C.llama_prepare_pred_vars(params, l.state)
	remainCOUNT := C.llama_get_remain_count(predVARs)

	for remainCOUNT > 0 {
		idsSIZE := int(C.llama_get_embedding_ids(params, predVARs))

		for i := 0; i < idsSIZE; i++ {
			id := C.llama_get_id(predVARs, C.int(i))
			embedCSTR := C.llama_get_embed_string(predVARs, id)
			embedSTR := C.GoString(embedCSTR)
			fmt.Print(embedSTR)
		}

		remainCOUNT = C.llama_get_remain_count(predVARs)
		isTokenEND := C.llama_check_token_end(predVARs)
		if bool(isTokenEND) {
			break
		}
	}
	fmt.Println()

	C.llama_default_signal_action()
	C.llama_free_params(params)

	return nil
}
