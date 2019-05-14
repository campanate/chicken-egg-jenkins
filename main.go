package main

import (
	"chicken-egg-jenkins/cli"
	"chicken-egg-jenkins/jenkins"
	"flag"
	"fmt"

)

func main() {

	args := flag.Args()
	_, err := cli.Parse(args, "")

	err = jenkins.CreateJenkins()

	if err != nil {
		fmt.Println(err)
	}
}