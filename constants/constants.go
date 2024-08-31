package utils

const (
	SPLITTER_PUBLISHER  = uint8(0x23)
	SPLITTER_SUBSCRIBER = uint8(0x32)
)

const (
	PUBLISHER_REGISTERED         = uint8(0x67)
	ERROR_UNKNOWN                = uint8(0x56)
	ERROR_PUBKEY_ALRREADY_EXISTS = uint8(0x64)
	PUBLISHER_NOT_FOUND          = uint8(0x46)
	SUBSCRIBE_DONE               = uint8(0x31)
)

func GetMessage(t uint8) string {
	switch t {
	case PUBLISHER_REGISTERED:
		return "PUBLISHER_REGISTERED"
	case ERROR_UNKNOWN:
		return "ERROR_UNKNOWN"
	case ERROR_PUBKEY_ALRREADY_EXISTS:
		return "ERROR_PUBKEY_ALRREADY_EXISTS"
	case PUBLISHER_NOT_FOUND:
		return "PUBLISHER_NOT_FOUND"
	case SUBSCRIBE_DONE:
		return "SUBSCRIBE_DONE"
	default:
		return "UNKNOWN ERROR MESSAGE"
	}
}
