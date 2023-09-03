package protocol

import (
	"encoding/binary"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/bet"
)

const OP_CODE = 1

func EncodeBet(bet *bet.Bet) []byte {
	//longitud first y last name, unicas variables
	firstNameLen := len(bet.FirstName)
	lastNameLen := len(bet.LastName)

	//longitud total(bytes fijos + long first y last name + byte del fin del nombre)
	totalLen := 14 + firstNameLen + lastNameLen + 1
	betBytes := make([]byte, totalLen)

	//header --> 2 bytes del len del payload
	binary.BigEndian.PutUint16(betBytes[:2], uint16(totalLen-2))

	//1 byte codigo de operacion
	betBytes[2] = byte(OP_CODE)

	//1 byte de la agencia
	betBytes[3] = byte(bet.Agency)

	//4 bytes DNI
	binary.BigEndian.PutUint32(betBytes[4:8], uint32(bet.Dni))

	//2 bytes numero de apuesta
	binary.BigEndian.PutUint16(betBytes[8:10], uint16(bet.Number))

	//2 bytes del año
	binary.BigEndian.PutUint16(betBytes[10:12], uint16(bet.Year))

	//1 byte del mes
	betBytes[12] = byte(bet.Month)

	//1 byte del dia
	betBytes[13] = byte(bet.Day)

	//nombre con un 0 al final
	copy(betBytes[14:], []byte(bet.FirstName))
	betBytes[14+firstNameLen] = 0x00

	//last name
	copy(betBytes[14+firstNameLen+1:], []byte(bet.LastName))

	return betBytes

}
