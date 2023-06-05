# TA 3
### Introducción
El informe de la Tarea Académica es sobre la recreación del juego Hoop hop Showdown creado en el lenguaje de progrmación Go. El juego trata de la división de una clase en grupos entre 3 a 4 personas y se utiliza hula hoops, conos, cubeta y fichas (siendo lo que se gana cuando terminas). Cada cubeta contiene entre 20 a 25 tokens.

Entra un alumno de cada equipo a la vez. Yo les digo a mis alumnos que una persona de su equipo va en el tablero a la vez. Los alumnos deben saltar en cada aro para moverse por el tablero. Su objetivo es llegar al cono de otro equipo para ganar una de sus fichas y llevarla al cubeta de su equipo. Mientras el alumno salta, puede enfrentarse a otro alumno. Estos dos estudiantes juegan a piedra, papel o tijera (RPS en inglés). El alumno que pierde sale del aro y vuelve corriendo a su equipo. En cuanto un alumno sale de un aro (del tablero de juego), la siguiente persona de su equipo puede empezar. Lo mismo ocurre si un estudiante llega al cono de otro equipo; ese estudiante sale del aro (lo que permite a su siguiente compañero de equipo comenzar) y entonces puede coger una ficha y correr de vuelta al cubo de su equipo.

### Funcionamiento del código
El código crea dos estructuras, una para el jugador y la otra para el juego (que sirve como una especie de controlador). Luego se crea una función del Nuevo Juego para que cree las cantidades necesarias de hula hoops, conos, fichas, número de equipos, número de jugadores  y las fichas que lo dejamos en 20 por cubeta.

Cuando empieza el juego en la función Run() se le da identificadores a todos los jugadores y equipos para una mejor lectura en los resultados, y se busca si hay aros disponibles para saltar, si lo hay solo salta hasta que se encuentre con otra persona y juegan a piedra, papel o tijera (RPS en inglés). 

![image](https://github.com/GustavoBecerra1109/Hoop-Hop-Showdown-/assets/54639476/663b2174-b156-4053-b0b1-088ea59d0c7b)

En la imagen muestra que si el jugador no gana, regresa a la primera posición y muestra el mensaje que perdió el piedra, papel o tijera. Luego hay un condicional que si ganan, o en este caso si es que llega al final, sale el segundo mensaje y agarra un token o ficha y regresa al primer hula hoop. 

Luego se crea un semáforo para aumentar la cantidad de tokens en cada equipo que haya llegado al final.

Se crea la función IsWinner() para determinar el ganador del piedra, papel y tijera usando un random para ver que número es el ganador de los jugadores. Al final en la función main se agrega el los equipos y jugadores y se empieza el juego.

### Conclusiones
Aunque el código funciona, creo que existe algunos problemas con los tiempos de ejecucción. Muchos de los problemas presentados en el código eran los deadlocks, así que el uso de un semáforo fue la solución óptima que hemos tenido.
