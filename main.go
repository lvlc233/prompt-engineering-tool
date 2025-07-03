package main

import (
	"context"
	// "fmt"
	// "prompt/evaluate"
	"prompt/iterate"
)

var ctx=context.Background()

func main() {

	// 评估
	// metaEvaluatePrompt  :=evaluate.Test()
	// fmt.Println(metaEvaluatePrompt.ToJSON())

	// 迭代
	iterate.Test()


}





