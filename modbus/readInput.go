package modbus

import (
	"encoding/hex"
)

const (
	FUNCCODE_READ_COILS  = 1
	FUNCCODE_READ_INPUTS = 2
)

func Uint16GetLoByte(aIn uint16) uint8 {
	return uint8(aIn & 0xFF)
}

func Uint16GetHiByte(aIn uint16) uint8 {
	return uint8((aIn >> 8) & 0xFF)
}

func InsertLRC(aIn []byte) {
	if len(aIn) < 2 {
		panic("Expects atleast 2 bytes!")
	}

	aIn[len(aIn)-1] = 0

	lLRC := uint8(0)
	for _, cData := range aIn {
		lLRC += cData
	}

	aIn[len(aIn)-1] = -lLRC
}

type (
	ReadInput struct {
		slave    uint8
		function uint8
		addr     uint16
		quant    uint16
	}
)

func InitReadInput(aSlave, aFunction uint8, aAddr, aQuant uint16) ReadInput {
	return ReadInput{
		slave:    aSlave,
		function: aFunction,
		addr:     aAddr,
		quant:    aQuant,
	}
}

func (meRI ReadInput) MakeCore() (rMod []byte) {
	rMod = make([]byte, 14, 14)

	lBytes := [7]byte{}

	lBytes[0] = meRI.slave
	lBytes[1] = meRI.function
	lBytes[2] = Uint16GetHiByte(meRI.addr)
	lBytes[3] = Uint16GetLoByte(meRI.addr)
	lBytes[4] = Uint16GetHiByte(meRI.quant)
	lBytes[5] = Uint16GetLoByte(meRI.quant)
	InsertLRC(lBytes[:])

	hex.Encode(rMod, lBytes[:])

	return
}

func (meRI ReadInput) MakeWhole() (rMod []byte) {
	rMod = make([]byte, 17, 17)

	rMod[0] = ':'
	rMod[len(rMod)-2] = '\r'
	rMod[len(rMod)-1] = '\n'

	copy(rMod[1:len(rMod)-2], meRI.MakeCore())

	return
}
