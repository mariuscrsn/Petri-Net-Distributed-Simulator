# PetriSim
A distributed and conservative Petri net simulator

## TODO
 - [x] Entender todos los campos del json, listas tb.
 - [x] Como se envían transiciones... Ficheros lefs, transition & simulation_engine
 - [ ] Que envío de un nodo a otro?? Como se saben los idglobales, como añadir info de red
 al fichero json o donde sea.
 - [x] Probar que funciona en centralizado. Ejemplo de 
 3rx2t y 2rx2t
 - [x] Entender funcionamiento de las lefs
 - [ ] Diseñar el sistema, como voy a pasar las IPs y 
 puertos a cada subred para identificar a que máquinas 
 envían los lugares de salida
 
### Tests
- [ ] Definir las subredes a particionar y ampliar con 
ejemplos que incrementen el nº de rams a ejecutar en paralelo; n<= 5
- [ ] Cambiar los tiempos de transiciones para observar comportamientos
 temporales
- [ ] Incrementar la longitud de las trans de cada rama, hasta 4txr y hacer
una prueba adicional con tiempos dispares en las transiciones para observar el 
comportamiento en distribuido
 
 
### Dudas código
- [ ] Sirve para algo la variable ii_Tiempo de cada LEF. No se podría hacer sin ella? 