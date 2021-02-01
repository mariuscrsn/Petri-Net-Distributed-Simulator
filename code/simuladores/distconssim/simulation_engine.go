/*Package distconssim with several files to offer a distributed conservative simulation
PROPOSITO: Tipo abstracto para realizar la simulacion de una (sub)RdP.
COMENTARIOS:
	- El resultado de una simulacion local sera un slice dinamico de
	componentes, de forma que cada una de ella sera una structura estatica de
	dos enteros, el primero de ellos sera el codigo de la transiciondisparada y
	el segundo sera el valor del reloj local para el que se disparo.
*/
package distconssim

import (
	"distconssim/utils"
	"time"
)

// TypeClock defines integer size for holding time.
type TypeClock int64

// ResultadoTransition holds fired transition id and time of firing
type ResultadoTransition struct {
	CodTransition     IndTrans
	ValorRelojDisparo TypeClock
}

// SimulationEngine is the basic data type for simulation execution
type SimulationEngine struct {
	Node               Node      // Estructura para la comunicación con el resto de nodos
	ilMisLefs          Lefs      // Estructura de datos del simulador
	iiRelojLocal       TypeClock // Valor de mi reloj local
	iiFinClk           TypeClock
	IlEventosPend      EventList             //Lista de eventos a procesar
	ivTransResults     []ResultadoTransition // slice dinámico con los resultados
	EventNumber        float64               // cantidad de eventos ejecutados
	MapTransitionsNode MapTransitionNode     // diccionario con el nombre del nodo en el que se encuentra cada transición
	ChRecvEv           chan Event
	ChFinish           chan bool
	Logger             *utils.Logger
}

// MakeMotorSimulation : inicializar SimulationEngine struct
func MakeMotorSimulation(node *Node, alLaLef Lefs, transDistr MapTransitionNode, finClk TypeClock, logger *utils.Logger) *SimulationEngine {
	m := SimulationEngine{}
	m.Node = *node
	m.iiRelojLocal = 0
	m.iiFinClk = finClk
	m.ilMisLefs = alLaLef
	m.IlEventosPend = MakeEventList(utils.MaxEventsQueueCap) //aun siendo dinámicos...
	m.ivTransResults = make([]ResultadoTransition, 0, utils.MaxEventsQueueCap)
	m.EventNumber = 0
	m.MapTransitionsNode = transDistr
	m.Logger = logger
	m.ChRecvEv = make(chan Event, 0)
	m.ChFinish = make(chan bool, 1)
	go m.Node.waitEvent(m.ChRecvEv, m.ChFinish)
	return &m
}

// disparar una transicion. Esto es, generar todos los eventos
//	   ocurridos por el disparo de una transicion
//   RECIBE: Indice en el vector de la transicion a disparar
func (se *SimulationEngine) dispararTransicion(ilTr IndTrans) {
	// Prepare 5 local variables
	trList := se.ilMisLefs.IaRed              // transition list
	timeTrans := trList[ilTr].IiTiempo        // time to spread to new events
	timeDur := trList[ilTr].IiDuracionDisparo // firing time length
	listPul := trList[ilTr].TransConstPul     // Pul list of pairs Trans, Ctes
	listIul := trList[ilTr].TransConstIul     // Iul list of pairs Trans, Ctes

	// First apply Iul propagations (Inmediate : 0 propagation time)
	for _, trCo := range listIul {
		trIndLocal := trList.getLocalIndTrans(IndTrans(trCo[0]))
		// check if IUL transition belongs to local transition
		if trIndLocal == TRANS_IND_ERROR {
			//TODO: it always should be true, delete when it works
			se.Logger.Error.Printf("Error: IUL transition not belongs to local transition."+
				"Node: [%s] - IUL trans: %v\n", se.Node.Name, trCo)
		}
		trList[trIndLocal].updateFuncValue(TypeConst(trCo[1]))
	}
	// Generamos eventos ocurridos por disparo de transicion ilTr
	for _, trCo := range listPul {
		// tiempo = tiempo de la transicion + coste disparo
		evDstNode := se.MapTransitionsNode[IndTrans(trCo[0])]
		ev := Event{timeTrans + timeDur,
			IndTrans(trCo[0]),
			TypeConst(trCo[1]), 0, se.Node.Name}
		if evDstNode != se.Node.Name { // la transición destino está en otro nodo
			se.Logger.Info.Printf("Sending event %s to node [%s]\n", ev, evDstNode)
			se.Node.sendEvent(&ev, evDstNode)
		} else {
			se.Logger.Trace.Printf("Appending local event: %s\n", ev)
			se.IlEventosPend.inserta(ev)
		}
	}
}

