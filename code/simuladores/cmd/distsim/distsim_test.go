package main

import (
	"distconssim"
	"distconssim/utils"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"sync"
	"testing"
)

var sshConn map[string]*ssh.Client

func CreateMotorSimulation(nodeName string, filesPrefix string) *distconssim.SimulationEngine {
	// init logger, create files and build files names
	err := os.Mkdir(utils.RelOutputPath+"results/"+filesPrefix, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating log dir: %s\n", err)
	}
	netFile, lefsFile := utils.ParseFilesNames(nodeName, filesPrefix)
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

func terminate() {
	// var killed_once bool = false
	for _, conn := range sshConn {
		_ = conn.Close()
	}
}

// Source: http://networkbit.ch/golang-ssh-client/
func startNodes(filesPrefix string, finishClk distconssim.TypeClock, partners *distconssim.Partners, wg *sync.WaitGroup) {
	sshConn = make(map[string]*ssh.Client, 0)
	for nodeName, node := range *partners {
		sshConn[nodeName] = utils.ConnectSSH(node.Username, node.IP)

		// Execute program
		fmt.Println("Starting: " + nodeName)
		var cmd = utils.AbsWorkPath + utils.BinFilePath + fmt.Sprintf(" %s %s %d", nodeName, filesPrefix, finishClk)
		fmt.Printf("Node [%s]:$ %s\n", nodeName, cmd)
		go utils.RunCommandSSH(cmd, sshConn[nodeName], wg)
	}
}

func Test2SubNets(t *testing.T) {
	// WaitGroup for syncronization goroutines, test wait until each net finish
	var wg sync.WaitGroup
	wg.Add(1)

	// Setup Motor Simulation of root net
	filesPrefix := "2subredes"
	var finalClk distconssim.TypeClock = 4
	ms := CreateMotorSimulation("P0", filesPrefix)
	startNodes(filesPrefix, finalClk, ms.Node.Partners, &wg)
	fmt.Println("[P0] Simulating net...")
	ms.SimularPeriodo(0, finalClk)

	wg.Wait()
	terminate()
}
