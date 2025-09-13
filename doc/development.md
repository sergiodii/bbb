# 💻 BBB Development Guide - Globo.com Voting System

## 1. Quick Start - BBB Development Environment

### 1.1. Prerequisites - Production Grade Setup

```bash
# Required tools for BBB development
✅ Go 1.23+                    # Latest Go for performance
✅ Docker 24.0+               # Containerization  
✅ Docker Compose v2          # Local orchestration
✅ Redis 7.0+                 # High-performance storage
✅ Make                       # Build automation
✅ Git 2.40+                  # Version control

# Optional but recommended
⚡ VS Code + Go extension     # IDE with debugging
🔍 Redis CLI                 # Database inspection  
📊 Prometheus + Grafana      # Local monitoring
🧪 Postman/Insomnia         # API testing
```

### 1.2. BBB Project Setup

```bash
# Clone BBB voting system
git clone https://github.com/sergiodii/bbb.git
cd bbb

# Download dependencies
go mod tidy
go mod verify

# Verify installation
go version                    # Should be 1.23+
docker --version             # Should be 24.0+
docker-compose --version     # Should be v2+

# Build CLI tools
go build -o bbb-voting .
./bbb-voting --help          # Verify CLI works
```

### 1.3. Environment Configuration - BBB Context

#### **Development Environment Variables**
```bash
# Create .env file for local development
cat > .env << EOF
REDIS_ADDR
EOF

# Load environment
source .env
export $(cat .env | grep -v '^#' | xargs)
```

#### **Production Environment Template**
```bash
# Production .env template (.env.production)
cat > .env.production << EOF

# Redis Cluster (Production)
REDIS_ADDR=redis-cluster.bbb.internal:6379

# Anti-Bot (Stricter in production)
BLOCKED_IP_RANGES=\${BLOCKED_RANGES_CONFIG}

EOF
```

## 2. Comandos Disponíveis

### 2.1. Makefile
```bash
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
    useCase := NewQueryVote(map[HandlerFuncEnum]Pipe[QueryDTO]{
        HandlerFuncGetTotalVotes: mockPipe,
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
# Executar teste de carga
go run . loadtest

```

## 7. Deploy e Distribuição

### 7.1. Build para Produção
```bash
# Build otimizado
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

# Build com Docker
docker build -t bbb-voting:latest .
```