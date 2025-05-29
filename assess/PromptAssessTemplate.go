package assess

import (
	// "context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"github.com/xuri/excelize/v2"
)

//ps,若想看作者的心路历程可以看下v0.1~最新版前(v0.5)的部分,若不想的话,直接看v0.6的部分即可

// v0.1
// 评价标准和目标任务之间存在一定的耦合关系,但是也不能是完全耦合...这里就用组合的方式交给使用者吧
type PromptAssessTemplateV01 struct {
	//提示词测评模板:包含以下内容
	//目的:输入待测目标提示词和目标任务,根据评价标准,给出评价结果和分数
	//用于可靠的提示词优化
	PromptToTest       string  //待测提示词
	TargetTask         string  //目标任务
	EvaluationCriteria string  //评价标准
	Output             string  //输出
	Score              float64 //分数
	Basis              string  //依据:即根据评价标准给出的评价依据,例如为什么是这个分数...
}

// New methods for all structs
func NewPromptAssessTemplateV01() *PromptAssessTemplateV01 {
	return &PromptAssessTemplateV01{}
}

// v0.2
// 这里考虑到目标任务和标准标准作为切片可能更加方便,也更加容易进行单元测试
// 而提示词的输入最终不论输入会是变成一个字符串,就保持这样子的处理好了,同理输出也是
type PromptAssessTemplateV02 struct {
	//提示词测评模板:包含以下内容
	//目的:输入待测目标提示词和目标任务,根据评价标准,给出评价结果和分数
	//用于可靠的提示词优化
	PromptToTest       string   //待测提示词
	TargetTask         []string //目标任务
	EvaluationCriteria []string //评价标准
	Output             string   //输出
	Score              float64  //分数
	Basis              string   //依据:即根据评价标准给出的评价依据,例如为什么是这个分数...
}

func NewPromptAssessTemplateV02() *PromptAssessTemplateV02 {
	return &PromptAssessTemplateV02{}
}

// v0.3
// 或许我们需要一个目标任务和评价标准的结构体...
type PromptAssessTemplateV03 struct {
	//提示词测评模板:包含以下内容
	//目的:输入待测目标提示词和目标任务,根据评价标准,给出评价结果和分数
	//用于可靠的提示词优化
	PromptToTest       string               //待测提示词
	TargetTask         []TargetTask         //目标任务
	EvaluationCriteria []EvaluationCriteria //评价标准
	Output             string               //输出
	Score              float64              //分数
	Basis              string               //依据:即根据评价标准给出的评价依据,例如为什么是这个分数...
}

func NewPromptAssessTemplateV03() *PromptAssessTemplateV03 {
	return &PromptAssessTemplateV03{}
}

// 我们在该版本中尝试将"任务的完成质量"和"提示词质量"分开,诚然,任务的完成质量是提示词质量中非常重要的一部分
// 但是,不同的任务的完成方式,目的,标准却是截然不同,若将"任务的完成质量"同"提示词质量"混为一谈,则很容易导致工程上的灾难
// 因此,我们将"任务的完成质量"和"提示词质量"分开,这样子的设计的好处是有以下至少 三点
// 1.资源利用更高效:
//
//	1-1:混合的情况下,提示词本身的质量考虑"任务的完成质量"和"提示词质量"两个方面,而当"任务的完成质量"不达标的时候,仍然需要进行"提示词质量"的评估,显然,一定程度上会浪费一定的资源
//	1-2:分开的情况下,我们可以根据"任务的完成质量"来判断是否需要进行"提示词质量"的评估,这样可以更加高效地利用资源
//
// 2.测试的边界更加清晰,我们可以进行更加精细的单元测试
// 3.可维护性与可扩展性更强
//
//	3-1:可以想想下,如果我们将"任务的完成质量"和"提示词质量"合并,当涉及的内容在"提示词质量"上重复,而仅在"任务"的维度上有所不同,那么使用分开的方式可以更好进行重复利用,
//		而无需重复的进行编写,尽管我们可以复制粘贴.但是那样子工程将会非常糟糕
//
// 那么?什么是"提示词质量",我们可以很轻松的找出一些指标,例如:准确性,相关性,可读性,创造性...等等,就是类似于此的指标
// 当然,实际上,您完完全全可以将"提示词质量"和"任务的完成质量"进行任意的转换或者是混合,本质上,它们只是为了完成评测任务而存在的概念
// 您也可以将"任务的完成质量"视作与系统无关的部分,例如输出的json,文档,图片等,又或者具体的说,是LLM输出的文本内容,
// 而"提示词质量"则是与系统相关的部分,例如响应时间,token使用情况,模型名称...等,又或者具体的说,是无法在LLM中得到的内容
// 当然,您怎么方便,怎么来。
type EvaluationCriteria struct {
	Id          int
	Description string
}

