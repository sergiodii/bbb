package incrementtest

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

func IncrementTestCommand() *cobra.Command {
	c := cobra.Command{
		Use:   "increment-test",
		Short: "Inicia o teste de incremento chamando a api de comandos várias vezes",
	}

	c.Flags().StringP("command-api-url", "c", "localhost:8082", "URL da API de comandos")
	c.Run = func(cmd *cobra.Command, args []string) {
		maxIncrements := 1000
		commandApiUrl, _ := cmd.Flags().GetString("command-api-url")

		participants := []string{}
		for i := 0; i < maxIncrements; i++ {
			participants = append(participants, getRandomStringFromSlice())
		}

		start := time.Now()
		fmt.Printf("\n[STARTING INCREMENT TEST] Iniciando teste de incremento com %d requisições para a API de comandos %s...\n", maxIncrements, commandApiUrl)

		for _, l := range slice.TransformSliceToMultipleSlices(participants, 100) {
			wg := sync.WaitGroup{}
			for _, participant := range l {
				wg.Add(1)
				go func(p string) {
					defer wg.Done()
					body := fmt.Sprintf(`{"participant_id": "%v"}`, p)
					_, err := http.Post("http://"+commandApiUrl+"/round2", "application/json", strings.NewReader(body))
					if err != nil {
						fmt.Println("Participant ID:", body, ", Error:", err)
					}
				}(participant)

			}
			wg.Wait()
		}

		elapsed := time.Since(start)

		// Se o teste demorar mais de 1 segundo, exibe um aviso.
		if elapsed.Seconds() > 1 {
			fmt.Printf("[WARNING][FINISHED INCREMENT TEST] O teste de incremento demorou mais de 1 segundo, finalizando em: %f segundos\n", elapsed.Seconds())
			return
		}
		fmt.Printf("[FINISHED INCREMENT TEST] Teste de incremento finalizado em %f segundos\n", elapsed.Seconds())

	}

	return &c
}
