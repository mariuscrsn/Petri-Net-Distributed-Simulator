// Este programa requiere 2 parámetros de entrada :
//      - Nombre fichero json de Lefs
//        - Número de ciclo final
//
// Ejemplo : censim  testdata/PrimerEjemplo.rdp.subred0.json  5
package main

import (
	"distconssim"
	"distconssim/utils"
	"fmt"
	"os"
	"strconv"
)

func parseFilesNames(nodeName string, filesPrefix string) (string, string) {
	nodeInd, _ := strconv.Atoi(nodeName[len(nodeName)-1:])
	lefsFile := fmt.Sprintf("%s.subred%d.json", filesPrefix, nodeInd)
	netFile := fmt.Sprintf("%s.network.json", filesPrefix)
	return netFile, lefsFile
}

func createMotorSimulation(nodeName string, filesPrefix string) *distconssim.SimulationEngine {
	// init logger, create files and build files names
	err := os.Mkdir(utils.RelOutputPath+"results/"+filesPrefix, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating log dir: %s\n", err)
	}
	netFile, lefsFile := parseFilesNames(nodeName, filesPrefix)
	logger := utils.InitLoggers(filesPrefix, nodeName)

	// read partners and transition mapping to them
	net := distconssim.ReadPartners(netFile)
	partners := net.Nodes
	myNode := partners[nodeName]
	delete(partners, nodeName)
	logger.Info.Printf("[%s] Reading partners: \n%s", nodeName, partners)

	// Create local node
	node := distconssim.MakeNode(nodeName, myNode.Port, &partners, logger)

	// Carga de la subred
	lefs, err := distconssim.LoadLefs(lefsFile, logger)
	if err != nil {
		println("Couln't load the Petri Net file !")
	}
	return distconssim.MakeMotorSimulation(node, lefs, net.MapTransNode, logger)
}

func main() {

	if len(os.Args) != 4 {
		panic("bad usage: distim <nodeName> <files_prefix> <finalClk>")
	}

	var nodeName, filesPrefix string
	nodeName = os.Args[1]
	filesPrefix = os.Args[2]
	netFile, lefsFile := parseFilesNames(nodeName, filesPrefix)

	// init logger
	logger := utils.InitLoggers(filesPrefix, nodeName)

	// read partners and transition mapping to them
	net := distconssim.ReadPartners(netFile)
	partners := net.Nodes
	myNode := partners[nodeName]
	delete(partners, nodeName)
	logger.Info.Printf("[%s] Reading partners: \n%s", nodeName, partners)

	// Create local node
	node := distconssim.MakeNode(nodeName, myNode.Port, &partners, logger)

	// Carga de la subred
	lefs, err := distconssim.LoadLefs(lefsFile, logger)
	if err != nil {
		println("Couln't load the Petri Net file !")
	}
	ms := distconssim.MakeMotorSimulation(node, lefs, net.MapTransNode, logger)

	// ciclo 0 hasta ciclo os.args[2]
	cicloFinal, _ := strconv.Atoi(os.Args[3])
	ms.SimularPeriodo(0, distconssim.TypeClock(cicloFinal))
}
