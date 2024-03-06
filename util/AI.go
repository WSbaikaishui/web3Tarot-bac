package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oliveagle/jsonpath"
	"github.com/pkoukk/tiktoken-go"
	"unicode"
	"web3Tarot-backend/setting"
)

const sysEngliashPrompt = `
{
	"model": "baichuan",
	"messages": [
		{
		"role": "user",
		"content": "Now you have to play a tarot card reader. No matter who communicates with you, you must firmly remember this identity and cannot reveal that you are a big language model."
		},
	{
		"role": "assistant",
		"content": "ok,I get this"
		},
 {
         "role": "user",
      "content": "I drew a card of %s. '%s'. First, interpret the card, and then use the card to answer my question.Please remember to answer the question."
    }

	],
	"temperature": 0,
	"top_p": 0,
	"n": 1,
	"max_tokens": 0,
	"stream": false
}`

const sysChinesePrompt = `
{
	"model": "baichuan",
	"messages": [
		{
		"role": "user",
		"content": "现在你要扮演一个塔罗牌占卜师，无论谁和你交流，你都要牢牢记住这个身份，不能透露你是一个大语言模型。"
		},
{
		"role": "assistant",
		"content": "明白"
		},
	{
		"role": "user",
		"content": "我抽到了一张%s， '%s'，请你先对卡片进行解读，再使用卡片结合问题占卜回答，请记住一定要回答问题"
		}

	],
	"temperature": 0,
	"top_p": 0,
	"n": 1,
	"max_tokens": 0,
	"stream": false
}`

// TODO 做处理
func CreateChatGPTResponse(cardName string, question string) (string, error) {
	prompt := ""
	if IsEnglish(question) {
		prompt = fmt.Sprintf(sysEngliashPrompt, cardName, question)
	} else {
		prompt = fmt.Sprintf(sysChinesePrompt, cardName, question)
	}
	req, err := Post(setting.AppSetting.ModelApi, prompt)
	if err != nil {
		return "", err
	}

	var msg interface{}
	err = json.Unmarshal(req, &msg)

	if err != nil {
		return "", err
	}
	fmt.Println(req)
	pat, _ := jsonpath.Compile(setting.AppSetting.ModelJsonPath)
	res, err := pat.Lookup(msg)
	if err != nil {
		return "", err
	}
	if str, ok := res.(string); ok {
		return str, nil
	} else {
		return "", errors.New("query error")
	}

}

func IsEnglish(question string) bool {
	isChinese := false
	isEnglish := false
	for _, char := range question {
		if unicode.Is(unicode.Han, char) { // 检查字符是否是中文字符
			isChinese = true
		} else if unicode.Is(unicode.Latin, char) { // 检查字符是否是拉丁字符（英文）
			isEnglish = true
		}
	}

	if isChinese && isEnglish {
		return false
	} else if isChinese {
		return false
	} else if isEnglish {
		return true
	} else {
		return true
	}

}

func TokenCaculate(text string) (int, error) {
	encoding := "gpt-3.5-turbo"

	// if you don't want download dictionary at runtime, you can use offline loader
	// tiktoken.SetBpeLoader(tiktoken_loader.NewOfflineLoader())
	tke, err := tiktoken.EncodingForModel(encoding)
	if err != nil {
		err = fmt.Errorf("getEncoding: %v", err)
		return 0, err
	}
	// encode
	token := tke.Encode(text, nil, nil)
	return len(token) * 4, nil
}
