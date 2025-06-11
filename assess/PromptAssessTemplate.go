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

//ps:从该版本(v0.6.1)版本开始,将去除非最新版的信息,之前的非主要信息,可以去github上查看
// v0.6.1
/*
提示词评估模板(PromptAssessTemplate)负责组织和构建执行评估的过程,是整个"评估器"的最核心的一部分
通过它,我们可以得到一个可靠的提示词评估结果
*/
type PromptAssessTemplate struct {
    //评估模板id,您可以用任意的您喜欢的方式创建id,这里默认使用uuid,用于唯一标记一个评估模板,无其他任何作用,允许相同的id
	PromptAssessTemplateId string
    //一个评估模板的描述信息,可以是任意的字符串,用于描述该评估模板的用途
	Description  string        
    //待测试的提示词,同ActualOutput一起说好了,看Message就可以看出,这就是模型的输入和输出,当然,这里只是定义上的
    //输入这两个字段,不一定要实际的调用,使用您收集数据集进行加载仍然是可以的
    //PromptToTest 是输入的提示词,虽然没有明确的定义,但是这里的Message的Role推荐为User和System
    //您可以使用assess.UserMessage()和assess.SystemMessage()来分别创建User和System消息
	PromptToTest []*Message    //待测提示词
    //ActualOutput 是模型的提示词,虽然没有明确的定义,但是这里的Message的Role推荐为Assistant和Tool
    //您可以使用assess.AssistantMessage()和assess.ToolMessage()来分别创建User和System消息
	ActualOutput []*Message    //实际输出
    //静态的...
    //标签
    //相似性回归....
	Evaluation   []*Evaluation //评价
	Evaluator    Evaluator     //评价器
    //确实应该会有一些静态的...或者说,可以通过简单程序的判断的,那这一点可能需要在评价器上动手了
}

//参数绑定
func NewPromptAssessTemplateWithOptions(
	promptToTest []*Message,
    actualOutput []*Message,
	evaluations []*Evaluation,
    evaluator Evaluator,
    opts ...PromptAssessTemplateOption,
) *PromptAssessTemplate {
    e := NewPromptAssessTemplate(promptToTest, actualOutput, evaluations,evaluator)
    for _, opt := range opts {
        opt(e)
    }
    return e
}

//默认创建一个使用uuid的PromptAssessTemplate,要求至少有输入输出和一个评价
func NewPromptAssessTemplate(
	promptToTest    []*Message,
    actualOutput    []*Message,
	evaluations     []*Evaluation,
    evaluator       Evaluator,
) *PromptAssessTemplate {
	PromptAssessTemplate := &PromptAssessTemplate{
		PromptAssessTemplateId: generateUUID(),
		PromptToTest:           promptToTest,
        ActualOutput :          actualOutput,
		Evaluation:             evaluations,
        Evaluator:              evaluator,
	}
	return PromptAssessTemplate
}

//参数包含id,描述
type PromptAssessTemplateOption func(*PromptAssessTemplate)

func WithPromptAssessTemplateId(id string) PromptAssessTemplateOption {
    return func(p *PromptAssessTemplate) {
        p.PromptAssessTemplateId = id
    }
}

func WithDescription(description string) PromptAssessTemplateOption {
    return func(p *PromptAssessTemplate) {
        p.Description = description
    }
}

//提供Set方法
//设置PromptAssessTemplateId
func (p *PromptAssessTemplate) SetPromptAssessTemplateId(id string) {
    p.PromptAssessTemplateId = id
}

//设置描述
func (p *PromptAssessTemplate) SetDescription(description string) {
    p.Description = description
}

//设置待测试提示词
func (p *PromptAssessTemplate) SetPromptToTest(promptToTest []*Message) {
    p.PromptToTest = promptToTest
}

//设置实际输出
func (p *PromptAssessTemplate) SetActualOutput(actualOutput []*Message) {
    p.ActualOutput = actualOutput
}

//添加评价
func (p *PromptAssessTemplate) AddEvaluation(evaluation *Evaluation) {
    p.Evaluation = append(p.Evaluation, evaluation)
}

//添加评价列表
func (p *PromptAssessTemplate) AddEvaluationList(evaluations []*Evaluation) {
    for _, eval := range evaluations {
        p.AddEvaluation(eval)
        
    }
}
//移除评价通过id
func (p *PromptAssessTemplate) RemoveEvaluation(id string) {
    for i, eval := range p.Evaluation {
        if eval.EvaluationId == id {
            p.Evaluation = append(p.Evaluation[:i], p.Evaluation[i+1:]...)
            return
        }
    }    
}
//批次移除评价
func (p *PromptAssessTemplate) RemoveEvaluationList(ids []string) {
    for _, id := range ids {
        p.RemoveEvaluation(id)
    }
}

//设置评价
func (p *PromptAssessTemplate) SetEvaluation(evaluations []*Evaluation) {
    p.Evaluation = evaluations
}

//设置评价器
func (p *PromptAssessTemplate) SetEvaluator(evaluator Evaluator) {
    p.Evaluator = evaluator
}

