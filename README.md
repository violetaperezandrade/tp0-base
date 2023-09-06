### Ejercicio N°7:
Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo, no podrá responder consultas por la lista de ganadores.
Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

Al protocolo del ejercicio anterior se le agregan mas mensajes.

- Header de 2 bytes que indica la longitud, en bytes, del batch(sin contar el header)
- 1 Byte indicando la cantidad de apuestas en el batch

    - Header de 2 bytes que indican la longitud del payload
    - 1 byte indicando el codigo de operacion

    Código de operacion 
        1 --> enviar una apuesta:
            - 1 byte agencia de la apuesta 1
            - 4 bytes dni de la apuesta 1
            - 2 bytes numero de la apuesta 1
            - 2 bytes año de la apuesta 1
            - 1 byte mes de la apuesta 1
            - 1 byte dia de la apuesta 1
            - Los siguientes bytes, hasta encontrar un cero son el nombre de la apuesta 1
            - Los bytes restantes, apellido de la apuesta 1

            - El siguiente byte es el codigo de operacion de la apuesta 2
            ... etc
        
        2 --> envio de apuestas finalizado
        3 --> consultar ganadores
            1 byte: agencia

Respuestas del servidor

- 2 bytes indicando la longitud del payload(sin contar estos dos bytes):

     - 1 Byte indicando codigo del mensaje:
        - 0 --> Aun no hay ganadores
        - 2 --> ACK
        - 3 --> ganadores
            - primeros 4 bytes DNI ganador 1
            - siguientes 4 bytes DNI ganador 2
            ...

Para ver el funcionamiento basta con ejecutar 

```console
$ make docker-compose-up
```
