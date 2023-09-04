### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento. La cantidad de apuestas dentro de cada _batch_ debe ser configurable.
El servidor, por otro lado, deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

Al protocolo del ejercicio anterior se le agrega un header mas. Para enviar un batch:

- Header de 2 bytes que indica la longitud, en bytes, del batch(sin contar el header)
- 1 Byte indicando la cantidad de apuestas en el batch

    - Header de 2 bytes que indican la longitud del payload
    - 1 byte indicando el codigo de operacion

    Código de operacion 1 --> enviar una apuesta:
        - 1 byte agencia de la apuesta 1
        - 4 bytes dni de la apuesta 1
        - 2 bytes numero de la apuesta 1
        - 2 bytes año de la apuesta 1
        - 1 byte mes de la apuesta 1
        - 1 byte dia de la apuesta 1
        - Los siguientes bytes, hasta encontrar un cero son el nombre de la apuesta 1
        - Los bytes restantes, apellido de la apuesta 1

    - Los siguientes dos bytes son el header de la apuesta 2
    ... etc

Respuesta del servidor

1 byte:
    - 1, un ACK indicando que la apuesta fue recibida de forma correcta

Para ver el funcionamiento basta con ejecutar 

```console
$ make docker-compose-up
```
