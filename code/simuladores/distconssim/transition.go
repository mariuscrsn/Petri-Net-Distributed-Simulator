package distconssim

import (
	"distconssim/utils"
)

//--------------------------------------------------------------------------

const TRANS_IND_ERROR IndTrans = -1

// IndTrans is a index of a transition in the local lefs list
type IndTrans int

//TypeConst is the constant to propagate in lefs
type TypeConst int

//------------------------------------------------------------------------

// Transition : Tipo abstracto  para guardar la informacion de una transicion
type Transition struct {
	// IiIdLocal en la tabla global de transiciones
	IiIndGlobal IndTrans `json:"ii_idglobal"`

	// iiValorLef es el valor que tiene la funcion de
	// sensibilizacion en el instante de tiempo que nos da
	// la variable ii_tiempo
	IiValorLef TypeConst `json:"ii_valor"`
	IiTiempo   TypeClock `json:"ii_tiempo"`

	// tiempo que dura el disparo de la transicion
	IiDuracionDisparo TypeClock `json:"ii_duracion_disparo"`

	// vector con parejas :
	//		transicion junto con cte a actualizarle de forma inmediata
	//TransConstIul [][2]int `json:"ii_listactes_IUL"`
	TransConstIul [][]int `json:"ii_listactes_IUL"`
	// vector con parejas :
	//		de transiciones a las que tengo que propagar cte
	// 		en el tiempo de disparo de esta transicion, junto con la cte que
	// 		tengo que propagar
	TransConstPul [][2]int `json:"ii_listactes_PUL"`

	// True si y solo si es una trans de salida hacia otro nodo
	DeSalida bool `json:"ib_desalida"`
}

// actualizaTiempo modifica el tiempo de la transicion dada
func (t *Transition) actualizaTiempo(aiTi TypeClock) {
	// Modificacion del tiempo
	t.IiTiempo = aiTi
}

// updateFuncValue modifica valor funcion de sensibilizacion de transicion dada
// RECIBE: Codigo de la transicion y valor con el que modificar
//		OJO, no es el valor definitivo, sino la CTE a a�adir al valor que tenia
//		antes la funcion
func (t *Transition) updateFuncValue(aiValLef TypeConst) {
	// Modificacion del valor de la funcion lef
	t.IiValorLef += aiValLef
}

// Imprime los atributos de una transicion para depurar errores
func (t *Transition) Imprime(logger *utils.Logger) {
	logger.Trace.Println("Dato Transicion:")
	logger.Trace.Println("IDLOCALTRANSICION: ", t.IiIndGlobal)
	logger.Trace.Println("\tVALOR LEF: ", t.IiValorLef)
	logger.Trace.Println("\tTIEMPO: ", t.IiTiempo)
	logger.Trace.Println("\tDURACION DISPARO: ", t.IiDuracionDisparo)
	logger.Trace.Println("\tLISTA DE CTES IUL: ")
	for _, v := range t.TransConstIul {
		logger.Trace.Println("\tTRANSICION: ", v[0], "\t\tCTE: ", v[1])
	}
	logger.Trace.Printf("\tLISTA DE CTES PUL: \n")
	for _, v := range t.TransConstPul {
		logger.Trace.Printf("\t\tTRANSICION: %d\t\tCTE: %d\n", v[0], v[1])
	}
	logger.Trace.Println()
}

// ImprimeValores de la transición
func (t *Transition) ImprimeValores(logger *utils.Logger) {
	logger.Trace.Println("Transicion -> ")
	logger.Trace.Println("\tIDLOCALTRANSICION: ", t.IiIndGlobal)
	logger.Trace.Println("\t\tVALOR LEF: ", t.IiValorLef)
	logger.Trace.Println("\t\tTIEMPO: ", t.IiTiempo)
}

//--------------------------------------------------------------------------

// TransitionList is a list of transitions themselves
type TransitionList []Transition //Slice de transiciones como Lista

// length return length of ListTransitions with type adapted to IndTrans
func (tl TransitionList) length() IndTrans {
	return IndTrans(len(tl))
}

// Get local IndTrans of global transition
func (tl *TransitionList) getLocalIndTrans(indGlob IndTrans) IndTrans {
	for indLoc, tr := range *tl {
		if tr.IiIndGlobal == indGlob {
			return IndTrans(indLoc)
		}
	}
	return TRANS_IND_ERROR
}

func (tl TransitionList) ImprimeTL(logger *utils.Logger) {
	logger.Trace.Println("Transition list: [ ")
	for _, tr := range tl {
		logger.Trace.Println(tr)
	}
}

//----------------------------------------------------------------------

// TransitionStack is a Stack of transition indices
type TransitionStack []IndTrans

// MakeTransitionStack crea lista de tamaño aiLongitud
func MakeTransitionStack(capacidad int) TransitionStack {
	// cero length and capacidad capacity
	return make(TransitionStack, 0, capacidad)
}

// push transition id to stack
func (st *TransitionStack) push(iTr IndTrans) {
	*st = append(*st, iTr)
}

// pop transition id from stack
func (st *TransitionStack) pop() IndTrans {
	if (*st).isEmpty() {
		return -1
	}

	iTr := (*st)[len(*st)-1] // obtener dato de lo alto de la pila
	*st = (*st)[:len(*st)-1] //desempilar

	return iTr
}

// isEmpty  the transition stack ?
func (st TransitionStack) isEmpty() bool {
	return len(st) == 0
}

func (st TransitionStack) ImprimeTransStack(logger *utils.Logger) {
	if st.isEmpty() {
		logger.Trace.Println("\tStack TRANSICIONES VACIA")
	} else {
		for _, iTr := range st {
			logger.Trace.Println("\t\t\t", iTr)
		}
	}
}
