package main

import (
	"github.com/tarm/serial"
	"log"
	"modbus/comm"
	"modbus/modbus"
	"time"
)

func main() {
	log.SetFlags(log.Lshortfile)

	/*
		lRI := modbus.InitReadInput(1, modbus.FUNCCODE_READ_INPUTS, 1240, 64)
		lBuf := lRI.MakeWhole()
	*/

	lBuf := modbus.CoverWithEnvelop([]byte("010204D80040E1"))

	log.Println("Sending '" + string(lBuf) + "' ...")

	lComConf := serial.Config{Name: "COM10", Baud: 57600, ReadTimeout: time.Second * 5}
	lCom := comm.MakeClient(lComConf)

	if bErr := lCom.Open(); bErr != nil {
		log.Println("Comm open error: " + bErr.Error())
		return
	}

	defer func() {
		if bErr := lCom.Close(); bErr != nil {
			log.Println("Comm close error: " + bErr.Error())
		}
	}()

	lRep, lErr := lCom.Request(lBuf, true)
	if lErr != nil {
		log.Println("Comm request error: " + lErr.Error())
		log.Println("Data received: " + string(lRep))
		return
	}

	log.Println(string(lRep))
}
