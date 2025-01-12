package main

import (
	"fmt"
	"log"
	"net"
)

type VlessResponse struct {
	Version   string `json:"version"`
	SessionID string `json:"session_id"`
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"`
}

func main() {
	// Créer un serveur Vless qui écoute sur le port 443
	ln, err := net.Listen("tcp", ":7777")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()

	fmt.Println("Serveur Vless démarré sur le port 443")

	for {
		// Accepter une nouvelle connexion
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer conn.Close()

		// Lire la requête du client Vless
		buf := make([]byte, 1024)
			for {
				n, err := conn.Read(buf)
				if err != nil || n < 1 {

				}
				
					data := string(buf[:n])
					log.Println(data)
					conn.Write([]byte("HTTP/1.1 200 OK"))
				
			}
		}
	}

type VlessRequest struct {
	// Champs de la requête Vless
}
