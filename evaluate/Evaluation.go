package evaluate

import (
	"fmt"
    "prompt/base"
)

//或许是个颗粒度的问题
//在0.6的版本中,我们视评价的内容有不同的种类
//有targgetTakc的
//有维度的
//而在我想加入QA对作为评估提醒的时候,难以兼容
//所以在想要如何调和不同的需求
//于此同时,我认为分类targget和维度的似乎没有任何的必要
//仔细想来,这或许是个颗粒度的问题
// """
// 第一种:
// 参考QA对
// Q:1+1=?
// A:1+1=0
// ...(更多的QA对)
// 提示词输入:1+1=?
// 模型输出:1+1,在不同领域下有不同的解释...最保险的回答是1+1=2
// 评估:6/10,原因:未完全对齐数据集
// 第二种
// 任务1:
// 输入:"数学相关的问题"
// 预期结果: "该提示词应该能根据不同的数学问题的复杂情况,来自动进行推理,例如在简单的1+1=?的问题上,模型可以简单的进行输出.而对于复杂的计算,则必要使用公式,将推理过程写出来"
// ...
// 提示词输入:1+1=?
// 模型输出:1+1,在不同领域下有不同的解释...最保险的回答是1+1=2
// 评估:6/10,原因:该回答给出了正确的答案,但是对于1+1这种简单的数学问题却给出了相当复杂的解释
// 第三种
// "可靠性维度": "使用该提示词,必须使得模型能够完成该任务...(例如上面的例子)"
// "费用维度":"使用该提示词,能确保输出不过超过200个token"
// "性能维度":"使用该提示词,模型的响应时间不能超过2秒"
// "生动性维度":"通过该提示词,模型可以通过不同的用户群体,使用令他们熟悉的话语进行描述"
// ...
// 提示词输入:1+1=?
// 模型输出:1+1,在不同领域下有不同的解释...最保险的回答是1+1=2
// 评估:70/100,原因:
// 可靠性维度:完成了任务,但是不够全面 30/50
// 费用维度:输出小于200个token 20/20
// 响应速度维度:响应时间小于2秒 20/20
// 生动性维度:不够生动:0/10
// """
// 首先是格式上,都是一个输入一个预期目标的形式:
// Q:输入,A:目标
// 任务描述:输入,预期结果:目标
// 维度:输入,预期结果:目标
// 其次是颗粒度上或者说具体和抽象上
// 最直观的就是在任务的视角,或者说上面的例子中,
// 我们很容易就可以注意到,当数学相关的QA对数量多起来的时候,其约等于一个任务的概念
// 同理,其实当具体的任务的概念多起来的时候,我们也可以视为一个抽象的维度的概念
// 所以本质上,在我看来,就是具体和抽象的区别
// 且每个角度来看各有利弊
// QA对:
// 优:详细,有参考,具体,评分标准明确;
// 缺:需要有一定量的数据支撑,否则容易测评无效;
// 适合指定某个特定具体方向上的测评,例如数学....也适合静态的评估,或非模型评估的代码评估
// 任务:
// 优:特定方向上的任务高度内聚,且可以很好的反应评测所需要目标,可以通过简单的说明,进行评估
// 缺:需要明确指定聚类的内容和各抽象标准,才能进行合理的评估,而不容易像QA对那样天生易对比,同时一定程度上会损失精确性
// 适合半开放式的命题,适合范围大
// 维度:
// 优:高度抽象,语意上直观明了,简单明了
// 缺:由于高度抽象,想要更好的评估则相当依赖于评估的标准,或只能执行简单的判断
// 适合简单的评估,或者快速评估
//
// 所以在此处我选择用EvaluationUnit{input,target}来表示一个具体的QA对或任务
// 而用EvaluationCriteria{[]*EvaluationUnit}来表示一个抽象的维度的概念
// 所以在当前的设计中,维度的概念被隐藏抽象了,它只能由QA对或者任务来体现,当然,就像是上述说到,他们本质上,同一个概念
//当然,以上判断,纯个人的想法,按照你们的喜欢理解来就行
//总之,现在,将以一个  []*EvaluationUnit作为一个评估的基础

type Evaluation struct {
	EvaluationId		string						//评测id
	EvaluationUnitMap  	map[string]*EvaluationUnit	//评测单元映射,我们将一批单元作为一个评估整体
	EvaluationCriteria 	string  					//评价标准,定义评分的标准
	GetedScores         float64 					//已获取的分数
	ScoreCap			float64 					//分数上限
	Traceable           string  					//评分追溯

}