type Evaluator interface {
	Evaluate(promptAssessTemplateV06 *PromptAssessTemplate) 
}

// 模板的评测方法（调用接口实现）
func (p *PromptAssessTemplate) RunEvaluation() {
	p.Evaluator.Evaluate(p)
}

// 计算总分数
func (p *PromptAssessTemplate) GetTotalScore() float64 {
    total := 0.0
    for _, eval := range p.Evaluation {
        total += eval.GetedScores
    }
    return total
}

// 计算分数上限
func (p *PromptAssessTemplate) GetScoreCap() float64 {
    total := 0.0
    for _, eval := range p.Evaluation {
        total += eval.ScoreCap
    }
    return total
}

// 获取分数百分比
func (p *PromptAssessTemplate) GetScorePercentage() float64 {
    cap := p.GetScoreCap()
    if cap == 0 {
        return 0
    }
    return (p.GetTotalScore() / cap) * 100
}

//以下A的,不用细看了(),我只能说,能用
func (p *PromptAssessTemplate) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("marshal PromptAssessTemplateV06 failed: %v", err)
	}
	return string(jsonBytes), nil
}

// func (p *PromptAssessTemplate) SaveToCSV(filename string) error {
//     // 检查文件是否存在，如果不存在则创建并写入表头
//     var file *os.File
//     var err error
//     var writeHeader bool

//     if _, err := os.Stat(filename); os.IsNotExist(err) {
//         // 文件不存在，创建并写入表头
//         file, err = os.Create(filename)
//         writeHeader = true
//     } else {
//         // 文件存在，以追加模式打开
//         file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
//         writeHeader = false
//     }
    
//     if err != nil {
//         return err
//     }
//     defer file.Close()

//     writer := csv.NewWriter(file)
//     defer writer.Flush()

//     // 如果文件是新创建的，写入表头
//     if writeHeader {
//         header := []string{"评估ID", "描述", "待测试提示词", "实际输出", "目标任务", "评价", "分数占比", "总分数", "分数上限"}
//         if err := writer.Write(header); err != nil {
//             return err
//         }
//     }

//     // 获取最长的数组
//     maxRows := max(
//         len(p.PromptToTest),
//         len(p.ActualOutput),
//         len(p.TargetTask),
//         len(p.Evaluation),
//     )
    
//     // 写入数据行
//     for i := 0; i < maxRows; i++ {
// 		var record []string
// 		if(i==0){
// 			record = []string{
// 				strconv.Itoa(p.EvaluationId),
// 				p.Description,
// 				getMessageToString(p.PromptToTest, i),
// 				getMessageToString(p.ActualOutput, i),
// 				getTargetTaskToString(p.TargetTask, i),
// 				getEvaluationToStringNotScore(p.Evaluation, i),
// 				getEvaluationToStringOnlySroce(p.Evaluation, i),
// 				strconv.FormatFloat(p.TotalScore, 'f', 2, 64),
// 				strconv.FormatFloat(p.ScoreCap, 'f', 2, 64),
// 			}
// 		}else{
// 			record = []string{
// 				"",
// 				"",
// 				getMessageToString(p.PromptToTest, i),
// 				getMessageToString(p.ActualOutput, i),
// 				getTargetTaskToString(p.TargetTask, i),
// 				getEvaluationToStringNotScore(p.Evaluation, i),
// 				getEvaluationToStringOnlySroce(p.Evaluation, i),
// 				"",
// 				"",
// 			}
// 		}
//         if err := writer.Write(record); err != nil {
//             return err
//         }
//     }

//     return nil
// }

// // 辅助函数：获取 Message 内容
// func getMessageToString(messages []*Message, index int) string {
// 	if index < len(messages) && messages[index] != nil {
// 		return messages[index].toString()
// 	}
// 	return ""
// }

// // 辅助函数：获取 Evaluation
// func getEvaluationToStringOnlySroce(evaluations []*Evaluation, index int) string {
// 	if index < len(evaluations) && evaluations[index] != nil {
// 		return evaluations[index].ToStringOnlyScores()
// 	}
// 	return ""
// }

// // 辅助函数：获取 Evaluation
// func getEvaluationToStringNotScore(evaluations []*Evaluation, index int) string {
// 	if index < len(evaluations) && evaluations[index] != nil {
// 		return evaluations[index].ToStringNotScores()
// 	}
// 	return ""
// }

// // 辅助函数：获取 Evaluation
// func getEvaluationToString(evaluations []*Evaluation, index int) string {
// 	if index < len(evaluations) && evaluations[index] != nil {
// 		return evaluations[index].ToString()
// 	}
// 	return ""
// }

// func max(values ...int) int {
// 	maxVal := 0
// 	for _, v := range values {
// 		if v > maxVal {
// 			maxVal = v
// 		}
// 	}
// 	return maxVal
// }


// func (p *PromptAssessTemplate) SaveToExcel(filename string) error {
//     // 1. 检查文件是否存在
//     _, err := os.Stat(filename)
//     fileExists := !os.IsNotExist(err)