func NewEvaluationCriteria(
	id int,
	description string,
) *EvaluationCriteria {
	return &EvaluationCriteria{
		Id:          id,
		Description: description,
	}
}

// v0.4
// 或许我们需要考虑加入一个mate信息,用于记录除了标准之外的一些信息...例如token使用情况,任务执行时间等...考虑到可扩展性....
type PromptAssessTemplateV04 struct {
	//提示词测评模板:包含以下内容
	//目的:输入待测目标提示词和目标任务,根据评价标准,给出评价结果和分数
	//用于可靠的提示词优化
	PromptToTest       string                 //待测提示词
	TargetTask         []TargetTask           //目标任务
	EvaluationCriteria []EvaluationCriteria   //评价标准
	Output             string                 //输出
	Score              float64                //分数
	Basis              string                 //依据:即根据评价标准给出的评价依据,例如为什么是这个分数...
	Mate               map[string]interface{} //mate信息,就这样子吧(乐)
}

func NewPromptAssessTemplateV04(
	promptToTest string,
	targetTask []TargetTask,
	evaluationCriteria []EvaluationCriteria,
) *PromptAssessTemplateV04 {
	return &PromptAssessTemplateV04{
		PromptToTest:       promptToTest,
		TargetTask:         targetTask,
		EvaluationCriteria: evaluationCriteria,
		Mate:               make(map[string]interface{}),
	}
}

// ToJSON converts PromptAssessTemplateV04 to JSON string
func (p *PromptAssessTemplateV04) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("marshal PromptAssessTemplateV04 failed: %v", err)
	}
	return string(jsonBytes), nil
}

// FromJSON creates PromptAssessTemplateV04 from JSON string
func (p *PromptAssessTemplateV04) FromJSON(jsonStr string) error {
	err := json.Unmarshal([]byte(jsonStr), p)
	if err != nil {
		return fmt.Errorf("unmarshal to PromptAssessTemplateV04 failed: %v", err)
	}
	return nil
}

// v0.5
// 将依据和分数加入到评价标准中
type PromptAssessTemplateV05 struct {
	//提示词测评模板:包含以下内容
	//目的:输入待测目标提示词和目标任务,根据评价标准,给出评价结果和分数
	//用于可靠的提示词优化
	Id           int                    //id
	description  string                 //描述
	PromptToTest string                 //待测提示词
	Output       string                 //输出
	TargetTask   []TargetTask           //目标任务
	Evaluation   []Evaluation           //评价
	TotalScore   float64                //总分数
	Mate         map[string]interface{} //mate信息,就这样子吧(乐)
}

func (p *PromptAssessTemplateV05) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("marshal PromptAssessTemplateV05 failed: %v", err)
	}
	return string(jsonBytes), nil
}

func NewPromptAssessTemplateV05(
	promptToTest string,
	targetTask []TargetTask,
	evaluationCriteria []Evaluation,
) *PromptAssessTemplateV05 {
	return &PromptAssessTemplateV05{
		PromptToTest: promptToTest,
		TargetTask:   targetTask,
		Evaluation:   evaluationCriteria,
	}
}

