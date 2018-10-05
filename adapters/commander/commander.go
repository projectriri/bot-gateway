package main

import (
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/cmd"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	"github.com/projectriri/bot-gateway/utils"
	log "github.com/sirupsen/logrus"
	"strings"
)

var (
	BuildTag      string
	BuildDate     string
	GitCommitSHA1 string
	GitTag        string
)

type CommanderPlugin struct {
	config           Config
	allowEmptyPrefix bool
}

var manifest = types.Manifest{
	BasicInfo: types.BasicInfo{
		Name:    "commander",
		Author:  "Project Riri Staff",
		Version: "v0.1",
		License: "MIT",
		URL:     "https://github.com/projectriri/bot-gateway/adapters/commander",
	},
	BuildInfo: types.BuildInfo{
		BuildTag:      BuildTag,
		BuildDate:     BuildDate,
		GitCommitSHA1: GitCommitSHA1,
		GitTag:        GitTag,
	},
}

func (p *CommanderPlugin) GetManifest() types.Manifest {
	return manifest
}

func (p *CommanderPlugin) Init(filename string, configPath string) {
	// load toml config
	_, err := toml.DecodeFile(configPath+"/"+filename+".toml", &p.config)
	if err != nil {
		panic(err)
	}
	p.allowEmptyPrefix = p.checkAllowEmptyPrefix()
}

