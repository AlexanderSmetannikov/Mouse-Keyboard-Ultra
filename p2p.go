package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"p2p/ShareLogic"
	"p2p/Utility"
	"strings"
)

type Node struct {
	Connections map[string]bool
	Address     Address
}

type Address struct {
	IPv4 string
	Port string
}

type Package struct {
	To        string
	From      string
	Data      string
	MousePosX int
	MousPosY  int
}

func init() {
	if len(os.Args) != 2 {
		panic("Usage: ./program <IP Address>:<Port>")
	}
}

func main() {
	ShareLogic.DisplayCoord()
	node := NewNode(os.Args[1])
	if node == nil {
		panic("Invalid address format. Use IPv4:Port.")
	}
	node.Run(handleServer, handleClient)
}

func NewNode(address string) *Node {
	splited := strings.Split(address, ":")
	if len(splited) != 2 {
		return nil
	}
	return &Node{
		Connections: make(map[string]bool),
		Address: Address{
			IPv4: Utility.GetLocalIp().String(),
			Port: ":" + splited[1],
		},
	}
}

func (node *Node) Run(handleServer func(*Node), handleClient func(*Node)) {
	go handleServer(node)
	handleClient(node)
}

func handleConnection(node *Node, conn net.Conn) {
	defer conn.Close()
	var (
		buffer  = make([]byte, 512)
		message string
		pack    Package
	)
	for {
		length, err := conn.Read(buffer)
		if err != nil {
			break
		}
		message += string(buffer[:length])
	}
	err := json.Unmarshal([]byte(message), &pack)
	if err != nil {
		return
	}
	node.ConnectTo([]string{pack.From})
	fmt.Println(pack.Data)
}

func handleServer(node *Node) {
	listen, err := net.Listen("tcp", node.Address.Port)
	if err != nil {
		panic("Listen error: " + err.Error())
	}
	defer listen.Close()
	fmt.Println("Listening on", node.Address.IPv4+node.Address.Port)
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Connection error:", err)
			continue
		}
		go handleConnection(node, conn)
	}
}

func handleClient(node *Node) {
	for {
		message := InputString()
		splited := strings.Split(message, " ")
		switch splited[0] {
		case "/exit":
			os.Exit(0)
		case "/connect":
			node.ConnectTo(splited[1:])
		case "/network":
			node.PrintNetwork()
		default:
			node.SendMessageToAll(message)
		}
	}
}

func (node *Node) PrintNetwork() {
	fmt.Println("Connected nodes:")
	for addr := range node.Connections {
		fmt.Println("|", addr)
	}
}

func (node *Node) ConnectTo(addresses []string) {
	for _, addr := range addresses {
		node.Connections[addr] = true
	}
}

func (node *Node) SendMessageToAll(message string) {
	var new_pack = Package{
		From: node.Address.IPv4 + node.Address.Port,
		Data: message,
	}
	for addr := range node.Connections {
		new_pack.To = addr
		node.Send(&new_pack)
	}
}

func (node *Node) Send(pack *Package) {
	conn, err := net.Dial("tcp", pack.To)
	if err != nil {
		fmt.Println("Dial error:", err)
		delete(node.Connections, pack.To)
		return
	}
	defer conn.Close()
	json_pack, _ := json.Marshal(*pack)
	conn.Write(json_pack)
}

func InputString() string {
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return strings.Replace(msg, "\n", "", -1)
}
