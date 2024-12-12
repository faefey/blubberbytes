/*
This is a SOCKS proxy using go. It logs the total number of ingoing and outgoing bytes
for each user (1 user = 1 IP address) and every 5 minutes this information is logged to
a txt file in the format [IP, bytes]\n[IP, bytes]\n[IP, bytes]
*/

package proxy

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"

	"server/database/operations"

	"github.com/armon/go-socks5"
	"github.com/libp2p/go-libp2p/core/host"
)

var paymentInformation = make(map[string]int64)
var mutex sync.Mutex

type trafficInterceptor struct {
	conn     net.Conn
	clientIP string
	read     int64
	written  int64
}

func (t *trafficInterceptor) Read(b []byte) (n int, err error) {
	n, err = t.conn.Read(b)
	if err == nil {
		t.read += int64(n)
	}
	//log.Printf("Total received: %d", t.read)
	return
}

func (t *trafficInterceptor) Write(b []byte) (n int, err error) {
	n, err = t.conn.Write(b)
	if err == nil {
		t.written += int64(n)
	}
	//log.Printf("Total sent %d", t.written)
	return
}

func (t *trafficInterceptor) Close() error {
	log.Printf("Final bytes received: %d", t.read)
	log.Printf("Final bytes sent: %d", t.written)
	log.Printf("IP Of the bytes above: %s", t.clientIP)

	updatePaymentInfo(strings.Split(t.clientIP, ":")[0], t.read+t.written)

	return t.conn.Close()
}

func (t *trafficInterceptor) LocalAddr() net.Addr {
	return t.conn.LocalAddr()
}

func (t *trafficInterceptor) RemoteAddr() net.Addr {
	return t.conn.RemoteAddr()
}

func (t *trafficInterceptor) SetDeadline(deadline time.Time) error {
	return t.conn.SetDeadline(deadline)
}

func (t *trafficInterceptor) SetReadDeadline(deadline time.Time) error {
	return t.conn.SetReadDeadline(deadline)
}

func (t *trafficInterceptor) SetWriteDeadline(deadline time.Time) error {
	return t.conn.SetWriteDeadline(deadline)
}

func (t *trafficInterceptor) GetBytesSent() int64 {
	return t.written
}

func (t *trafficInterceptor) GetBytesReceived() int64 {
	return t.read
}

type clientAddressRuleset struct {
	socks5.RuleSet
}

func updatePaymentInfo(key string, value int64) {
	mutex.Lock()
	defer mutex.Unlock()

	log.Printf("Key: %s value: %d\n", key, value)

	if currentBytes, exists := paymentInformation[key]; exists {
		paymentInformation[key] = currentBytes + value
	} else {
		paymentInformation[key] = value
	}

}

func (r *clientAddressRuleset) Connect(ctx context.Context, conn net.Conn, target *socks5.AddrSpec) (*socks5.AddrSpec, error) {
	//clientAddr := conn.RemoteAddr()
	//log.Printf("Client address: %s", clientAddr)
	return target, nil
}

func (r *clientAddressRuleset) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	if req.RemoteAddr != nil {
		clientIP := req.RemoteAddr.String()
		log.Printf("Client IP: %s", clientIP)
		return context.WithValue(ctx, "clientIP", clientIP), true
	}

	return ctx, true
}

func customDial(ctx context.Context, network, addr string) (net.Conn, error) {

	//remoteSocksProxy := "23.239.12.179:8000"

	conn, err := net.Dial(network, addr)
	if err != nil {
		return nil, err
	}

	clientIP, _ := ctx.Value("clientIP").(string)
	// Wrap the connection to intercept traffic
	return &trafficInterceptor{conn: conn, clientIP: clientIP}, nil
}

func Proxy(node host.Host, db *sql.DB) {
	dial := customDial
	conf := &socks5.Config{Dial: dial, Rules: &clientAddressRuleset{}}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			mutex.Lock()

			for key, value := range paymentInformation {
				//log.Println("Hello!")
				log.Printf("%s : %d", key, value)
				operations.AddProxyLogs(db, key, value, time.Now().Unix())
			}

			// timeBefore := time.Now().Unix() - (5 * time.Minute).Milliseconds()
			// log.Println("Sending request for proxy payment...")
			// err := p2p.SendProxyBillWithConfirmation(node, )

			for key := range paymentInformation {
				delete(paymentInformation, key)
			}

			mutex.Unlock()
		}
	}()

	fmt.Println("Proxy is running on http://localhost:8000.")

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", "0.0.0.0:8000"); err != nil {
		panic(err)
	}
}
