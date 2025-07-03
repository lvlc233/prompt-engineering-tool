package evaluate

import (
	"prompt/base"
)

// // 评估器接口,需要自己实现
type Evaluator interface {
	Evaluate(evaluation *Evaluation, evaluatePrompt []*base.Message, actualOutput []*base.Message)
}
