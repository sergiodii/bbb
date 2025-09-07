# Guia de Desenvolvimento

## 1. Configuração do Ambiente

### 1.1. Pré-requisitos
- Go 1.23+ instalado
- Docker e Docker Compose
- Make (para automação)
- Redis (para desenvolvimento local - opcional, pois pode usar Docker)

### 1.2. Clonando e Configurando
```bash
git clone <repository-url>
cd github.com/sergiodii/bbb
go mod tidy
```

### 1.3. Variáveis de Ambiente
```bash
# Configuração do Redis
REDIS_ADDR=localhost:6379

# Porta da API
PORT=8080
```

## 2. Comandos Disponíveis

### 2.1. Makefile
```bash
# Instalar dependências Cobra
make install-cobra

# Build da aplicação
make build

# Executar aplicação
make run

# Executar testes
make test

# Subir ambiente Docker
make docker-up

# Parar ambiente Docker
make docker-down

# Limpar binários
make clean
```

### 2.2. Comandos CLI da Aplicação
```bash
# API unificada (comando + query)
go run . api --port 8080

# Apenas API de consultas
go run . query-api --port 8081

# Apenas API de comandos  
go run . command-api --port 8082

# Teste de performance
go run . increment-test
```

## 3. Estrutura de Desenvolvimento

### 3.1. Adicionando Novas Features

**1. Domínio (Domain)**
- Adicione novas entidades em `internal/domain/entity/`
- Defina interfaces de repositório em `internal/domain/repository/`

**2. Caso de Uso (UseCase)**
- Para operações de escrita: `internal/usecase/vote/command/`
- Para operações de leitura: `internal/usecase/vote/query/`
- Use o padrão de pipeline para composição de operações

**3. Interface (API)**
- Adicione rotas em `cmd/api/route/`
- Registre rotas nos arquivos `query_api.go` ou `command_api.go`

**4. Infraestrutura**
- Implemente repositórios concretos em `pkg/`
- Use interfaces para manter desacoplamento

### 3.2. Exemplo: Adicionando Nova Consulta

**1. Interface do Use Case**
```go
// internal/usecase/vote/query/interface.go
type QueryVoteUseCase interface {
    // ... métodos existentes
    GetVotesByTime(ctx context.Context, roundID string, startTime, endTime int64) ([]Vote, error)
}
```

**2. Implementação**
```go
// internal/usecase/vote/query/query.go
func (q *queryVote) GetVotesByTime(ctx context.Context, roundID string, startTime, endTime int64) ([]Vote, error) {
    // Implementação usando pipeline
}
```

**3. Rota da API**
```go
// cmd/api/route/vote/query.go
router.GET("/votes-by-time/:roundId", func(c *gin.Context) {
    // Handler implementation
})
```

## 4. Padrões de Código

### 4.1. Nomenclatura
- **Interfaces**: Sufixo com propósito (ex: `Repository`, `UseCase`)
- **Implementações**: Nome descritivo sem sufixos genéricos
- **DTOs**: Sufixo `DTO` para objetos de transferência
- **Mocks**: Sufixo `Mock` em package separado

### 4.2. Estrutura de Funções
```go
// Construtores sempre começam com New
func NewQueryVote(pipes map[HandlerFuncEnum]OrderedExecutionPipeDTO) QueryVoteUseCase {
    return &queryVote{
        orderedExecutionPipes: pipes,
    }
}

// Métodos públicos com documentação
// GetTotalVotes retorna o total de votos para um round específico
func (q *queryVote) GetTotalVotes(ctx context.Context, roundID string) (int, error) {
    // implementação
}
```

### 4.3. Tratamento de Erros
```go
// Use erros específicos do domínio
var ErrRoundNotFound = errors.New("round not found")

// Para pipeline, use pipe.ONF quando objeto não for encontrado
if len(votes) == 0 {
    return dto, pipe.ONF
}

// Propague erros com contexto
if err != nil {
    return dto, fmt.Errorf("failed to get votes for round %s: %w", roundID, err)
}
```

## 5. Testes

### 5.1. Executando Testes
```bash
# Todos os testes
go test ./...

# Testes específicos com verbose
go test -v ./internal/usecase/vote/query/

# Testes de integração
go test -v ./pkg/redis/

# Com coverage
go test -cover ./...
```

### 5.2. Estrutura de Testes

**Testes Unitários**
```go
func TestQueryVote_GetTotalVotes(t *testing.T) {
    // Arrange
    mockPipe := mock.NewPipeMock[QueryDTO]()
    useCase := NewQueryVote(map[HandlerFuncEnum]OrderedExecutionPipeDTO{
        HandlerFuncGetTotalVotes: {
            ExecutionType: "SEQUENTIAL",
            Pipe: mockPipe,
        },
    })
    
    // Act & Assert
    // ...
}
```

**Testes de Integração**
```go
func TestRedisRepository_Integration(t *testing.T) {
    // Setup miniredis
    s, err := miniredis.Run()
    require.NoError(t, err)
    defer s.Close()
    
    repo := NewRedisRepository(s.Addr())
    // ...
}
```

### 5.3. Mocks
```bash
# Gerar mocks (já configurado com go:generate)
go generate ./...

# Ou específico
mockgen -source=internal/usecase/vote/query/interface.go -destination=internal/usecase/vote/query/mock/query_mock.go
```

## 6. Performance e Monitoramento

### 6.1. Teste de Performance
```bash
# Executar teste de incremento
go run . increment-test

# Profile de CPU
go test -cpuprofile=cpu.prof -bench=.

# Profile de memória  
go test -memprofile=mem.prof -bench=.
```

### 6.2. Logs e Debugging
```go
// Use contexto para rastreamento
ctx = context.WithValue(ctx, "requestID", requestID)

// Logs estruturados (considere usar logrus ou zap)
log.Printf("[%s] Processing vote for round: %s", requestID, roundID)
```

## 7. Deploy e Distribuição

### 7.1. Build para Produção
```bash
# Build otimizado
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Build com Docker
docker build -t bbb-voting:latest .
```

### 7.2. Kubernetes
```bash
# Deploy
kubectl apply -f chart/

# Verificar status
kubectl get pods -l app=bbb-voting

# Logs
kubectl logs -f deployment/bbb-voting
```

## 8. Boas Práticas

### 8.1. Git
- Commits pequenos e focados
- Mensagens descritivas
- Branch por feature/bugfix
- Pull requests para code review

### 8.2. Código
- Sempre execute `go fmt` antes de commit
- Use `go vet` para verificar problemas
- Mantenha cobertura de testes > 80%
- Documente APIs públicas

### 8.3. Performance
- Use pipeline com moderação (evite over-engineering)
- Monitor uso de goroutines
- Profile aplicação em ambiente similar à produção
- Cache dados quando apropriado (Redis)

### 8.4. Segurança
- Validação de entrada em todas as APIs
- Rate limiting implementado
- Logs não devem expor dados sensíveis
- Use HTTPS em produção
