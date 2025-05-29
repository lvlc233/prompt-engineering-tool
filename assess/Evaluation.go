package assess

import "fmt"

//要不要加个状态呢?
type Evaluation struct {
	EvaluationId       	int
	EvaluationCriteria 	string  //评价标准
	GetedScores         float64 //获取的分数
	ScoreCap			float64 //分数上限
	Basis              	string  //依据:即根据评价标准给出的评价依据,例如为什么是这个分数...
	
}
func (e *Evaluation) ToString() string {
	return fmt.Sprintf("评测id: %d\n评测标准: %s\n获取分数: %.2f\n分数上限: %.2f\n判断依据: %s\n",
	e.EvaluationId,e.EvaluationCriteria,e.GetedScores,e.ScoreCap,e.Basis)
}
func (e *Evaluation) ToStringNotScores() string {
	return fmt.Sprintf("评测id: %d\n评测标准: %s\n判断依据: %s\n",
	e.EvaluationId,e.EvaluationCriteria,e.Basis)
}
func (e *Evaluation) ToStringOnlyScores() (string){
	return fmt.Sprintf("%.2f/%.2f\n",e.GetedScores,e.ScoreCap)
}

func (e *Evaluation) SetScore(GetedScores float64 ) {
	e.GetedScores=GetedScores
}
func (e *Evaluation) SetBasis(Basis string) {
	e.Basis=Basis
}

func NewEvaluation(
	id int,
	evaluationCriteria string,
	ScoreCap float64,
) *Evaluation {
	return &Evaluation{
		EvaluationId:                 id,
		EvaluationCriteria: evaluationCriteria,
		ScoreCap: ScoreCap,
	}
}