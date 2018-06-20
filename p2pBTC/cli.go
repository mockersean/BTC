package main

import (
	"os"
	"fmt"
	"flag"
)

const Usage  = `
	createchain -address ADDRESS	"创建区块链"
	send -from FROM -to To -amount AMOUNT 			"生成交易"
	printchain						"查看区块"
	getbalance -address ADDRESS 	"余额"
`


type CLI struct {
}

func (cli *CLI) Run() {
	if len(os.Args) < 2{
		fmt.Println(Usage)
		os.Exit(1)
	}

	creactChainCmd := flag.NewFlagSet("createchain",flag.ExitOnError)
	printCmd := flag.NewFlagSet("printCmd",flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance",flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send",flag.ExitOnError)


	creactChainCmdPara := creactChainCmd.String("address","","address data")
	getbalanceCmdPara := getbalanceCmd.String("address","","address data")

	fromPara := sendCmd.String("from","","from data")
	toPara := sendCmd.String("to","","to data")
	amountPara := sendCmd.Float64("amount",0,"amount data")


	switch os.Args[1] {
	case "createchain":
		err := creactChainCmd.Parse(os.Args[2:])
		CheckErr(err)
	case "printchain":
		err := printCmd.Parse(os.Args[2:])
		CheckErr(err)
	case "getbalance":
		err := getbalanceCmd.Parse(os.Args[2:])
		CheckErr(err)

	case "send":
		err := sendCmd.Parse(os.Args[2:])
		CheckErr(err)

	default:
		fmt.Println("invalid cmd\n",Usage)
		os.Exit(1)
	}

	if creactChainCmd.Parsed() {
		if *creactChainCmdPara == ""{
			fmt.Println(Usage)
			os.Exit(1)
		}else {
			cli.CreactChain(*creactChainCmdPara)
		}
	}

	if getbalanceCmd.Parsed() {
		if *getbalanceCmdPara == "" {
			fmt.Println(Usage)
			os.Exit(1)
		}else {
			cli.GetBalance(*getbalanceCmdPara)
		}
	}
	if sendCmd.Parsed() {
		if *fromPara == "" || *toPara == "" || *amountPara == 0 {
			fmt.Println(Usage)
			os.Exit(1)
		}else {
			cli.Send(*fromPara,*toPara,*amountPara)
		}
	}

	if printCmd.Parsed() {
		cli.PrintChain()
	}
}