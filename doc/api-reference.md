# Documentação da API

## 1. Visão Geral

A API do sistema de votação BBB implementa padrão CQRS (Command Query Responsibility Segregation) com três formas de acesso:

- **API Unificada**: Porta 8080 (comandos + consultas)
- **API de Comando**: Porta 8082 (apenas operações de escrita)
- **API de Consulta**: Porta 8081 (apenas operações de leitura)

## 2. Endpoints de Comando (Escrita)

### 2.1. Registrar Voto

**POST** `/command/{{ roundId }}`

Registra um novo voto para um participante em um round específico.

**Parâmetros:**
- `roundId` (path): ID do round

**Request:**
```json
{
  "participant_id": "participant-123"
}
```

**Response (201 OK):**
```json
{
  "status":"vote created"
}
```

**Response (400 Bad Request):**
```json
{
  "error": "invalid request body", 
  "details": "the error text"
}
```

**Response (500 Internal Server Error):**
```json
{
  "error": "Internal server error",
}
```

**Exemplo cURL:**
```bash
curl -X POST http://localhost:8080/command/{{ roundId }} \
  -H "Content-Type: application/json" \
  -d '{
    "participant_id": "participant-123"
  }'
```

## 3. Endpoints de Consulta (Leitura)

### 3.1. Total de Votos por Round

**GET** `/query/{{ roundId }}`

Retorna o total de votos registrados para um round específico.

**Parâmetros:**
- `roundId` (path): ID do round

**Response (200 OK):**
```json
{
  "total": 15420
}
```

**Response (500 Internal Server Error):**
```json
{
  "error": "Round not found"
}
```

**Exemplo cURL:**
```bash
curl http://localhost:8081/query/round-001
```

### 3.2. Votos por Participante

**GET** `/query/{{ roundId }}/participant`

Retorna o total de votos por participante em um round específico.

**Parâmetros:**
- `roundId` (path): ID do round

**Response (200 OK):**
```json
{
  "participant1": 1500, 
  "participant2": 750,
  "participant3": 750
}

```

**Exemplo cURL:**
```bash
curl http://localhost:8081/query/round-001/participant
```

### 3.3. Votos por Hora

**GET** `/query/{{ roundId }}/hour`

Retorna a distribuição de votos por hora (timestamp unix) para um round específico.

**Parâmetros:**
- `roundId` (path): ID do round

**Response (200 OK):**
```json
{
    "451411": 1500, // Timestamp Unix da hora
    "451412": 750,
    "451413": 750,
    "488134": 1,
    "488272": 18616
}
```

**Exemplo cURL:**
```bash
curl http://localhost:8081/query/round-001/hour
```

## 4. Códigos de Status HTTP

| Código | Significado | Quando Ocorre |
|--------|-------------|---------------|
| 200 | OK | Operação realizada com sucesso |
| 201 | Created | Voto criado com sucesso |
| 400 | Bad Request | Dados de entrada inválidos |
| 500 | Internal Server Error | Erro interno do servidor |

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


## 7. Formato de Dados

### 7.1. Timestamps

Todos os timestamps são retornados como Unix timestamp (segundos desde 1970-01-01).

```json
{
  "timestamp": 1694518800
}
```

### 7.2. IDs

- **Round ID**: String alfanumérica (ex: "round-001", "paredao-final")
- **Participant ID**: String alfanumérica (ex: "participant-123", "joao-silva")

### 7.3. Encoding

- **Content-Type**: `application/json`
- **Charset**: UTF-8
- **Formato de data**: ISO 8601 quando aplicável