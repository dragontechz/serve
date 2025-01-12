package main

import (
	"encoding/json"
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
		conn.Read(buf)
		// Renvoyer la réponse au client Vless
		resp := VlessResponse{
			Version:   "1.0",
			SessionID: "1234567890",
			Port:      443,
			Protocol:  "vless",
		}
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			fmt.Println(err)
			continue
		}
		_, err = conn.Write(jsonResp)
		if err != nil {
			fmt.Println(err)
			for {
				n, err := conn.Read(buf)
				if err != nil {

				}
				if n > 0 {
					data := string(buf[:n])
					log.Println(data)
					conn.Write([]byte("HTTP/2.2 200 ok"))
				}
			}
		}
	}
}

type VlessRequest struct {
	// Champs de la requête Vless
}
