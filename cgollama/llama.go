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

func (l *LLama) Predict(text string, po PredictOptions) (string, error) {
	input := C.CString(text)
	// if po.Tokens == 0 {
	// 	po.Tokens = 99999999
	// }
	// out := make([]byte, po.Tokens)

	params := C.llama_allocate_params(input, C.int(po.Seed), C.int(po.Threads), C.int(po.Tokens), C.int(po.TopK),
		C.float(po.TopP), C.float(po.Temperature), C.float(po.Penalty), C.int(po.Repeat), C.bool(po.IgnoreEOS), C.bool(po.F16KV))

	// pred_vars := C.llama_get_prediction(params, l.state, (*C.char)(unsafe.Pointer(&out[0])))
	pred_vars := C.llama_get_prediction(params, l.state)
	n_remains := C.llama_get_remain_count(pred_vars)

	for n_remains > 0 {
		embed_size := int(C.llama_loop_prediction(params, pred_vars))

		for i := 0; i < embed_size; i++ {
			id := C.llama_get_id(pred_vars, C.int(i))
			embedCSTR := C.llama_get_embed_string(pred_vars, id)
			embedSTR := C.GoString(embedCSTR)
			fmt.Print(embedSTR)
		}

		n_remains = C.llama_get_remain_count(pred_vars)
	}
	fmt.Println()

	// res := C.GoString((*C.char)(unsafe.Pointer(&out[0])))
	// res = strings.TrimPrefix(res, " ")
	// res = strings.TrimPrefix(res, text)
	// res = strings.TrimPrefix(res, "\n")

	C.llama_default_signal_action()

	C.llama_free_params(params)

	// return res, nil
	return "", nil
}

// func (l *LLama) Predict(text string, opts ...PredictOption) (string, error) {
// 	po := NewPredictOptions(opts...)

// 	input := C.CString(text)
// 	if po.Tokens == 0 {
// 		po.Tokens = 99999999
// 	}
// 	out := make([]byte, po.Tokens)

// 	params := C.llama_allocate_params(input, C.int(po.Seed), C.int(po.Threads), C.int(po.Tokens), C.int(po.TopK),
// 		C.float(po.TopP), C.float(po.Temperature), C.float(po.Penalty), C.int(po.Repeat), C.bool(po.IgnoreEOS), C.bool(po.F16KV))
// 	ret := C.llama_predict(params, l.state, (*C.char)(unsafe.Pointer(&out[0])))
// 	if ret != 0 {
// 		return "", fmt.Errorf("inference failed")
// 	}
// 	res := C.GoString((*C.char)(unsafe.Pointer(&out[0])))

// 	res = strings.TrimPrefix(res, " ")
// 	res = strings.TrimPrefix(res, text)
// 	res = strings.TrimPrefix(res, "\n")

// 	C.llama_free_params(params)

// 	return res, nil
// }
