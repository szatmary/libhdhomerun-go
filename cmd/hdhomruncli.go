package main

import (
	"fmt"
	"net"
	"os"

	"github.com/szatmary/libhdhomerun-go/hdhomerun"
)

// func MulticastAddrs() []net.Addr {
// 	var addrs []net.Addr
// 	ifaces, _ := net.Interfaces()
// 	for _, iface := range ifaces {
// 		if iface.Flags&net.FlagMulticast != 0 {
// 			fmt.Printf("%v\n", iface)
// 			iface.
// 		}
// 		// fmt.Printf("%v\n", iface)
// 		// addrs, _ := iface.MulticastAddrs()
// 		// for _, addr := range addrs {
// 		// 	// fmt.Printf("%v\n", addr)
// 		// 	addrs = append(addrs, addr)
// 		// }
// 	}
// 	return addrsxxwxw
// }

func main() {
	// for _, addr := range MulticastAddrs() {
	// 	fmt.Printf("%v\n", addr)
	// }
	// find tuners
	listen, _ := net.ListenUDP("udp", &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 0, Zone: ""})
	// udpreader := bufio.NewReader(listen)
	bcast, _ := net.ResolveUDPAddr("udp", "239.255.255.250:65001")

	pkt, _ := hdhomerun.Discover(hdhomerun.DEVICE_TYPE_TUNER, hdhomerun.DEVICE_ID_WILDCARD)
	listen.WriteTo(pkt, bcast)
	fmt.Printf("Pkt %v\n", pkt)
	for {
		fmt.Print("Waiting\n")
		var buf [1500]byte
		n, addr, err := listen.ReadFromUDP(buf[:])
		if err != nil {
			os.Exit(0)
		}
		fmt.Printf("Got bytes %v from %v\n", buf[:n], addr)
		var pkt hdhomerun.Packet
		err = pkt.UnmarshalBinary(buf[:n])
		fmt.Printf("Parsed %s (%v)\n", pkt.String(), err)
		switch pkt.FrameType {
		case hdhomerun.TYPE_DISCOVER_RPY:
			for _, x := range pkt.Tags {
				fmt.Printf("tag: %v\n", x)
			}

			fmt.Printf("Getting tuner status\n")
			dev, err := hdhomerun.NewDevice(addr.String())
			if err != nil {
				fmt.Printf("%v\n", err)
			}
			stat, err := dev.GetTunerStatus(1, 0)
			if err != nil {
				fmt.Printf("Err: %v\n", err)
			}
			fmt.Printf("Tuner status %v\n", stat)
			dev.GetStreamInfo(1)
		}
	}
}
