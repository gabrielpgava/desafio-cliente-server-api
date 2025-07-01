package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)



func main (){

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err:=http.NewRequestWithContext(ctx, http.MethodGet, `http://localhost:8080/cotacao`, nil)
	if err != nil {
		if err == context.DeadlineExceeded {
			fmt.Println("Tempo limite excedido")
			return
		}

		fmt.Println("Erro ao criar requisição:", err)
		return
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Erro ao fazer requisição:", err)
		return 
	}
	
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Erro ao ler o corpo da resposta:", err)
		return 
	}


	var response map[string]string
	if err := json.Unmarshal(body, &response); err != nil {
		fmt.Println("Erro ao fazer o unmarshal:", err)
		return
	}
	
	fmt.Println("Cotação do Dolar: ", response["bid"])

}