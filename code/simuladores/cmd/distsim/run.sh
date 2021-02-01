#!/bin/bash

TestCases="2SubNets2Br 3SubNets2Br 6SubNets5BrHomogen 6SubNets5Br1BrSlow 6SubNets5BrLA"

rm -r bin/distsim results/*
echo "Building project..."
go build -i -o bin/distsim distsim.go

for i in ${TestCases}; do
  echo "Running test ${i} ->->->->"
  go test -v  -run "$i" distsim_test.go
done