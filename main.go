package main

import "healthcheckr/cmd/healthcheckr"

func main() {
	if err := healthcheckr.GetCommand().Execute(); err != nil {
		panic(err)
	}
}
