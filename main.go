package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/alteamc/minequery/v2"
	"github.com/zan8in/masscan"
)

// BESTPORTS "21,22,80,U:137,U:161,443,445,U:1900,3306,3389,U:5353,8080"
const BESTPORTS = `
	25565
	`

func main() {
	//context, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	//defer cancel()

	var (
		scannerResult []masscan.ScannerResult
		errorBytes    []byte
	)

	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets("176.9.20.0-176.9.20.255"),
		masscan.SetParamPorts(BESTPORTS),
		masscan.EnableDebug(),
		masscan.SetParamWait(0),
		masscan.SetParamRate(1000),
	)

	if err != nil {
		log.Fatalf("unable to create masscan scanner: %v", err)
	}

	if err := scanner.RunAsync(); err != nil {
		panic(err)
	}

	stdout := scanner.GetStdout()

	stderr := scanner.GetStderr()

	go func() {
		for stdout.Scan() {
			srs := masscan.ParseResult(stdout.Bytes())
			scannerResult = append(scannerResult, srs)
			if i, err := strconv.Atoi(srs.Port); err == nil {
				pinger := minequery.NewPinger(
					minequery.WithUseStrict(true),
				)
				res, err := pinger.Ping17(srs.IP, int(i))
				if err == nil {
					motd := res.Description.(map[string]interface{})["text"]
					if motd == "" {
						fmt.Printf("%s (%s) (%s:%s) Unable to get motd\n", res.VersionName, fmt.Sprint(res.ProtocolVersion), srs.IP, srs.Port)
					} else if motd != "" {
						fmt.Printf("%s (%s) (%s:%s) (%s/%s) %s\n", res.VersionName, fmt.Sprint(res.ProtocolVersion), srs.IP, srs.Port, fmt.Sprint(res.OnlinePlayers), fmt.Sprint(res.MaxPlayers), res.Description.(map[string]interface{})["text"])
					}
				} else if err != nil {
					fmt.Println(err)
				}
			}
		}
	}()

	go func() {
		for stderr.Scan() {
			fmt.Println(stderr.Text())
			errorBytes = append(errorBytes, stderr.Bytes()...)
		}
	}()

	if err := scanner.Wait(); err != nil {
		panic(err)
	}

	fmt.Println("masscan result count : ", len(scannerResult))

}
