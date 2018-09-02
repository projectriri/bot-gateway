package router

import (
	"time"
)

type RouterConfig struct {
	BufferSize      uint
	ChannelLifeTime time.Duration
	GCInterval      time.Duration
}
