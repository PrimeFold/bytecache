package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

type Server struct{
	listeningAddress string
	cache *Cache
}

func NewServer(listenAddr string , c *Cache ) *Server{
	return &Server{
		listeningAddress: listenAddr,
		cache: c,
	}
}



func (s *Server) handleClient(conn net.Conn){
	defer conn.Close()

	scanner:= bufio.NewScanner(conn)


	for scanner.Scan(){
		line:= strings.TrimSpace(scanner.Text())
		if line == ""{
			continue
		}

		parts:= strings.Fields(line)

		command:= strings.ToUpper(parts[0])

		switch command{
		case "SET":
			if len(parts) < 3{
				fmt.Println(conn,"ERR wrong number of arguments for SET command.. ")
				continue
			}
			key:= parts[1]
			value:= []byte(parts[2])
			var ttl time.Duration

			if len(parts) >= 4{
				d,err:= time.ParseDuration(parts[3])
				if err != nil{
					fmt.Printf("Wrong ttl duration / format ")
					return
				}
				ttl = d;
			}

			s.cache.Set(key,value,ttl)
			fmt.Fprintln(conn,"OK")

		case "GET":
			if len(parts) < 2{
				fmt.Println(conn,"ERR .. Least number of arguments for GET command is 2 , Try again..")
				continue
			}
			value,ok:=s.cache.Get(parts[1])
			
			if !ok || value == nil{
				fmt.Println(conn,"(nill)")
			}
			fmt.Fprintf(conn, "%s\n", string(value))

		case "DELETE":
			if len(parts) < 2{
				fmt.Println(conn,"ERR .. Least number of arguments for DELETE command is 2 , Try again..")
				continue
			}

			s.cache.Delete(parts[1])
			fmt.Fprintln(conn, "OK")

		default:
			fmt.Fprintf(conn, "ERR unknown command '%s'\n", command)
		}




	}
}


func (s *Server) Start() error {
	clientConn,err:= net.Listen("tcp",s.listeningAddress)

	if err != nil{
		log.Fatalf("Couldn't listen on the port [%s] : %v",s.listeningAddress,err)
	}

	ctx,cancel:= context.WithCancel(context.Background());
	defer cancel()

	defer clientConn.Close()

	for{
		conn,err:=clientConn.Accept();

		if err != nil{
			log.Printf("Error accepting the connection : %v",err)
			return err
		}

		go s.cache.StartEvictionWorker(ctx,1*time.Second)
		go s.handleClient(conn)


	}
}