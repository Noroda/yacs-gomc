package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/alteamc/minequery/v2"
	"github.com/zan8in/masscan"
)

// BESTPORTS "21,22,80,U:137,U:161,443,445,U:1900,3306,3389,U:5353,8080"
const BESTPORTS = `
	25565
	`

func cleartext(text string) string {
	if strings.ContainsAny(text, "§l") == true {
		text := strings.ReplaceAll(text, "§l", "")
		return text
	} else if strings.ContainsAny(text, "§0") == true {
		text := strings.ReplaceAll(text, "§0", "")
		return text
	} else if strings.ContainsAny(text, "§1") == true {
		text := strings.ReplaceAll(text, "§1", "")
		return text
	} else if strings.ContainsAny(text, "§2") == true {
		text := strings.ReplaceAll(text, "§2", "")
		return text
	} else if strings.ContainsAny(text, "§3") == true {
		text := strings.ReplaceAll(text, "§3", "")
		return text
	} else if strings.ContainsAny(text, "§4") == true {
		text := strings.ReplaceAll(text, "§4", "")
		return text
	} else if strings.ContainsAny(text, "§5") == true {
		text := strings.ReplaceAll(text, "§5", "")
		return text
	} else if strings.ContainsAny(text, "§6") == true {
		text := strings.ReplaceAll(text, "§6", "")
		return text
	} else if strings.ContainsAny(text, "§7") == true {
		text := strings.ReplaceAll(text, "§7", "")
		return text
	} else if strings.ContainsAny(text, "§8") == true {
		text := strings.ReplaceAll(text, "§8", "")
		return text
	} else if strings.ContainsAny(text, "§9") == true {
		text := strings.ReplaceAll(text, "§9", "")
		return text
	} else if strings.ContainsAny(text, "§a") == true {
		text := strings.ReplaceAll(text, "§a", "")
		return text
	} else if strings.ContainsAny(text, "§b") == true {
		text := strings.ReplaceAll(text, "§b", "")
		return text
	} else if strings.ContainsAny(text, "§c") == true {
		text := strings.ReplaceAll(text, "§c", "")
		return text
	} else if strings.ContainsAny(text, "§d") == true {
		text := strings.ReplaceAll(text, "§d", "")
		return text
	} else if strings.ContainsAny(text, "§e") == true {
		text := strings.ReplaceAll(text, "§e", "")
		return text
	} else if strings.ContainsAny(text, "§f") == true {
		text := strings.ReplaceAll(text, "§f", "")
		return text
	} else if strings.ContainsAny(text, "\n") == true {
		text := strings.ReplaceAll(text, "\n", " ")
		return text
	}
	return text
}

func main() {
	//context, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	//defer cancel()

	var (
		scannerResult []masscan.ScannerResult
		errorBytes    []byte
	)

	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets("176.9.0.0-176.9.255.255"),
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
					f, err := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
					if err == nil {
						if motd == "" {
							fmt.Printf("%s (%s) (%s:%s) (%s/%s) Unable to get motd\n", cleartext(res.VersionName), fmt.Sprint(res.ProtocolVersion), srs.IP, srs.Port, fmt.Sprint(res.OnlinePlayers), fmt.Sprint(res.MaxPlayers))
							f.WriteString(fmt.Sprintf("%s (%s) (%s:%s) (%s:%s) Unable to get motd\n", cleartext(res.VersionName), fmt.Sprint(res.ProtocolVersion), srs.IP, srs.Port, fmt.Sprint(res.OnlinePlayers), fmt.Sprint(res.MaxPlayers)))
						} else if motd != "" {
							fmt.Printf("%s (%s) (%s:%s) (%s/%s) %s\n", cleartext(res.VersionName), fmt.Sprint(res.ProtocolVersion), srs.IP, srs.Port, fmt.Sprint(res.OnlinePlayers), fmt.Sprint(res.MaxPlayers), cleartext(fmt.Sprint(res.Description.(map[string]interface{})["text"])))
							f.WriteString(fmt.Sprintf("%s (%s) (%s:%s) (%s/%s) %s\n", cleartext(res.VersionName), fmt.Sprint(res.ProtocolVersion), srs.IP, srs.Port, fmt.Sprint(res.OnlinePlayers), fmt.Sprint(res.MaxPlayers), cleartext(fmt.Sprint(res.Description.(map[string]interface{})["text"]))))
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
