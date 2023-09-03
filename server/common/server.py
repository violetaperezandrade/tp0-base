import socket
import logging
import signal

from protocol.protocol import decode
from .utils import store_bets

ACK = 1


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self.running = True

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
        while self.running:
            try:
                client_sock = self.__accept_new_connection()
            except OSError:
                if not self.running:
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

            bet = decode(payload)
            store_bets([bet])
            addr = client_sock.getpeername()
            logging.info(
                f'action: apuesta_enviada | result: success | ip: {addr[0]} | dni: {bet.document} | numero: {bet.number}')
            # TODO: Modify the send to avoid short-writes
            client_sock.send(ACK.to_bytes(1, byteorder='big'))
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
        self.running = False
        logging.info(f'action: close_server | result: success')
        return
