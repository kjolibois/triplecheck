package gemini_manager

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

var thinking_prompt string = "think with step  by step before you answer using <thinking> tags and then answer using <answer> tags. make sure to close the tags with </"

func CallAI(chat_message string) (string, error) {

	ctx := context.Background()
	// Access your API key as an environment variable (see "Set up your API key" above)
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	fmt.Println(thinking_prompt + " " + chat_message)

	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	// The Gemini 1.5 models are versatile and work with both text-only and multimodal prompts
	model := client.GenerativeModel("gemini-1.5-flash")
	resp, err := model.GenerateContent(ctx, genai.Text(thinking_prompt+" "+chat_message))
	if err != nil {
		log.Fatal(err)
	}
	printResponse(resp)
	return getResponseAsString(resp), nil
}
func getResponseAsString(resp *genai.GenerateContentResponse) string {
	var builder strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				switch v := part.(type) {
				case genai.Text:
					builder.WriteString(string(v))
				default:
					builder.WriteString(fmt.Sprintf("Unexpected part type: %T\n", v))
				}
			}
		}
	}
	return builder.String()
}
func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)

			}
		}
	}
	fmt.Println("---")
}
