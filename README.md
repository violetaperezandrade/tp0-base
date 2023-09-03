# TP0: Docker + Comunicaciones + Concurrencia
### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

Para ejecutar se debe correr

```console
$ make docker-compose-up
$ make docker-compose-logs
```
Luego, en otra terminal correr
```console
$ make docker-compose-down
```

Y en la primera terminal se podra ver como tanto el cliente como el servidor terminan de forma _graceful_
