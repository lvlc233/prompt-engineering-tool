package base

import (
	"context"
	"log"
	"os"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/schema"
	"github.com/joho/godotenv"
)




func _init(ctx context.Context) (*openai.ChatModel, error) {
	// 加载.env文件
	evn_load_err := godotenv.Load()
	if evn_load_err != nil {
		log.Fatal(" .env文件加载失败")
	}

	// 从环境变量获取配置
	var ChatModel, new_chat_err = openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   os.Getenv("model_name"),
		APIKey:  os.Getenv("api_key"),
		BaseURL: os.Getenv("base_url"),
	})
	return ChatModel, new_chat_err
}


//这里就暂时用eino的接口进行模型的调用了
//后期可以换成接口,用于适配不同的框架
//暂时的想法是定义一个userModel()(*[]string)的接口
//因为已经定义了自己的Message(虽然是抄的enio的(欸嘿)),也许可能换成userModel()(*[]Message)
//后面再计划吧
func UseModel(ctx context.Context,in []*schema.Message)(outMsg *schema.Message) {
	 model,new_chat_err :=_init(ctx)
	 if new_chat_err!=nil {
		log.Fatal("LLM模型加载失败")
	 }
	 out,model_generate_err := model.Generate(ctx,in)
	 if model_generate_err!=nil {
		log.Fatal("LLM模型消息生成失败")
	 }
	 return out

}


// 使用 Option

func Test() {
 
}
