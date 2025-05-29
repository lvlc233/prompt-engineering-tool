package assess

import "fmt"

//目标，目标可以被看作是，一个输入输出对，一个任务作为一个条件输入和输出对
//因此我们可以从评价标准上加以区分
//可以根据预期的执行描述来查看目标
// 例如:任务:--将一个字符串转换为拼音--
// 将可以视为接收一个字符串,期望获得拼音:因此我们便可以根据基础的输入输出,
// 来衡量该模型在该提示词上,是否能完成"内容"活动
type TargetTask struct {
	TargetTaskId   int//任务id
	Input          string//输入,或者说任务描述
	ExpectedOutput string//预期输出,期望得到的结果
}
func (t *TargetTask) toString() string {
	return fmt.Sprintf("任务id: %d\n任务: %s\n目标: %s\n", t.TargetTaskId,t.Input, t.ExpectedOutput)
}

func NewTargetTask(
	id int,//任务id
	input string,//输入,或者说任务描述
	expectedOutput string,//预期输出,期望得到的结果
) *TargetTask {
	return &TargetTask{
		TargetTaskId:   id,
		Input:          input,
		ExpectedOutput: expectedOutput,
	}
}




