package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/quic-go/quic-go"
	"log"
)

func main() {
	addr := ":4242"
	listener, err := quic.ListenAddr(addr, generateTLSConfig(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Serveur QUIC en attente de connexion sur", addr)

	for {
		ctx := context.Background()
		session, err := listener.Accept(ctx)
		if err != nil {
			log.Fatal(err)
		}

		for {
			stream, err := session.AcceptStream(ctx)
			if err != nil {
				log.Println(err)
				return
			}

			go handleconn(stream)
		}
	}
}

func handleconn(stream quic.Stream) {
	for {
		buf := make([]byte, 1024)
		n, err := stream.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}

		message := string(buf[:n])
		fmt.Println("Message reçu:", message)

		if message == "salut" {
			_, err := stream.Write([]byte("HTTP/1.1 200 OK"))
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func generateTLSConfig() *tls.Config {
	// Remplace ces chemins par ceux de tes certificats
	certFile := "server.crt" // Chemin vers ton certificat
	keyFile := "server.key"  // Chemin vers ta clé privée

	// Créer une configuration TLS
	return &tls.Config{
		Certificates: []tls.Certificate{loadCertificate(certFile, keyFile)},
		NextProtos:   []string{"h3", "http/1.1"}, // Protocoles ALPN supportés
		MinVersion:   tls.VersionTLS13,
	}
}

func loadCertificate(certFile, keyFile string) tls.Certificate {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("failed to load certificate: %v", err)
	}
	return cert
}
