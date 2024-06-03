package democracy

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"sort"
)

type NodeData struct {
	Id       int
	Source   string
	Weight   string
	State    string
	Last     int
	Voters   []NodeData
	Channels []string
}

type Message struct {
	Content   string `json:"content"`
	Candidate string `json:"candidate"`
	State     string `json:"state"`
}

type Chunk struct {
	Chunk   string `json:"chunk"`
	Id      string `json:"id"`
	Total   int    `json:"total"`
	Counter int    `json:"counter"`
}

type NodeConfig struct {
	Interval      int
	Conn          *net.UDPConn
	SourceAddr    *net.UDPAddr
	Timeout       int
	Source        string
	Peers         []string
	Weight        int
	Id            string
	Channels      []string
	MaxPacketSize int
}

type Democracy struct {
	Nodes     []NodeData
	ChunkData map[string][]Chunk
	Config    *NodeConfig
}

func NewDemocracy() *Democracy {
	nc := NewConfig()
	d := &Democracy{
		Config:    nc,
		Nodes:     []NodeData{},
		ChunkData: map[string][]Chunk{},
	}
	return d
}

func (d *Democracy) SendData(event string) {
	message := Message{
		Content:   event,
		Candidate: d.Config.Source,
		State:     "leader",
	}

	encoded, _ := json.Marshal(message)
	// fmt.Println(string(encoded))
	chunks := d.GenerateChunks(encoded)
	addr, err := net.ResolveUDPAddr("udp", d.Config.Source)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		os.Exit(1)
	}

	// conn, err := net.ListenUDP("udp", addr)
	// if err != nil {
	// 	fmt.Println("Error creating UDP connection:", err)
	// 	os.Exit(1)
	// }
	// defer conn.Close()
	for _, chunk := range chunks {
		jsonData, _ := json.Marshal(chunk)

		_, err = d.Config.Conn.WriteToUDP(jsonData, addr)
		if err != nil {
			fmt.Println("Error sending data:", err)
			os.Exit(1)
		}

		fmt.Println("Sent message:", chunk)
	}
}

func (d *Democracy) Start() chan bool {
	addr, err := net.ResolveUDPAddr("udp", d.Config.Source)
	if err != nil {
		fmt.Println("Error resolving UDP address:", err)
		os.Exit(1)
	}

	d.Config.Conn, err = net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error creating UDP connection:", err)
		os.Exit(1)
	}
	ch := make(chan bool)
	go func(ch chan<- bool) {
		defer d.Config.Conn.Close()
		fmt.Println("Starting to listen on the socket")
		buf := make([]byte, 1024)
		for {
			n, addr, err := d.Config.Conn.ReadFromUDP(buf)
			if err != nil {
				fmt.Println("Error reading from UDP:", err)
				break
			}

			// Decode JSON data
			var receivedMsg Chunk
			err = json.Unmarshal(buf[:n], &receivedMsg)
			if err != nil {
				fmt.Println("Error decoding JSON:", err)
				continue
			}

			fmt.Printf("Received message from %s: %v\n", addr.String(), receivedMsg)
			if _, ok := d.ChunkData[receivedMsg.Id]; !ok {
				d.ChunkData[receivedMsg.Id] = []Chunk{}
			}
			d.ChunkData[receivedMsg.Id] = append(d.ChunkData[receivedMsg.Id], receivedMsg)

			if len(d.ChunkData[receivedMsg.Id]) == receivedMsg.Total {
				sort.Slice(d.ChunkData[receivedMsg.Id], func(i, j int) bool {
					return d.ChunkData[receivedMsg.Id][i].Id < d.ChunkData[receivedMsg.Id][j].Id
				})

				fmt.Println("Decoded chunks", d.ChunkData[receivedMsg.Id])

				content := []byte{}
				for _, ch := range d.ChunkData[receivedMsg.Id] {
					content = append(content, []byte(ch.Chunk)...)
				}

				var msg Message
				err = json.Unmarshal([]byte(content), &msg)
				if err != nil {
					fmt.Println("Error decoding Final JSON:", err)
				}

				fmt.Println("Final decoded msg ", msg)
				delete(d.ChunkData, receivedMsg.Id)
			}
		}
		ch <- true
	}(ch)

	return ch
}

// func (d *Democracy) Send(event string, id int) {
// 	fmt.Println("Sending now", event, id)
// 	addr, err := net.ResolveUDPAddr("udp", d.Source)
// 	if err != nil {
// 		fmt.Println("Error resolving UDP address:", err)
// 		os.Exit(1)
// 	}

// 	// conn, err := net.ListenUDP("udp", addr)
// 	// if err != nil {
// 	// 	fmt.Println("Error creating UDP connection:", err)
// 	// 	os.Exit(1)
// 	// }
// 	// defer conn.Close()

// 	message := Message{Content: event}
// 	jsonData, _ := json.Marshal(message)

// 	_, err = d.Conn.WriteToUDP(jsonData, addr)
// 	if err != nil {
// 		fmt.Println("Error sending data:", err)
// 		os.Exit(1)
// 	}

// 	fmt.Println("Sent message:", message.Content)
// }

func (d *Democracy) GenerateChunks(data []byte) []Chunk {
	chunks := []Chunk{}
	shortId, _ := GenerateShortID()
	if len(data) > d.Config.MaxPacketSize {
		count := Ceil(len(data), d.Config.MaxPacketSize)
		for i := 0; i < count; i++ {
			last := (i + 1) * d.Config.MaxPacketSize
			if last > len(data) {
				last = len(data)
			}
			chunks = append(chunks, Chunk{
				Id:      shortId,
				Chunk:   string(data[i*d.Config.MaxPacketSize : last]),
				Total:   count,
				Counter: i,
			})
		}
	} else {
		chunks = append(chunks, Chunk{
			Id:      shortId,
			Chunk:   string(data),
			Total:   1,
			Counter: 0,
		})
	}
	return chunks
}
