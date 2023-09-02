# TP0: Docker + Comunicaciones + Concurrencia
### Ejercicio N°3:
Crear un script que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un EchoServer, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado. Netcat no debe ser instalado en la máquina _host_ y no se puede exponer puertos del servidor para realizar la comunicación (hint: `docker network`).

Para ejecutar la prueba se debe correr

```console
$ make docker-compose-up
$ docker exec netcat /netcat_script.sh
```