package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	_ "modernc.org/sqlite"
)

type Dolar struct {
	Name 	  string `json:"name"`
	Bid         string `json:"bid"`
}

func main() {
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		bid,err := GetDolar()
		if err != nil {
			http.Error(w, "Erro ao requisitar o dado", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(
			map[string]string{"bid": bid,})
				if err != nil {
					http.Error(w, "Erro ao criar o json", http.StatusInternalServerError)
					return
		}

		w.Write(response)
	})

	fmt.Println("Servidor rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
	
}

func GetDolar() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", `https://economia.awesomeapi.com.br/json/last/USD-BRL`, nil)
	if err != nil {
		if err == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout na requisição da API: %w", err)
		}
		return "", err
	}
	
	data, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer data.Body.Close()

	bodyBytes, err := io.ReadAll(data.Body)
	if err != nil {
		return "", err
	}

	var result map[string]Dolar
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return "", err
	}


	salvarNoDB(result["USDBRL"].Name, result["USDBRL"].Bid)


	return result["USDBRL"].Bid, nil
}





func salvarNoDB(name, bid string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	db, err := sql.Open("sqlite", "../dolar.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS dolar (name TEXT, bid TEXT)")
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout ao criar tabela: %w", err)
		}
		return err
	}

	_, err = db.ExecContext(ctx, "INSERT INTO dolar (name,bid) VALUES (?,?)", name, bid)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return fmt.Errorf("timeout ao inserir dados: %w", err)
		}
		return err
	}

	return nil
}

