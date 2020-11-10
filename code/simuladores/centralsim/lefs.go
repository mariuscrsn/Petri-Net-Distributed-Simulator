//Package centralsim with several files to offer a centralized simulation
// This file deals with the low level lefs encoding of a petri net
package centralsim

import (
	"encoding/json"
	"fmt"
	"os"
)

//type TypeIndexSubnet int32

//----------------------------------------------------------------------------

// Lefs es el tipo de datos principal que gestiona el disparo de transiciones.
type Lefs struct {
	// Slice de transiciones de esta subred
	IaRed TransitionList `json:"ia_red"`
	//ii_indice int32	// Contador de transiciones agnadidas, Necesario ???
	// Identificadores de las transiciones sensibilizadas para
	// T = Reloj local actual. Slice que funciona como Stack
	IsTransSensib TransitionStack
}

// Load obtains Lefs from a json file
func Load(filename string) (Lefs, error) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open json lefs file: %v\n", err)
		return Lefs{}, err
	}
	defer file.Close()

	result := Lefs{}
	if err := json.NewDecoder(file).Decode(&result); err != nil {
		fmt.Fprintf(os.Stderr, "Decode json 		file: %v\n", err)
		return Lefs{}, err
	}

	result.IsTransSensib = MakeTransitionStack(100) //aun siendo dinamicos...

	return result, nil
}

/*
	-----------------------------------------------------------------
	   METODO: MakeLefs
	   RECIBE: numero de transiciones en la red
	   DEVUELVE: Nada
	   PROPOSITO: crear los datos iniciales
	   HISTORIA DE CAMBIOS:
   	COMENTARIOS:

func MakeLefs (int ai_ntrans) Lefs {
		// creamos los arrays de la dimension que nos indican
		ia_red = make(ListaTransiciones, ai_ntrans)
		ii_indice = 0;
		isTransSensib = nil
		il_eventos = nil
	}


func NewLefs (listaTransiciones ListaTransiciones) Lefs {
	l := Lefs{}
	l.ia_red = listaTransiciones
	l.ii_indice = 0
	l.isTransSensib = nil
	l.il_eventos = nil
}


// TranslateTablesToLefs takes a structure of a Petri Net of a
// subnet ai_subred from the global net in Tables format
// and translates it to a Lefs structure (antigua metodo transforma
// en antigua clase transformar
func TranslateTablesToLefs(aT_red Tables, ai_subred TypeIndexSubnet) {
	for
}
*/

/*
-----------------------------------------------------------------
   METODO: agnade_transicion
   RECIBE: indice en la tabla global de transiciones,
	  valor de la funcion de sensibilizacion, tiempo para el que
	  ese valor es valido, y el coste de disparo de esa transicion
   DEVUELVE: ii_indice de la transicion insertada
   PROPOSITO: Crea la instancia de la clase Transicion y la inserta
	  en la lista de transiciones para esta subred
   HISTORIA DE CAMBIOS:
COMENTARIOS:
-----------------------------------------------------------------
*/
/*
func (me *Lefs) Agnade_transicion (int ai_id,int ai_valor,int ai_tiempo,int ai_duracion,int ctes[][]) int32 {
	me.ia_red[ii_indice]=Transition{}
	me.ia_red[ii_indice].ii_idglobal=ai_id;
	me.ia_red[ii_indice].iiValorLef=ai_valor;
	me.ia_red[ii_indice].ii_tiempo=ai_tiempo;
	me.ia_red[ii_indice].ii_duracion_disparo=ai_duracion;
	me.ia_red[ii_indice].ii_listactes=ctes;
	me.ii_indice++;
	return ii_indice - 1;
}
*/

/*
-----------------------------------------------------------------
   METODO: agnade_sensibilizada
   RECIBE: Transicion sensibilizada a a�adir
   DEVUELVE: OK si todo va bien o ERROR en caso contrario
   PROPOSITO: A�ade a la lista de transiciones sensibilizadas
   HISTORIA DE CAMBIOS:
COMENTARIOS:
-----------------------------------------------------------------
*/
func (l *Lefs) agnadeSensibilizada(aiTransicion IndLocalTrans) bool {
	l.IsTransSensib.push(aiTransicion)
	return true // OK
}

// haySensibilizadas permite saber si tenemos transiciones sensibilizadas;
// se supone que previamente se ha llamado a actualizaSensibilizadas(relojLocal)
func (l Lefs) haySensibilizadas() bool {
	return !l.IsTransSensib.isEmpty()
}

// getSensibilizada coge el primer identificador de la lista de transiciones
//	 		sensibilizadas
func (l *Lefs) getSensibilizada() IndLocalTrans {
	if (*l).IsTransSensib.isEmpty() {
		return -1
	}

	return (*l).IsTransSensib.pop()
}

// actualizaSensibilizadas recorre toda la lista de transiciones
//	   e inserta trans sensibilizadas, con el mismo tiempo que el reloj local,
//  en la pila de transiciones sensibilizadas
func (l *Lefs) actualizaSensibilizadas(aiRelojLocal TypeClock) bool {
	for IndT, t := range (*l).IaRed {
		if t.IiValorLef <= 0 && t.IiTiempo == aiRelojLocal {
			(*l).IsTransSensib.push(IndLocalTrans(IndT))
		}
	}
	return true
}

// ImprimeTransiciones para depurar errores
func (l Lefs) ImprimeTransiciones() {
	fmt.Println(" ")
	fmt.Println("------IMPRIMIMOS LA LISTA DE TRANSICIONES---------")
	for _, tr := range l.IaRed {
		tr.ImprimeValores()
	}
	fmt.Println("------FINAL DE LA LISTA DE TRANSICIONES---------")
	fmt.Println(" ")
}

// ImprimeLefs : Imprimir los atributos de la clase para depurar errores
func (l Lefs) ImprimeLefs() {

	fmt.Println("STRUCT LEFS")
	//fmt.Println ("\tNº transiciones: ", self.ii_indice)
	fmt.Println("\tNº transiciones: ", l.IaRed.length())

	fmt.Println("------Lista transiciones---------")
	for _, tr := range l.IaRed {
		tr.Imprime()
	}
	fmt.Println("------Final lista transiciones---------")

	fmt.Println("FINAL ESTRUCTURA LEFS")
}
