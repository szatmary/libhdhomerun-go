package hdhomerun

import (
	"errors"
	"fmt"
	"io"
	"net"
	"regexp"
)

// TODO move the networking in here

func Discover(deviceType, deviceId uint32) ([]byte, error) {
	pkt := Packet{
		FrameType: TYPE_DISCOVER_REQ,
		Tags: []Tag{{
			Type:  TAG_DEVICE_TYPE,
			Value: WriteUint32(deviceType),
		}, {
			Type:  TAG_DEVICE_ID,
			Value: WriteUint32(deviceId),
		}},
	}
	return pkt.MarshalBinary()
}

type Device struct {
	DeviceId uint32
	Addr     string
	Conn     net.Conn
	Tuners   int
}

func NewDevice(addr string) (*Device, error) {
	Conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Device{
		Conn: Conn,
	}, nil
}
func (dev *Device) SendReceive(pkt *Packet) (*Packet, error) {
	data, err := pkt.MarshalBinary()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Sending: %v\n", data)
	if _, err = dev.Conn.Write(data); err != nil {
		return nil, err
	}
	head := make([]byte, 4)
	if _, err = io.ReadFull(dev.Conn, head); err != nil {
		return nil, err
	}
	length := 4 + int(head[2])<<8 | int(head[3])
	fmt.Printf("length: %v\n", length)
	body := make([]byte, length)
	if _, err = io.ReadFull(dev.Conn, body); err != nil {
		return nil, err
	}
	data = []byte{}
	data = append(data, head...)
	data = append(data, body...)
	resp := Packet{}
	fmt.Printf("UnmarshalBinary: %v\n", data)
	if err = resp.UnmarshalBinary(data); err != nil {
		return nil, err
	}
	if message := pkt.Find(TAG_ERROR_MESSAGE); len(message) != 0 {
		return nil, errors.New(string(message))
	}
	return &resp, nil
}

type TunerStatus struct {
	Tuner                int
	Channel              string
	Lock                 string
	SignalStrength       int
	SignalToNoiseQuality int
	SymbolErrorQuality   int
	RawBitsPerSecond     uint32
	PacketsPerSecond     uint32
}

func (dev *Device) GetTunerStatus(tuner int, lockkey uint32) (*TunerStatus, error) {
	pkt := &Packet{
		FrameType: TYPE_GETSET_REQ,
		Tags:      []Tag{{TAG_GETSET_NAME, []byte(fmt.Sprintf("/tuner%d/status", tuner))}},
	}
	if lockkey != 0 {
		pkt.Tags = append(pkt.Tags, Tag{TAG_GETSET_LOCKKEY, WriteUint32(lockkey)})
	}
	pkt, err := dev.SendReceive(pkt)
	if err != nil {
		return nil, err
	}
	// TODO check return FrameType
	tunerStatus := &TunerStatus{}
	fmt.Sscanf(string(pkt.Find(TAG_GETSET_NAME)), "/tuner%d", &tunerStatus.Tuner)
	value := pkt.Find(TAG_GETSET_VALUE)
	if match := regexp.MustCompile(`ch=([^ ]+)`).FindSubmatch(value); len(match) == 2 {
		tunerStatus.Channel = string(match[1])
	}
	if match := regexp.MustCompile(`lock=([^ ]+)`).FindSubmatch(value); len(match) == 2 {
		tunerStatus.Lock = string(match[1])
	}
	if match := regexp.MustCompile(`ss=(\d+)`).FindSubmatch(value); len(match) == 2 {
		fmt.Sscanf(string(match[1]), "%d", &tunerStatus.SignalStrength)
	}
	if match := regexp.MustCompile(`snq=(\d+)`).FindSubmatch(value); len(match) == 2 {
		fmt.Sscanf(string(match[1]), "%d", &tunerStatus.SignalToNoiseQuality)
	}
	if match := regexp.MustCompile(`seq=(\d+)`).FindSubmatch(value); len(match) == 2 {
		fmt.Sscanf(string(match[1]), "%d", &tunerStatus.SymbolErrorQuality)
	}
	if match := regexp.MustCompile(`bps=(\d+)`).FindSubmatch(value); len(match) == 2 {
		fmt.Sscanf(string(match[1]), "%d", &tunerStatus.RawBitsPerSecond)
	}
	if match := regexp.MustCompile(`pps=(\d+)`).FindSubmatch(value); len(match) == 2 {
		fmt.Sscanf(string(match[1]), "%d", &tunerStatus.PacketsPerSecond)
	}
	return tunerStatus, nil
}

