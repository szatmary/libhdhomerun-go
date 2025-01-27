package hdhomerun

const (
	DISCOVER_UDP_PORT = 65001
	CONTROL_TCP_PORT  = 65001

	TAG_DEVICE_TYPE                = 0x01
	TAG_DEVICE_ID                  = 0x02
	TAG_GETSET_NAME                = 0x03
	TAG_GETSET_VALUE               = 0x04
	TAG_GETSET_LOCKKEY             = 0x15
	TAG_ERROR_MESSAGE              = 0x05
	TAG_TUNER_COUNT                = 0x10
	TAG_LINEUP_URL                 = 0x27
	TAG_STORAGE_URL                = 0x28
	TAG_DEVICE_AUTH_BIN_DEPRECATED = 0x29
	TAG_BASE_URL                   = 0x2A
	TAG_DEVICE_AUTH_STR            = 0x2B
	TAG_STORAGE_ID                 = 0x2C

	TYPE_DISCOVER_REQ = 0x0002
	TYPE_DISCOVER_RPY = 0x0003
	TYPE_GETSET_REQ   = 0x0004
	TYPE_GETSET_RPY   = 0x0005
	TYPE_UPGRADE_REQ  = 0x0006
	TYPE_UPGRADE_RPY  = 0x0007

	DEVICE_TYPE_WILDCARD = 0xFFFFFFFF
	DEVICE_TYPE_TUNER    = 0x00000001
	DEVICE_TYPE_STORAGE  = 0x00000005

	DEVICE_ID_WILDCARD = 0xFFFFFFFF
)
