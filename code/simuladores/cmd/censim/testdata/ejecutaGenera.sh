#!/bin/sh

# Previamente se han convertido algunos ficheros java antiguo de :
# ISO-8859-15 a UTF-8 con herramienta iconv de la siguiente forma :
# iconv -f ISO-8859-15 -t UTF-8 -o f.new f

# Ademas se han compilado algunos ficheros fuente en diferentes directorios
# con cuidado ya que, en la compilacion, habia que incluir
# un -classpath adecuado (mirando en imports) estilo :
#	javac -classpath ../.. *java

# Parametros
#	$1: Nombre fichero de salida de la RdP Textual
#	$2: numero de lugares horizontales
#	$3: numero de lugares verticales

java -jar genera.jar $1 $2 $3
