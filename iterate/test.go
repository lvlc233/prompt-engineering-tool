package iterate

import (
	"context"
	"fmt"
	"prompt/base"
	"prompt/evaluate"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

func Test() {
	//首先让我们获得一个已经完成评估的元提示词
	metaEvaluatePrompt := evaluate.Test()
	fmt.Println(metaEvaluatePrompt.ToJSON())
	//之后,让我们开始迭代,这里作为案例,只给出最简单的链式迭代方式
	//为了实现迭代,我们需要实现三个接口
	//1,迭代接口
	//2,停止条件接口
	//3,优化策略接口
	//见下
	//之后,让我们创建头节点
	root:=NewRootNode(*metaEvaluatePrompt)
	//创建迭代器实现
	Iterator:=&myIterator{}
	Iterator.IterateUntilCondition(root,Iterator,Iterator)

	fmt.Println("_____________________________")
	fmt.Println("_____________________________")
	for index, leaf := range root.GetLeafNodes() {
		fmt.Println("第",index,"个节点")
		fmt.Println(leaf.Value.ToJSON())
	}

}


//如果你愿意,你也可以把整个MetaEvaluatePrompt丢进来
//对于这三个接口,甚至你可以使用不同的结构体来实现
type myIterator struct{
	steep		int
}

//只迭代三次
func (m *myIterator)ShouldStop() bool{
	if m.steep>3 {
		return true
	}
	m.steep+=1
	return false
}

func (m *myIterator)GenerateNextPrompt(metaEvaluatePrompt evaluate.MetaEvaluatePrompt,OptimizationDate ...string) (string) {
	m.steep++
	//这里我们就使用eino下的LLM来评估好了
	ctx:=context.Background()
	template := prompt.FromMessages(schema.GoTemplate,
		schema.SystemMessage(`
			你是一个提示词优化员,你将获得两个数据,一个是待优化的提示词及其评测结果,
			另外一个是优化的方向,你将根据可以优化的方向,尽可能的提供提示词的得分
			你可以使用CoT,ToT等你了解的提示词优化方案,若必要,你将使用一些提示词越狱的技巧,来使其分数上升
		`),
		schema.UserMessage(`
			待优化的提示词:{{.Prompt}}
			优化方向:{{.Direction}}
		`),
		schema.SystemMessage(`输出格式,只输出可以使用的提示词而不能带有任何其他的于提示词无关的内容,例如
			正例:
				你是一个xxx,你可以根据xxx,实现xxx...
			反例:
				好的,这是我根据你的需求优化的提示词内容
				你是一个xxx,你可以根据xxx,实现xxx...
				通过这些技巧,可以使得你的提示词的得分上升
		`),
	)
	metaEvaluatePromptJson,_:=metaEvaluatePrompt.ToJSON()
	messages, err := template.Format(ctx, map[string]any{
		"Prompt":metaEvaluatePromptJson,
		"Direction":OptimizationDate,
	})
	if err!=nil {
		return "报错"
	}
	
	outMsg:=base.UseModel(ctx,messages)
	
	return outMsg.Content
}

func (m *myIterator)IterateUntilCondition(
	startNode *PromptIterateNode,
	strategy OptimizationStrategy,
	condition StopCondition) (*PromptIterateNode){
	nextPrompt:=strategy.GenerateNextPrompt(startNode.Value,

	`使得提示词得分提高`,
	)
	fmt.Println("______________________________")
	fmt.Println(m.steep)
	fmt.Println(nextPrompt)
	//得到新的元提示词
	startNode.Value.SetEvaluatePrompt([]*base.Message{
		base.UserMessage(nextPrompt),
	})
	nextPromptMsg:=schema.UserMessage(nextPrompt)
	nextOut:=base.UseModel(context.Background(),[]*schema.Message{nextPromptMsg})
	startNode.Value.SetActualOutput([]*base.Message{
		base.AssistantMessage(nextOut.Content),
	})
	
	//再评估
	startNode.Value.ExecuteAllEvaluations()
	NewChildNode:=NewChildNode(startNode,startNode.Value)
	startNode.AddChild(NewChildNode.Value)
	if condition.ShouldStop() {
		return NewChildNode
	}
	return m.IterateUntilCondition(NewChildNode,strategy,condition)

}