package ubmapi

import "fmt"

func (m *CID) FormatString() string {
	return fmt.Sprintf("%s://%s:%s", m.Messenger, m.ChatType, m.ChatId)
}

func (m *UID) FormatString() string {
	return fmt.Sprintf("%s://%s", m.Messenger, m.Id)
}
