package client

import (
	. "github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/jsonrpc-any"
	"github.com/projectriri/bot-gateway/types"
	log "github.com/sirupsen/logrus"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sync"
	"time"
)

var (
	UUID             string
	Timeout          = time.Hour
	Limit            = 100
	Addr             string
	MaxRetryInterval = 5 * time.Minute
)

var (
	r       *rpc.Client
	conn    net.Conn
	ready   = false
	dialing sync.Mutex
	timer   = time.Second
)

func Init(addr string, uuid string) {
	Addr = addr
	UUID = uuid
	Dial(UUID)
}

func Dial(uuid string) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[RiriSDK-Dial] Error %v", err)
			go recoverNetwork()
			return
		}
		timer = time.Second
		ready = true
		log.Infof("[RiriSDK-Dial] Ready!")
	}()
	log.Infof("[RiriSDK-Dial]")
	ConnectRPC(Addr)
	InitChannel(uuid)
}

func Close() {
	log.Infof("[RiriSDK-Close]")
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[RiriSDK-Close] Error %v", err)
			return
		}
	}()
	conn.Close()
}

func ConnectRPC(addr string) {
	log.Infof("[RiriSDK-ConnectRPC]")
	var err error
	conn, err = net.Dial("tcp", addr)
	log.Infof("[RiriSDK-ConnectRPC] Connected to RPC at %v", addr)
	if err != nil {
		log.Errorf("[RiriSDK-ConnectRPC] Error %v", err)
		panic(err)
	}
	r = jsonrpc.NewClient(conn)
}

func recoverNetwork() {
	log.Warning("[RiriSDK-recoverNetwork]")
	ready = false
	dialing.Lock()
	if ready == false {
		log.Warning("[RiriSDK-recoverNetwork] Recovering")
		Close()
		countDown()
		Dial(UUID)
	} else {
		log.Infof("[RiriSDK-recoverNetwork]")
	}
	dialing.Unlock()
}

func countDown() {
	log.Infof("[RiriSDK-countDown] %s", timer)
	t := time.NewTimer(timer)
	<-t.C
	timer = timer * 2
	if timer > MaxRetryInterval {
		timer = MaxRetryInterval
	}
}

func InitChannel(key string) (msg string, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[RiriSDK-InitChannel] error: %v", err)
			recoverNetwork()
			return
		}
	}()
	args := &ChannelInitRequest{
		UUID:     key,
		Producer: true,
		Consumer: true,
	}
	reply := ChannelInitResponse{}

	err = r.Call("Broker.InitChannel", args, &reply)
	if err != nil {
		panic(err)
	}
	UUID = reply.UUID
	log.Debugf("[RiriSDK-InitChannel] %v", UUID)
	return
}

func GetUpdates() (packets []types.Packet, err error) {
	packets, err = GetChannelUpdates(UUID, Timeout, Limit)
	return
}

func GetChannelUpdates(uuid string, timeout time.Duration, limit int) (packets []types.Packet, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[RiriSDK-GetUpdates] error: %v", err)
			recoverNetwork()
			return
		}
	}()
	log.Debugf("[RiriSDK-GetUpdates] Preparing updates")
	args := &ChannelConsumeRequest{
		UUID:    uuid,
		Timeout: timeout,
		Limit:   limit,
	}
	reply := ChannelConsumeResponse{}
	err = r.Call("Broker.GetUpdates", args, &reply)
	if err != nil {
		panic(err)
	}
	packets = reply.Packets
	log.Debugf("[RiriSDK-GetUpdates] %v", packets)
	return
}

func GetUpdatesChan(bufferSize int) (UpdatesChannel, error) {
	ch := make(chan *types.Packet, bufferSize)

	go func() {
		for {
			updates, err := GetUpdates()
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

func MakeRequest(request ChannelProduceRequest) (reply ChannelProduceResponse, err error) {
	defer func() {
		if err := recover(); err != nil {
			log.Errorf("[RiriSDK-MakeRequest] error: %v", err)
			recoverNetwork()
			return
		}
	}()
	args := &request
	err = r.Call("Broker.Send", args, &reply)
	if err != nil {
		panic(err)
	}
	log.Debugf("[RiriSDK-PushMessage]")
	return
}
