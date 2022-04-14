package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"

	"github.com/lucas-clemente/quic-go"
)

// 这份代码copy自：https://github.com/lucas-clemente/quic-go/blob/master/example/echo/echo.go
// 做了一些简单的注释，代码相对来说还是比较简单的，难的是底层的实现原理

const addr = "localhost:4242"

const message = "foobar"

func main() {
	go func() { log.Fatal(echoServer()) }()

	err := clientMain()
	if err != nil {
		panic(err)
	}
}

// 开启一个server，第一个流收到什么内容都全部原样写回去
func echoServer() error {
	// 在 quic 中监听的都是udp的
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		return err
	}
	// 可以看到这里和UDP的监听是不一样的，流的建立过程是分成了两步
	conn, err := listener.Accept(context.Background())
	if err != nil {
		return err
	}

	// 在建立的连接上监听一个流
	stream, err := conn.AcceptStream(context.Background())
	if err != nil {
		panic(err)
	}

	// loggingWriter只是封装了下io.Writer，为了加入server端的打印信息，好像也没必要怎么做
	_, err = io.Copy(loggingWriter{stream}, stream)
	return err
}

func clientMain() error {
	// 设置tls
	tlsConf := &tls.Config{
		InsecureSkipVerify: true,                          // 不验证服务端的证书和域名信息
		NextProtos:         []string{"quic-echo-example"}, // 设置协商协议ALPN中可选的协议
	}
	conn, err := quic.DialAddr(addr, tlsConf, nil)
	if err != nil {
		return err
	}
	// 连接成功后打开一个流，打开流分成了好几个api，OpenStreamSync, OpenStream, OpenUniStream, OpenUniStreamSync
	// OpenStreamSync 阻塞等待代码一个双向的 QUIC 流
	stream, err := conn.OpenStreamSync(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Client: Sending '%s'\n", message)
	_, err = stream.Write([]byte(message))
	if err != nil {
		return err
	}

	// server端收到client端发送的流会echo回来
	buf := make([]byte, len(message))
	_, err = io.ReadFull(stream, buf)
	if err != nil {
		return err
	}
	fmt.Printf("Client: Got '%s'\n", buf)

	// 注意到，这里并没有关闭流的操作和关闭conn的操作，可能是示例代码的问题
	return nil
}

// 封装 io.Writer 同时输出相应的log信息
type loggingWriter struct{ io.Writer }

func (w loggingWriter) Write(b []byte) (int, error) {
	fmt.Printf("Server: Got '%s'\n", string(b))
	return w.Writer.Write(b)
}

// 生成一个简单的TLS配置给server端
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}