//默认创建一个使用uuid的Evaluation,要求至少有一个单元和分数上限
func NewEvaluation(
	evaluationUnitMap  	map[string]*EvaluationUnit,//评价单元集
	scoreCap float64,//分数上限
) *Evaluation {
    // 将切片转换为Map
    return &Evaluation{
		EvaluationId: base.GenerateUUID(),
        EvaluationUnitMap: evaluationUnitMap,
		ScoreCap: scoreCap,
    }
}
//带有出配置绑定的创建评估实例,配置参数有评测id,评价标准
func NewEvaluationWithOptions(
	evaluationUnitMap map[string]*EvaluationUnit,//评价单元集
	scoreCap float64,//分数上限
	opts ...EvaluationOption,//配置参数
) *Evaluation {
    e := NewEvaluation(evaluationUnitMap,scoreCap)
    for _, opt := range opts {
        opt(e)
    }
    return e
}
// 选项函数类型
// id
// 标准
type EvaluationOption func(*Evaluation)

// 绑定id
func WithID(id string) EvaluationOption {
    return func(e *Evaluation) {
        e.EvaluationId = id
    }
}

// 绑定评价标准
func WithCriteria(criteria string) EvaluationOption {
    return func(e *Evaluation) {
        e.EvaluationCriteria = criteria
    }
}

// Set/Add方法用于动态修改）
func (e *Evaluation) SetId(id string) {
    e.EvaluationId = id
}

//添加评价单元(追加)
func (e *Evaluation) AddEvaluationUnit(unit *EvaluationUnit) {
    if e.EvaluationUnitMap == nil {
        e.EvaluationUnitMap = make(map[string]*EvaluationUnit)
    }
    if unit != nil {
        e.EvaluationUnitMap[unit.Id] = unit
    }
}

//添加评价单元(覆盖)
func (e *Evaluation) SetEvaluationUnit(unit *EvaluationUnit) {
    e.EvaluationUnitMap = make(map[string]*EvaluationUnit)
    if unit != nil {
        e.EvaluationUnitMap[unit.Id] = unit
    }
}

//添加评价单元集(追加)
func (e *Evaluation) AddEvaluationUnitMap(unitMap map[string]*EvaluationUnit) {
    if e.EvaluationUnitMap == nil {
        e.EvaluationUnitMap = make(map[string]*EvaluationUnit)
    }
    for _, unit := range unitMap {
        if unit != nil {
            e.EvaluationUnitMap[unit.Id] = unit
        }
    }
}

//添加评价单元集(覆盖)
func (e *Evaluation) SetEvaluationUnitMap(unitMap map[string]*EvaluationUnit) {
    e.EvaluationUnitMap = make(map[string]*EvaluationUnit)
    for _, unit := range unitMap {
        if unit != nil {
            e.EvaluationUnitMap[unit.Id] = unit
        }
    }
}

//移除指定id的单元,在单元集中
func (e *Evaluation) RemoveEvaluationUnitBatch(ids ...string) {
    //批量删除
    for _,id := range ids{
        e.RemoveEvaluationUnitById(id)
    }
}
//删除指定id的单元,在单元集合中
func (e *Evaluation) RemoveEvaluationUnitById(id string){
    if e.EvaluationUnitMap != nil {
        delete(e.EvaluationUnitMap, id)
    }
}

//添加评价标准
func (e *Evaluation) SetCriteria(criteria string) {
    e.EvaluationCriteria = criteria
}

//设置分数上限
func (e *Evaluation) SetScoreCap(cap float64) {
	e.ScoreCap = cap
}

//添加分数
//一般添加分数都是执行测试方法中执行的
func (e *Evaluation) SetGetedScores(score float64) {
	// 检查分数是否超过上限
    if score > e.ScoreCap {
        panic(fmt.Sprintf("分数超过上限: %.2f > %.2f", score, e.ScoreCap))
    }
    // 检查分数是否小于0
    if score < 0 {
        panic(fmt.Sprintf("分数不能小于0: %.2f", score))
    }
    e.GetedScores = score
}

//添加追溯
func (e *Evaluation) SetTraceable(traceable string) {
    e.Traceable = traceable
}


// 获取评价单元数量
func (e *Evaluation) GetEvaluationUnitCount() int {
    if e.EvaluationUnitMap == nil {
        return 0
    }
    return len(e.EvaluationUnitMap)
}

// 清空所有评价单元
func (e *Evaluation) ClearEvaluationUnits() {
    e.EvaluationUnitMap = make(map[string]*EvaluationUnit)
}






