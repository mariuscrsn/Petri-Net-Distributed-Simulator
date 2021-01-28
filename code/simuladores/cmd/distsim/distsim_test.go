package main

import (
	"distconssim"
	"distconssim/utils"
	"fmt"
	"golang.org/x/crypto/ssh"
	"sync"
	"testing"
)

var sshConn map[string]*ssh.Client

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
		fmt.Printf("Node [%s]:$ %s", nodeName, cmd)
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
	ms := createMotorSimulation("P0", filesPrefix)
	startNodes(filesPrefix, finalClk, ms.Node.Partners, &wg)
	ms.SimularPeriodo(0, finalClk)

	wg.Wait()
	terminate()
}
