//Package centralsim with several files to offer a centralized simulation
package centralsim

import "fmt"

// Event define el evento básico de simulación
type Event struct {
	// Tiempo para el que debemos considerar el evento
	IiTiempo TypeClock
	// A que transicion (indice transicion en subred)
	IiTransicion IndLocalTrans
	// Constante que mandamos
	IiCte TypeConst
}

/*
-----------------------------------------------------------------
   METODO: NewEvento
   RECIBE: Tiempo, transicion y cte del evento a crear
   DEVUELVE: Event
   PROPOSITO: Crear evento con todos los datos del nuevo evento creados
   HISTORIA DE CAMBIOS:
   COMENTARIOS:
-----------------------------------------------------------------

func NewEvento(ai_tiempo TypeClock, ai_transicion IndLocalTrans, ai_cte TypeConst) *Event {
	e := new(Event)
	set_tiempo(e.ai_tiempo)
	set_transicion(e.ai_transicion)
	set_cte(e.ai_cte)
}
*/

// SetTiempo mmodifica el tiempo del evneto
func (e *Event) SetTiempo(aiTiempo TypeClock) {
	e.IiTiempo = aiTiempo
}

// SetTransicion modifica la transicion del evento
func (e *Event) SetTransicion(aiTransicion IndLocalTrans) {
	e.IiTransicion = aiTransicion
}

// SetCte modifica la cte del evento
func (e *Event) SetCte(aiCte TypeConst) {
	e.IiCte = aiCte
}

// getTiempo obtiene el tiempo del evento
func (e Event) getTiempo() TypeClock {
	return e.IiTiempo
}

// getTransicion obtiene la trasicion del evento
func (e Event) getTransicion() IndLocalTrans {
	return e.IiTransicion
}

// getCte obtiene la cte del evento a aplicar a la transicion
func (e Event) getCte() TypeConst {
	return e.IiCte
}

// Imprime atributos de evento para depurar errores
func (e Event) Imprime(i int) {
	fmt.Println("  Evento -> ", i)
	fmt.Println("    TIEMPO: ", e.IiTiempo)
	fmt.Println("    TRANSICION: ", e.IiTransicion)
	fmt.Println("    CONSTANTE: ", e.IiCte)
}

//----------------------------------------------------------------------------

// EventList es el tipo que almacena la lista de eventos necesaria
// para los motores de	simulacion.
type EventList []Event

// MakeEventList crea lista de tamaño aiLongitud
func MakeEventList(capacidad int) EventList {
	// cero length and capacidad capacity
	return make(EventList, 0, capacidad)
}

// longitud : numero de elementos de la lista eventos
func (el EventList) longitud() int {
	return len(el)
}

// inserta evento en la lista de eventos con ordenación de tiempo
func (el *EventList) inserta(aeEvento Event) {
	var i int // INITIALIZED to 0 !!!

	//fmt.Println("Insertar evento en lista : ", ae_evento, *self)

	// Obtengo la posicion ordenada del evento en slice con i
	for _, e := range *el {
		if e.getTiempo() >= aeEvento.getTiempo() {
			break
		}
		i++
	}

	//fmt.Println("POSICION a INSERTAR en lista de eventos : ", i)
	*el = append((*el)[:i], append([]Event{aeEvento}, (*el)[i:]...)...)

	//fmt.Println("DESPUES de insertar : ", *self)
}

// recogePrimerEvento encolado
func (el EventList) leePrimerEvento() Event {
	if len(el) > 0 {
		return el[0]
	}

	return Event{} //sino devuelve el tipo Event, zeroed
}

// eliminaPrimerEvento encolado
func (el *EventList) eliminaPrimerEvento() {
	if len(*el) > 0 {
		//suprimir con posibilidad de liberacion de memoria
		copy(*el, (*el)[1:])
		(*el)[len(*el)-1] = Event{} //pongo a zero el previo último Event
		(*el) = (*el)[:len(*el)-1]
	}
}

// getPrimerEvento toma el primer evento de la lista de eventos
func (el *EventList) popPrimerEvento() Event {
	/* fmt.Println("Lista antes de eliminar primer evento :")
	(*self).il_eventos.Imprime()
	*/
	leEvento := el.leePrimerEvento()
	el.eliminaPrimerEvento()
	/*fmt.Println("Lista DESPUES de eliminar primer evento :")
	(*self).il_eventos.Imprime()
	*/
	return leEvento
}

/* tiempoPrimerEvento : valor temporal del primer evento para conocer
	   posteriormente si debemos avanzar el reloj local
   			DEVUELVE: El valor del tiempo del primer evento de lista de eventos.
	  				 *** -1 si ocurrio un error o no hay eventos.

*/
func (el *EventList) tiempoPrimerEvento() TypeClock {
	if el.longitud() > 0 {
		return el.leePrimerEvento().IiTiempo
	}

	return -1
}

// hayEventos permite saber si quedan eventos disponible para este tiempo
func (el *EventList) hayEventos(aiTiempo TypeClock) bool {
	if el.tiempoPrimerEvento() == aiTiempo {
		return true
	}

	return false
}

// Imprime la lista de eventos para depurar errores
func (el EventList) Imprime() {
	fmt.Println("Estructura EventList")
	for i, e := range el {
		e.Imprime(i)
	}
}

// FIN DEL TIPO ABSTRACTO EventList
