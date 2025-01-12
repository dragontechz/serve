package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	proxyAddr := "170.205.31.126:443" // Adresse du proxy
	sshHost := "ssh.example.com:22"   // Adresse du serveur SSH
	sshUser := "vip2"                 // Nom d'utilisateur SSH
	sshPassword := "vip2"

	handleClient(proxyAddr, sshHost, sshUser, sshPassword) // Mot de passe SSH

}

func handleClient(proxyAddr, sshHost, sshUser, sshPassword string) {
	key := []byte("0123456789abcdef0123456789abcdef")
	conn, err := net.Dial("tcp", proxyAddr)
	if err != nil {
		log.Fatalf("Erreur de connexion au proxy: %v", err)
	}
	defer conn.Close()
	// Envoyer une requête CONNECT au proxy
	data, err := encrypt([]byte("HTTP/1.1 200\r\nHost: www.google.com\r\n\r\n"), key)
	conn.Write([]byte(data))
	if err != nil {
		log.Fatalf("Erreur lors de l'envoi de la requête CONNECT: %v", err)
	}
	// Lire la réponse du proxy pour vérifier si la connexion a réussi
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture de la réponse du proxy: %v", err)
	}

	response := string(buf[:n])
	log.Println("Réponse du proxy:", response)

	// Établir une session SSH sur la même connexion
	sshConfig := &ssh.ClientConfig{
		User: sshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(sshPassword),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	// Créer un client SSH à partir de la connexion existante
	sshClientConn, chans, reqs, err := ssh.NewClientConn(conn, sshHost, sshConfig)
	if err != nil {
		log.Fatalf("Erreur lors de la création du client SSH: %v", err)
	}
	ssh.NewClient(sshClientConn, chans, reqs)
	conn.Write([]byte("HTTP/1.1 200\r\nHost: www.google.com\r\n\r\n"))

	buff := make([]byte, 4096)
	n2, err := conn.Read(buff)
	if err != nil {
		log.Fatalf("Erreur lors de la lecture de la réponse du proxy: %v", err)
	}

	response2 := string(buf[:n2])
	log.Println("Réponse du proxy:", response2)

}

func pad(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padtext...)
}

func unpad(data []byte) ([]byte, error) {
	padding := data[len(data)-1]
	if int(padding) > len(data) {
		return nil, fmt.Errorf("invalid padding")
	}
	return data[:len(data)-int(padding)], nil
}

func encrypt(payload []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	payload = pad(payload)

	ciphertext := make([]byte, aes.BlockSize+len(payload))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], payload)

	return hex.EncodeToString(ciphertext), nil
}

func decrypt(ciphertextHex string, key []byte) ([]byte, error) {
	ciphertext, err := hex.DecodeString(ciphertextHex)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	return unpad(ciphertext)
}
