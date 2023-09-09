package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	// "os"
	//"strings"
	"sync"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/zan8in/masscan"

	"github.com/Tnze/go-mc/bot"
	"github.com/Tnze/go-mc/chat"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	_ "go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

/*
type ServerDB struct {
	ServerIP       string `bson:"serverIP"`
	Description    string `bson:"description"`
	Version        string `bson:"version"`
	OnlinePlayers  int    `bson:"onlinePlayers"`
	MaxPlayers     int    `bson:"maxPlayers"`
	FoundAt        string `bson:"foundAt"`
	CurrentPlayers []string `bson:"currentPlayers"`
}
*/

type ServerDB struct {
	ServerIP    string `bson:"serverIP"`
	Description string `bson:"description"`
	Version     string `bson:"version"`
	Players     struct {
		Max    int `bson:"max"`
		Online int `bson:"online"`
		List   []struct {
			ID   string `bson:"id"`
			Name string `bson:"name"`
		} `bson:"list"`
	} `bson:"players"`
	FoundAt string `bson:"foundAt"`
	Favicon Icon   `bson:"favicon"`
}

type ServerDBbutMf struct {
	ServerIP    string `bson:"serverIP"`
	Description string `bson:"description"`
	Version     string `bson:"version"`
	Players     struct {
		Max    int `bson:"max"`
		Online int `bson:"online"`
		List   []struct {
			ID   uuid.UUID `bson:"id"`
			Name string    `bson:"name"`
		} `bson:"list"`
	} `bson:"players"`
	FoundAt string `bson:"foundAt"`
	Favicon Icon   `bson:"favicon"`
}

func convertServerDB(server ServerDBbutMf) ServerDB {
	converted := ServerDB{
		ServerIP:    server.ServerIP,
		Description: server.Description,
		Version:     server.Version,
		FoundAt:     server.FoundAt,
	}

	converted.Players.Max = server.Players.Max
	converted.Players.Online = server.Players.Online
	converted.Favicon = server.Favicon

	converted.Players.List = make([]struct {
		ID   string `bson:"id"`
		Name string `bson:"name"`
	}, len(server.Players.List))

	for i, player := range server.Players.List {
		converted.Players.List[i].ID = fmt.Sprint(player.ID)
		converted.Players.List[i].Name = player.Name
	}

	return converted
}

func main() {
	//context, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	//defer cancel()

	//godotenv.Load()

	//PORTRANGE := os.Getenv("PORT_RANGE")
	//IPRANGE := os.Getenv("IP_RANGE")
	IPRANGE1 := flag.String("range", "127.0.0.1", "IP range to scan")
	PORTRANGE1 := flag.String("port-range", "25565", "Port range to scan")
	//OUTFILE1 := flag.String("output", "output.txt", "You can't disable it")
	RATE1 := flag.Int("rate", 1000, "masscan rate")
	TIMEOUT := flag.Duration("timeout", 5*time.Second, "pinger timeout time")
	flag.Parse()
	TIMEOUT2 := *TIMEOUT
	IPRANGE2 := *IPRANGE1
	PORTRANGE2 := *PORTRANGE1
	//OUTFILE2 := *OUTFILE1
	RATE2 := *RATE1

	var (
		scannerResult []masscan.ScannerResult
		errorBytes    []byte
	)

	scanner, err := masscan.NewScanner(
		masscan.SetParamTargets(IPRANGE2),
		masscan.SetParamPorts(PORTRANGE2),
		masscan.EnableDebug(),
		masscan.SetParamWait(0),
		masscan.SetParamRate(RATE2),
		masscan.SetParamExclude("255.255.255.255"),
	)

	if err != nil {
		log.Fatalf("unable to create masscan scanner: %v", err)
	}

	if err := scanner.RunAsync(); err != nil {
		fmt.Println(err)
	}

	stdout := scanner.GetStdout()

	stderr := scanner.GetStderr()
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://<username>:<password>@<host>:<port>/<database>"))
	if err != nil {
		fmt.Println(err)
	}
	scanAndInsert := func(ip string, port string) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered!")
			}
		}()
		var wg sync.WaitGroup
		db := client.Database("YourDatabase").Collection("YourCollection")
		resp, delay, err := bot.PingAndListTimeout(fmt.Sprint(ip+":"+port), TIMEOUT2)
		if err != nil {
			fmt.Printf("Ping and list server fail: %v\n", err)
		}
		var s status
		err = json.Unmarshal(resp, &s)
		if err != nil {
			fmt.Printf("Parse json response fail:%v\n", err)
		}
		s.Delay = delay
		if err == nil {
			if err != nil {
				fmt.Println(err)
			}
			serverbutmf := ServerDBbutMf{
				ServerIP:    ip + ":" + port,
				Description: stripansi.Strip(fmt.Sprint(s.Description)),
				Version:     s.Version.Name + " (" + fmt.Sprint(s.Version.Protocol) + ")",
				Players: struct {
					Max    int `bson:"max"`
					Online int `bson:"online"`
					List   []struct {
						ID   uuid.UUID `bson:"id"`
						Name string    `bson:"name"`
					} `bson:"list"`
				}{
					Max:    s.Players.Max,
					Online: s.Players.Online,
					List: []struct {
						ID   uuid.UUID `bson:"id"`
						Name string    `bson:"name"`
					}(s.Players.Sample),
				},
				FoundAt: time.Now().Format("2006-01-02"),
				Favicon: s.Favicon,
			}
			server := convertServerDB(serverbutmf)
			var existingDocument ServerDB
			filter := bson.M{"serverIP": ip +":"+ fmt.Sprint(port)}
			err := db.FindOne(context.TODO(), filter).Decode(&existingDocument)
			if err == mongo.ErrNoDocuments {
				ins, err := db.InsertOne(context.TODO(), server)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Inserted new document with ID:", ins.InsertedID)
			} else if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Document already exists, skipping insertion.")
			}
			//defer insrow.Close()
		}
		//db.Close()
		wg.Wait()
	}

	go func() {
		for stdout.Scan() {
			srs := masscan.ParseResult(stdout.Bytes())
			scannerResult = append(scannerResult, srs)
			go scanAndInsert(srs.IP, srs.Port)
		}
	}()

	go func() {
		for stderr.Scan() {
			fmt.Println(stderr.Text())
			errorBytes = append(errorBytes, stderr.Bytes()...)
		}
	}()

	if err := scanner.Wait(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("masscan result count : ", len(scannerResult))
}
