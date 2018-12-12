package chat

import (
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

/*------------------------------通信连接管理器--------------------------------*/

var connCache map[string]*grpc.ClientConn
var mutex sync.RWMutex

type statshandler struct {
}

func (h *statshandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	//fmt.Println("tagconn:")
	return context.WithValue(ctx, "address", info.RemoteAddr.String()) //往ctx存储该链接的对端地址
}

func (h *statshandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	//fmt.Println("tagrpc:")
	return ctx
}

func (h *statshandler) HandleConn(ctx context.Context, s stats.ConnStats) {
	//fmt.Println("HandleConn")
	switch s.(type) {
	case *stats.ConnEnd:
		mutex.Lock()
		defer mutex.Unlock()
		value := ctx.Value("address") //从ctx取对端地址
		address := value.(string)

		delete(connCache, address) //从connCache中删除对应的链接
		//fmt.Println("Conn End-conn-address:", connCache[address], address)
	}
}

func (h *statshandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	//Do Nothing
}

/*
*获取通信连接
*address: 对端节点的IP地址
*timeout： 超时时间
*返回值： 通信连接
 */
func getConn(ip string, timout time.Duration) (*grpc.ClientConn, error) {

	targetHost := fmt.Sprintf("%s:%s", ip, chatPort) //将IP和通信port进行连接

	mutex.RLock()
	conn, ok := connCache[targetHost]
	mutex.RUnlock()
	if ok { //如果在缓存中找到，则直接返回conn
		//fmt.Println("getConn ok", conn, ip)
		return conn, nil
	} else { //如果没有找到，则新建立连接，并加入缓存
		var opts []grpc.DialOption
		h := &statshandler{}
		opts = append(opts, grpc.WithInsecure())
		opts = append(opts, grpc.WithTimeout(timout))
		opts = append(opts, grpc.WithBlock())
		opts = append(opts, grpc.WithStatsHandler(h))

		conn, err := grpc.Dial(targetHost, opts...)
		if err == nil {
			mutex.Lock()
			connCache[targetHost] = conn
			mutex.Unlock()
			//fmt.Println("conn-address:", conn, address)
			return conn, nil
		} else {
			return nil, err
		}
	}
}

/*
*释放通信连接
*conn: 通信连接
 */
func releaseConn(conn *grpc.ClientConn) {
	//conn.Close()
}
