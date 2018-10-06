package cmd

import (
	"github.com/projectriri/bot-gateway/types/ubm-api"
)

type Command struct {
	Cmd     []ubm_api.RichTextElement   `json:"cmd,omitempty"`
	CmdStr  string                      `json:"cmd_str,omitempty"`
	Args    [][]ubm_api.RichTextElement `json:"args,omitempty"`
	ArgsTxt []string                    `json:"args_txt,omitempty"`
	ArgsStr string                      `json:"args_str,omitempty"`
	Message *ubm_api.Message            `json:"message"`
}
