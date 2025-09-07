# Documentação da API

## 1. Visão Geral

A API do sistema de votação BBB implementa padrão CQRS (Command Query Responsibility Segregation) com três formas de acesso:

- **API Unificada**: Porta 8080 (comandos + consultas)
- **API de Comando**: Porta 8082 (apenas operações de escrita)
- **API de Consulta**: Porta 8081 (apenas operações de leitura)

## 2. Endpoints de Comando (Escrita)

### 2.1. Registrar Voto

**POST** `/command/vote`

Registra um novo voto para um participante em um round específico.

**Request:**
```json
{
  "round_id": "round-001",
  "participant_id": "participant-123"
}
```

**Response (200 OK):**
```json
{
  "success": true,
  "message": "Vote registered successfully",
  "timestamp": 1694518800
}
```

**Response (400 Bad Request):**
```json
{
  "error": "Invalid request",
  "details": "round_id is required"
}
```

**Response (500 Internal Server Error):**
```json
{
  "error": "Internal server error",
  "details": "Failed to register vote"
}
```

**Exemplo cURL:**
```bash
curl -X POST http://localhost:8080/command/vote \
  -H "Content-Type: application/json" \
  -d '{
    "round_id": "round-001",
    "participant_id": "participant-123"
  }'
```

## 3. Endpoints de Consulta (Leitura)

### 3.1. Total de Votos por Round

**GET** `/query/total/{roundId}`

Retorna o total de votos registrados para um round específico.

**Parâmetros:**
- `roundId` (path): ID do round

**Response (200 OK):**
```json
{
  "round_id": "round-001",
  "total_votes": 15420
}
```

**Response (404 Not Found):**
```json
{
  "error": "Round not found",
  "round_id": "round-001"
}
```

**Exemplo cURL:**
```bash
curl http://localhost:8081/query/total/round-001
```

### 3.2. Votos por Participante

**GET** `/query/participant/{roundId}`

Retorna o total de votos por participante em um round específico.

**Parâmetros:**
- `roundId` (path): ID do round

**Response (200 OK):**
```json
{
  "round_id": "round-001",
  "participants": {
    "participant-123": 8500,
    "participant-456": 4200,
    "participant-789": 2720
  },
  "total_votes": 15420
}
```

**Exemplo cURL:**
```bash
curl http://localhost:8081/query/participant/round-001
```

### 3.3. Votos por Hora

**GET** `/query/hour/{roundId}`

Retorna a distribuição de votos por hora para um round específico.

**Parâmetros:**
- `roundId` (path): ID do round

**Response (200 OK):**
```json
{
  "round_id": "round-001",
  "hours": {
    "2023-09-12T20": 3200,
    "2023-09-12T21": 8500,
    "2023-09-12T22": 3720
  },
  "total_votes": 15420
}
```

**Exemplo cURL:**
```bash
curl http://localhost:8081/query/hour/round-001
```

## 4. Códigos de Status HTTP

| Código | Significado | Quando Ocorre |
|--------|-------------|---------------|
| 200 | OK | Operação realizada com sucesso |
| 400 | Bad Request | Dados de entrada inválidos |
| 404 | Not Found | Recurso não encontrado |
| 500 | Internal Server Error | Erro interno do servidor |
| 503 | Service Unavailable | Serviço temporariamente indisponível |

## 5. Rate Limiting

A API implementa rate limiting para prevenir abuso:

- **Limite**: 100 requests por minuto por IP
- **Headers de resposta**:
  - `X-RateLimit-Limit`: Limite total por janela
  - `X-RateLimit-Remaining`: Requests restantes
  - `X-RateLimit-Reset`: Timestamp do reset

**Response (429 Too Many Requests):**
```json
{
  "error": "Rate limit exceeded",
  "retry_after": 60
}
```

## 6. Autenticação e Autorização

**Versão Atual**: A API não requer autenticação (conforme requisitos do BBB).