/* fireEnabledTransitions dispara todas las transiciones sensibilizadas
   		PROPOSITO: Accede a lista de transiciones sensibilizadas y procede con
	   	su disparo, lo que generara nuevos eventos y modificara el marcado de
		transicion disparada. Igualmente anotara en resultados el disparo de
		cada transicion para el reloj actual dado
*/
func (se *SimulationEngine) fireEnabledTransitions(aiLocalClock TypeClock) {
	for se.ilMisLefs.haySensibilizadas() { //while
		liCodTrans := se.ilMisLefs.getSensibilizada()
		se.dispararTransicion(liCodTrans)

		// Anotar el Resultado que disparo la liCodTrans en tiempoaiLocalClock
		se.ivTransResults = append(se.ivTransResults,
			ResultadoTransition{liCodTrans, aiLocalClock})
	}
}

// tratarEventos : Accede a lista eventos y trata todos con tiempo aiTiempo
func (se *SimulationEngine) tratarEvento(ev *Event) {
	if ev.getTiempo() == se.iiRelojLocal {
		se.Logger.Info.Printf("Processing next event: %s .....\n", ev)
		trList := se.ilMisLefs.IaRed // obtener lista de transiciones de Lefs
		localIndTrans := trList.getLocalIndTrans(ev.IiTransicion)
		if localIndTrans == TRANS_IND_ERROR {
			se.Logger.Error.Panicf("Cannot find transition id [%d] in >> %v <<", ev.IiTransicion, trList)
		}
		// Establecer nuevo valor de la funcion
		trList[localIndTrans].updateFuncValue(ev.IiCte)
		// Establecer nuevo valor del tiempo
		trList[localIndTrans].actualizaTiempo(ev.IiTiempo)

		se.EventNumber++
	} else {
		se.Logger.Error.Panicf("Processing event in other time: event time: %d - local time: %d", ev.IiTiempo, se.iiRelojLocal)
	}
}

// devolverResultados : Mostrar los resultados de la simulacion
func (se SimulationEngine) devolverResultados() {
	se.Logger.Info.Println("----------------------------------------")
	se.Logger.Info.Println("Resultados del simulador local")
	se.Logger.Info.Println("----------------------------------------")
	if len(se.ivTransResults) == 0 {
		se.Logger.Info.Println("No esperes ningun resultado...")
	}

	for _, liResult := range se.ivTransResults {
		se.Logger.Info.Printf("TIEMPO: %d  -> TRANSICION: %d\n", liResult.ValorRelojDisparo, liResult.CodTransition)
	}

	se.Logger.Info.Printf("========== TOTAL DE TRANSICIONES DISPARADAS = %d\n", len(se.ivTransResults))
}

