package modbus

func CoverWithEnvelop(aIn []byte) []byte {
	aIn = append([]byte{58}, aIn...)
	aIn = append(aIn, 13)
	aIn = append(aIn, 10)

	return aIn
}