**Versão Futura**: Considerações para implementação:
- JWT tokens para administradores
- API keys para integrações
- Rate limiting diferenciado por nível de acesso

## 7. Monitoramento e Health Check

### 7.1. Health Check

**GET** `/health`

Verifica a saúde da aplicação e suas dependências.

**Response (200 OK):**
```json
{
  "status": "healthy",
  "timestamp": "2023-09-12T21:30:00Z",
  "version": "1.0.0",
  "dependencies": {
    "redis": "connected",
    "database": "connected"
  }
}
```

### 7.2. Métricas

**GET** `/metrics`

Endpoint para coleta de métricas (formato Prometheus).

```
# HELP votes_total Total number of votes
# TYPE votes_total counter
votes_total{round_id="round-001"} 15420

# HELP api_requests_total Total number of API requests
# TYPE api_requests_total counter
api_requests_total{method="POST",endpoint="/command/vote",status="200"} 15420
```

## 8. Formato de Dados

### 8.1. Timestamps

Todos os timestamps são retornados como Unix timestamp (segundos desde 1970-01-01).

```json
{
  "timestamp": 1694518800
}
```

### 8.2. IDs

- **Round ID**: String alfanumérica (ex: "round-001", "paredao-final")
- **Participant ID**: String alfanumérica (ex: "participant-123", "joao-silva")

### 8.3. Encoding

- **Content-Type**: `application/json`
- **Charset**: UTF-8
- **Formato de data**: ISO 8601 quando aplicável

## 9. Exemplos de Integração

### 9.1. JavaScript (Frontend)

```javascript
// Registrar voto
async function registerVote(roundId, participantId) {
  try {
    const response = await fetch('/command/vote', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        round_id: roundId,
        participant_id: participantId
      })
    });
    
    const result = await response.json();
    
    if (response.ok) {
      console.log('Voto registrado:', result);
      // Atualizar interface
      updateResults(roundId);
    } else {
      console.error('Erro:', result.error);
    }
  } catch (error) {
    console.error('Erro de rede:', error);
  }
}

// Consultar resultados
async function getResults(roundId) {
  try {
    const response = await fetch(`/query/participant/${roundId}`);
    const results = await response.json();
    
    if (response.ok) {
      return results.participants;
    } else {
      throw new Error(results.error);
    }
  } catch (error) {
    console.error('Erro ao buscar resultados:', error);
    return {};
  }
}
```

### 9.2. Go (Cliente)

```go
package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type VoteRequest struct {
    RoundID       string `json:"round_id"`
    ParticipantID string `json:"participant_id"`
}

func registerVote(roundID, participantID string) error {
    vote := VoteRequest{
        RoundID:       roundID,
        ParticipantID: participantID,
    }
    
    jsonData, err := json.Marshal(vote)
    if err != nil {
        return err
    }
    
    resp, err := http.Post("http://localhost:8080/command/vote", 
                          "application/json", 
                          bytes.NewBuffer(jsonData))
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to register vote: %d", resp.StatusCode)
    }
    
    return nil
}
```

## 10. Versionamento da API

**Versão Atual**: v1

**Estratégia**: Path-based versioning (`/v1/command/vote`)

**Backward Compatibility**: Mantida por pelo menos 6 meses após nova versão.

## 11. Tratamento de Erros

### 11.1. Estrutura Padrão de Erro

```json
{
  "error": "Error message",
  "details": "Additional error details",
  "code": "ERROR_CODE",
  "timestamp": 1694518800,
  "request_id": "req-123456"
}
```

### 11.2. Códigos de Erro Customizados

| Código | Descrição |
|--------|-----------|
| `INVALID_ROUND_ID` | Round ID inválido ou não encontrado |
| `INVALID_PARTICIPANT_ID` | Participant ID inválido |
| `RATE_LIMIT_EXCEEDED` | Rate limit excedido |
| `INTERNAL_ERROR` | Erro interno do servidor |
| `SERVICE_UNAVAILABLE` | Serviço temporariamente indisponível |
