package protocol

import (
	"encoding/binary"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/bet"
)

const OP_CODE = 1

func EncodeBet(bet bet.Bet) []byte {
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

func EncodeBets(bets []bet.Bet) []byte {
	encodedBets := make([][]byte, 0)
	for i := 0; i < len(bets); i++ {
		encodedBet := EncodeBet(bets[i])
		encodedBets = append(encodedBets, encodedBet)
	}

	totalLen := 1 // 1 byte para la cantidad de apuestas
	for _, encodedBet := range encodedBets {
		totalLen += len(encodedBet)
	} //los bytes de cada apuesta

	batch := make([]byte, totalLen+2) //2 bytes de header

	// la long del payload
	binary.BigEndian.PutUint16(batch[:2], uint16(totalLen))

	// 1 byte indicando cantidad de apuestas
	batch[2] = byte(len(encodedBets))

	offset := 3 //
	for _, encodedBet := range encodedBets {
		copy(batch[offset:], encodedBet)
		offset += len(encodedBet)
	}

	return batch

}
