package main

import (
	"encoding/json"
	"fmt"
)

func main() {

	signatures, err := GetSignatures()
	if err != nil {
		ErrorLog.Println(err)
	}

	resolvers, err := GetResolvers()
	if err != nil {
		ErrorLog.Println(err)
	}

	Signatures = signatures
	Resolvers = resolvers

	data, _ := Dig("0x4f.in")

	outputByte, _ := json.MarshalIndent(data, "", "  ")
	fmt.Println(string(outputByte))

}
