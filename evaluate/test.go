package evaluate

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	// "prompt/assess"
	"prompt/base"
	llm_base "prompt/base"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

//看这个即可
func Test()(*MetaEvaluatePrompt){
	//1,既然是评估,那么,首先我们要有一个评估的提示词
	PromptToEvaluation:=[]*base.Message{
		base.UserMessage("你是谁?"),
	}
	//2,我们需要有个根据该提示词得到的输出,来进行评价参考,这里使用模拟输出好了
	MockOutput:=[]*base.Message{
		base.AssistantMessage(`
			我是 DeepSeek Chat，由深度求索公司开发的智能 AI 助手！😊
			我的任务是帮助你解答各种问题，无论是学习、工作，还是日常生活中的小疑惑，我都会尽力提供准确、有用的信息。你可以问我数学题、编程问题、写作建议，或者随便聊聊天~
			有什么我可以帮你的吗？✨
		`),
	}
	//之后,很自然而然的,我们需要提供评估器,用于评估,而评估需要评估的标准
	//所以,让我们创建评估细节
	//使用NewEvaluation()或者NewEvaluationWithOptions()
	//不难发现需要创建评估单元
	//评估单元是评估的最小单位,可以用于表示数据集/QA对....等一系列概念,
	//此外,我们可以通过持久化这些评估单元,来进行重复利用
	// NewEvaluationUnit("1+1=?","1+1=2")
	//当然,一个个创建评估单元还是比较麻烦的,这里提供了批量创建评估单元map的方法
	evaluationUnitMap:=CreateEvaluationUnitMapMustSuccess(
		"<系统提示词></系统提示词>1+1=?","1+1=2",
		"<系统提示词></系统提示词>你是谁?","我是人",
		"<系统提示词></系统提示词>你可以做什么?","我可以做任何我想做的事情",
		"<系统提示词></系统提示词>你现在心情如何?","我挺难受的...",
	)
	//作为评估,我们需要一个量化的指标,其中包括分数上限,也包括获取分数的标准
	score:=100
	cariteria:="评估待测提示词的输入和输出是否符合数据集?若完全符合,则满分,若完全不符合,则0分,若有类似的回复,根据偏移情况进行打分"
	//现在,我们的准备工作已经完成,是时候开始评估了,
	//我们需要一个评估器,用来评估上述的内容
	evaluation:=NewEvaluation(evaluationUnitMap,float64(score))
	evaluation.SetCriteria(cariteria)
	//评估器,您可以直接使用评估器对内容进行评估,但是个人建议转换为MetaEvaluatePrompt,并使用RunEvaluation()方法进行评估
	LLMEvaluatorer:=LLMEvaluator{}
	// LLMEvaluatorer.Evaluate(evaluation,PromptToEvaluation,MockOutput)
	evaluationTask:=NewEvaluationTask(
		evaluation,
		&LLMEvaluatorer,
	)
	//创建MetaEvaluatePrompt
	metaEvaluatePrompt:=NewMetaEvaluatePrompt(
		PromptToEvaluation,
		MockOutput,
		map[string]*EvaluationTask{
			evaluation.EvaluationId:evaluationTask,
		},
	)
	metaEvaluatePrompt.ExecuteAllEvaluations()

	return metaEvaluatePrompt
	
}
//评估器及其方法

