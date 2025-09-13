package loadtest

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/sergiodii/bbb/extension/slice"
	"github.com/spf13/cobra"
)

func getRandomStringFromSlice() string {
	slice := []string{"apple", "banana", "cherry", "date", "elderberry"}

	rand.Seed(time.Now().UnixNano())

	randomIndex := rand.Intn(len(slice))

	// Retorna a string no índice aleatório.
	return slice[randomIndex]
}

func LoadTestCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "loadtest",
		Short: "Inicia o teste de carga chamando a api de inserção de dados N vezes",
	}

	c.Flags().StringP("url", "c", "localhost:8082", "URL da API de Inserção de votos")
	c.Flags().IntP("requests", "r", 1000, "Número de requisições a serem feitas")
	c.Flags().IntP("concurrent", "n", 100, "Número de requisições concorrentes")
	c.Flags().String("round-id", "round1", "ID da rodada para a qual os votos serão enviados")

	c.Run = func(cmd *cobra.Command, args []string) {

		// Get flags
		maxIncrements, _ := cmd.Flags().GetInt("requests")
		url, _ := cmd.Flags().GetString("url")
		roundID, _ := cmd.Flags().GetString("round-id")
		concurrent, _ := cmd.Flags().GetInt("concurrent")

		// Generate random participant IDs
		participants := []string{}
		for i := 0; i < maxIncrements; i++ {
			participants = append(participants, getRandomStringFromSlice())
		}

		// Start load test time tracking
		start := time.Now()
		fmt.Printf("\n 🏁 [STARTING LOAD TEST] Iniciando teste de carga com %d requisições para a API de inserção de votos %s...\n", maxIncrements, url)

		// Client HTTP otimizado
		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		for _, l := range slice.TransformSliceToMultipleSlices(participants, concurrent) {
			wg := sync.WaitGroup{}
			finalUrl := "http://" + url + "/" + roundID
			for _, participant := range l {
				wg.Add(1)
				go func(p string) {
					defer wg.Done()
					body := fmt.Sprintf(`{"participant_id": "%v"}`, p)
					resp, err := client.Post(finalUrl, "application/json", strings.NewReader(body))
					if err != nil {
						fmt.Println("Participant ID:", body, ", Error:", err)
					} else if resp.StatusCode >= 400 {
						fmt.Printf("Participant ID: %s, HTTP Status: %d\n", body, resp.StatusCode)
						resp.Body.Close()
					} else {
						resp.Body.Close()
					}
				}(participant)

			}
			wg.Wait()
			// Small delay between batches to avoid overwhelming the server
			time.Sleep(10 * time.Millisecond)
		}

		// End load test time tracking
		elapsed := time.Since(start)

		// Se o teste demorar mais de 1 segundo, exibe um aviso.
		if elapsed.Seconds() > 1 {
			fmt.Printf("\n⚠️ - [FINISHED LOAD TEST] O teste de carga demorou mais de 1 segundo:\n\n- ⏰ Tempo: %f segundos\n- 📈 Throughput: %f req/s\n\n", elapsed.Seconds(), calculateThroughput(maxIncrements, elapsed))
			return
		}
		fmt.Printf("\n✅ [FINISHED LOAD TEST] Teste de carga finalizado:\n\n- ⏰ Tempo: %f segundos\n- 📈 Throughput: %f req/s\n\n", elapsed.Seconds(), calculateThroughput(maxIncrements, elapsed))

	}

	return &c
}

func calculateThroughput(requests int, duration time.Duration) float64 {
	seconds := duration.Seconds()
	if seconds == 0 {
		return 0
	}
	return (float64(requests) / seconds)
}
