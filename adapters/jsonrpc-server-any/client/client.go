package client

import (
	. "github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/jsonrpc-any"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	log "github.com/sirupsen/logrus"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"
)

type Client struct {
	UUID             string
	Timeout          time.Duration
	Limit            int
	Addr             string
	MaxRetryInterval time.Duration
	Accept           []router.RoutingRule

	r       *rpc.Client
	conn    net.Conn
	ready   bool
	dialing sync.Mutex
	timer   time.Duration
}

func (c *Client) Init(addr string, uuid string) {
	c.Timeout = time.Hour
	c.Limit = 100
	c.Addr = addr
	c.MaxRetryInterval = 5 * time.Minute
	c.UUID = uuid
	c.ready = false
	c.timer = time.Second
}

func (c *Client) Dial() {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[RiriSDK-Dial] Error %v", err)
			go c.recoverNetwork()
			return
		}
		c.timer = time.Second
		c.ready = true
		log.Infof("[RiriSDK-Dial] Ready!")
	}()
	log.Infof("[RiriSDK-Dial]")
	c.ConnectRPC()
	c.InitChannel(c.UUID)
}

func (c *Client) Close() {
	log.Infof("[RiriSDK-Close]")
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[RiriSDK-Close] Error %v", err)
			return
		}
	}()
	c.conn.Close()
}

func (c *Client) ConnectRPC() {
	log.Infof("[RiriSDK-ConnectRPC]")
	var err error
	c.conn, err = net.Dial("tcp", c.Addr)
	log.Infof("[RiriSDK-ConnectRPC] Connected to RPC at %v", c.Addr)
	if err != nil {
		log.Errorf("[RiriSDK-ConnectRPC] Error %v", err)
		panic(err)
	}
	c.r = jsonrpc.NewClient(c.conn)
}

func (c *Client) recoverNetwork() {
	log.Warning("[RiriSDK-recoverNetwork]")
	c.ready = false
	c.dialing.Lock()
	if c.ready == false {
		log.Warning("[RiriSDK-recoverNetwork] Recovering")
		c.Close()
		c.countDown()
		c.Dial()
	} else {
		log.Infof("[RiriSDK-recoverNetwork]")
	}
	c.dialing.Unlock()
}

func (c *Client) countDown() {
	log.Infof("[RiriSDK-countDown] %s", c.timer)
	t := time.NewTimer(c.timer)
	<-t.C
	c.timer = c.timer * 2
	if c.timer > c.MaxRetryInterval {
		c.timer = c.MaxRetryInterval
	}
}

func (c *Client) InitChannel(key string) (msg string, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[RiriSDK-InitChannel] error: %v", err)
			c.recoverNetwork()
			return
		}
	}()
	args := &ChannelInitRequest{
		UUID:     key,
		Producer: true,
		Consumer: true,
		Accept:   c.Accept,
	}
	reply := ChannelInitResponse{}

	err = c.r.Call("Broker.InitChannel", args, &reply)
	if err != nil {
		panic(err)
	}
	c.UUID = reply.UUID
	log.Debugf("[RiriSDK-InitChannel] %v", c.UUID)
	return
}

func (c *Client) GetUpdates() (packets []types.Packet, err error) {
	packets, err = c.GetChannelUpdates(c.UUID, c.Timeout.String(), c.Limit)
	return
}

func (c *Client) GetChannelUpdates(uuid string, timeout string, limit int) (packets []types.Packet, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[RiriSDK-GetUpdates] error: %v", err)
			c.recoverNetwork()
			return
		}
	}()
	log.Debugf("[RiriSDK-GetUpdates] Preparing updates")
	args := &ChannelConsumeRequest{
		UUID:       uuid,
		TimeoutStr: timeout,
		Limit:      limit,
	}
	reply := ChannelConsumeResponse{}
	err = c.r.Call("Broker.GetUpdates", args, &reply)
	if err != nil {
		panic(err)
	}
	packets = reply.Packets
	log.Debugf("[RiriSDK-GetUpdates] %v", packets)
	return
}

func (c *Client) GetUpdatesChan(bufferSize int) (UpdatesChannel, error) {
	ch := make(chan *types.Packet, bufferSize)

	go func() {
		for {
			updates, err := c.GetUpdates()
			if err != nil {
				log.Println(err)
				log.Println("Failed to get updates, retrying in 3 seconds...")
				time.Sleep(time.Second * 3)

				continue
			}

			for _, update := range updates {
				ch <- &update
			}
		}
	}()

	return ch, nil
}

func (c *Client) MakeRequest(request ChannelProduceRequest) (reply ChannelProduceResponse, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[RiriSDK-MakeRequest] error: %v", err)
			c.recoverNetwork()
			return
		}
	}()
	args := &request
	err = c.r.Call("Broker.Send", args, &reply)
	if err != nil {
		panic(err)
	}
	log.Debugf("[RiriSDK-MakeRequest]")
	return
}
