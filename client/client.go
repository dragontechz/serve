package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/quic-go/quic-go"
)

func main() {
	addr := "170.205.31.126:4242"
	ctx := context.Background()
	config := &tls.Config{
		InsecureSkipVerify: true,                       // À ne pas utiliser en production
		NextProtos:         []string{"h3", "http/1.1"}, // Spécifie les protocoles souhaités
	}
	quicConfig := &quic.Config{} // Tu peux personnaliser cette configuration si nécessaire

	session, err := quic.DialAddr(ctx, addr, config, quicConfig)
	if err != nil {
		log.Fatal(err)
	}

	stream, err := session.OpenStream()
	if err != nil {
		log.Fatal(err)
	}
	for {
		_, err = stream.Write([]byte("HTTP/1.1 200 ok\r\nHost:www.google.com\r\n\r\n"))
		if err != nil {
			log.Fatal(err)
		}

		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Réponse du serveur:", string(buf[:n]))
		continue
	}
}

func generateTLSConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{"quic-echo-example"},
	}
}
