package iterate

import (
	"fmt"
	"prompt/base"
	"prompt/evaluate"
)

// 支持分支的迭代节点
type PromptIterateNode struct {
	Value    	evaluate.MetaEvaluatePrompt 	// 元评估提示词
	Version  	string                      	// 版本
	Parent   	*PromptIterateNode        		// 父节点
	Children 	[]*PromptIterateNode      		// 子分支列表
}


// 迭代接口
// 接受一个待迭代的节点,停止迭代条件,和优化策略
//其中,你可以选择任意的迭代节点进行迭代,并最终输出下一个节点,下一个节点可以是原本的节点,从而实现对同一个节点多次的迭代
//也可以指向下一个不同的节点,使其成为链式,这都随意
//您可以指定任意的停止条件,例如迭代次数,例如评分...也随您
//此外,虽然并非硬性要求,但是我个人推荐在优化的时候也提供一定的数据集,指导其优化方向,
//当然,您仍然可以只使用MetaEvaluatePrompt作为优化的依据
type IterationInterface interface {
	IterateUntilCondition(
		startNode *PromptIterateNode,
		strategy OptimizationStrategy,
		condition StopCondition) (*PromptIterateNode, error)
}


// 停止条件接口
type StopCondition interface {
	ShouldStop() bool

}

// 优化策略接口
type OptimizationStrategy interface {
	GenerateNextPrompt(metaEvaluatePrompt evaluate.MetaEvaluatePrompt,OptimizationDate ...string) (string)
}





//基础的关于树的方法
// 创建根节点
func NewRootNode(value evaluate.MetaEvaluatePrompt) *PromptIterateNode {
	return &PromptIterateNode{
		Value:    	value,
		Version:  	base.GenerateUUID(),
		Parent:   	nil,
		Children: 	make([]*PromptIterateNode, 0),
	}
}

// 创建子节点
func NewChildNode(parent *PromptIterateNode, value evaluate.MetaEvaluatePrompt) *PromptIterateNode {
	if parent == nil {
		return nil
	}
	
	child := &PromptIterateNode{
		Value:    	value,
		Version:  	base.GenerateUUID(),
		Parent:   parent,
		Children: make([]*PromptIterateNode, 0),
	}
	
	parent.Children = append(parent.Children, child)
	return child
}

// 添加子节点到指定父节点
func (node *PromptIterateNode) AddChild(value evaluate.MetaEvaluatePrompt) *PromptIterateNode {
	return NewChildNode(node, value)
}

// 删除节点（会同时删除其所有子节点）
func (node *PromptIterateNode) Delete() bool {
	if node == nil {
		return false
	}
	
	// 如果是根节点，不能删除
	if node.Parent == nil {
		return false
	}
	
	// 从父节点的子节点列表中移除
	parent := node.Parent
	for i, child := range parent.Children {
		if child == node {
			// 删除该子节点
			parent.Children = append(parent.Children[:i], parent.Children[i+1:]...)
			break
		}
	}
	
	// 递归删除所有子节点
	node.deleteAllChildren()
	
	// 清空当前节点的引用
	node.Parent = nil
	node.Children = nil
	
	return true
}

// 递归删除所有子节点
func (node *PromptIterateNode) deleteAllChildren() {
	for _, child := range node.Children {
		child.deleteAllChildren()
		child.Parent = nil
		child.Children = nil
	}
	node.Children = make([]*PromptIterateNode, 0)
}

// 根据版本号查找节点
func (node *PromptIterateNode) FindByVersion(version string) *PromptIterateNode {
	if node == nil {
		return nil
	}
	
	if node.Version == version {
		return node
	}
	
	// 递归查找子节点
	for _, child := range node.Children {
		if found := child.FindByVersion(version); found != nil {
			return found
		}
	}
	
	return nil
}

// 根据 MetaEvaluatePromptId 查找节点
func (node *PromptIterateNode) FindById(id string) *PromptIterateNode {
	if node == nil {
		return nil
	}
	
	if node.Value.MetaEvaluatePromptId == id {
		return node
	}
	
	// 递归查找子节点
	for _, child := range node.Children {
		if found := child.FindById(id); found != nil {
			return found
		}
	}
	
	return nil
}

// 获取所有叶子节点（没有子节点的节点）
func (node *PromptIterateNode) GetLeafNodes() []*PromptIterateNode {
	if node == nil {
		return nil
	}
	
	var leaves []*PromptIterateNode
	
	// 如果没有子节点，则是叶子节点
	if len(node.Children) == 0 {
		leaves = append(leaves, node)
	} else {
		// 递归获取所有子节点的叶子节点
		for _, child := range node.Children {
			leaves = append(leaves, child.GetLeafNodes()...)
		}
	}
	
	return leaves
}

// 获取从根节点到当前节点的路径
func (node *PromptIterateNode) GetPath() []*PromptIterateNode {
	if node == nil {
		return nil
	}
	
	var path []*PromptIterateNode
	current := node
	
	// 从当前节点向上遍历到根节点
	for current != nil {
		path = append([]*PromptIterateNode{current}, path...)
		current = current.Parent
	}
	
	return path
}

// 获取节点深度（根节点深度为0）
func (node *PromptIterateNode) GetDepth() int {
	if node == nil {
		return -1
	}
	
	depth := 0
	current := node
	
	for current.Parent != nil {
		depth++
		current = current.Parent
	}
	
	return depth
}

// 获取子节点数量
func (node *PromptIterateNode) GetChildrenCount() int {
	if node == nil {
		return 0
	}
	return len(node.Children)
}

// 判断是否为叶子节点
func (node *PromptIterateNode) IsLeaf() bool {
	return node != nil && len(node.Children) == 0
}

// 判断是否为根节点
func (node *PromptIterateNode) IsRoot() bool {
	return node != nil && node.Parent == nil
}

// 打印节点树结构（用于调试）
func (node *PromptIterateNode) PrintTree(prefix string, isLast bool) {
	if node == nil {
		return
	}
	
	// 打印当前节点
	connector := "├── "
	if isLast {
		connector = "└── "
	}
	
	fmt.Printf("%s%sV%d: %s\n", prefix, connector, node.Version, node.Value.MetaEvaluatePromptId)
	
	// 打印子节点
	newPrefix := prefix
	if isLast {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}
	
	for i, child := range node.Children {
		isLastChild := i == len(node.Children)-1
		child.PrintTree(newPrefix, isLastChild)
	}
}



