package hdhomerun

import (
	"fmt"
	"io"
)

func crc(pkt []byte) uint32 {
	s := func(x, mask, crc, coef uint32) uint32 {
		if x&mask != 0 {
			return crc ^ coef
		}
		return crc
	}
	crc := ^uint32(0)
	for _, v := range pkt {
		x := crc ^ uint32(v)
		crc = crc >> 8
		crc = s(x, 0x01, crc, 0x77073096)
		crc = s(x, 0x02, crc, 0xEE0E612C)
		crc = s(x, 0x04, crc, 0x076DC419)
		crc = s(x, 0x08, crc, 0x0EDB8832)
		crc = s(x, 0x10, crc, 0x1DB71064)
		crc = s(x, 0x20, crc, 0x3B6E20C8)
		crc = s(x, 0x40, crc, 0x76DC4190)
		crc = s(x, 0x80, crc, 0xEDB88320)
	}
	return crc ^ 0xFFFFFFFF
}

// type Packet []byte

// Packet writing
func WriteVarLen(v int) []byte {
	if v <= 127 {
		return []byte{byte(v & 0x7f)}
	}
	return []byte{0x80 | byte(v&0x7f), byte(v >> 7)}
}

// func (pkt *Packet) WriteUint8(v uint8) {
// 	*pkt = append(*pkt, v)
// }
// func (pkt *Packet) WriteUint16(v uint16) {
// 	*pkt = append(*pkt, []byte{byte(v >> 8), byte(v)}...)
// }
func WriteUint32(v uint32) []byte {
	return []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}
}

// func (pkt *Packet) WriteUint32(v uint32) {
// 	*pkt = append(*pkt, []byte{byte(v >> 24), byte(v >> 16), byte(v >> 8), byte(v)}...)
// }
// func (pkt *Packet) WriteTag8(t byte, v uint8) {
// 	pkt.WriteUint8(t)
// 	// pkt.WriteVarLen(1)
// 	pkt.WriteUint8(v)
// }
// func (pkt *Packet) WriteTag16(t byte, v uint16) {
// 	pkt.WriteUint8(t)
// 	// pkt.WriteVarLen(2)
// 	pkt.WriteUint16(v)
// }
// func (pkt *Packet) WriteTag32(t byte, v uint32) {
// 	pkt.WriteUint8(t)
// 	// pkt.WriteVarLen(4)
// 	pkt.WriteUint32(v)
// }

// func (pkt *Packet) SealFrame(frameType uint16) {
// 	length := len(*pkt) - 4
// 	(*pkt)[0] = byte(frameType >> 8)
// 	(*pkt)[1] = byte(frameType >> 0)
// 	(*pkt)[2] = byte(length >> 8)
// 	(*pkt)[3] = byte(length >> 0)
// 	crc := crc(*pkt)
// 	*pkt = append(*pkt, []byte{byte(crc), byte(crc >> 8), byte(crc >> 16), byte(crc >> 24)}...)
// }

type Tag struct {
	Type  byte
	Value []byte
}

type Packet struct {
	FrameType uint16
	Tags      []Tag
}

func (pkt *Packet) Find(tagType byte) []byte {
	for _, tag := range pkt.Tags {
		if tag.Type == tagType {
			return tag.Value
		}
	}
	return nil
}

func (pkt *Packet) String() string {
	s := fmt.Sprintf("%d\n", pkt.FrameType)
	for _, tag := range pkt.Tags {
		s += fmt.Sprintf("%s\n", tag.String())
	}
	return s
}