// Devuelve el evento menor entre la FIFO con tiempo menor y el local.
// Devuelve nil si el menor tiempo no tiene eventos pendientes.
// Devuelve true si el evento es local y falso si es remoto
func (se *SimulationEngine) getLowerEvent() (*Event, bool) {
	_, lowestTimeNode := se.Node.getLowerTimeFIFO()
	// There is any local event
	if se.IlEventosPend.isEmpty() {
		// No events from retarded node
		if lowestTimeNode.IncomingEvFIFO.isEmpty() {
			return nil, false
		}

		// Event in lazy node FIFO
		ev := lowestTimeNode.IncomingEvFIFO.leePrimerEvento()
		return &ev, false
	}

	// Get local event
	localEv := se.IlEventosPend.leePrimerEvento()
	se.Logger.Trace.Printf("getLowerEvent: LOCAL->%s, lowerFIFOtime: %d\n", localEv, lowestTimeNode.RemoteSafeTime)
	// No events in lazy node FIFO
	if lowestTimeNode.IncomingEvFIFO.isEmpty() {
		if localEv.IiTiempo > lowestTimeNode.RemoteSafeTime {
			return nil, false // I should return remote ev, but it not exist
		} else {
			return &localEv, true
		}
	}

	// Events in lazy node FIFO
	remoteEv := lowestTimeNode.IncomingEvFIFO.leePrimerEvento()
	if localEv.IiTiempo <= remoteEv.IiTiempo {
		return &localEv, true
	} else {
		return &remoteEv, false
	}
}

// Get the lowerEvent. If it not exists, blocks until new event arrive. If it is an event and is the lowest return it.
// If not, blocks again until receive the lowset and all FIFO have at least one event. If the recv event is null with
// lower time, blocks again until receives the correct one.
func (se *SimulationEngine) getNextEvent() *Event {
	// Iterate until get a processable event or finish event
	for {
		ev, isLocalEv := se.getLowerEvent()
		// Not blocked, I've get an event to process
		if ev != nil {
			// delete event for list
			if isLocalEv {
				se.IlEventosPend.eliminaPrimerEvento()
				se.Logger.Info.Printf("Lower event is local: %s\n", ev)
			} else {
				name, remoteNode := se.Node.getLowerTimeFIFO()
				remoteNode.IncomingEvFIFO.eliminaPrimerEvento()
				se.Node.Partners[name] = remoteNode
				se.Logger.Info.Printf("Lower event is remote: %s\n", ev)
			}
			return ev
		}

		// I'm gonna to block, send before it an NULL event
		_, lowestNodeTime := se.Node.getLowerTimeFIFO()
		lowestTime := lowestNodeTime.RemoteSafeTime
		clkFstEvLocal := se.IlEventosPend.leePrimerEvento().IiTiempo
		if lowestTime > clkFstEvLocal && !se.IlEventosPend.isEmpty() { // Time on NULL event depend on local events
			lowestTime = clkFstEvLocal
		}
		// Send null event or finish
		nextTime := lowestTime + lowestNodeTime.LookAhead
		if lowestTime >= se.iiFinClk || nextTime >= se.iiFinClk {
			return se.FinishSim()
		}
		nullEv := Event{
			IiTransicion: 0,
			IiCte:        0,
			Is_Sender:    se.Node.Name,
			IiTiempo:     nextTime,
			Ib_IsNULL:    1,
		}
		se.Logger.Trace.Printf("Sending NULL event: %s\n", nullEv)
		if nullEv.getTiempo() <= se.iiFinClk {
			se.Node.sendEv2All(&nullEv)
		}

		// Wait for a event or a null message
		se.Logger.Trace.Printf("	Waiting for incoming event...\n")
		recvEv := <-se.ChRecvEv

		// Process event
		if recvEv.IsClosingEvent() {
			se.Logger.Warning.Printf("Received closing event %s\n", recvEv)
			time.Sleep(utils.TimeWaitStop)
			se.ChFinish <- true
			return &recvEv
		} else if recvEv.IsNullEvent() {
			// Update RemoteSafeTime
			sender := se.Node.Partners[recvEv.getSender()]
			sender.RemoteSafeTime = recvEv.IiTiempo
			se.Node.Partners[recvEv.getSender()] = sender

			// Check if any event has been unblocked
			se.Logger.Trace.Printf("NULL received: %s\n", recvEv)
		} else { // Event can be processed
			// Insert event in remote node FIFO
			se.Logger.Trace.Printf("Adding received event to incFIFO: %s\n", recvEv)
			senderNode := se.Node.Partners[recvEv.getSender()]
			senderNode.RemoteSafeTime = recvEv.IiTiempo
			senderNode.IncomingEvFIFO.inserta(recvEv)
			se.Node.Partners[recvEv.getSender()] = senderNode
		}
	}
}