// type Evaluation struct {
// 	Id                 int
// 	EvaluationCriteria string  //评价标准
// 	Score              float64 //分数
// 	Basis              string  //依据:即根据评价标准给出的评价依据,例如为什么是这个分数...
// }

// func NewEvaluation(
// 	id int,
// 	evaluationCriteria string,
// ) *Evaluation {
// 	return &Evaluation{
// 		Id:                 id,
// 		EvaluationCriteria: evaluationCriteria,
// 	}
// }

// 定义执行评价的接口,用于执行评价
// type executeEvaluation interface {
// 	executeEvaluationUseLLM(evaluator func(*PromptAssessTemplateV05) error)
// }

// 定义一个默认的使用LLM进行评价的实现
func (p *PromptAssessTemplateV05) executeEvaluationUseLLM(evaluator func(*PromptAssessTemplateV05) error) {
	err := evaluator(p)
	if err != nil {
		fmt.Println("执行评估发生错误")
		fmt.Println(err)
		return
	}
}

// v0.6
//可参考v0.5的版本解释
type PromptAssessTemplateV06 struct {
	//提示词测评模板:包含以下内容
	//目的:输入待测目标提示词和目标任务,根据评价标准,给出评价结果和分数
	//用于可靠的提示词优化
	EvaluationId int           //id
	Description  string        //描述
	PromptToTest []*Message    //待测提示词
	ActualOutput []*Message    //实际输出
	TargetTask   []*TargetTask //目标任务
	Evaluation   []*Evaluation //评价
	TotalScore   float64       //总分数
	ScoreCap     float64       //分数上限
	Evaluator    Evaluator     //评价器
}

func NewPromptAssessTemplateV06(
	evaluationId int,
	description string,
	promptToTest []*Message,
	targetTasks []*TargetTask,
	evaluations []*Evaluation,
	Evaluator Evaluator,
) *PromptAssessTemplateV06 {
	PromptAssessTemplateV06 := &PromptAssessTemplateV06{
		EvaluationId: evaluationId,
		Description:  description,
		PromptToTest: promptToTest,
		TargetTask:   targetTasks,
		Evaluation:   evaluations,
		Evaluator:    Evaluator,
	}
	// 解引用指针并遍历切片
	for _, eval := range evaluations {
		// 计算总体分数上限
		//这里会有bug吗?
		PromptAssessTemplateV06.ScoreCap += eval.ScoreCap
	}
	return PromptAssessTemplateV06
}
type Evaluator interface {
	Evaluate(promptAssessTemplateV06 *PromptAssessTemplateV06) 
}

// 模板的评测方法（调用接口实现）
func (p *PromptAssessTemplateV06) RunEvaluation() {
	p.Evaluator.Evaluate(p)
}

//以下A的,不用细看了(),我只能说,能用
func (p *PromptAssessTemplateV06) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		return "", fmt.Errorf("marshal PromptAssessTemplateV06 failed: %v", err)
	}
	return string(jsonBytes), nil
}

