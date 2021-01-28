#!/bin/bash

FINISH_CLK=4
LEFS_PATH="2subredes.subred"
NET_FILE_PATH="2subredes.network.json"
GO_CMD="/usr/local/go/bin/go"
GO_TARGET="bin/go_build_distsim_go_script"
WORKDIR="/home/cms/Escritorio/Datos/UNI/redes/practicas/miniproyecto/code/simuladores/cmd/distsim/"
GO_MAIN_FILE="distsim.go"
echo "Running script on: $(pwd)"

rm "${WORKDIR}results/Log_P"*
netstat -tlpun  2>/dev/null | grep 1609  | tr -s ' ' | cut -d' ' -f7 | cut -d'/' -f1 | xargs -i kill {}

echo "Building package ..."
${GO_CMD} build -i -o "${WORKDIR}${GO_TARGET}" "${WORKDIR}${GO_MAIN_FILE}"
for i in {1..2} ; do
  nodeIdx="$((i-1))"
  nodeName="P${nodeIdx}"
  echo "Starting node: ${nodeName} ..."
  "${WORKDIR}${GO_TARGET}" "${nodeName}" "${LEFS_PATH}${nodeIdx}.json" "${NET_FILE_PATH}" "${FINISH_CLK}" &
done
