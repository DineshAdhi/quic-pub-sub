package utils

import (
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/quic-go/quic-go"
)

var SessionMap map[string]quic.Connection
var smap sync.RWMutex

func RegisterSession(conn quic.Connection) string {
	smap.Lock()
	defer smap.Unlock()

	id := uuid.New().String()

	SessionMap[id] = conn

	log.Printf("[New Session Registered][%s]", id)

	return id
}

func DeleteSession(id string) {
	smap.Lock()
	defer smap.Unlock()

	delete(SessionMap, id)

	log.Printf("[Session Deleted][%s]", id)
}

func GetSession(id string) quic.Connection {
	smap.RLock()
	defer smap.RUnlock()

	return SessionMap[id]
}

var PubMap map[string][]string
var mutex sync.RWMutex

func RegisterPublisher(pubkey string) uint8 {
	mutex.Lock()
	defer mutex.Unlock()

	_, ok := PubMap[pubkey]

	if ok {
		return ERROR_PUBKEY_ALRREADY_EXISTS
	}

	PubMap[pubkey] = []string{}

	return PUBLISHER_REGISTERED
}

func DeletePublisher(pubkey string) error {
	mutex.Lock()
	defer mutex.Unlock()

	_, ok := PubMap[pubkey]

	if ok {
		delete(PubMap, pubkey)
		log.Printf("[Publisher Deleted][%s]", pubkey)
		return nil
	}

	return fmt.Errorf("[Delete][Publisher does not exist][%s]", pubkey)
}

func AddSubscriber(pubkey string, sid string) uint8 {
	mutex.Lock()
	defer mutex.Unlock()

	if _, ok := PubMap[pubkey]; ok {
		PubMap[pubkey] = append(PubMap[pubkey], sid)
		return SUBSCRIBE_DONE
	}

	return PUBLISHER_NOT_FOUND
}

func DeleteSubscriber(pubkey string, sid string) uint8 {
	mutex.Lock()
	defer mutex.Unlock()

	slist := PubMap[pubkey]
	updatedList := []string{}

	for _, id := range slist {
		if id == sid {
			continue
		}

		updatedList = append(updatedList, id)
	}

	PubMap[pubkey] = updatedList

	return PUBLISHER_NOT_FOUND
}

func PublishData(pubkey string, data []byte) {
	mutex.RLock()
	defer mutex.RUnlock()

	slist := PubMap[pubkey]

	for _, sid := range slist {
		conn := GetSession(sid)

		if conn != nil {
			conn.SendDatagram(data)
		}
	}
}

// Read - Write Utils

func WriteInt(stream quic.Stream, data uint8) error {
	d := []byte{data}
	_, err := stream.Write(d)

	return err
}

func WriteString(stream quic.Stream, str string) error {
	l := uint8(len(str))
	WriteInt(stream, l)

	data := []byte(str)

	_, err := stream.Write(data)

	return err
}

func ReadInt(stream quic.Stream) (uint8, error) {
	d := make([]byte, 1)
	n, err := stream.Read(d)

	if n < 1 || err != nil {
		return 0, fmt.Errorf("[Error Reading Int][%s]", err)
	}

	return d[0], nil
}

func ReadString(stream quic.Stream) (string, error) {
	l, err := ReadInt(stream)

	if err != nil {
		return "", err
	}

	data := make([]byte, l)

	_, err = stream.Read(data)

	if err != nil {
		return "", err
	}

	return string(data), nil
}