func (p *PromptAssessTemplateV06) SaveToCSV(filename string) error {
    // 检查文件是否存在，如果不存在则创建并写入表头
    var file *os.File
    var err error
    var writeHeader bool

    if _, err := os.Stat(filename); os.IsNotExist(err) {
        // 文件不存在，创建并写入表头
        file, err = os.Create(filename)
        writeHeader = true
    } else {
        // 文件存在，以追加模式打开
        file, err = os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
        writeHeader = false
    }
    
    if err != nil {
        return err
    }
    defer file.Close()

    writer := csv.NewWriter(file)
    defer writer.Flush()

    // 如果文件是新创建的，写入表头
    if writeHeader {
        header := []string{"评估ID", "描述", "待测试提示词", "实际输出", "目标任务", "评价", "分数占比", "总分数", "分数上限"}
        if err := writer.Write(header); err != nil {
            return err
        }
    }

    // 获取最长的数组
    maxRows := max(
        len(p.PromptToTest),
        len(p.ActualOutput),
        len(p.TargetTask),
        len(p.Evaluation),
    )
    
    // 写入数据行
    for i := 0; i < maxRows; i++ {
		var record []string
		if(i==0){
			record = []string{
				strconv.Itoa(p.EvaluationId),
				p.Description,
				getMessageToString(p.PromptToTest, i),
				getMessageToString(p.ActualOutput, i),
				getTargetTaskToString(p.TargetTask, i),
				getEvaluationToStringNotScore(p.Evaluation, i),
				getEvaluationToStringOnlySroce(p.Evaluation, i),
				strconv.FormatFloat(p.TotalScore, 'f', 2, 64),
				strconv.FormatFloat(p.ScoreCap, 'f', 2, 64),
			}
		}else{
			record = []string{
				"",
				"",
				getMessageToString(p.PromptToTest, i),
				getMessageToString(p.ActualOutput, i),
				getTargetTaskToString(p.TargetTask, i),
				getEvaluationToStringNotScore(p.Evaluation, i),
				getEvaluationToStringOnlySroce(p.Evaluation, i),
				"",
				"",
			}
		}
        if err := writer.Write(record); err != nil {
            return err
        }
    }

    return nil
}

// 辅助函数：获取 Message 内容
func getMessageToString(messages []*Message, index int) string {
	if index < len(messages) && messages[index] != nil {
		return messages[index].toString()
	}
	return ""
}

// 辅助函数：获取 TargetTask 内容
func getTargetTaskToString(tasks []*TargetTask, index int) string {
	if index < len(tasks) && tasks[index] != nil {
		return tasks[index].toString()
	}
	return ""
}

// 辅助函数：获取 Evaluation
func getEvaluationToStringOnlySroce(evaluations []*Evaluation, index int) string {
	if index < len(evaluations) && evaluations[index] != nil {
		return evaluations[index].ToStringOnlyScores()
	}
	return ""
}

// 辅助函数：获取 Evaluation
func getEvaluationToStringNotScore(evaluations []*Evaluation, index int) string {
	if index < len(evaluations) && evaluations[index] != nil {
		return evaluations[index].ToStringNotScores()
	}
	return ""
}

// 辅助函数：获取 Evaluation
func getEvaluationToString(evaluations []*Evaluation, index int) string {
	if index < len(evaluations) && evaluations[index] != nil {
		return evaluations[index].ToString()
	}
	return ""
}

func max(values ...int) int {
	maxVal := 0
	for _, v := range values {
		if v > maxVal {
			maxVal = v
		}
	}
	return maxVal
}


