package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	quicsplitter "quic-splitter/constants"
	utils "quic-splitter/constants"

	"github.com/quic-go/quic-go"
)

func main() {

	utils.PubMap = make(map[string][]string)
	utils.SessionMap = make(map[string]quic.Connection)

	Port := flag.Int("port", 5000, "Listening Port")
	CertPath := flag.String("cert", "./certs/", "Certificate Path")

	flag.Parse()

	ListenAddr := fmt.Sprintf("0.0.0.0:%d", *Port)
	certs, err := tls.LoadX509KeyPair(*CertPath+"/localhost.crt", *CertPath+"/localhost.key")

	flag.Parse()

	if err != nil {
		log.Fatalf("[Error Loading Certificates. Location : %s]", *CertPath)
	}

	TLSConfig := tls.Config{
		Certificates: []tls.Certificate{certs},
		NextProtos:   []string{"quic-splitter"},
	}

	QuicConfig := quic.Config{
		EnableDatagrams: true,
	}

	Listener, err := quic.ListenAddr(ListenAddr, &TLSConfig, &QuicConfig)

	if err != nil {
		log.Fatalf("[Error Listening to Address : %s][%s]", ListenAddr, err)
	}

	log.Printf("[QUIC Server Listening : %s]", ListenAddr)

	for {
		Conn, err := Listener.Accept(context.TODO())

		if err != nil {
			log.Fatal("Error Acceping Connection")
		}

		go handleConn(Conn)
	}
}

func handleConn(Conn quic.Connection) {
	log.Printf("[New Quic Connection][%s]", Conn.RemoteAddr())

	ControlStream, err := Conn.AcceptStream(context.TODO())

	if err != nil {
		log.Printf("Errpr Accepign Control Stream %s", err)
		return
	}

	ctype, err := utils.ReadInt(ControlStream)

	if err != nil {
		log.Printf("[Error Reading control Stream][%s]", err)
		return
	}

	switch ctype {
	case quicsplitter.SPLITTER_PUBLISHER:
		handlePublisher(Conn, ControlStream)
	case quicsplitter.SPLITTER_SUBSCRIBER:
		handleSubscriber(Conn, ControlStream)
	default:
		log.Print("Unknown Tyoe")
		return
	}
}

func handlePublisher(Conn quic.Connection, ControlStream quic.Stream) {

	sid := utils.RegisterSession(Conn)

	pubkey, err := utils.ReadString(ControlStream)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	log.Printf("[New Publisher][%s][Pub key - %s]", sid, pubkey)

	res := utils.RegisterPublisher(pubkey)
	utils.WriteInt(ControlStream, res)

	switch res {
	case utils.PUBLISHER_REGISTERED:
		log.Printf("[Pubsliher Registration Successful][Key - %s]", pubkey)
	case utils.ERROR_PUBKEY_ALRREADY_EXISTS:
		log.Printf("[Publisher Registration Rejected][%s]", utils.GetMessage(res))
		return
	default:
		return
	}

	for {
		data, err := Conn.ReceiveDatagram(context.TODO())

		if err != nil {
			log.Printf("Error Receiving Datagrams %s", err)
			utils.DeletePublisher(pubkey)
			utils.DeleteSession(sid)
			return
		}

		utils.PublishData(pubkey, data)
	}
}

func handleSubscriber(Conn quic.Connection, ControlStream quic.Stream) {
	sid := utils.RegisterSession(Conn)

	pubkey, err := utils.ReadString(ControlStream)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	log.Printf("[New Subscriber][%s]", sid)

	res := utils.AddSubscriber(pubkey, sid)
	utils.WriteInt(ControlStream, res)

	switch res {
	case utils.SUBSCRIBE_DONE:
		log.Printf("[Subscription Successful][%s]", sid)
	default:
		log.Printf("[Error Subscribing][%s]", utils.GetMessage(res))
		return
	}
}
