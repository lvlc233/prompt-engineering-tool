
# 项目技术栈与结构说明

## 技术栈
- Go
- Eino

## 项目结构
```
assess/                  # 评估模块
    Evaluation/          # 评估模块
    Message/             # 消息模块
    PromptAssesTempalte/ # 提示词评估模板模块
    TargetTack/          # 目标任务模块
    test/                # 测试模块
base/                    # 基础模块
    model/               # 模型模块
```

## 项目简介
本项目专注于提示词工程，涵盖提示词的生成、评估、优化迭代等全流程。工程化的核心在于开发、测试与迭代的闭环。

### 测试的重要性
提示词的测试无法依赖传统工具，必须通过数据集结合人工/AI评估。直观的评估数据是优化提示词的关键，因此本项目优先开发评估器，后续将逐步完善优化器、生成器等模块。

### 评估器设计
评估器用于量化提示词的质量，其核心要素包括：
1. **提示词**
   - 输入提示词（系统/用户输入）
   - 输出提示词（模型的实际输出）
2. **评分标准**
   - 任务指标（监督学习式参考）
   - 性能指标（专家评分式维度）
3. **评分备注**
   - 各维度的性能评价
4. **分数**
   - 各维度分数
   - 总分

> 注：输出提示词、评分标准、评分备注和维度分数为可选项，可根据复杂度调整。

### 模块功能
- **Message**  
  定义消息结构（角色：系统/用户/assistant/tool + 内容）。
- **Evaluation**  
  定义评价结构（评分标准、备注、分数、总分）。
- **TargetTack**  
  定义目标任务结构（任务描述、示例）。
- **PromptAssesTempalte**  
  提供评价框架，支持JSON序列化、CSV/Excel导出及评分执行。

#### 评估接口示例
```go
type PromptAssessTemplateV06 struct {
    Evaluator Evaluator // 评价器接口
}

type Evaluator interface {
    Evaluate(promptAssessTemplateV06 *PromptAssessTemplateV06) 
}

// 执行评估
func (p *PromptAssessTemplateV06) RunEvaluation() {
    p.Evaluator.Evaluate(p)
}
```
用户需实现`Evaluator`接口的`Evaluate`方法，绑定后调用`RunEvaluation`即可。参考`test/`中的示例代码。

## 使用方法
1. 创建`PromptAssessTemplateV06`实例。
2. 绑定自定义的`Evaluator`实现。
3. 调用`RunEvaluation()`执行评估。
4. 使用`ToJson()`输出结果，或通过`SaveToCsv()`/`SaveToExcel()`保存至本地。

# 后续计划
对assess继续优化,使其更加规范,易于使用。
开发其他模块。