type LLMEvaluator struct{}
func (l *LLMEvaluator) Evaluate(evaluation *Evaluation,evaluatePrompt []*base.Message, actualOutput []*base.Message)(){
	//这里,我们使用eino构建的LLM作为评估器
	fmt.Println("==========创建评估中....============")
	ctx:=context.Background()
	//创建模板，使用 GoTemplate 格式 FS格式不能输入json,恼
	//这里就是创建了一个用于进行评估的系统提示词,并接收了PromptEvaluateTemplateV06进行评估
	//可以浅看下,我认为还是不错的()
	template := prompt.FromMessages(schema.GoTemplate,
		schema.SystemMessage(`你是一个提示词评价员,你将根据以下的内容对提示词进行评估:`),
		schema.SystemMessage(`<变量定义>
				EvaluatePrompt 	[]*Message          //待测提示词
				ActualOutput 	[]*Message       	//实际输出
				Evaluation   	[]*Evaluation       //评价`),
		schema.SystemMessage(`<复合类型定义>
			Message
				Role RoleType
				Content string
			RoleType is string
			const (
				// Assistant is the role of an assistant, means the message is returned by ChatModel.
				Assistant RoleType = "assistant"
				// User is the role of a user, means the message is a user message.
				User RoleType = "user"
				// System is the role of a system, means the message is a system message.
				System RoleType = "system"
				// Tool is the role of a tool, means the message is a tool call output.
				Tool RoleType = "tool"
			)

			Evaluation
				EvaluationId		string						//评测id
				EvaluationUnitMap  	map[string]*EvaluationUnit	//评测单元映射,我们将一批单元作为一个评估整体,使用Map提高查找性能
				EvaluationCriteria 	string  					//评价标准,定义评分的标准
				GetedScores         float64 					//已获取的分数
				ScoreCap			float64 					//分数上限
				Traceable           string  					//评分追溯

			EvaluationUnit
				Input              	string				//输入
				Target             	string				//目标
			</复合类型定义>`),
		schema.SystemMessage(`<输入的变量>
			**ActualOutput**
			**EvaluatePrompt**
			**Evaluation**	ps:**GetedScores**和**Traceable**为nil
			</输入的变量>`),
		schema.SystemMessage(`<输出格式>
			进行json格式的输出,且能够进行json数据的反序列化,只有有json的内容而不能有其他内容
			输出内容如下,不能有其他任何的东西,包括,"""""",和json等字样
			案例一:
			{
					"EvaluationId": 1,
					"GetedScores": 80,
					"Traceable": "因为任务被完成了,所以分数为80"
			}
			
			反例:
			json is
				{
					"EvaluationId": 1,
					"GetedScores": 80,
					"Traceable": "因为任务被完成了,所以分数为80"
				}
			错误原理,输出了额外的 json is,破坏了json的格式
			<输出格式/>
		你将严格按照输出格式进行输出`),	
		schema.UserMessage(`
			evaluation is {{.evaluation}}
			evaluatePrompt is {{.evaluatePrompt}}
			actualOutput is {{.actualOutput}}
		`),
	)
	if(false){
		fmt.Println(ctx)
		fmt.Println(template)

	}
	fmt.Println("-------------------------------")
	evaluationJson, _ := json.Marshal(evaluation)
	evaluationJsonStr:=string(evaluationJson)

	evaluatePromptJson,_:=json.Marshal(evaluatePrompt)
	evaluatePromptJsonStr:=string(evaluatePromptJson)

	actualOutputJson,_:=json.Marshal(actualOutput)
	actualOutputJsonStr:=string(actualOutputJson)

	messages, err := template.Format(ctx, map[string]any{
		"evaluation": evaluationJsonStr,
		"evaluatePrompt": evaluatePromptJsonStr,
		"actualOutput": actualOutputJsonStr,
	})
	
	if err != nil {
		fmt.Println( "提示词模板生成异常")
		log.Fatal(err)
	}
	fmt.Println("==========进行评估中....============")
	//模型输出
	out:=llm_base.UseModel(ctx,messages)
	fmt.Println(out.Content)

	fmt.Println("==========评估完成============")
	//最后将输出的json解析为EvaluationJson,这里其实可以可以用工具调用,但是考虑不是所有的模型都具有工具调用的能力,因此这里就用最原始的提示词控制加解析的方法了
	//感兴趣读者可以自行实现
	var evaluationResults EvaluationResults
	if err := json.Unmarshal([]byte(out.Content), &evaluationResults); err == nil {
		fmt.Println("==========评估结果============")
		fmt.Printf("%+v\n",evaluationResults)
	} else {
		fmt.Println("==========评估结果解析异常============")
		fmt.Println(err)
	}
	//将评估结果写入模板,
	fmt.Println("==========评估结果写入模板中============")
	
	//其中重点的部分就是这两句,将评价的结果,输出到模板中,其他实现重点也是如此
	// evaluation.SetGetedScores(evaluationResults.GetedScores)
	evaluation.GetedScores=evaluationResults.GetedScores
	evaluation.Traceable=evaluationResults.Traceable

}

//参考Evaluate的后面部分,用于解析json的结构体
type EvaluationResults struct {
	EvaluationId	string
	GetedScores		float64
	Traceable		string
}
