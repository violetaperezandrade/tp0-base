# TP0: Docker + Comunicaciones + Concurrencia
### Ejercicio N°1.1:
Definir un script (en el lenguaje deseado) que permita crear una definición de DockerCompose con una cantidad configurable de clientes.

### Requerimentos:
-Se requiere tener Python3 instalado para correr el script.
### Para correrlo:

```console
$ cd scripts
$ chmod +x generate_docker_compose
$ ./generate_docker_compose [-h] [-n CLIENTS] [-f FILEPATH]
```

Donde:
- `-h`: Mostrar la ayuda
- `-n`: La cantidad de clientes que se desean agregar
- `-f`: Ruta al archivo docker compose

El script sobreescribe el archivo original y asume que hay al menos un cliente
