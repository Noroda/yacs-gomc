package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/zan8in/masscan"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
)

type status struct {
	Description chat.Message
	Players     struct {
		Max    int
		Online int
		Sample []struct {
			ID   uuid.UUID
			Name string
		}
	}
	Version struct {
		Name     string
		Protocol int
	}
	Favicon Icon
	Delay   time.Duration
}
type Icon string

var outTemp = template.Must(template.New("output").Parse(`
Version: {{ .Version.Name }} ({{ .Version.Protocol }})
Description:
{{ .Description }}
Players: {{ .Players.Online }}/{{ .Players.Max }}{{ range .Players.Sample }}
- {{ .Name }} {{ end }}
`))

func (s *status) String() string {
	var sb strings.Builder
	err := outTemp.Execute(&sb, s)
	if err != nil {
		panic(err)
	}
	return sb.String()
}

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
			resp, delay, err := bot.PingAndList(fmt.Sprint(srs.IP + ":" + srs.Port))
			if err != nil {
				fmt.Printf("Ping and list server fail: %v", err)
				os.Exit(1)
			}
			var s status
			err = json.Unmarshal(resp, &s)
			if err != nil {
				fmt.Print("Parse json response fail:", err)
				os.Exit(1)
			}
			s.Delay = delay
			//ss := *&s
			fmt.Println("IP: " + srs.IP + ":" + srs.Port + fmt.Sprint(&s))
			f, err := os.OpenFile(fmt.Sprint(OUTFILE2), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
			if err == nil {
				f.WriteString("IP: " + srs.IP + ":" + srs.Port + fmt.Sprint(&s))
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
