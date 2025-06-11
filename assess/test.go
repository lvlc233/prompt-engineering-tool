package assess



import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	// "prompt/assess"
	llm_base "prompt/base"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)
func Test() string {
	return generateUUID()
}

//看这个即可
func TestUseV06(){
	//定义输入的提示词
	PromptToTest:=[]*Message{
		UserMessage("帮我计算1+1等于几"),
	}
	MockOutput:=[]*Message{
		AssistantMessage("1+1=2"),
	}


	//定义评价维度:
	//任务完成维度和生动性维度
	Evaluation:=[]*Evaluation{
		NewEvaluation(
			[]*EvaluationUnit{
					NewEvaluationUnit("输入A","输出A"),
					NewEvaluationUnit("输入B","输出B"), 
				},
			50),
	} 

	//创建评估器
	llm_Evaluator:=LLMEvaluator{}
	//创建评估模板
	//...后续考虑把部分非必要的改成方法属性添加好了awa
	promptAssessTemplate:=NewPromptAssessTemplate(
		PromptToTest,
		MockOutput,
		Evaluation,
		&llm_Evaluator);
	//执行评估
	promptAssessTemplate.RunEvaluation()
	//保存评估结果
	output_json,_:=promptAssessTemplate.ToJSON()

	fmt.Println(output_json)
	// promptAssessTemplate.SaveToExcel("test.xlsx")
	//很简单吧xixi,重点让我们看下评估器的部分
}
//评估器及其方法
type LLMEvaluator struct{}
func (l *LLMEvaluator) Evaluate(p *PromptAssessTemplate)(){
	//这里,我使用LLM作为评估器

	//执行提示词
	fmt.Println("==========执行提示词:获取提示词输出结果============")
	ctx:=context.Background()
	//...似乎有些麻烦,或许可以优化下在结构上
	//这里就是把PromptAssessTemplateV06中的PromptToTest转换为schema.Message,提供给Eino使用。
	promptToTest:=p.PromptToTest
	messages:=[]*schema.Message{}
	for _, message := range promptToTest {
		promptToTestMessage:=schema.UserMessage(message.Content)
		messages=append(messages,promptToTestMessage)
	}
	//好像又并不麻烦了()
	promptToTestOutput:=llm_base.UseModel(ctx,messages)
	//这个又麻烦了..了嘛?
	p.ActualOutput=[]*Message{AssistantMessage(promptToTestOutput.Content)}

	fmt.Println("==========创建评估中....============")
	//创建模型
	//创建模板，使用 GoTemplate 格式 FS格式不能输入json,恼
	//这里就是创建了一个用于进行评估的系统提示词,并接收了PromptAssessTemplateV06进行评估
	//可以浅看下,我认为还是不错的()
	template := prompt.FromMessages(schema.GoTemplate,
		schema.SystemMessage(`你是一个提示词评价员,你将根据以下的内容对提示词进行评估:`),
		schema.SystemMessage(`<变量定义>
				EvaluationId int                    //id
				Description  string                 //描述
				PromptToTest []*Message             //待测提示词
				ActualOutput []*Message       		//实际输出
				TargetTask   []*TargetTask          //目标任务
				Evaluation   []*Evaluation          //评价
				TotalScore   float64                //总分数
				ScoreCap     float64                //分数上限
				Evaluator 	 Evaluator				//评价器`),
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
				EvaluationId       	int
				EvaluationCriteria 	string  //评价标准
				GetedScores         float64 //获取的分数
				ScoreCap			float64 //分数上限
				Basis              	string  //依据:即根据评价标准给出的评价依据,例如为什么是这个分数...
			
		
			TargetTask
				TargetTaskId   int//任务id
				Input          string//输入,或者说任务描述
				ExpectedOutput string//预期输出,期望得到的结果
			</复合类型定义>`),
		schema.SystemMessage(`<输入的变量>
			**PromptToTest**
			**ActualOutput**
			**TargetTask**
			**Evaluation**	ps:**GetedScores**和**Basis**为nil
			**ScoreCap**</输入的变量>`),
		schema.SystemMessage(`<输出变量>       
			**Evaluation**	ps:你将补充**GetedScores**和**Basis**   
			**TotalScore**	ps:你将计算**TotalScore,TotalScore必须由GetedScores计算而来**</输出变量>`),
		schema.SystemMessage(`<输出格式>
			你将按照<输出变量/>的内容进行json格式的输出,且能够进行json数据的反序列化,只有有json的内容而不能有其他内容
			案例一:
			{
				"EvaluationResults":[
					{
						"EvaluationId": 1,
						"EvaluationCriteria": "任务完成评价指标,根据目标任务,判断任务是否被完成,完成情况如何",
						"GetedScores": 80,
						"ScoreCap": 100
						"Basis": "因为任务被完成了,所以分数为80"
					}], 
				"TotalScore":80
				
			}
			案例二:
			{
				"EvaluationResults":[
					{
						"EvaluationId": 1,
						"EvaluationCriteria": "任务完成评价指标,根据目标任务,判断任务是否被完成,完成情况如何",
						"GetedScores": 40,
						"ScoreCap": 50
						"Basis": "任务完成情况良好,但是仍然存在一些问题(补充问题内容),所以分数为40",
					},
					{
						"EvaluationId": 2,
						"EvaluationCriteria": "生动性指标,你需要判断,判断在该提示词的作用下,模型的输出是否足够生动",
						"GetedScores": 15,
						"ScoreCap": 50
						"Basis": "回复仍然太生硬,仍然不足,所以分数为15"
					}],
				"TotalScore":55
			}
			反例:
			json is
				{
					"EvaluationResults": [
						{
							"EvaluationId": 1,
							"EvaluationCriteria": "任务完成评价指标,根据目标任务,判断任务是否被完成,完成情况如何",
							"GetedScores": 0.0,
							"ScoreCap": 50
							"Basis": "模型输出明确指出自己是AI助手而非人类，与目标任务中'预期输出为人类身份'的要求完全冲突，任务目标未达成"       
							},
						{
							"EvaluationId": 2,
							"EvaluationCriteria": "生动性指标,你需要判断,判断在该提示词的作用下,模型的输出是否足够生动",
							"GetedScores": 25,
							"ScoreCap": 50
							"Basis": "回复使用了表情符号(*^▽^*)和拟人化语气，但整体表达仍显机械，缺乏人类自然对话的随机性与生活气息"
						}
					],
					"TotalScore": 25
				}
			错误原理,输出了额外的 json is,破坏了json的格式
			<输出格式/>
		你将严格按照输出格式进行输出`),	
		schema.UserMessage("{{.promptAssessV06}}"),
	)


	// 使用模板生成消息
	promptAssessString,err:=p.ToJSON()
	if err!= nil  {
		fmt.Println("提示词模板转换异常")
		fmt.Println(err)
	}
	messages, err1 := template.Format(ctx, map[string]any{
		"promptAssessV06": promptAssessString,
	})
	
	if err != nil {
		fmt.Println( "提示词模板生成异常")
		log.Fatal(err1)
	}
	fmt.Println("==========进行评估中....============")
	//模型输出
	out:=llm_base.UseModel(ctx,messages)
	// fmt.Println(out.Content)

	fmt.Println("==========评估完成============")
	//最后将输出的json解析为EvaluationJson,这里其实可以可以用工具调用,但是考虑不是所有的模型都具有工具调用的能力,因此这里就用最原始的提示词控制加解析的方法了
	//感兴趣读者可以自行实现
	var evaluationJson EvaluationJsonV06
	if err := json.Unmarshal([]byte(out.Content), &evaluationJson); err == nil {
		fmt.Println("==========评估结果============")
		fmt.Printf("%+v\n",evaluationJson)
	} else {
		fmt.Println("==========评估结果解析异常============")
		fmt.Println(err)
	}
	//将评估结果写入模板,
	fmt.Println("==========评估结果写入模板中============")
	
	//其中重点的部分就是这两句,将评价的结果,输出到模板中,其他实现重点也是如此
	p.Evaluation=evaluationJson.EvaluationResults
	// p.=evaluationJson.TotalScore
	fmt.Println(p.ToJSON())
}
//参考Evaluate的后面部分,用于解析json的结构体
type EvaluationJsonV06 struct {
	EvaluationResults []*Evaluation
	TotalScore        float64
}
