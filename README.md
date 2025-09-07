# Sistema de Votação BBB

Sistema de votação do paredão do BBB desenvolvido em Go, seguindo **Clean Architecture**, **SOLID**, **Clean Code** e implementando um padrão análogo ao **CQRS** (Command Query Responsibility Segregation).

## 🏗️ Arquitetura

O sistema implementa CQRS com APIs separadas para escalabilidade independente:
- **Command API** (Porta 8082): Operações de escrita (registrar votos)
- **Query API** (Porta 8081): Operações de leitura (consultar resultados)
- **API Unificada** (Porta 8080): Comandos + consultas em uma única API

### Estrutura do Projeto
```
github.com/sergiodii/bbb/
├── cmd/                    # Comandos CLI e APIs
├── internal/domain/        # Entidades e interfaces de domínio
├── internal/usecase/       # Casos de uso (CQRS: command/query)
├── pkg/                    # Repositórios (Redis, SQL)
├── extension/              # Pipeline, SafeChannel, utilitários
├── chart/                  # Manifests Kubernetes
└── doc/                    # Documentação completa
```

## 🚀 Como Rodar

### Opção 1: Docker Compose (Recomendado)
```bash
# Subir todas as APIs + Redis
docker-compose up

# APIs disponíveis:
# - Command API: http://localhost:8082
# - Query API:   http://localhost:8081
# - Redis:       localhost:6379
```

### Opção 2: Make (Desenvolvimento)
```bash
# Build e executar
make build && make run

# Apenas testes
make test

# Docker completo
make docker-up
```

### Opção 3: Go CLI (Desenvolvimento)
```bash
# API unificada (comando + consulta)
go run . api --port 8080

# Apenas API de consultas
go run . query-api --port 8081

# Apenas API de comandos
go run . command-api --port 8082
```

## 📡 Endpoints da API

### Command API (Escrita) - Porta 8082

#### Registrar Voto
```http
POST /command/{round_id}
Content-Type: application/json

{
  "participant_id": "participante-123"
}
```

**Resposta (201):**
```json
{
  "status": "vote created"
}
```

### Query API (Leitura) - Porta 8081

#### 1. Total de Votos por Round
```http
GET /query/{round_id}
```

**Resposta:**
```json
{
  "total": 15420
}
```

#### 2. Votos por Participante
```http
GET /query/{round_id}/participant
```

**Resposta:**
```json
{
  "participant-123": 8500,
  "participant-456": 4200,
  "participant-789": 2720
}
```

#### 3. Votos por Hora
```http
GET /query/{round_id}/hour
```

**Resposta:**
```json
{
  "1694518800": 3200,
  "1694522400": 8500,
  "1694526000": 3720
}
```

### Exemplos de Uso

#### Registrar Voto
```bash
curl -X POST http://localhost:8082/command/round1 \
  -H "Content-Type: application/json" \
  -d '{"participant_id": "alice"}'
```

#### Consultar Resultados
```bash
# Total de votos
curl http://localhost:8081/query/round1

# Votos por participante
curl http://localhost:8081/query/round1/participant

# Votos por hora
curl http://localhost:8081/query/round1/hour
```

## 🛠️ Comandos CLI

O sistema possui uma interface CLI robusta usando **Cobra**:

### 1. `api` - API Unificada
```bash
go run . api --port 8080
```
- **Função**: Inicia uma API que combina comandos e consultas
- **Uso**: Ambiente de desenvolvimento ou quando não há necessidade de separação
- **Portas**: Personalizável via flag `--port`

### 2. `query-api` - API de Consultas
```bash
go run . query-api --port 8081
```
- **Função**: Inicia apenas a API de leitura/consultas
- **Uso**: Escalabilidade independente, cache dedicado
- **Endpoints**: Apenas rotas de consulta (`GET`)

### 3. `command-api` - API de Comandos
```bash
go run . command-api --port 8082
```
- **Função**: Inicia apenas a API de escrita/comandos
- **Uso**: Escalabilidade independente, otimizações de escrita
- **Endpoints**: Apenas rotas de comando (`POST`)

### 4. `increment-test` - Teste de Performance
```bash
go run . increment-test --command-api-url localhost:8082
```
- **Função**: Executa teste de carga com 1000 requisições concorrentes
- **Uso**: Validar performance e capacidade do sistema
- **Configurável**: URL da API de comandos via flag

### Exemplos Avançados
```bash
# API unificada em porta customizada
go run . api --port 9000

# Teste de performance em API remota
go run . increment-test --command-api-url production-api:8082

# APIs separadas para microserviços
go run . command-api --port 8082 &  # Background
go run . query-api --port 8081 &    # Background
```

## 🧪 Testes

### Executar Testes
```bash
# Todos os testes
go test ./...

# Testes com verbose
go test -v ./...

# Testes com coverage
go test -cover ./...

# Teste de performance
go run . increment-test
```

### Tipos de Teste
- **Unitários**: Mocks gerados com `mockgen`
- **Integração**: Usando `miniredis` para Redis
- **Performance**: 1000 requisições concorrentes
- **API**: Testes de endpoints

## 🔧 Desenvolvimento

### Pré-requisitos
- Go 1.23+
- Docker e Docker Compose
- Make (opcional)
- Redis (para desenvolvimento local)

### Estrutura de Desenvolvimento
```bash
# Instalar dependências
go mod tidy

# Gerar mocks
go generate ./...

# Build
make build

# Limpar cache
go clean -cache -testcache -modcache
```

## 📊 Performance

### Resultados de Teste
- **Throughput**: 1000+ requisições/segundo
- **Concorrência**: Goroutines controladas (máx 10 por lote)
- **Latência**: Sub-100ms para operações Redis
- **Escalabilidade**: APIs independentes permitem escala horizontal

### Estratégias de Execução
- **SEQUENTIAL**: Para consistência
- **CONCURRENT**: Para performance
- **SEQUENTIAL_WITH_FIRST_RESULT**: Para failover
- **SEQUENTIAL_BLOCKING_ONLY_FIRST**: Para replicação

## 📚 Documentação

Documentação completa disponível em `/doc/`:
- **[architecture.md](./doc/architecture.md)**: Arquitetura detalhada
- **[api-reference.md](./doc/api-reference.md)**: Referência completa da API
- **[development.md](./doc/development.md)**: Guia do desenvolvedor
- **[infra.md](./doc/infra.md)**: Deploy e infraestrutura

## 🎯 Features

- ✅ **CQRS**: Separação de comando e consulta
- ✅ **Pipeline**: Sistema flexível de execução
- ✅ **Agregadores**: Múltiplas fontes de dados
- ✅ **Failover**: Recuperação automática
- ✅ **Escalabilidade**: APIs independentes
- ✅ **Performance**: Teste de 1000 req/s
- ✅ **Docker**: Containerização completa
- ✅ **Kubernetes**: Deploy em produção
- ✅ **CLI**: Interface de linha de comando
- ✅ **Testes**: Cobertura >80%

## 📋 Histórico e Decisões

Para detalhes sobre decisões arquiteturais, implementação e histórico completo, consulte:
- **[HISTORY.md](./HISTORY.md)**: Histórico detalhado do desenvolvimento
