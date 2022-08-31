package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/alteamc/minequery/v2"
	"github.com/zan8in/masscan"
)

func main() {
	//context, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	//defer cancel()

	//godotenv.Load()

	//PORTRANGE := os.Getenv("PORT_RANGE")
	//IPRANGE := os.Getenv("IP_RANGE")
	IPRANGE1 := flag.String("range", "127.0.0.1", "IP range to scan")
	PORTRANGE1 := flag.String("port-range", "25565", "Port range to scan")
	OUTFILE1 := flag.String("output", "output.txt", "You can't disable it")
	flag.Parse()
	IPRANGE2 := *IPRANGE1
	PORTRANGE2 := *PORTRANGE1
	OUTFILE2 := *OUTFILE1

	var (
		scannerResult []masscan.ScannerResult
		errorBytes    []byte
	)

	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets(IPRANGE2),
		masscan.SetParamPorts(PORTRANGE2),
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
					f, err := os.OpenFile(fmt.Sprint(OUTFILE2), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
					if err == nil {
						if motd == "" {
							fmt.Printf("%s (%s) (%s:%s) (%s/%s) Unable to get motd\n", res.VersionName, fmt.Sprint(res.ProtocolVersion), srs.IP, srs.Port, fmt.Sprint(res.OnlinePlayers), fmt.Sprint(res.MaxPlayers))
							f.WriteString(fmt.Sprintf("%s (%s) (%s:%s) (%s:%s) Unable to get motd\n", res.VersionName, fmt.Sprint(res.ProtocolVersion), srs.IP, srs.Port, fmt.Sprint(res.OnlinePlayers), fmt.Sprint(res.MaxPlayers)))
						} else if motd != "" {
							fmt.Printf("%s (%s) (%s:%s) (%s/%s) %s\n", res.VersionName, fmt.Sprint(res.ProtocolVersion), srs.IP, srs.Port, fmt.Sprint(res.OnlinePlayers), fmt.Sprint(res.MaxPlayers), fmt.Sprint(res.Description.(map[string]interface{})["text"]))
							f.WriteString(fmt.Sprintf("%s (%s) (%s:%s) (%s/%s) %s\n", res.VersionName, fmt.Sprint(res.ProtocolVersion), srs.IP, srs.Port, fmt.Sprint(res.OnlinePlayers), fmt.Sprint(res.MaxPlayers), fmt.Sprint(res.Description.(map[string]interface{})["text"])))
						}
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
