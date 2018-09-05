package ubm_api

import "net/url"

type UBM struct {
	Type     string    `json:"type"`
	Message  *Message  `json:"message,omitempty"`
	Notice   *Notice   `json:"notice,omitempty"`
	Response *Response `json:"response,omitempty"`
	Action   *Action   `json:"action,omitempty"`
}

type Message struct {
	ID              string    `json:"id,omitempty"`
	From            *User     `json:"from,omitempty"`
	Chat            *Chat     `json:"chat,omitempty"`
	UID             *UID      `json:"uid,omitempty"`
	CID             *CID      `json:"cid,omitempty"`
	Type            string    `json:"type"`
	ReplyID         string    `json:"reply_id,omitempty"`
	EditID          string    `json:"edit_id,omitempty"`
	DeleteID        string    `json:"delete_id,omitempty"`
	ForwardFromChat *Chat     `json:"forward_from_chat,omitempty"`
	ForwardID       string    `json:"forward_id,omitempty"`
	ForwardFrom     *User     `json:"forward_from"`
	RichText        *RichText `json:"rich_text,omitempty"`
	Sticker         *Sticker  `json:"sticker,omitempty"`
	Audio           *Audio    `json:"audio,omitempty"`
	Location        *Location `json:"location,omitempty"`
}

type CID struct {
	Messenger string `json:"messenger"`
	ChatID    string `json:"chat_id"`
	ChatType  string `json:"chat_type"`
}

type Chat struct {
	CID         CID    `json:"cid"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UID struct {
	Messenger string `json:"messenger"`
	ID        string `json:"id"`
	Username  string `json:"username"`
}

type User struct {
	DisplayName string `json:"display_name"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	UID         UID    `json:"uid"`
	PrivateChat CID    `json:"private_chat"`
}

type RichText []RichTextElement

type RichTextElement struct {
	Type       string      `json:"type"`
	Text       string      `json:"text,omitempty"`
	StyledText *StyledText `json:"styled_text,omitempty"`
	At         *At         `json:"at,omitempty"`
	Image      *Image      `json:"image,omitempty"`
}

type StyledText struct {
	Text   string `json:"text"`
	Format string `json:"format"`
}

type At struct {
	DisplayName string `json:"display_name"`
	UID         UID    `json:"uid"`
}

type Image struct {
	Format string   `json:"format"`
	Width  int      `json:"width"`
	Height int      `json:"height"`
	Data   *[]byte  `json:"data"`
	URL    *url.URL `json:"url"`
}

type Sticker struct {
	ID     string `json:"id"`
	PackID string `json:"pack_id"`
	Image  *Image `json:"image"`
}

type Audio struct {
	Format   string   `json:"format"`
	Duration int      `json:"duration"`
	Data     *[]byte  `json:"data"`
	URL      *url.URL `json:"url"`
}

type Location struct {
	Content   string  `json:"content"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Title     string  `json:"title"`
}

type Notice struct {
}

type Action struct {
}

type Response struct {
}
