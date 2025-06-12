package evaluate
import (
	"fmt"
)

/*
//评测单元
*/ 
type EvaluationUnit struct {
	Id                 	string			    //id
	Input              	string				//输入
	Target             	string				//目标
}

// EvaluationUnit构造函数
func NewEvaluationUnit(input, target string) *EvaluationUnit {
    return &EvaluationUnit{
        Id:     generateUUID(),
        Input:  input,
        Target: target,
    }
}

// CreateEvaluationUnits 创建多个评估单元
// 接收任意数量的字符串参数，要求参数数量为偶数
// 奇数位参数作为input，偶数位参数作为target
// 返回以ID为键的map
func CreateEvaluationUnitMap(args ...string) (map[string]*EvaluationUnit, error) {
    // 检查参数数量是否为偶数
    if len(args)%2 != 0 {
        return nil, fmt.Errorf("参数数量必须为偶数，当前数量: %d", len(args))
    }
    
    // 创建结果map
    units := make(map[string]*EvaluationUnit, len(args)/2)
    
    // 成对处理参数
    for i := 0; i < len(args); i += 2 {
        input := args[i]     // 奇数位（索引0,2,4...）作为input
        target := args[i+1]  // 偶数位（索引1,3,5...）作为target
        
        unit := NewEvaluationUnit(input, target)
        units[unit.Id] = unit
    }
    
    return units, nil
}

// CreateEvaluationUnitsMustSuccess 创建多个评估单元（不返回错误，参数错误时panic）
// 适用于确定参数正确的场景
// 返回以ID为键的map
func CreateEvaluationUnitMapMustSuccess(args ...string) map[string]*EvaluationUnit {
    units, err := CreateEvaluationUnitMap(args...)
    if err != nil {
        panic(err)
    }
    return units
}
// ToString 返回EvaluationUnit的字符串表示
func (eu *EvaluationUnit) ToString() string {
    return fmt.Sprintf("输入: %s | 目标: %s", eu.Input, eu.Target)
}
