package main

import (
	"flag"
	// "fmt"
	// "log"
	"fmt"
	"time"

	"github.com/simonvetter/modbus"
)

func main() {

	var cycle uint16 = 0

	var singleRegAddr uint16 = 4
	var singleRegValue uint16 = 0

	var twoRegAddr1 uint16 = 0
	var twoRegAddr2 uint16 = 2
	var twoRegValue uint32 = 0

	// get the device serial port from the command line
	var (
		// Serial (Hardware)
		serialPort string
		baudRate   uint
		// Serial (Protocol)
		//     dataBits int
		parity   uint
		stopBits uint
		// Modbus (Device)
		slaveDevice   int
		startAddr     int
		responsePause int
	)

	const (
		// Serial (Hardware)
		defaultPort     = ""
		defaultBaudRate = 9600

		// Serial (Protocol)
		defaultDataBits = 8
		defaultParity   = modbus.PARITY_NONE
		defaultStopBits = 1

		// Modbus (Device)
		defaultSlave         = 1
		defaultStartAddress  = 0
		defaultResponsePause = 5
	)

	// Serial (Hardware)
	flag.StringVar(&serialPort, "serial", defaultPort, "Serial port (RS485) to use, e.g., /dev/ttyS0 (try \"dmesg | grep tty\" to find)")
	flag.UintVar(&baudRate, "baud", defaultBaudRate, fmt.Sprintf("Baud Rate (default is %d)", defaultBaudRate))

	// Serial (Protocol)
	// flag.IntVar(&dataBits, "dataBits", defaultDataBits, fmt.Sprintf("Data bits (default is %d)", defaultDataBits))
	flag.UintVar(&parity, "parity", defaultParity, fmt.Sprintf("Set the parity value \"PARITY NONE\"=%d, \"PARITY EVEN\"=%d or \"PARITY_ODD\"=%d. (default is %d)", modbus.PARITY_NONE, modbus.PARITY_EVEN, modbus.PARITY_ODD, defaultParity))
	flag.UintVar(&stopBits, "stopBits", defaultStopBits, fmt.Sprintf("Set the stop bits value 1 or 2. (default is %d)", defaultParity))

	// Modbus (Device)
	flag.IntVar(&slaveDevice, "slave", defaultSlave, fmt.Sprintf("Slave device number (default is %d)", defaultSlave))
	flag.IntVar(&startAddr, "start", defaultStartAddress, fmt.Sprintf("Start address (default is %d)", defaultStartAddress))
	flag.IntVar(&responsePause, "pause", defaultResponsePause, fmt.Sprintf("Pause between write cycle in milliseconds (default is %d)", defaultResponsePause))

	// End of the Args setup
	flag.Parse()

	if len(serialPort) == 0 { // If the serial port isn't set
		// Show the commands and close the program and close the program
		flag.PrintDefaults()
		return
	}

	var client *modbus.ModbusClient
	var err error

	// for an RTU (serial) device/bus
	client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:      "rtu://" + serialPort, // need a value
		Speed:    baudRate,              // default
		DataBits: defaultDataBits,       // default, optional
		Parity:   parity,                // default, optional
		StopBits: stopBits,              // default if no parity, optional
		Timeout:  time.Millisecond * 300,
	})
	if err != nil {
		// error out if client creation failed
	}

	// now that the client is created and configured, attempt to connect
	err = client.Open()
	if err != nil {
		// error out if we failed to connect/open the device
		// note: multiple Open() attempts can be made on the same client until
		// the connection succeeds (i.e. err == nil), calling the constructor again
		// is unnecessary.
		// likewise, a client can be opened and closed as many times as needed.
	}

	// Switch to unit ID (a.k.a. slave ID) #1
	client.SetUnitId(1)

	// Writing Cycle
	for true {
		// write -200 to 16-bit (holding) register 100, as a signed integer
		err = client.WriteRegister(singleRegAddr, singleRegValue)
		if err != nil {
			fmt.Printf("failed to read register 0x4000: %v\n", err)
		} else {
			fmt.Printf("register 0x4000: 0x%04x\n", singleRegValue)
		}

		// write a int32 value between 2 registers
		err = client.WriteUint32(twoRegAddr1, twoRegValue)
		if err != nil {
			fmt.Printf("failed to read registers 0x4000 and 0x4001: %v\n", err)
		} else {
			fmt.Printf("register 0x4000 and 0x4001: 0x%08x\n", twoRegValue)
		}

		// write a int32 value between 2 registers
		err = client.WriteUint32(twoRegAddr2, twoRegValue)
		if err != nil {
			fmt.Printf("failed to read registers 0x4000 and 0x4001: %v\n", err)
		} else {
			fmt.Printf("register 0x4002 and 0x4003: 0x%08x\n", twoRegValue)
		}

		time.Sleep(time.Duration(responsePause) * time.Second)

		cycle++

		if cycle > 60 { // whait for a moment and increment the registers
			cycle = 0
			singleRegValue++
			twoRegValue++

			if singleRegValue > 1000 || twoRegValue > 1000 {
				singleRegValue = 0
				twoRegValue = 0
			}
		}

	}

	// close the TCP connection/serial port
	client.Close()
}
