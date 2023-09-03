from common.utils import Bet

RECEIVE_BET_CODE = 1


def decode(payload):
    op_code = int(payload[0])
    return STRATEGY_MAP[op_code](payload)


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


STRATEGY_MAP = {
    1: decode_bet_received
}
