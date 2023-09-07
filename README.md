## Parte 3: Repaso de Concurrencia

### Ejercicio N°8:
Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo.
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

Para este ejercicio, se creó un grupo de hilos (conocido como un "pool de threads") que puede tener hasta cinco hilos ejecutándose simultáneamente. Cuando el hilo principal recibe una nueva conexión, llama a la función `handle_new_thread`. Esta función verifica si en el grupo de hilos hay algún hilo disponible para ser utilizado. Si encuentra un hilo disponible, como en Python no es posible reiniciar un hilo que ya ha sido utilizado, se crea un hilo completamente nuevo en su lugar pisando el anterior.

En el caso de que todos los hilos estén ocupados manejando otras conexiones, el hilo principal espera hasta que al menos uno de los hilos se libere para poder ser utilizado nuevamente.

Para ver el funcionamiento basta con ejecutar 

```console
$ make docker-compose-up
```
