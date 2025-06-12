package assess

import (
	// "context"
	// "encoding/csv"
	"encoding/json"
	"fmt"
	// "os"
	// "strconv"
	// "github.com/xuri/excelize/v2"
)

// v0.6.2
/*
元评估提示词(MetaEvaluatePrompt),用于构建一个元评估提示词,用于评估
通过它,我们可以得到一个可靠的提示词评估结果
*/
type MetaEvaluatePrompt struct {
    //评估模板id,您可以用任意的您喜欢的方式创建id,这里默认使用uuid,用于唯一标记一个评估模板,无其他任何作用,允许相同的id
	MetaEvaluatePromptId string
    //一个评估模板的描述信息,可以是任意的字符串,用于描述该评估模板的用途
	Description  string        
    //待测试的提示词,同ActualOutput一起说好了,看Message就可以看出,这就是模型的输入和输出,当然,这里只是定义上的
    //输入这两个字段,不一定要实际的调用,使用您收集数据集进行加载仍然是可以的
    //PromptToTest 是输入的提示词,虽然没有明确的定义,但是这里的Message的Role推荐为User和System
    //您可以使用evaluate.UserMessage()和evaluate.SystemMessage()来分别创建User和System消息
	PromptToTest []*Message    //待测提示词
    //ActualOutput 是模型的提示词,虽然没有明确的定义,但是这里的Message的Role推荐为Assistant和Tool
    //您可以使用evaluate.AssistantMessage()和evaluate.ToolMessage()来分别创建Assistant和Tool消息
	ActualOutput []*Message    //实际输出
    //这里是用list还是map好?
    //看来得从DDD或者是持久化的角度去想了
	Evaluation   []*Evaluation //评价
}

//默认创建一个使用uuid的PromptEvaluateTemplate,要求至少有输入输出和一个评价
func NewMetaEvaluatePrompt(
	prompToTest    []*Message,
    actualOutput    []*Message,
	evaluations     []*Evaluation,
    evaluator       Evaluator,
) *MetaEvaluatePrompt {
	MetaEvaluatePrompt := &MetaEvaluatePrompt{
		MetaEvaluatePromptId: generateUUID(),
		PromptToTest:           promptToTest,
        ActualOutput :          actualOutput,
		Evaluation:             evaluations,
	}
	return MetaEvaluatePrompt
}

//参数绑定
func NewMetaEvaluatePromptWithOptions(
	prompToTest []*Message,
    actualOutput []*Message,
	evaluations []*Evaluation,
    evaluator Evaluator,
    opts ...MetaEvaluatePromptOption,
) *MetaEvaluatePrompt {
    e := NewMetaEvaluatePrompt(promptToTest, actualOutput, evaluations,evaluator)
    for _, opt := range opts {
        opt(e)
    }
    return e
}



//参数包含id,描述
type MetaEvaluatePromptOption func(*MetaEvaluatePrompt)

func WithMetaEvaluatePromptOptionId(id string) MetaEvaluatePromptOption {
    return func(p *MetaEvaluatePrompt) {
        p.MetaEvaluatePromptId = id
    }
}

func WithDescription(description string) MetaEvaluatePromptOption {
    return func(p *MetaEvaluatePrompt) {
        p.Description = description
    }
}

//提供Set方法
//设置PromptEvaluateTemplateId
func (p *MetaEvaluatePrompt) SetMetaEvaluatePromptId(id string) {
    p.MetaEvaluatePromptId = id
}

//设置描述
func (p *MetaEvaluatePrompt) SetDescription(description string) {
    p.Description = description
}

//设置待测试提示词
func (p *MetaEvaluatePrompt) SetPromptToTest(promptToTest []*Message) {
    p.PromptToTest = promptToTest
}

//设置实际输出
func (p *MetaEvaluatePrompt) SetActualOutput(actualOutput []*Message) {
    p.ActualOutput = actualOutput
}

//添加评价
func (p *MetaEvaluatePrompt) AddEvaluation(evaluation *Evaluation) {
    p.Evaluation = append(p.Evaluation, evaluation)
}

//添加评价列表
func (p *MetaEvaluatePrompt) AddEvaluationList(evaluations []*Evaluation) {
    for _, eval := range evaluations {
        p.AddEvaluation(eval)
        
    }
}
//移除评价通过id
func (p *MetaEvaluatePrompt) RemoveEvaluation(id string) {
    for i, eval := range p.Evaluation {
        if eval.EvaluationId == id {
            p.Evaluation = append(p.Evaluation[:i], p.Evaluation[i+1:]...)
            return
        }
    }    
}
//批次移除评价
func (p *MetaEvaluatePrompt) RemoveEvaluationList(ids []string) {
    for _, id := range ids {
        p.RemoveEvaluation(id)
    }
}

//设置评价
func (p *MetaEvaluatePrompt) SetEvaluation(evaluations []*Evaluation) {
    p.Evaluation = evaluations
}


// 计算总分数
func (p *MetaEvaluatePrompt) GetTotalScore() float64 {
    total := 0.0
    for _, eval := range p.Evaluation {
        total += eval.GetedScores
    }
    return total
}

// 计算分数上限
func (p *MetaEvaluatePrompt) GetScoreCap() float64 {
    total := 0.0
    for _, eval := range p.Evaluation {
        total += eval.ScoreCap
    }
    return total
}

// 获取分数百分比
func (p *MetaEvaluatePrompt) GetScorePercentage() float64 {
    cap := p.GetScoreCap()
    if cap == 0 {
        return 0
    }
    return (p.GetTotalScore() / cap) * 100
}

//以下A的,不用细看了(),我只能说,能用
func (p *MetaEvaluatePrompt) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("marshal PromptEvaluateTemplateV06 failed: %v", err)
	}
	return string(jsonBytes), nil
}