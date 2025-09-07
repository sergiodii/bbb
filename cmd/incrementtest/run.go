package incrementtest

import (
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

func getRandomStringFromSlice() string {
	slice := []string{"apple", "banana", "cherry", "date", "elderberry"}

	rand.Seed(time.Now().UnixNano())

	randomIndex := rand.Intn(len(slice))

	// Retorna a string no índice aleatório.
	return slice[randomIndex]
}

func IncrementTestCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "increment-test",
		Short: "Inicia o teste de incremento chamando a api de comandos várias vezes",
	}

	c.Flags().StringP("command-api-url", "c", "localhost:8082", "URL da API de comandos")
	c.Run = func(cmd *cobra.Command, args []string) {
		maxIncrements := 1000
		commandApiUrl, _ := cmd.Flags().GetString("command-api-url")

		start := time.Now()
		fmt.Printf("\n[STARTING INCREMENT TEST] Iniciando teste de incremento com %d requisições para a API de comandos %s...\n", maxIncrements, commandApiUrl)

		wg := sync.WaitGroup{}

		for i := 0; i < maxIncrements; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				body := fmt.Sprintf(`{"participant_id": "%v"}`, getRandomStringFromSlice())
				_, err := http.Post("http://"+commandApiUrl+"/round2", "application/json", strings.NewReader(body))
				if err != nil {
					fmt.Println("Participant ID:", body, ", Error:", err)
				}
			}()
		}
		wg.Wait()

		elapsed := time.Since(start)
		fmt.Printf("\n[FINISHED INCREMENT TEST] Teste de incremento finalizado em %f segundos\n", elapsed.Seconds())

	}

	return &c
}