func (p *CommanderPlugin) Start() {
	log.Infof("[commander] registering consumer channel %v", p.config.ChannelUUID)
	cc := router.RegisterConsumerChannel(p.config.ChannelUUID, []router.RoutingRule{
		{
			From: ".*",
			To:   ".*",
			Formats: []types.Format{
				{
					API:     "ubm-api",
					Version: "1.0",
					Method:  "receive",
				},
			},
		},
	})
	defer cc.Close()
	log.Infof("[commander] registered consumer channel %v", cc.UUID)
	log.Infof("[commander] registering producer channel %v", p.config.ChannelUUID)
	pc := router.RegisterProducerChannel(p.config.ChannelUUID, false)
	defer pc.Close()
	log.Infof("[commander] registered producer channel %v", pc.UUID)
	for {
		packet := cc.Consume()
		req := ubm_api.UBM{}
		err := json.Unmarshal(packet.Body, &req)
		if err != nil {
			log.Errorf("[commander] message %v has an incorrect body type %v", packet.Head.UUID, err)
		}
		if req.Type == "message" && req.Message != nil {
			if req.Message.Type != "rich_text" || req.Message.RichText == nil {
				continue
			}

			richTexts := *req.Message.RichText
			// Trim all leading white characters
			for i := 0; i < len(richTexts) && richTexts[i].Type == "text"; i++ {
				richTexts[i].Text = strings.TrimLeftFunc(richTexts[i].Text, p.isWhiteChar)
				if len(richTexts[i].Text) == 0 {
					richTexts = richTexts[1:]
					i--
				} else {
					break
				}
			}
			if len(richTexts) == 0 {
				continue
			}
			if richTexts[0].Type == "text" {
				// If the first rich text element is text, trim the command prefix
				if pfx, ok := p.checkPrefix(richTexts[0].Text); !ok {
					continue
				} else {
					if richTexts[0].Text == pfx {
						richTexts = richTexts[1:]
					} else {
						richTexts[0].Text = richTexts[0].Text[len(pfx):]
					}
				}
			} else {
				// else check allowEmptyPrefix
				if !p.allowEmptyPrefix {
					continue
				}
			}

			// Process this command
			parsedCommand := make([][]ubm_api.RichTextElement, 1)
			parsedCommand[0] = make([]ubm_api.RichTextElement, 0)
			// Process rich text array
			lastEscape := false
			lastWhiteChar := false
			inQuote := false
			var lastQuoteChar rune
			buffer := make([]rune, 0)
			nowP := 0
			for _, elem := range richTexts {
				if elem.Type == "text" {
					// Text needs to be parsed
					for _, r := range elem.Text {
						// state operations
						if lastEscape {
							lastEscape = false
							buffer = append(buffer, r)
							continue
						}
						if r == ESCAPE_CHAR {
							lastEscape = true
							continue
						}
						if lastWhiteChar {
							if p.isWhiteChar(r) {
								continue
							} else {
								lastWhiteChar = false
							}
						}
						if inQuote {
							if r == lastQuoteChar {
								// end of quote
								inQuote = false
							} else {
								buffer = append(buffer, r)
							}
							continue
						}
						// state transfer
						if p.isWhiteChar(r) {
							lastWhiteChar = true
							// append and clear buffer
							if len(buffer) > 0 {
								parsedCommand[nowP] = append(parsedCommand[nowP],
									ubm_api.RichTextElement{
										Type: "text",
										Text: string(buffer),
									})
								buffer = make([]rune, 0)
							}
							// append parsedCommand
							parsedCommand = append(parsedCommand, make([]ubm_api.RichTextElement, 0))
							nowP++
						} else if p.isQuoteChar(r) {
							inQuote = true
							lastQuoteChar = r
						} else {
							// normal char
							buffer = append(buffer, r)
						}
					}
				} else {
					// Other type of message, append buffer
					if len(buffer) > 0 {
						parsedCommand[nowP] = append(parsedCommand[nowP],
							ubm_api.RichTextElement{
								Type: "text",
								Text: string(buffer),
							})
						buffer = make([]rune, 0)
					}
					// and append cur elem to the end
					parsedCommand[nowP] = append(parsedCommand[nowP], elem)
					// and clear some states
					lastEscape = false
					lastWhiteChar = false
				}
			}
			// Append buffer in the end
			if len(buffer) > 0 {
				parsedCommand[nowP] = append(parsedCommand[nowP],
					ubm_api.RichTextElement{
						Type: "text",
						Text: string(buffer),
					})
				buffer = make([]rune, 0)
			}

			// compose response according to config.ResponseMode in bit mask
			c := cmd.Command{}
			if p.config.ResponseMode&RESPONSE_CMD != 0 {
				c.Cmd = parsedCommand[0]
			}
			if p.config.ResponseMode&RESPONSE_CMDSTR != 0 {
				for _, elem := range parsedCommand[0] {
					if elem.Type == "text" {
						c.CmdStr += elem.Text
					}
				}
			}
			if p.config.ResponseMode&RESPONSE_ARGS != 0 {
				c.Args = parsedCommand[1:]
			}
			if p.config.ResponseMode&RESPONSE_ARGSTXT != 0 ||
				p.config.ResponseMode&RESPONSE_ARGSSTR != 0 {
				tmpArgsTxt := make([]string, 0)
				for _, aCmd := range parsedCommand[1:] {
					tmp := ""
					for _, elem := range aCmd {
						if elem.Type == "text" {
							tmp += elem.Text
						}
					}
					if len(tmp) != 0 {
						tmpArgsTxt = append(tmpArgsTxt, tmp)
					}
				}
				if p.config.ResponseMode&RESPONSE_ARGSTXT != 0 {
					c.ArgsTxt = tmpArgsTxt
				}
				if p.config.ResponseMode&RESPONSE_ARGSSTR != 0 {
					c.ArgsStr = strings.Join(tmpArgsTxt, " ")
				}
			}
			c.Message = *req.Message

			b, _ := json.Marshal(c)
			pc.Produce(types.Packet{
				Head: types.Head{
					From: packet.Head.From,
					To:   packet.Head.To,
					UUID: utils.GenerateUUID(),
					Format: types.Format{
						API:      "cmd",
						Method:   "cmd",
						Version:  "1.0",
						Protocol: "",
					},
				},
				Body: b,
			})
		}

	}
}

var PluginInstance types.Adapter = &CommanderPlugin{}
