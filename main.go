package main

import (
	"encoding/json"
	"fmt"
	"os"
	"tabular_method/quinemccluskey"
)

type logicFunctionOutput struct {
	S []uint64 `json:"s"`
	D []uint64 `json:"d"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	functionFilePath := os.Args[1]
	jsonStream, err := os.ReadFile(functionFilePath)
	check(err)

	var inLabels quinemccluskey.InputLabels
	var outLabels quinemccluskey.OutputLabels

	var topLevel map[string]json.RawMessage
	e := json.Unmarshal(jsonStream, &topLevel)
	check(e)

	if inputLabelBlob, ok := topLevel["inputs"]; ok {
		var inputLabels map[int]string
		e := json.Unmarshal(inputLabelBlob, &inputLabels)
		check(e)

		for label := range inputLabels {
			inLabels.Set(label, inputLabels[label])
		}
	}

	outputs := []logicFunctionOutput{}
	for item := range topLevel {
		if item == "inputs" {
			continue
		}

		outputs = append(outputs, logicFunctionOutput{})
		e := json.Unmarshal(topLevel[item], &outputs[len(outputs)-1])
		check(e)
		outLabels.Add(item)
	}

	var logicFunction quinemccluskey.LogicFunction
	logicFunction.Init(true)

	for _, output := range outputs {
		logicFunction.AddOutput(output.S, output.D)
	}

	fmt.Printf(logicFunction.GetMinimumCostCover(inLabels, outLabels))

}
