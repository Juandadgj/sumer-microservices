# sumer-microservices
Ejemplo de sencillo de implementación de colas con Golang y RabbitMQ  
## Sobre el ejemplo
Para este caso encontrarás una de las tantas soluciones a un problema común:
¿Qué pasa cuándo una request conlleva muchos procesos cómo para simplemente
manejarlo a través de HTTP?  
Para este ejemplo implementamos un sencillo enpoint que recibe una petición para
la confirmación de un pago. Ahora bien, una verificación de pago puede contener también
relacionado un proceso de notificación, este es nuestro caso. El problema radica en que
el servicio encargado de realizar las noticicaciones al email y telefono movil del cliente
es un poco tardío, no podemos hacer esperar al cliente tardo porque afectaría la experiencia.
# ¿Solución?
¿Por qué no entonces en vez me manejar consultas síncronas implementamos un método asíncrono?
Así se podría enviar una respuesta rápida al cliente de que el pago fue verificado mientras nuestro
servicio se encarga de notificarle a su información de contacto. Es para esto que implementamos las
colas con RabbitMQ. Te explico brevemente los dos servicios que ves en el repositorio.
## verify-payment
Al ejecutar este servicio podrás hacer una simple petición POST a la ruta localhost:3000/payment y enviarle
en forma de JSON los siguiente datos: {"ID": int, "Client": string y "Amount": float}. Este servicio a travez 
de sleep simula estar trabajando y nos entrega una respuesta en 700ms aproximadamente con la confirmación, una 
velocidad aceptable. Hecha la verificación envía un evento a RabbitMQ con el payment y nuestro servicio de 
notificaciones hace el resto.
## notifications
Este servicio por su parte va a estar atento a los eventos que ingresen a la cola he intentará resolverlos de inmediato. 
Una de las ventajas que tiene implementar este sistema es que si llegásemos a tener cuellos de botella podríamos simplemente 
ejecutar más workers, es decir, instancias de este mismo servicio y se repartiran la carga de forma sencilla gracias a RabbitMQ. 
Esto puedes verificarlo tú mismo viendo los logs en la terminal.
# Ejecucion
Es muy sencillo, simplemente ejecutaremos con go run o build los main.go que se encuentran en los directorios. Para el caso de 
verify-payment sólo tendremos que ejecutarlo una vez y conectarnos atravez de un cliente cómo Postman por ejemplo. Pero para 
notifications podemos ejecutar multiples intancias y ver cómo se comportan cuándo hacemos las peticiones.
