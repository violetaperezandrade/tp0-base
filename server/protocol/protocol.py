import struct
from common.utils import Bet


RECEIVE_BETS_CODE = 1
SEND_ACK_RECEIVED = 1
RECEIVE_BET_FINISHED = 2
RECEIVE_WINNERS_QUERY = 3


def decode(payload):
    if len(payload) > 5:  # batch
        op_code = int(payload[3])
    else:
        op_code = int(payload[0])
    return op_code, DECODE_MAP[op_code](payload)


def encode(op_code):
    return ENCODE_MAP[op_code]()


def decode_bets(payload):
    bets = []
    i = 1

    while i < len(payload):

        bet_len = int.from_bytes(payload[i:i+2], byteorder='big')
        bet_i = payload[i+2:i+bet_len+2]
        bets.append(decode_bet_received(bet_i))

        i += bet_len+2

    return bets


def decode_bet_received(payload):
    agency = int(payload[1])
    dni = int.from_bytes(payload[2:6], byteorder='big')
    bet_number = int.from_bytes(payload[6:8], byteorder='big')
    year = int.from_bytes(payload[8:10], byteorder='big')
    month = int(payload[10])
    day = int(payload[11])
    first_last_name = payload[12:]
    end_name_index = first_last_name.index(b'\x00')
    name = first_last_name[:end_name_index].decode('utf-8')
    last_name = first_last_name[end_name_index+1:].decode('utf-8')
    birthdate = f"{year}-{month:02}-{day:02}"
    bet = Bet(agency, name, last_name, str(dni), birthdate, str(bet_number))
    return bet


def decode_bets_end(payload):
    return int(payload[0])


def encode_ack():
    msg = bytearray(3)
    # struct.pack_into('>H', msg, 0, 1)
    msg[0] = 0x00
    msg[1] = 0x01
    msg[2] = 0x02
    return bytes(msg)


def decode_winners_query(payload):
    return int(payload[1])


def encode_winners_not_ready():
    msg = bytearray(3)
    msg[0] = 0x00
    msg[1] = 0x01
    msg[2] = 0x00
    return msg


def encode_winners(winners):
    total_len = 1 + len(winners)*4  # 1 byte de opcode + 4 bytes de cada dni

    msg = bytearray(total_len + 2)  # los 2 bytes de header
    msg[0:2] = total_len.to_bytes(2, byteorder='big')  # header

    op_code = 3
    msg[2:3] = op_code.to_bytes(1, byteorder='big')

    i = 3
    for dni in winners:
        msg[i:i+4] = int(dni).to_bytes(4, byteorder='big')
        i += 4

    return msg


DECODE_MAP = {
    RECEIVE_BETS_CODE: decode_bets,
    RECEIVE_BET_FINISHED: decode_bets_end,
    RECEIVE_WINNERS_QUERY: decode_winners_query
}

ENCODE_MAP = {
    SEND_ACK_RECEIVED: encode_ack
}