func (tag *Tag) String() string {
	switch tag.Type {
	case TAG_DEVICE_TYPE:
		if len(tag.Value) != 4 {
			return ""
		}
		deviceType := uint32(tag.Value[0]<<24) | uint32(tag.Value[2]<<8) | uint32(tag.Value[2]<<8) | uint32(tag.Value[3])
		switch deviceType {
		case DEVICE_TYPE_WILDCARD:
			return "TAG_DEVICE_TYPE: DEVICE_TYPE_WILDCARD"
		case DEVICE_TYPE_TUNER:
			return "TAG_DEVICE_TYPE: DEVICE_TYPE_TUNER"
		case DEVICE_TYPE_STORAGE:
			return "TAG_DEVICE_TYPE: DEVICE_TYPE_STORAGE"
		default:
			return "TAG_DEVICE_TYPE: DEVICE_TYPE_UNKNOWN"
		}
	case TAG_DEVICE_ID:
		deviceId := uint32(tag.Value[0]<<24) | uint32(tag.Value[2]<<8) | uint32(tag.Value[2]<<8) | uint32(tag.Value[3])
		return fmt.Sprintf("TAG_DEVICE_ID: 0x%08X", deviceId)
	case TAG_GETSET_NAME:
		return fmt.Sprintf("TAG_GETSET_NAME: %v", string(tag.Value))
	case TAG_GETSET_VALUE:
		return fmt.Sprintf("TAG_GETSET_VALUE: %v", string(tag.Value))
	case TAG_GETSET_LOCKKEY:
		return fmt.Sprintf("TAG_GETSET_LOCKKEY: %d", tag.Value[0])
	case TAG_ERROR_MESSAGE:
		return fmt.Sprintf("TAG_ERROR_MESSAGE: %v", string(tag.Value))
	case TAG_TUNER_COUNT:
		return fmt.Sprintf("TAG_TUNER_COUNT: %v", tag.Value[0])
	case TAG_LINEUP_URL:
		return fmt.Sprintf("TAG_LINEUP_URL: %v", string(tag.Value))
	case TAG_STORAGE_URL:
		return fmt.Sprintf("TAG_STORAGE_URL: %v", string(tag.Value))
	case TAG_DEVICE_AUTH_BIN_DEPRECATED:
		return fmt.Sprintf("TAG_DEVICE_AUTH_BIN_DEPRECATED: %v", tag.Value)
	case TAG_BASE_URL:
		return fmt.Sprintf("TAG_BASE_URL: %v", string(tag.Value))
	case TAG_DEVICE_AUTH_STR:
		return fmt.Sprintf("TAG_DEVICE_AUTH_STR: %v", string(tag.Value))
	case TAG_STORAGE_ID:
		return fmt.Sprintf("TAG_STORAGE_ID: %v", tag.Value)
	default:
		return "DEVICE_TAG_UNKNOWN"
	}
}

func (pkt *Packet) UnmarshalBinary(data []byte) error {
	// TODO check CRC
	if len(data) < 4 {
		return io.ErrUnexpectedEOF
	}
	pkt.FrameType = uint16(data[0])<<8 | uint16(data[1])
	length, data := 4+int(uint16(data[2])<<8|uint16(data[3])), data[4:] // the crc is not counted in the size
	if length != len(data) {
		return io.ErrUnexpectedEOF
	}
	for len(data) > 4 {
		tagType, varLen := data[0], uint16(data[1])
		data = data[2:]
		if varLen&0x0080 != 0 {
			if len(data) == 0 {
				return io.ErrUnexpectedEOF
			}
			varLen = (varLen & 0x7f) | uint16(data[0])<<7
			data = data[1:]
		}
		if len(data) < int(varLen) {
			return io.ErrUnexpectedEOF
		}
		pkt.Tags = append(pkt.Tags, Tag{tagType, data[:varLen]})
		data = data[varLen:]
	}
	// The last 4 bytes are the CRC, we should hav already checked them
	if len(data) != 4 {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func (tag *Tag) MarshalBinary() ([]byte, error) {
	var data []byte
	data = append(data, tag.Type)
	data = append(data, WriteVarLen(len(tag.Value))...)
	return append(data, tag.Value...), nil
}

func (pkt *Packet) MarshalBinary() ([]byte, error) {
	data := []byte{byte(pkt.FrameType >> 8), byte(pkt.FrameType), 0, 0}
	for _, tag := range pkt.Tags {
		bin, err := tag.MarshalBinary()
		if err != nil {
			return nil, err
		}
		data = append(data, bin...)
	}
	length := len(data) - 4
	data[2] = byte(length >> 8)
	data[3] = byte(length >> 0)
	crc := crc(data)
	return append(data, []byte{byte(crc), byte(crc >> 8), byte(crc >> 16), byte(crc >> 24)}...), nil
}
