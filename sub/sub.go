package main

import (
	"context"
	"crypto/tls"
	"flag"
	"log"
	utils "quic-splitter/constants"

	"github.com/quic-go/quic-go"
)

func main() {

	RelayAddr := flag.String("relay", "localhost:5000", "Relay Address")
	PubKey := flag.String("pubkey", "testdata", "Publish Key")

	flag.Parse()

	tlsConfig := tls.Config{
		NextProtos: []string{"quic-splitter"},
	}

	QuicConfig := quic.Config{
		EnableDatagrams: true,
	}

	conn, err := quic.DialAddr(context.TODO(), *RelayAddr, &tlsConfig, &QuicConfig)

	if err != nil {
		log.Fatalf("[Error connecting to Relay : %s", err)
	}

	log.Printf("[Connected to Relay : %s]", *RelayAddr)

	ControlStream, err := conn.OpenStream()

	utils.WriteInt(ControlStream, utils.SPLITTER_SUBSCRIBER)
	utils.WriteString(ControlStream, *PubKey)

	res, err := utils.ReadInt(ControlStream)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	switch res {
	case utils.SUBSCRIBE_DONE:
		log.Printf("[Subscription Successful][%s]", *PubKey)
	default:
		log.Printf("[Error Subscribing][%s]", utils.GetMessage(res))
		return
	}

	for {
		data, err := conn.ReceiveDatagram(context.TODO())

		if err != nil {
			log.Printf("[Error Receiving Datagram][%s]", err)
			return
		}

		str := string(data)

		log.Printf("Subscriber - %s", str)
	}

}
