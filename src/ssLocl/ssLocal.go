package ssLocal

import (
	"net"
	"io"
	"log"
	"os"
	"time"
	"strconv"
	"math/rand"
	ss "github.com/shadowsocks/shadowsocks-go/shadowsocks"
	"errors"
	"encoding/binary"
	"helper"
)

var debug ss.DebugLog
var cipher *ss.Cipher

var (
	errAddrType      = errors.New("socks addr type not supported")
	errVer           = errors.New("socks version not supported")
	errMethod        = errors.New("socks only support 1 method now")
	errAuthExtraData = errors.New("socks authentication get extra data")
	errReqExtraData  = errors.New("socks request get extra data")
	errCmd           = errors.New("socks command not supported")
)

var (
	//  ss password
	ssPassword []string = []string{
	}
	// domain.com:3308
	ssServer []string = []string{
	}
	// :1080  127.0.0.1:1080
	localPort []string = []string{
	}
)

const (
	socksVer5       = 5
	socksCmdConnect = 1
)

type Server struct {
	password      string
	encryptMethod string
	port          int
	listenIP      string
}

func (cls *Server) Listen(password string, encryptMethod string) {
	cls.password = password
	cls.encryptMethod = encryptMethod

	listenAddress := cls.getListenIP() + ":" + strconv.Itoa(cls.getPort())

	ln, err := net.Listen("tcp", listenAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("starting local socks5 server at %v ...\n", listenAddress)
	for {
		// 等待每次请求
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		//异步处理每次请求
		go handleConnection(conn, 0)
	}
}

func (cls *Server) getListenIP() string {
	cls.listenIP = "127.0.0.1"
	return cls.listenIP
}

func (cls *Server) getPort() int {
	cls.port = 1082 //todo comment out
	return cls.port

	if cls.port == 0 {
		var err error
		cls.port, err = helper.GetFreePort()
		if err != nil {
			log.Fatal(err)
		}
	}
	return cls.port
}

func init() {
	rand.Seed(time.Now().Unix())
}

func handShake(conn net.Conn) (err error) {
	const (
		idVer     = 0
		idNmethod = 1
	)
	// version identification and method selection message in theory can have
	// at most 256 methods, plus version and nmethod field in total 258 bytes
	// the current rfc defines only 3 authentication methods (plus 2 reserved),
	// so it won't be such long in practice

	buf := make([]byte, 258)

	var n int
	ss.SetReadTimeout(conn)
	// make sure we get the nmethod field
	if n, err = io.ReadAtLeast(conn, buf, idNmethod+1); err != nil {
		return
	}
	if buf[idVer] != socksVer5 {
		return errVer
	}
	nmethod := int(buf[idNmethod])
	msgLen := nmethod + 2
	if n == msgLen { // handshake done, common case
		// do nothing, jump directly to send confirmation
	} else if n < msgLen { // has more methods to read, rare case
		if _, err = io.ReadFull(conn, buf[n:msgLen]); err != nil {
			return
		}
	} else { // error, should not get extra data
		return errAuthExtraData
	}
	// send confirmation: version 5, no authentication required
	_, err = conn.Write([]byte{socksVer5, 0})
	return
}

func getRequest(conn net.Conn) (rawaddr []byte, host string, err error) {
	const (
		idVer   = 0
		idCmd   = 1
		idType  = 3 // address type index
		idIP0   = 4 // ip addres start index
		idDmLen = 4 // domain address length index
		idDm0   = 5 // domain address start index

		typeIPv4 = 1 // type is ipv4 address
		typeDm   = 3 // type is domain address
		typeIPv6 = 4 // type is ipv6 address

		lenIPv4   = 3 + 1 + net.IPv4len + 2 // 3(ver+cmd+rsv) + 1addrType + ipv4 + 2port
		lenIPv6   = 3 + 1 + net.IPv6len + 2 // 3(ver+cmd+rsv) + 1addrType + ipv6 + 2port
		lenDmBase = 3 + 1 + 1 + 2           // 3 + 1addrType + 1addrLen + 2port, plus addrLen
	)
	// refer to getRequest in server.go for why set buffer size to 263
	buf := make([]byte, 263)
	var n int
	ss.SetReadTimeout(conn)
	// read till we get possible domain length field
	if n, err = io.ReadAtLeast(conn, buf, idDmLen+1); err != nil {
		return
	}
	// check version and cmd
	if buf[idVer] != socksVer5 {
		err = errVer
		return
	}
	if buf[idCmd] != socksCmdConnect {
		err = errCmd
		return
	}

	reqLen := -1
	switch buf[idType] {
	case typeIPv4:
		reqLen = lenIPv4
	case typeIPv6:
		reqLen = lenIPv6
	case typeDm:
		reqLen = int(buf[idDmLen]) + lenDmBase
	default:
		err = errAddrType
		return
	}

	if n == reqLen {
		// common case, do nothing
	} else if n < reqLen { // rare case
		if _, err = io.ReadFull(conn, buf[n:reqLen]); err != nil {
			return
		}
	} else {
		err = errReqExtraData
		return
	}

	rawaddr = buf[idType:reqLen]

	if debug {
		switch buf[idType] {
		case typeIPv4:
			host = net.IP(buf[idIP0: idIP0+net.IPv4len]).String()
		case typeIPv6:
			host = net.IP(buf[idIP0: idIP0+net.IPv6len]).String()
		case typeDm:
			host = string(buf[idDm0: idDm0+buf[idDmLen]])
		}
		port := binary.BigEndian.Uint16(buf[reqLen-2: reqLen])
		host = net.JoinHostPort(host, strconv.Itoa(int(port)))
	}

	return
}

func handleConnection(conn net.Conn, serverId int) {
	if debug {
		debug.Printf("socks connect from %s\n", conn.RemoteAddr().String())
	}
	closed := false
	defer func() {
		if !closed {
			conn.Close()
		}
	}()

	var err error = nil
	// borrow & ss-local 握手
	if err = handShake(conn); err != nil {
		log.Println("socks handshake:", err)
		return
	}

	// 从socket中获取请求地址
	rawaddr, addr, err := getRequest(conn)
	if err != nil {
		log.Println("error getting request:", err)
		return
	}
	// Sending connection established message immediately to client.
	// This some round trip time for creating socks connection with the client.
	// But if connection failed, the client will get connection reset error.
	_, err = conn.Write([]byte{0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x08, 0x43})
	if err != nil {
		debug.Println("send connection confirmation:", err)
		return
	}

	cipher, err = ss.NewCipher("rc4-md5", ssPassword[serverId])
	if err != nil {
		os.Exit(100)
	}

	//fmt.Println(rand.Intn(2))

	//remote, err := connectToServer(0, rawaddr, addr)
	remote, err := ss.DialWithRawAddr(rawaddr, ssServer[serverId], cipher.Copy())
	if err != nil {
		log.Println("error connecting to shadowsocks server:", err)
		return
	}
	defer func() {
		if !closed {
			remote.Close()
		}
	}()

	go ss.PipeThenClose(conn, remote)
	ss.PipeThenClose(remote, conn)
	closed = true
	debug.Println("closed connection to", addr)
}
