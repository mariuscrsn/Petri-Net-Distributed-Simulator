package utils

import "time"

//const HomePath = "/home/a721609/"
//const AbsWorkPath = HomePath + "Desktop/redes/miniproyecto/"
const HomePath = "/home/cms/"
const AbsWorkPath = HomePath + "Escritorio/Datos/UNI/redes/practicas/miniproyecto/code/simuladores/"
const RelOutputPath = "results/"
const LoggerPath = AbsWorkPath + "cmd/distsim/results/"
const RelTestDataPath = "cmd/distsim/testdata/"
const MaxAttempsConnect = 5
const PeriodRetry = 2 * time.Second // second between connection retries
const BinFilePath = "cmd/distsim/bin/distsim"
const MaxEventsQueueCap = 1000
const TimeWaitStop = 3 * time.Second
