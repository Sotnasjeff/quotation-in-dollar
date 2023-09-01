package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("Error in sending request to server: %v", err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("Error in taking response from server: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalf("Error in reading json response from server: %v", err)
	}

	var bid string
	err = json.Unmarshal(body, &bid)
	if err != nil {
		log.Fatalf("Error in parsing json response from server: %v", err)
	}

	newFile, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatalf("Error in creating file: %v", err)
	}
	defer newFile.Close()

	_, err = newFile.WriteString(fmt.Sprintf("Dolar: %s", bid))
	if err != nil {
		log.Fatalf("Error in writing in file: %v", err)
	}

}
