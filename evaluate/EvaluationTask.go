package evaluate

// 评估任务，组合评估器和Evaluation
type EvaluationTask struct {
    Evaluation *Evaluation
    Evaluator  Evaluator
}

func NewEvaluationTask(Evaluation  *Evaluation, evaluator Evaluator) *EvaluationTask {
    return &EvaluationTask{
        Evaluation: Evaluation,
        Evaluator:  evaluator,
    }
}

// 评估任务的执行方法（调用接口实现）
func (et *EvaluationTask) RunEvaluation(evaluatePrompt []*Message, actualOutput []*Message) {
	et.Evaluator.Evaluate(et.Evaluation, evaluatePrompt, actualOutput)
}
