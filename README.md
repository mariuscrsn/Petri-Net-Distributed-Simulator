# PetriSim
A distributed and conservative Petri net simulator

## TODO
 - [ ] Probar que funciona en centralizado. Ejemplo de 
 3rx2t y 2rx2t
 - [ ] Entender funcionamiento de las lefs
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
 