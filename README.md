# TP0: Docker + Comunicaciones + Concurrencia
### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera un nuevo build de las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida afuera de la imagen (hint: `docker volumes`).

Para correr, se debe primero correr buildeando las imagenes

```console
$ make docker-compose-up
```

Luego, cambiar algun archivo de confiuracion, como ``server/config.ini`` o ``client/config.yaml`` y volver a correr pero sin buildear mediante

```console
$ docker compose -f docker-compose-dev.yaml up
```

y verificar que impacten los cambios hechos en el archivo de configuracion