type StreamInfo struct {
}

func (dev *Device) GetStreamInfo(tuner int) (*StreamInfo, error) {
	pkt := &Packet{
		FrameType: TYPE_GETSET_REQ,
		Tags:      []Tag{{TAG_GETSET_NAME, []byte(fmt.Sprintf("/tuner%d/streaminfo", tuner))}},
	}
	pkt, err := dev.SendReceive(pkt)
	if err != nil {
		return nil, err
	}
	// TODO check return FrameType
	// fmt.Sscanf(string(pkt.Find(TAG_GETSET_NAME)), "/tuner%d", &tunerStatus.Tuner)
	value := pkt.Find(TAG_GETSET_VALUE)
	fmt.Printf("value %v\n", string(value))
	return nil, nil
}

// extern LIBHDHOMERUN_API int hdhomerun_device_get_tuner_vstatus(struct hdhomerun_device_t *hd, char **pvstatus_str, struct hdhomerun_tuner_vstatus_t *vstatus);
// extern LIBHDHOMERUN_API int hdhomerun_device_get_tuner_channel(struct hdhomerun_device_t *hd, char **pchannel);
// extern LIBHDHOMERUN_API int hdhomerun_device_get_tuner_vchannel(struct hdhomerun_device_t *hd, char **pvchannel);
// extern LIBHDHOMERUN_API int hdhomerun_device_get_tuner_channelmap(struct hdhomerun_device_t *hd, char **pchannelmap);
// extern LIBHDHOMERUN_API int hdhomerun_device_get_tuner_filter(struct hdhomerun_device_t *hd, char **pfilter);
// extern LIBHDHOMERUN_API int hdhomerun_device_get_tuner_program(struct hdhomerun_device_t *hd, char **pprogram);
// extern LIBHDHOMERUN_API int hdhomerun_device_get_tuner_target(struct hdhomerun_device_t *hd, char **ptarget);
// extern LIBHDHOMERUN_API int hdhomerun_device_get_tuner_plotsample(struct hdhomerun_device_t *hd, struct hdhomerun_plotsample_t **psamples, size_t *pcount);
// extern LIBHDHOMERUN_API int hdhomerun_device_get_tuner_lockkey_owner(struct hdhomerun_device_t *hd, char **powner);
// extern LIBHDHOMERUN_API int hdhomerun_device_get_oob_status(struct hdhomerun_device_t *hd, char **pstatus_str, struct hdhomerun_tuner_status_t *status);
// extern LIBHDHOMERUN_API int hdhomerun_device_get_oob_plotsample(struct hdhomerun_device_t *hd, struct hdhomerun_plotsample_t **psamples, size_t *pcount);
// extern LIBHDHOMERUN_API int hdhomerun_device_get_ir_target(struct hdhomerun_device_t *hd, char **ptarget);
// extern LIBHDHOMERUN_API int hdhomerun_device_get_version(struct hdhomerun_device_t *hd, char **pversion_str, uint32_t *pversion_num);
// extern LIBHDHOMERUN_API int hdhomerun_device_get_supported(struct hdhomerun_device_t *hd, char *prefix, char **pstr);

// extern LIBHDHOMERUN_API uint32_t hdhomerun_device_get_tuner_status_ss_color(struct hdhomerun_tuner_status_t *status);
// extern LIBHDHOMERUN_API uint32_t hdhomerun_device_get_tuner_status_snq_color(struct hdhomerun_tuner_status_t *status);
// extern LIBHDHOMERUN_API uint32_t hdhomerun_device_get_tuner_status_seq_color(struct hdhomerun_tuner_status_t *status);

// extern LIBHDHOMERUN_API const char *hdhomerun_device_get_hw_model_str(struct hdhomerun_device_t *hd);
// extern LIBHDHOMERUN_API const char *hdhomerun_device_get_model_str(struct hdhomerun_device_t *hd);
