package utils

func LongTo8ByteArray(long int) []byte {
	// we want to represent the input as a 8-bytes array
	byteArray := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < len(byteArray); i++ {
		byt := long & 0xff
		byteArray[i] = byte(byt)
		long = (long - byt) / 256
	}
	return byteArray
}

func ShortTo2ByteArray(long int) []byte {
	byteArray := []byte{0, 0}
	for i := 0; i < len(byteArray); i++ {
		byt := long & 0xff
		byteArray[i] = byte(byt)
		long = (long - byt) / 256
	}
	return byteArray
}

func LongTo32ByteArray(long int) []byte {
	byteArray := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	for i := 0; i < len(byteArray); i++ {
		byt := long & 0xff
		byteArray[i] = byte(byt)
		long = (long - byt) / 256
	}
	return byteArray
}

func ByteArrayToLong(b []byte) int {
	value := 0
	for i := len(b) - 1; i >= 0; i-- {
		value = value*256 + int(b[i])
	}
	return value
}