// SimularUnpaso de una RdP con duración disparo >= 1. Devuelve true si se ha procesado el ultimo evento
func (se *SimulationEngine) simularUnpaso() bool {

	se.Logger.Info.Printf("####################### CLK = %d #######################\n", se.iiRelojLocal)
	se.ilMisLefs.ImprimeLefs()
	se.ilMisLefs.actualizaSensibilizadas(se.iiRelojLocal)
	se.Logger.Trace.Println(">>>>>>>>>> Stack de transiciones sensibilizadas <<<<<<<<")
	se.ilMisLefs.TransSensib.ImprimeTransStack(se.Logger)
	se.Logger.Trace.Println(">>>>>>>>>> Final Stack de transiciones <<<<<<<<<<<<<<<<<")

	// Fire enabled transitions and produce events
	if se.ilMisLefs.haySensibilizadas() {
		se.fireEnabledTransitions(se.iiRelojLocal)
	}

	se.Logger.Trace.Println("·········· Lista eventos después de disparos ········")
	se.Logger.Trace.Printf("Eventos locales: %s\n", se.IlEventosPend)
	se.Logger.Trace.Println(se.Node.Partners.StringFIFO())
	se.Logger.Trace.Println("·········· Final lista eventos ······················")

	ev := se.getNextEvent()
	if ev.IsClosingEvent() {
		return true
	} else if ev != nil {
		// advance local clock to soonest available event
		se.iiRelojLocal = ev.IiTiempo
		se.Logger.Trace.Printf("+++ NEXT CLOCK: %d +++\n", se.iiRelojLocal)

		// if events exist for current local clock, process them
		se.tratarEvento(ev)
		return false
	} else {
		se.Logger.Error.Panicf("Simulating nil event\n")
		return true
	}
}

func (se *SimulationEngine) FinishSim() *Event {
	// Send closing event to partners
	ev := Event{
		IiTiempo:     se.iiRelojLocal,
		IiTransicion: FINISH_EVENT,
		IiCte:        0,
		Ib_IsNULL:    0,
	}
	se.Logger.Info.Println("Sending closing event")
	se.Node.sendEv2All(&ev)
	time.Sleep(utils.TimeWaitStop)
	se.ChFinish <- true
	return &ev
}

// SimularPeriodo de una RdP
// RECIBE: - Ciclo inicial (por si marcado recibido no se corresponde al
//				inicial sino a uno obtenido tras simular ai_cicloinicial ciclos)
//		   - Ciclo con el que terminamos
func (se *SimulationEngine) SimularPeriodo() {
	ldIni := time.Now()

	//se.Node.Wait4PartnersSetup()
	finish := false
	initClk := se.iiRelojLocal
	for se.iiRelojLocal < se.iiFinClk && !finish {
		finish = se.simularUnpaso()
	}

	if !finish {
		se.FinishSim()
	}
	elapsedTime := time.Since(ldIni)

	se.Logger.Info.Printf("Eventos por segundo = %.4f\n",
		se.EventNumber/elapsedTime.Seconds())

	// Devolver los resultados de la simulacion
	se.devolverResultados()
	se.Logger.Info.Println("---------------------")
	se.Logger.Info.Printf("TIEMPO SIMULADO en ciclos: %d\n", se.iiFinClk-initClk)
	se.Logger.Info.Printf("TIEMPO ejecución REAL simulación: %s\n", elapsedTime.String())
}
