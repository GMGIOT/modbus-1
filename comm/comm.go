package comm

import (
	"github.com/tarm/serial"
	//"log"
)

type Err string

func (meErr Err) Error() string {
	return string(meErr)
}

const (
	ErrAttached    = Err("Comm has another serial port attached!")
	ErrNotAttached = Err("Comm has no serial port attached!")
	ErrError       = Err("Unknown error. Please file a report!")
	ErrRxTimeout   = Err("Reception timed out!")
	ErrRxNoBody    = Err("Empty body!")
	ErrRxNoCR      = Err("No carriage return in trailer!")
)

type (
	ModbusClient struct {
		config serial.Config
		port   *serial.Port
	}
)

func MakeClient(aConfig serial.Config) *ModbusClient {
	return &ModbusClient{
		config: aConfig,
	}
}

func (meCom *ModbusClient) Open() error {
	if meCom.port != nil {
		return ErrAttached
	}

	var lErr error
	meCom.port, lErr = serial.OpenPort(&meCom.config)

	return lErr
}

func (meCom *ModbusClient) Close() error {
	if meCom.port == nil {
		return ErrNotAttached
	}

	lRet := meCom.port.Close()
	meCom.port = nil

	return lRet
}

func (meCom *ModbusClient) Escape() error {
	if meCom.port == nil {
		return ErrNotAttached
	}

	//If escape, send some escape requests
	for cIdx := 0; cIdx < 3; cIdx++ {
		if _, bErr := meCom.port.Write([]byte{27}); bErr != nil {
			return bErr
		}
	}

	//Flush data
	if bErr := meCom.port.Flush(); bErr != nil {
		return bErr
	}

	//Read and discard all existing data
	for {
		bDummy := [1]byte{}
		bNum, _ := meCom.port.Read(bDummy[:])

		if bNum == 0 {
			break
		}
	}

	return nil
}

func (meCom *ModbusClient) Request(aReq []byte, aEscape bool) ([]byte, error) {
	if meCom.port == nil {
		return nil, ErrNotAttached
	}

	//If escape, send some escape requests
	if aEscape {
		for cIdx := 0; cIdx < 3; cIdx++ {
			if _, bErr := meCom.port.Write([]byte{27}); bErr != nil {
				return nil, bErr
			}
		}
	}

	//Flush data
	if bErr := meCom.port.Flush(); bErr != nil {
		return nil, bErr
	}

	//Read and discard all existing data
	for {
		bDummy := [1]byte{}
		bNum, _ := meCom.port.Read(bDummy[:])

		if bNum == 0 {
			break
		}
	}

	if _, bErr := meCom.port.Write(aReq); bErr != nil {
		return nil, bErr
	}

	return readLine(meCom.port)
}

// readLine reads a line from the given serial port
func readLine(aPort *serial.Port) ([]byte, error) {
	bRet := make([]byte, 0)

	for {
		bDummy := [1]byte{}
		bNum, bErr := aPort.Read(bDummy[:])

		if bErr != nil {
			return bRet, bErr
		}

		if bNum == 0 {
			return bRet, ErrRxTimeout
		}

		if bDummy[0] == 10 {
			if len(bDummy) == 0 {
				bRet = append(bRet, bDummy[0])
				return bRet, ErrRxNoBody
			}

			bPrev := bRet[len(bRet)-1]
			bRet = append(bRet, bDummy[0])

			if bPrev == 13 {
				break
			} else {
				return bRet, ErrRxNoCR
			}
		}

		bRet = append(bRet, bDummy[0])

	}

	return bRet, nil
}
