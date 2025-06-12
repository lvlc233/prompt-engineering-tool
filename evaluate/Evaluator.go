package evaluate

//// 评估器接口,需要自己实现
type Evaluator interface {
	Evaluate(evaluation *Evaluation, evaluatePrompt []*Message, actualOutput []*Message) 
}

