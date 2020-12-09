#!/bin/sh

# Previamente se han convertido algunos ficheros java antiguo de :
# ISO-8859-15 a UTF-8 con herramienta iconv de la siguiente forma :
# iconv -f ISO-8859-15 -t UTF-8 -o f.new f

# Ademas se han compilado algunos ficheros fuente en diferentes directorios
# con cuidado ya que, en la compilacion, habia que incluir
# un -classpath adecuado (mirando en imports) estilo :
#	javac -classpath ../.. *java

# Parametro $1: fichero con la red de petri textual

# Version no exportable
# java -classpath .:gson-2.8.6.jar Exportardp $1

java -jar ExportaLefs.jar $1
