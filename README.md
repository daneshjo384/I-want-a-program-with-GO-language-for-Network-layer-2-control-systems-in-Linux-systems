# I-want-a-program-with-GO-language-for-Network-layer-2-control-systems-in-Linux-systems
I want a program with GO language for Network layer 2 control systems in Linux systems

I want a program with GO language for Network layer 2 control systems in Linux systems

This program can include the following:

Layer 2 frame processing (Ethernet) Layer 2 traffic management (VLAN, MAC, Switching) Control of network devices through raw socket Implementation of a simple switch or layer 2 controller Ability to control traffic, filtering and analysis of layer 2 protocols Ability to connect to eBPF or DPDK in the future ✅ Application features:
Using gopacket for packet analysis, using afpacket or pcap to read layer 2 traffic, managing MAC table, filtering and controlling traffic, simple graphic display with termui or saving in a file, the ability to create and control VLANs, 

installing dependencies:
go mod init l2-controller
go get github.com/google/gopacket
go get github.com/google/gopacket/layers
go get github.com/gizak/termui/v3


RUN
go run main.go
sudo go run main.go  //Requires sudo access:

output sample:
L2 Frame: 12:34:56:78:90:ab -> cd:ef:gh:ij:kl:mn on eth0
L2 Frame: ab:cd:ef:gh:ij:kl -> 12:34:56:78:90:ab on eth1

and a terminal window showing the Layer 2 status.

Expandable features:
VLAN support Implement a simple switch (Switch) Traffic control based on layer 2 protocol (ARP, STP, ...) Connect to eBPF or netlink for advanced management DPDK support in high application mode Save traffic in file or database RESTful interface for remote control


