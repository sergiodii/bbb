# Índice de Documentação

## Documentação Disponível

Esta pasta contém toda a documentação técnica do sistema de votação BBB:

### 📚 Documentos Principais

1. **[architecture.md](./architecture.md)** - Documentação arquitetural completa
   - Visão geral da Clean Architecture
   - Estrutura de pastas e organização
   - Padrões arquiteturais aplicados (CQRS, Strategy, Factory, etc.)
   - Pipeline de execução e agregadores

2. **[api-reference.md](./api-reference.md)** - Referência completa da API
   - Endpoints de comando e consulta
   - Exemplos de request/response
   - Códigos de erro e rate limiting
   - Exemplos de integração

3. **[development.md](./development.md)** - Guia para desenvolvedores
   - Configuração do ambiente
   - Comandos disponíveis
   - Padrões de código e boas práticas
   - Estrutura de testes

4. **[infra.md](./infra.md)** - Documentação de infraestrutura
   - Componentes (Redis, Docker, Kubernetes)
   - Automação e CI/CD
   - Monitoramento e logging

## 🚀 Início Rápido

Para começar rapidamente:

1. **Setup**: Veja [development.md - Seção 1](./development.md#1-configuração-do-ambiente)
2. **Arquitetura**: Entenda a estrutura em [architecture.md](./architecture.md)
3. **API**: Teste os endpoints em [api-reference.md](./api-reference.md)
4. **Deploy**: Configure a infraestrutura com [infra.md](./infra.md)

## 🏗️ Visão Geral da Aplicação

Esta aplicação simula o sistema de votação do paredão do Big Brother Brasil (BBB), desenvolvida em Go seguindo:

- **Clean Architecture**: Separação clara de responsabilidades
- **SOLID**: Princípios de design de software
- **CQRS**: Separação de comando e consulta
- **Clean Code**: Código limpo e maintível

### Funcionalidades Principais

- ✅ Registro de votos em tempo real
- ✅ Consulta de resultados por participante
- ✅ Análise temporal (votos por hora)
- ✅ Rate limiting para proteção contra bots
- ✅ APIs separadas (comando/consulta)
- ✅ Escalabilidade horizontal
- ✅ Testes automatizados

### Tecnologias Utilizadas

- **Backend**: Go 1.23+
- **Banco**: Redis (alta performance)
- **API**: Gin Framework
- **CLI**: Cobra
- **Containers**: Docker + Kubernetes
- **Testes**: Testify + Mocks

## 🔧 Comandos Rápidos

```bash
# Build e execução
make build && make run

# Testes
make test

# Docker
make docker-up

# APIs específicas
go run . api --port 8080          # API unificada
go run . query-api --port 8081    # Apenas consultas
go run . command-api --port 8082  # Apenas comandos
```

## 📋 Estrutura de Arquivos

```
doc/
├── api-reference.md    # 📖 Referência completa da API
├── architecture.md     # 🏗️ Documentação arquitetural
├── development.md      # 👨‍💻 Guia do desenvolvedor
└── infra.md           # 🚀 Infraestrutura e deploy
```

Para contribuir com a documentação, edite os arquivos Markdown e mantenha a consistência com os padrões estabelecidos.
