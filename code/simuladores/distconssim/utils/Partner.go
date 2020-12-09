package utils

import (
	"distconssim"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Partner struct {
	IP              string `json:"IP"`
	Port            int    `json:"Port"`
	IncommingEvList distconssim.EventList
}

func ReadPartners() *map[string]Partner {

	// read file
	data, err := ioutil.ReadFile(AbsWorkPath + RelExecPath + "network.json")
	if err != nil {
		panic(err)
	}

	fmt.Println("Reading network configuration file...")

	var partners map[string]Partner
	// parse content of json file to Config struct
	err = json.Unmarshal(data, &partners)
	if err != nil {
		panic(err)
	}

	return &partners
}
