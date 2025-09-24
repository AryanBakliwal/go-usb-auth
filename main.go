package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jochenvg/go-udev"
)

const (
	// Hardcoded interface class/subclass/protocol to block
	BlockIFClass    = 8 // Mass storage
	BlockIFSubClass = 6
	BlockIFProtocol = 50 // example: Bulk-Only Transport (80)
)

func main() {
	u := udev.Udev{}

	monitor := u.NewMonitorFromNetlink("udev")
	if monitor == nil {
		log.Fatal("failed to create monitor")
	}

	// Create a context
	ctx, _ := context.WithCancel(context.Background())

	// Open the device channel (this internally calls udev_monitor_enable_receiving)
	ch, _, err := monitor.DeviceChan(ctx)
	if err != nil {
		log.Fatalf("failed to create device channel: %v", err)
	}

	fmt.Println("Monitoring USB add/remove events...")

	for d := range ch {

		if d == nil || d.Subsystem() != "usb" || d.Action() != "add" {
			continue
		}

		// dev.Syspath() is like /sys/bus/usb/devices/2-4.1
		sysPath := d.Syspath()
		base := filepath.Base(sysPath) // e.g., "2-4.1"
		devicePath := filepath.Join("/sys/bus/usb/devices", base)

		// list entries under the device path - find interfaces like "2-4.1:1.0"
		entries, err := ioutil.ReadDir(devicePath)
		if err != nil {
			log.Printf("cannot read %s: %v", devicePath, err)
			continue
		}

		for _, e := range entries {
			name := e.Name()
			// interface entries contain a colon
			if !strings.Contains(name, ":") {
				continue
			}
			ifPath := filepath.Join(devicePath, name)
			ifClass := readHexFile(ifPath + "/bInterfaceClass")
			ifSub := readHexFile(ifPath + "/bInterfaceSubClass")
			ifProto := readHexFile(ifPath + "/bInterfaceProtocol")

			fmt.Printf("Found interface %s: class=%d sub=%d proto=%d\n",
				name, ifClass, ifSub, ifProto)

			if ifClass == BlockIFClass && ifSub == BlockIFSubClass && ifProto == BlockIFProtocol {
				authorizedPath := filepath.Join(ifPath, "authorized")
				fmt.Printf("Blocking interface %s -> writing 0 to %s\n", name, authorizedPath)
				if err := ioutil.WriteFile(authorizedPath, []byte("0"), 0644); err != nil {
					log.Printf("failed to write authorized for %s: %v", name, err)
				} else {
					fmt.Printf("Blocked interface %s\n", name)
				}
			}
		}
	}

}

func readHexFile(path string) int {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return -1
	}
	s := strings.TrimSpace(string(b))
	// try parse as hex (common) then fallback to decimal
	if strings.HasPrefix(s, "0x") || strings.ContainsAny(s, "abcdefABCDEF") {
		v, err := strconv.ParseInt(strings.TrimPrefix(s, "0x"), 16, 64)
		if err == nil {
			return int(v)
		}
	}
	// fallback: try parse base 10
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1
	}
	return int(v)
}
