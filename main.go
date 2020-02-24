package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
)

type returnThingey struct {
	returnData string
	lastEntry  bool
}

func main() {

	argument := os.Args
	if len(argument) != 3 {
		fmt.Println("Usage: lookupblaster hostmask threads")
		return
	}

	maskPattern := regexp.MustCompile(`^([0-9.]+)/(\d+)$`)
	match := maskPattern.FindStringSubmatch(argument[1])
	if len(match) != 3 {
		fmt.Println("Invalid hostmask!\nUsage: lookupblaster hostmask threads")
		return
	}

	ipAddress := match[1]
	mask, _ := strconv.Atoi(match[2])
	goRoutines, _ := strconv.Atoi(argument[2])

	incomingMessageChannel := make(chan returnThingey)
	for t := 0; t < goRoutines; t++ {
		go doLookup(net.ParseIP(ipAddress), mask, t, goRoutines, incomingMessageChannel)
	}

	connectionsTerminated := 0
	for connectionsTerminated < goRoutines {
		fromConnection := <-incomingMessageChannel
		if fromConnection.lastEntry == false {
			fmt.Println(fromConnection.returnData)
		} else {
			connectionsTerminated++
		}
	}

}

func ip2int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

func doLookup(ip net.IP, mask int, instance int, instances int, incomingMessageChannel chan returnThingey) {
	ipv4Mask := net.CIDRMask(mask, 32)
	startNumeric := ip2int(ip.Mask(ipv4Mask))
	stopNumeric := startNumeric + 1<<uint32(32-mask)
	for startNumeric < stopNumeric {
		currentAddress := int2ip(startNumeric + uint32(instance))
		hosts, err := net.LookupAddr(currentAddress.String())
		if err == nil {
			for _, hostname := range hosts {
				hostname = hostname[:len(hostname)-1] // Remove trailing dot character
				incomingMessageChannel <- returnThingey{fmt.Sprint(currentAddress, "\t", hostname), false}
			}
		}
		startNumeric += uint32(instances)
	}
	incomingMessageChannel <- returnThingey{"", true}
}
