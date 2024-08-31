package main

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	utils "quic-splitter/constants"
	"time"

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

	if err != nil {
		log.Fatalf("[Error Opening Control Stream][%s]", err)
	}

	utils.WriteInt(ControlStream, utils.SPLITTER_PUBLISHER)
	utils.WriteString(ControlStream, *PubKey)

	response, err := utils.ReadInt(ControlStream)

	if err != nil {
		log.Printf("%s", err)
		return
	}

	switch response {
	case utils.PUBLISHER_REGISTERED:
		log.Printf("[Publisher Registered][Key - %s]", *PubKey)
	default:
		log.Printf("[Error Publishing][%s]", utils.GetMessage(response))
		return
	}

	var itr uint64 = 0

	for {

		str := fmt.Sprintf("%d", itr)
		data := []byte(str)
		itr++

		conn.SendDatagram(data)
		<-time.After(time.Second)
	}
}
