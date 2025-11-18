package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

// MAC table (MAC -> Interface)
type MacTable struct {
	entries map[string]string
}

func NewMacTable() *MacTable {
	return &MacTable{
		entries: make(map[string]string),
	}
}

func (mt *MacTable) Learn(mac string, iface string) {
	mt.entries[mac] = iface
}

func (mt *MacTable) Lookup(mac string) (string, bool) {
	iface, exists := mt.entries[mac]
	return iface, exists
}

// Layer 2 controller
type L2Controller struct {
	macTable *MacTable
	devices  []string
}

func NewL2Controller() *L2Controller {
	return &L2Controller{
		macTable: NewMacTable(),
		devices:  getNetworkInterfaces(),
	}
}

// Function to get all network devices
func getNetworkInterfaces() []string {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	var names []string
	for _, device := range devices {
		// Add only active devices
		if len(device.Addresses) > 0 {
			names = append(names, device.Name)
		}
	}
	return names
}

// Function to listen to a device
func (l2c *L2Controller) ListenOnDevice(deviceName string) {
	handle, err := pcap.OpenLive(deviceName, 1600, true, pcap.BlockForever)
	if err != nil {
		log.Printf("Error opening device %s: %v", deviceName, err)
		return
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range packetSource.Packets() {
		l2c.ProcessPacket(packet, deviceName)
	}
}

// Process the package
func (l2c *L2Controller) ProcessPacket(packet gopacket.Packet, iface string) {
	ethLayer := packet.Layer(layers.LayerTypeEthernet)
	if ethLayer != nil {
		eth, _ := ethLayer.(*layers.Ethernet)

		// Learn the source MAC
		srcMac := eth.SrcMAC.String()
		dstMac := eth.DstMAC.String()

		l2c.macTable.Learn(srcMac, iface)
		// Print information
		fmt.Printf("L2 Frame: %s -> %s on %s\n", srcMac, dstMac, iface)

		// You can decide here whether to forward the packet or not 
		// For example:
		// if dst, exists := l2c.macTable.Lookup(dstMac); exists { ... }
	}
}

// Display the status in the terminal
func (l2c *L2Controller) RunUI() {
	if err := termui.Init(); err != nil {
		log.Fatalf("Failed to initialize termui: %v", err)
	}
	defer termui.Close()

	p := widgets.NewParagraph("L2 Controller Status\n\nListening on all interfaces...")
	p.SetRect(0, 0, 60, 10)

	grid := termui.NewGrid()
	grid.SetRect(0, 0, 60, 10)
	grid.Set(
		termui.NewRow(1.0, p),
	)

	termui.Render(grid)

	go func() {
		for {
			// Update status
			status := fmt.Sprintf("L2 Controller Status\n\nActive Interfaces: %v\nMAC Entries: %d", l2c.devices, len(l2c.macTable.entries))
			p.Text = status
			termui.Render(grid)
			time.Sleep(1 * time.Second)
		}
	}()

	for e := range termui.PollEvents() {
		if e.Type == termui.KeyboardEvent {
			if e.ID == "q" || e.ID == "<C-c>" {
				break
			}
		}
	}
}

func main() {
	controller := NewL2Controller()

	// Listen to all network devices in a goroutine
	for _, dev := range controller.devices {
		if dev == "lo" {
			continue // Skip the loopback
		}
		go controller.ListenOnDevice(dev)
	}

	// Implement the UI
	controller.RunUI()
}