//     // 2. 初始化Excel文件对象
//     var f *excelize.File
//     if fileExists {
//         // 文件存在时打开现有文件
//         f, err = excelize.OpenFile(filename)
//         if err != nil {
//             return fmt.Errorf("打开现有Excel文件失败: %v", err)
//         }
//     } else {
//         // 文件不存在时创建新文件
//         f = excelize.NewFile()
//     }
//     defer f.Close()

//     // 3. 设置/获取工作表
//     sheetName := "评估结果"
//     index ,_:= f.GetSheetIndex(sheetName)
//     if index == -1 {
//         // 工作表不存在时创建
//         index,_= f.NewSheet(sheetName)
//         f.SetActiveSheet(index)
        
//         // 写入表头
//         headerStyle, _ := f.NewStyle(&excelize.Style{
//             Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
//             Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4F81BD"}, Pattern: 1},
//             Alignment: &excelize.Alignment{Horizontal: "center"},
//         })

//         headers := []string{"评估ID", "描述", "待测试提示词", "实际输出", "目标任务", "评价","分数占比", "总分数", "分数上限"}
//         for col, header := range headers {
//             cell, _ := excelize.CoordinatesToCellName(col+1, 1)
//             f.SetCellValue(sheetName, cell, header)
//             f.SetCellStyle(sheetName, cell, cell, headerStyle)
//         }
//     }

//     // 4. 获取现有数据的最后行号
//     rows, err := f.GetRows(sheetName)
//     if err != nil {
//         return fmt.Errorf("获取行数据失败: %v", err)
//     }
//     startRow := len(rows) + 1
//     if startRow == 1 {
//         startRow = 2 // 如果只有表头，从第2行开始
//     }

//     // 5. 写入新数据
//     maxRows := max(
//         len(p.PromptToTest),
//         len(p.ActualOutput),
//         len(p.TargetTask),
//         len(p.Evaluation),
//     )

//     for i := 0; i < maxRows; i++ {
//         row := startRow + i
//         f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), p.EvaluationId)
//         f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), p.Description)
//         f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), getMessageToString(p.PromptToTest, i))
//         f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), getMessageToString(p.ActualOutput, i))
//         f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), getTargetTaskToString(p.TargetTask, i))
//         f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), getEvaluationToStringNotScore(p.Evaluation, i))
// 		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), getEvaluationToStringOnlySroce(p.Evaluation, i))
// 		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), p.TotalScore)
//         f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), p.ScoreCap)
//     }

//     // 6. 设置列宽（如果是新文件）
//     if !fileExists {
//         f.SetColWidth(sheetName, "A", "A", 10)
//         f.SetColWidth(sheetName, "B", "B", 20)
//         f.SetColWidth(sheetName, "C", "F", 30)
//         f.SetColWidth(sheetName, "G", "I", 12)
//     }
// 	    // 合并单元格
// 		if err := mergeStaticColumns(f, sheetName, p, startRow, maxRows); err != nil {
// 			return err
// 		}

//     // 7. 保存文件
//     return f.SaveAs(filename)
// }
// // 合并单元格（同时处理新建文件和追加模式）
// func mergeStaticColumns(f *excelize.File, sheetName string, p *PromptAssessTemplate, startRow, maxRows int) error {
//     if maxRows <= 1 {
//         return nil // 不需要合并
//     }

//     // 需要合并的静态列
//     staticColumns := []string{"A", "B", "H", "I"}
    
//     for _, col := range staticColumns {
//         // 计算合并范围
//         startCell := fmt.Sprintf("%s%d", col, startRow)
//         endCell := fmt.Sprintf("%s%d", col, startRow+maxRows-1)
        
//         // 检查是否已存在合并单元格
//         mergedCells, _ := f.GetMergeCells(sheetName)
//         alreadyMerged := false
        
//         for _, merged := range mergedCells {
//             if merged[0] == startCell {
//                 alreadyMerged = true
//                 break
//             }
//         }

//         if !alreadyMerged {
//             // 执行合并
//             if err := f.MergeCell(sheetName, startCell, endCell); err != nil {
//                 return fmt.Errorf("合并列%s失败: %v", col, err)
//             }
            
//             // 设置居中对齐样式
//             style, _ := f.NewStyle(&excelize.Style{
//                 Alignment: &excelize.Alignment{
//                     Horizontal: "center",
//                     Vertical:   "center",
//                 },
//             })
//             if err := f.SetCellStyle(sheetName, startCell, startCell, style); err != nil {
//                 return fmt.Errorf("设置样式失败: %v", err)
//             }
//         } else {
//             // 如果是追加模式且已有合并单元格，需要解除合并后重新合并
//             if err := f.UnmergeCell(sheetName, startCell,endCell); err != nil {
//                 return fmt.Errorf("解除合并失败: %v", err)
//             }
//             if err := f.MergeCell(sheetName, startCell, endCell); err != nil {
//                 return fmt.Errorf("重新合并列%s失败: %v", col, err)
//             }
//         }
//     }
//     return nil
// }

