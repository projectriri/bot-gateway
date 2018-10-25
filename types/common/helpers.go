package common

import "net/http"

func (m *HTTPRequest) SetHeader(header http.Header) {
	h := make(map[string]string)
	for k, v := range header {
		if len(v) > 0 {
			h[k] = v[0];
		}
	}
	m.Header = h;
}