func (p *PromptAssessTemplateV06) SaveToExcel(filename string) error {
    // 1. 检查文件是否存在
    _, err := os.Stat(filename)
    fileExists := !os.IsNotExist(err)

    // 2. 初始化Excel文件对象
    var f *excelize.File
    if fileExists {
        // 文件存在时打开现有文件
        f, err = excelize.OpenFile(filename)
        if err != nil {
            return fmt.Errorf("打开现有Excel文件失败: %v", err)
        }
    } else {
        // 文件不存在时创建新文件
        f = excelize.NewFile()
    }
    defer f.Close()

    // 3. 设置/获取工作表
    sheetName := "评估结果"
    index ,_:= f.GetSheetIndex(sheetName)
    if index == -1 {
        // 工作表不存在时创建
        index,_= f.NewSheet(sheetName)
        f.SetActiveSheet(index)
        
        // 写入表头
        headerStyle, _ := f.NewStyle(&excelize.Style{
            Font:      &excelize.Font{Bold: true, Color: "#FFFFFF"},
            Fill:      excelize.Fill{Type: "pattern", Color: []string{"#4F81BD"}, Pattern: 1},
            Alignment: &excelize.Alignment{Horizontal: "center"},
        })

        headers := []string{"评估ID", "描述", "待测试提示词", "实际输出", "目标任务", "评价","分数占比", "总分数", "分数上限"}
        for col, header := range headers {
            cell, _ := excelize.CoordinatesToCellName(col+1, 1)
            f.SetCellValue(sheetName, cell, header)
            f.SetCellStyle(sheetName, cell, cell, headerStyle)
        }
    }

    // 4. 获取现有数据的最后行号
    rows, err := f.GetRows(sheetName)
    if err != nil {
        return fmt.Errorf("获取行数据失败: %v", err)
    }
    startRow := len(rows) + 1
    if startRow == 1 {
        startRow = 2 // 如果只有表头，从第2行开始
    }

    // 5. 写入新数据
    maxRows := max(
        len(p.PromptToTest),
        len(p.ActualOutput),
        len(p.TargetTask),
        len(p.Evaluation),
    )

    for i := 0; i < maxRows; i++ {
        row := startRow + i
        f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), p.EvaluationId)
        f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), p.Description)
        f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), getMessageToString(p.PromptToTest, i))
        f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), getMessageToString(p.ActualOutput, i))
        f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), getTargetTaskToString(p.TargetTask, i))
        f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), getEvaluationToStringNotScore(p.Evaluation, i))
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), getEvaluationToStringOnlySroce(p.Evaluation, i))
		f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), p.TotalScore)
        f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), p.ScoreCap)
    }

    // 6. 设置列宽（如果是新文件）
    if !fileExists {
        f.SetColWidth(sheetName, "A", "A", 10)
        f.SetColWidth(sheetName, "B", "B", 20)
        f.SetColWidth(sheetName, "C", "F", 30)
        f.SetColWidth(sheetName, "G", "I", 12)
    }
	    // 合并单元格
		if err := mergeStaticColumns(f, sheetName, p, startRow, maxRows); err != nil {
			return err
		}

    // 7. 保存文件
    return f.SaveAs(filename)
}
// 合并单元格（同时处理新建文件和追加模式）
func mergeStaticColumns(f *excelize.File, sheetName string, p *PromptAssessTemplateV06, startRow, maxRows int) error {
    if maxRows <= 1 {
        return nil // 不需要合并
    }

    // 需要合并的静态列
    staticColumns := []string{"A", "B", "H", "I"}
    
    for _, col := range staticColumns {
        // 计算合并范围
        startCell := fmt.Sprintf("%s%d", col, startRow)
        endCell := fmt.Sprintf("%s%d", col, startRow+maxRows-1)
        
        // 检查是否已存在合并单元格
        mergedCells, _ := f.GetMergeCells(sheetName)
        alreadyMerged := false
        
        for _, merged := range mergedCells {
            if merged[0] == startCell {
                alreadyMerged = true
                break
            }
        }

        if !alreadyMerged {
            // 执行合并
            if err := f.MergeCell(sheetName, startCell, endCell); err != nil {
                return fmt.Errorf("合并列%s失败: %v", col, err)
            }
            
            // 设置居中对齐样式
            style, _ := f.NewStyle(&excelize.Style{
                Alignment: &excelize.Alignment{
                    Horizontal: "center",
                    Vertical:   "center",
                },
            })
            if err := f.SetCellStyle(sheetName, startCell, startCell, style); err != nil {
                return fmt.Errorf("设置样式失败: %v", err)
            }
        } else {
            // 如果是追加模式且已有合并单元格，需要解除合并后重新合并
            if err := f.UnmergeCell(sheetName, startCell,endCell); err != nil {
                return fmt.Errorf("解除合并失败: %v", err)
            }
            if err := f.MergeCell(sheetName, startCell, endCell); err != nil {
                return fmt.Errorf("重新合并列%s失败: %v", col, err)
            }
        }
    }
    return nil
}

