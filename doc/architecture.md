# Documentação Arquitetural - Sistema de Votação BBB

## 1. Visão Geral da Arquitetura

O sistema foi desenvolvido seguindo os princípios de **Clean Architecture**, **SOLID** e **Clean Code**, organizando o código em camadas bem definidas que promovem baixo acoplamento e alta coesão.

### 1.1. Arquitetura Visual

```
┌─────────────────────────────────────────────────────────────┐
│                     🌐 INTERFACE LAYER                        │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ │
│  │   Command API   │ │   Query API     │ │  Unified API    │ │
│  │   (Port 8082)   │ │   (Port 8081)   │ │  (Port 8080)    │ │
│  │  🛡️ Rate Limit   │ │  📊 Optimized   │ │  🔧 Development │ │
│  │  POST /{round}  │ │  GET /{round}   │ │  All Endpoints  │ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                       ⚡ USE CASE LAYER                       │
│  ┌─────────────────────────────┐ ┌─────────────────────────┐ │
│  │    🔥 COMMAND SIDE          │ │    📊 QUERY SIDE        │ │
│  │  • RegisterVote             │ │  • GetTotal            │ │
│  │  • Anti-Bot Validation      │ │  • GetByParticipant    │ │
│  │  • Pipeline: Sequential     │ │  • GetByHour           │ │
│  │    Blocking First           │ │  • Pipeline: First     │ │
│  │  • Write-Through + Async    │ │    Result with Failover│ │
│  └─────────────────────────────┘ └─────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                        🎯 DOMAIN LAYER                        │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │  Vote(roundID, participant, timestamp, clientIP)       │ │
│  │  Round(id, participants[], startTime, endTime)         │ │
│  │  Participant(id, name, house)                          │ │
│  │  RoundRepository Interface + Business Rules            │ │
│  └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│                    ⚙️ INFRASTRUCTURE LAYER                    │
│  ┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐ │
│  │  🔴 Redis Repo  │ │  💾 SQL Repo    │ │  🛡️ Middleware  │ │
│  │  (Production)   │ │ (Development)   │ │  Rate Limit     │ │
│  │  INCR Atomics   │ │ SQLite Fallback │ │  IP Blocking    │ │
│  │  Connection Pool│ │ Prepared Stmts  │ │  Proxy Support  │ │
│  └─────────────────┘ └─────────────────┘ └─────────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## 2. Camada de Domínio (Domain Layer)

### 2.1. Entidades

**`internal/domain/entity/entities.go`**
- **`Vote`**: Representa um voto individual com roundID, participantID, timestamp e IP
- **`Round`**: Representa um paredão/round de votação
- **`Participant`**: Representa um participante do programa

### 2.2. Repositórios (Interfaces)

**`internal/domain/repository/repository.go`**
- **`RoundRepository`**: Interface que define operações de persistência:
  - `VoteRegister`: Registra um voto
  - `GetTotalVotes`: Retorna total de votos de um round
  - `GetTotalForParticipant`: Retorna votos por participante
  - `GetTotalForHour`: Retorna votos por hora

## 3. Camada de Caso de Uso (Use Case Layer)

### 3.1. CQRS (Command Query Responsibility Segregation)

O sistema implementa CQRS separando operações de escrita (Command) e leitura (Query):

**Command Side (`internal/usecase/vote/command/`)**
- **`CommandVoteUseCase`**: Interface para operações de escrita
- **`command.go`**: Implementação que registra votos usando pipeline de execução
- Suporte a diferentes tipos de execução (SEQUENTIAL, CONCURRENT, etc.)

**Query Side (`internal/usecase/vote/query/`)**
- **`QueryVoteUseCase`**: Interface para operações de leitura
- **`query.go`**: Implementação que consulta dados usando pipeline de execução
- Suporte a agregação de múltiplas fontes de dados

### 3.2. Pipeline de Execução

**`extension/pipe/`**
- Sistema flexível de pipeline que suporta diferentes estratégias de execução:
  - **SEQUENTIAL**: Execução sequencial das tarefas
  - **CONCURRENT**: Execução concorrente limitada
  - **SEQUENTIAL_WITH_FIRST_RESULT**: Para quando só precisamos do primeiro resultado válido
  - **SEQUENTIAL_BLOCKING_ONLY_FIRST**: Primeira tarefa bloqueia, demais são assíncronas

### 3.3. Agregadores

**`internal/usecase/vote/aggregator/`**
- **`QueryAggregator`**: Agrega dados de múltiplos repositórios usando padrão Singleton
- **`CommandAggregator`**: Distribui comandos para múltiplos repositórios
- Permite failover automático entre diferentes fontes de dados

## 4. Camada de Interface (Interface Layer)

### 4.1. APIs REST

**`cmd/api/`**
- **API Unificada**: Porta 8080 - Endpoints de comando e query
- **API de Query**: Porta 8081 - Apenas consultas
- **API de Command**: Porta 8082 - Apenas operações de escrita

**Rotas (`cmd/api/route/vote/`)**
- **`POST /command/vote`**: Registra um voto
- **`GET /query/total/{roundId}`**: Total de votos
- **`GET /query/participant/{roundId}`**: Votos por participante
- **`GET /query/hour/{roundId}`**: Votos por hora

### 4.2. CLI (Command Line Interface)

**`cmd/root.go`**
- Utiliza **Cobra** para criar interface de linha de comando
- Comandos disponíveis:
  - `api`: Inicia API unificada
  - `query-api`: Inicia apenas API de consultas
  - `command-api`: Inicia apenas API de comandos
  - `increment-test`: Executa teste de performance

## 5. Camada de Infraestrutura (Infrastructure Layer)

### 5.1. Repositórios Concretos

**`pkg/redis/`**
- Implementação usando Redis para alta performance
- Operações atômicas usando INCR para contadores
- Chaves estruturadas por contexto (round, participante, hora)

**`pkg/localsql/`**
- Implementação alternativa usando banco SQL local
- Útil para desenvolvimento e testes

### 5.2. Extensões Utilitárias

**`extension/channel/`**
- **`SafeChannel`**: Canal thread-safe com controle de estado
- Previne envio para canais fechados

**`extension/slice/`**
- Utilitários para manipulação de slices
- Divisão de slices para processamento em lotes

## 6. Padrões Arquiteturais Aplicados

### 6.1. Dependency Inversion (SOLID)
- Camadas superiores não dependem de implementações específicas
- Uso extensivo de interfaces para desacoplamento

### 6.2. Single Responsibility (SOLID)
- Cada classe/struct tem uma única responsabilidade
- Separação clara entre Command e Query

### 6.3. Strategy Pattern
- Pipeline de execução permite diferentes estratégias
- Agregadores permitem diferentes fontes de dados

### 6.4. Factory Pattern
- Funções `New*` para criação de instâncias
- Configuração centralizada de dependências

### 6.5. Singleton Pattern
- Agregadores implementam Singleton para evitar múltiplas instâncias
- Controle de concorrência com `sync.Once`

## 7. Escalabilidade e Performance

### 7.1. Concorrência
- Uso extensivo de goroutines controladas
- Pipeline com limitação de goroutines simultâneas (máximo 10 por lote)
- Channels thread-safe para comunicação

### 7.2. Agregação de Dados
- Múltiplos repositórios para redundância
- Failover automático usando `SEQUENTIAL_WITH_FIRST_RESULT`
- Cache distribuído com Redis

### 7.3. Separação de Responsabilidades
- APIs separadas para read/write permitem escalabilidade independente
- CQRS permite otimização específica para cada tipo de operação

## 8. Testes

### 8.1. Estratégia de Testes
- **Testes Unitários**: Mocks gerados com `mockgen`
- **Testes de Integração**: Usando `miniredis` para Redis
- **Testes de Performance**: Comando `increment-test`

### 8.2. Mocks
- Geração automática de mocks para interfaces
- Testes isolados de cada camada
- Validação de comportamento sem dependências externas
