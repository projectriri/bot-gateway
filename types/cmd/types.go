package cmd

import "github.com/projectriri/bot-gateway/types/ubm-api"

type Command struct {
	Cmd     []ubm_api.RichTextElement   `json:"cmd"`
	CmdStr  string                      `json:"cmd_str"`
	Args    [][]ubm_api.RichTextElement `json:"args"`
	ArgsTxt []string                    `json:"args_txt"`
	ArgsStr string                      `json:"args_str"`
}
