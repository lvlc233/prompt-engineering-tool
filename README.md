
# Prompt Engineering Tool

一个用于提示工程的全栈应用，包含 Go 后端和 React 前端。

## 快速启动

### 方式一：使用批处理脚本（推荐新手）

#### 开发模式启动
```bash
# 启动开发服务器（使用 go run）
start.bat
```

#### 生产模式启动
```bash
# 构建并启动生产服务器（使用编译的 exe）
start-production.bat
```

#### 手动构建
```bash
# 构建前后端为可执行文件
build.bat
```

### 方式二：使用 PowerShell 脚本（推荐开发者）

```powershell
# 查看所有可用命令
.\scripts.ps1 help

# 构建应用
.\scripts.ps1 build

# 启动开发服务器
.\scripts.ps1 start-dev

# 启动生产服务器（使用编译的 exe）
.\scripts.ps1 start-prod

# 清理构建文件
.\scripts.ps1 clean

# 运行测试
.\scripts.ps1 test
```

## 服务地址

- **前端**: http://localhost:3000
- **后端**: http://localhost:8593

## 开发 vs 生产模式

### 开发模式
- 使用 `go run` 启动后端，支持代码热重载
- 适合开发和调试
- 启动稍慢，但修改代码后重启快

### 生产模式
- 使用编译后的 `backend.exe` 启动
- 启动速度快，性能更好
- 适合部署和演示
- 修改代码后需要重新构建

## 依赖要求

- **Go**: 1.19+
- **Node.js**: 16+
- **npm**: 8+

## 项目技术栈
- **后端**: Go + Gin + SQLite + Eino
- **前端**: React + TypeScript + Vite
- **构建**: Go build + npm build
- **部署**: 可执行文件 + 静态文件

## 项目结构
```
prompt-engineering-tool/
├── bk/                 # Go 后端
│   ├── main.go        # 主入口文件
│   └── database.go    # 数据库相关
├── web/               # React 前端
│   ├── src/
│   ├── package.json
│   └── tsconfig.json
├── bin/               # 编译后的可执行文件
│   └── backend.exe    # 后端可执行文件
├── assess/            # 评估模块
│   ├── Evaluation/    # 评估模块
│   ├── Message/       # 消息模块（已移至base）
│   ├── PromptAssesTempalte/ # 提示词评估模板模块
│   ├── TargetTack/    # 目标任务模块
│   └── test/          # 测试模块
├── base/              # 基础模块
│   └── model/         # 模型模块
├── build.bat          # 构建脚本
├── start.bat          # 开发启动脚本
├── start-production.bat # 生产启动脚本
└── scripts.ps1        # PowerShell 管理脚本
```

## 常见问题

### 1. 编码问题
所有脚本都已配置 UTF-8 编码，如果仍有乱码，请确保终端支持 UTF-8。

### 2. 端口占用
如果端口被占用，请检查是否有其他服务在运行：
- 前端端口：3000
- 后端端口：8593

### 3. 构建失败
- 确保 Go 和 Node.js 已正确安装
- 检查网络连接（下载依赖时需要）
- 运行 `scripts.ps1 clean` 清理后重试

### 4. Go 编译为 exe 的优势
- **启动速度快**: 编译后的可执行文件启动比 `go run` 快很多
- **部署简单**: 单个 exe 文件，无需 Go 环境
- **性能更好**: 编译优化后的代码执行效率更高
- **生产就绪**: 适合生产环境部署

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

# 更新日记
## v0.6.1:
1,提供了一个简单的uuid生成器
2,优化了PromptAssesTemplate的结构
  具体的说:
  2-1:移除了PromptAssessTemplate中的总分数和分数上限,通过计算方法来进行计算,而非属性
  2-2:移除TaggetTack的概念:将维度,QA对,任务的概念合并为评估
3, 在Evaluation中添加了一个EvaluationUnit用于表示任意的QA对,任务,维度的概念,
4,添加了相关的方法

## v0.6.2:
1,将Assess统一为Evaluate
2,将asses文件夹,包名,修改为evaluate
3,添加了EvaluationTask在EvaluationTask.go中,用于确保Evaluation和Evaluator的一一对应关系,并在其中实现了Evaluator接口
4,移除可Evaluation中的Evaluator,使其只存储数据内容,并使用map来存储unit,用于提供性能
5,修改了EvaluationUnit中的批量创建的方法,
6,确定了Evaluation,EvaluationUnit,MetaEvaluatePrompt的id的必要性
7,Evaluator的函数签名修改
8,修改PromptAssessTemplate为MetaEvaluatePrompt,使其符合含义
9,在MetaEvaluatePrompt的数据结构中使用Map,存储不同的Evaluation
10, 添加批量执行评估的方法,移除持久化到cvs的方法(考虑后面实现的持久化方法)
11,更新了Test中的案例实现 

## v0.6.3:
1,将Message 和 uuidGen 移动至 base模块
2,初步搭建迭代器
    2-1:考虑到历史追踪的需求,这边将其数据结构定位为树的结构并提供相关的方法
    2-2:提供了迭代器执行接口及其相关停止迭代,优化策略结构