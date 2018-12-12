package chat

import (
	pb "chat/proto"
	"config"
	"errors"
	"fmt"
	"logging"
	"net"
	"reflect"
	"time"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//	"google.golang.org/grpc/reflection"
)

//peer communication port through chat module.
var chatPort string

//gRPC connection ttl.
var chatTimeout time.Duration

var logger = logging.MustGetLogger()

var server *grpc.Server

//it realized PeerServer interface.
type Server struct{}

//the type of handler
type Handler func(msg interface{}) (pb.Response_Type, interface{}, error)

//request map's value type.
type msgHandler struct {
	handler Handler      //the handler of handling msg from client
	reqType reflect.Type //the type of msg request content(it is a struct)
}

//request map: key is msg type, value is a struct that stores payload body type and callback function.
var reqMsgMap map[pb.Request_Type]*msgHandler = make(map[pb.Request_Type]*msgHandler)

//response map: key is msg type, value is payload body type.
var respMsgMap map[pb.Response_Type]reflect.Type = make(map[pb.Response_Type]reflect.Type)

func Initialize() error {
	chatTimeout = time.Second * time.Duration(config.GetChatTimeout())

	//init connection cache.
	connCache = make(map[string]*grpc.ClientConn)

	//listen chat port
	chatPort = config.GetChatPort()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", chatPort))
	if err != nil {
		logger.Errorf("%v", err)
		return err
	}

	//start grpc server
	concurrentStream := uint32(config.GetChatConcurrent())
	connWindowSize := int32(config.GetChatConnWindowSize())
	writeBufferSize := config.GetChatWriteBufferSize()
	server = grpc.NewServer(grpc.MaxConcurrentStreams(concurrentStream), grpc.InitialConnWindowSize(connWindowSize), grpc.WriteBufferSize(int(writeBufferSize)))
	pb.RegisterPeerServer(server, &Server{})
	go server.Serve(lis)
	return nil
}

//server handles msg from client.
func (server *Server) Chat(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	h, ok := reqMsgMap[req.Type] //find handler struct according msg type.
	if !ok {
		return nil, errors.New(fmt.Sprintf("handler not exist for %s", req.Type))
	}

	msg := reflect.New(h.reqType).Interface()
	err := proto.Unmarshal(req.Payload, msg.(proto.Message))
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	respType, respData, err := h.handler(msg)
	if err == nil {
		payload, err := proto.Marshal(respData.(proto.Message))
		if err == nil {
			resp := &pb.Response{
				Type:    respType,
				Payload: payload,
			}
			return resp, nil
		} else {
			return nil, err
		}
	} else {
		logger.Error("chat handler err:", err)
		return nil, err
	}
}

/*
*Description: register message.
*Parameters: request msg type, callback function, response msg type
*ReturnValue: null
 */

func RegisterMsg(reqMsgType pb.Request_Type, handler Handler, respMsgType pb.Response_Type) {
	var request interface{}
	var response interface{}

	switch reqMsgType {
	case pb.Request_CONSENSUS:
	//		request = pb.ConsensusRequest{}
	//		response = pb.ConsensusResponse{}

	//...
	default:
		logger.Error("Register a unknown msg type...")
		return
	}

	//expand request msg map
	var h *msgHandler = &msgHandler{}
	h.handler = handler
	h.reqType = reflect.TypeOf(request)
	reqMsgMap[reqMsgType] = h

	//expand response msg map
	respMsgMap[respMsgType] = reflect.TypeOf(response)
}

func SendMsg(msgType pb.Request_Type, arg interface{}, address string) (interface{}, error) {
	msg, err := msgPrepare(msgType, arg)

	resp, err := sendMsg(msg, address)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

/*
*package msg
*msgType: msg type
*arg: payload body to be serialized.
 */
func msgPrepare(msgType pb.Request_Type, arg interface{}) (*pb.Request, error) {
	//serialize the payload body.
	payload, err := proto.Marshal(arg.(proto.Message))
	if err != nil {
		logger.Error("Message marshal failed:%s.", msgType)
		return nil, err
	}

	// package the msg.
	msg := &pb.Request{
		Type:    msgType,
		Payload: payload,
	}

	return msg, nil
}

/*
*send smg to specified node
*msg:request body
*address:receiver's address
*retValue：response body, and error msg
 */
func sendMsg(msg *pb.Request, ip string) (interface{}, error) {
	//acquire the connection with a server that's address is "address"
	conn, err := getConn(ip, chatTimeout)
	if err != nil {
		logger.Errorf("Error creating connection to peer address %s: %s", ip, err)
		return nil, err
	}
	//create grpc client，and send message.
	client := pb.NewPeerClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	retMsg, err := client.Chat(ctx, msg, grpc.MaxCallRecvMsgSize(1024*1024*40))
	if err != nil {
		return nil, err
	}
	if retMsg == nil {
		return nil, errors.New("retMsg is nil")
	}

	respTT, ok := respMsgMap[retMsg.Type] //find payload type according msg type.
	if !ok {
		err := errors.New(fmt.Sprintf("resp msg type no found %s", retMsg.Type))
		logger.Error(err)
		return nil, err
	}

	resp := reflect.New(respTT).Interface()

	err = proto.Unmarshal(retMsg.Payload, resp.(proto.Message))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	//release the connection.
	defer releaseConn(conn)
	if err != nil {
		logger.Errorf("error: %v", err)
		return nil, err
	}
	return resp, nil
}
