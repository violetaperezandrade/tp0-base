import socket
import logging
import signal

from protocol.protocol import decode, encode, encode_winners_not_ready, encode_winners
from .utils import store_bets, load_bets, has_won

ACK = 1
AGENCIES = 5


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._running = True
        self._operations_map = {
            1: self.__store_bets_received,
            2: self.__receive_bets_end,
            3: self.__receive_winners_query
        }
        self._agencies_finished = 0
        self.agencies_finished_anounced = []

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        # TODO: Modify this program to handle signal to graceful shutdown
        # the server
        signal.signal(signal.SIGTERM, self.handle_sigterm)
        while self._running:
            try:
                client_sock = self.__accept_new_connection()
            except OSError:
                if not self._running:
                    logging.info(f'action: sigterm received')
                else:
                    raise
                return
            self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # TODO: Modify the receive to avoid short-reads
            header = self.__read_exact(2, client_sock)
            payload = self.__read_exact(int.from_bytes(
                header, byteorder='big'), client_sock)

            op_code, data = decode(payload)
            answer = self._operations_map.get(op_code, lambda _: None)(data)
            if answer == None:
                logging.error(
                    "action: receive_message | result: fail | error: received unkown operation code")

            self.__send_exact(answer, client_sock)
        except OSError as e:
            logging.error(
                "action: receive_message | result: fail | error: {e}")
        finally:
            client_sock.close()

    def __read_exact(self, bytes_to_read, client_sock):
        bytes_read = client_sock.recv(bytes_to_read)

        while len(bytes_read) != bytes_to_read:
            new_bytes_read = client_sock.recv(bytes_to_read - len(bytes_read))
            bytes_read += new_bytes_read

        return bytes_read

    def __send_exact(self, answer, client_sock):
        bytes_sent = 0
        while bytes_sent < len(answer):
            chunk_size = client_sock.send(answer[bytes_sent:])
            bytes_sent += chunk_size

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        logging.info(
            f'action: accept_connections | result: success | ip: {addr[0]}')
        return c

    def handle_sigterm(self, signum, frame):
        logging.info(
            f'action: sigterm received | signum: {signum}, frame:{frame}')
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()
        self._running = False
        logging.info(f'action: close_server | result: success')
        return

    def __store_bets_received(self, bets):
        store_bets(bets)
        print(
            f'sent BET: length: {len(bets)}')
        logging.info(
            f'action: apuesta_enviada | result: success | agency: {bets[0].agency} ')
        return encode(ACK)

    def __receive_bets_end(self, op_code):
        self._agencies_finished += 1
        logging.info(
            f'action: received_agency_finished | result: success | agencies finished: {self._agencies_finished}, remaining: {AGENCIES - self._agencies_finished}')
        return encode(ACK)

    def __receive_winners_query(self, agency_id):
        # print(f'RECEIVE WINNERS QUERY BEING CALLED FROM AGENCY: {agency_id}')
        # print(f'AAAAND, AGENCIES FINISHED ARE: {self._agencies_finished}')
        if self._agencies_finished != AGENCIES:
            return encode_winners_not_ready()
        else:
            return self.__get_winners(agency_id)

    def __get_winners(self, agency_id):
        #print(f'GGET WINNERS BEING CALLED FROM AGENCY: {agency_id}')
        bets = load_bets()
        winners = []
        for bet in bets:
            if bet.agency == agency_id and has_won(bet):
                logging.debug(
                    f"From agency: {agency_id} dni: {bet.document}, winner")
                winners.append(bet.document)
        return encode_winners(winners)
