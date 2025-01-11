package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
)

var key []byte = []byte("0123456789abcdef0123456789abcdef")

func main() {
	localAddr := ":9090"  // Port local sur lequel le proxy écoute
	remoteAddr := ":8888" // Adresse distante vers laquelle le trafic sera redirigé

	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		log.Fatalf("Erreur lors de l'écoute sur %s: %v", localAddr, err)
	}
	defer listener.Close()

	log.Printf("Proxy actif sur %s, redirige vers %s", localAddr, remoteAddr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Erreur lors de l'acceptation de la connexion: %v", err)
			continue
		}
		clientConn.Read(make([]byte, 1024))
		clientConn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))

		go handleConnection(clientConn, remoteAddr)
	}
}

func handleConnection(clientConn net.Conn, remoteAddr string) {
	defer clientConn.Close()

	// Connexion au serveur distant
	serverConn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Printf("Erreur lors de la connexion à %s: %v", remoteAddr, err)
		return
	}
	defer serverConn.Close()

	go recv_unencryted_send_encrypted(serverConn, clientConn) // Redirige les données du client vers le serveur
	recv_encryted_send_unencrypted(clientConn, serverConn)    // Redirige les données du serveur vers le client
}
func recv_encryted_send_unencrypted(src, dst net.Conn) error {
	// Créer un buffer pour les données
	buffer := make([]byte, 1024)

	for {
		// Lire des données depuis la connexion
		n, err := src.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break // Fin de flux
			}
			return fmt.Errorf("erreur lors de la lecture : %w", err)
		}
		data := string(buffer[:n])
		decryptedPayload, err := decrypt(data, key)
		if err != nil {
			fmt.Println("Erreur lors du déchiffrement:", err)
		}
		log.Println("Payload recu déchiffré :", string(decryptedPayload))

		// Écrire les données lues dans la même connexion
		if _, err := dst.Write([]byte(string(decryptedPayload))); err != nil {
			return fmt.Errorf("erreur lors de l'écriture : %w", err)
		}
	}

	return nil
}

func recv_unencryted_send_encrypted(src, dst net.Conn) error {
	// Créer un buffer pour les données
	buffer := make([]byte, 1024)

	for {
		// Lire des données depuis la connexion
		n, err := src.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break // Fin de flux
			}
			return fmt.Errorf("erreur lors de la lecture : %w", err)
		}
		data := string(buffer[:n])
		cryptedPayload, err := encrypt([]byte(data), key)
		if err != nil {
			fmt.Println("Erreur lors du chiffrement:", err)
		}
		log.Println("Payload a envoyer chiffré :", string(cryptedPayload))

		// Écrire les données lues dans la même connexion
		if _, err := dst.Write([]byte(cryptedPayload)); err != nil {
			return fmt.Errorf("erreur lors de l'écriture : %w", err)
		}
	}

	return nil
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

func t() {
	key := []byte("0123456789abcdef0123456789abcdef")

	payload := []byte("Ceci est un message secret!")

	ciphertextHex, err := encrypt(payload, key)
	if err != nil {
		fmt.Println("Erreur lors du chiffrement:", err)
		return
	}
	fmt.Println("Payload chiffré :", ciphertextHex)

	decryptedPayload, err := decrypt(ciphertextHex, key)
	if err != nil {
		fmt.Println("Erreur lors du déchiffrement:", err)
		return
	}
	fmt.Println("Payload déchiffré :", string(decryptedPayload))
}
