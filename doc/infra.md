# Documentação de Infraestrutura

## 1. Visão Geral

A infraestrutura da aplicação foi projetada para ser escalável, resiliente e automatizada, utilizando ferramentas open-source e práticas de DevOps.

## 2. Componentes

### 2.1. Redis

- **Função**: Armazenamento de votos e contadores.
- **Motivo da Escolha**: Redis é um banco de dados em memória extremamente rápido, ideal para operações de escrita intensiva e contagem, como no cenário de votação do BBB. Sua estrutura de dados `INCR` permite contagens atômicas e de alta performance.
- **Deploy**: A aplicação utiliza uma instância do Redis, que pode ser implantada como um serviço separado no Kubernetes ou via `docker-compose` em desenvolvimento.

### 2.2. Docker

- **Função**: Containerização da aplicação.
- **`Dockerfile`**: O `Dockerfile` utiliza um build multi-stage:
  1. **`builder`**: Compila a aplicação Go em um ambiente limpo.
  2. **`final`**: Copia apenas o binário compilado e os assets necessários para uma imagem final leve, baseada em `alpine`.
- **`docker-compose.yml`**: Facilita o ambiente de desenvolvimento, subindo a aplicação e o Redis com um único comando.

### 2.3. Kubernetes

- **Função**: Orquestração de contêineres em produção.
- **`chart/`**: O diretório contém os manifestos YAML para implantar a aplicação e o Redis no Kubernetes.
  - **`bbb-voting.yaml`**: Define o `Deployment` e o `Service` da aplicação.
  - **`redis.yaml`**: Define o `Deployment` e o `Service` do Redis.
- **Escalabilidade**: O `Deployment` da aplicação pode ser facilmente escalado para múltiplas réplicas para lidar com o aumento da carga.

## 3. Automação e CI/CD

- **`Makefile`**: Automatiza tarefas comuns de desenvolvimento, como build, teste e deploy local.
- **CI/CD (Sugestão)**: Para um ambiente de produção, recomenda-se a configuração de um pipeline de CI/CD (ex: GitHub Actions, GitLab CI) que automatize:
  1. **Build e Teste**: A cada push para o repositório.
  2. **Build da Imagem Docker**: E push para um registro de contêineres (ex: Docker Hub, GCR).
  3. **Deploy no Kubernetes**: Atualização do `Deployment` com a nova imagem.

## 4. Monitoramento e Logging (Sugestão)

Para garantir a saúde da aplicação em produção, recomenda-se a integração com ferramentas de monitoramento e logging:

- **Prometheus**: Para coletar métricas da aplicação (ex: número de votos, latência da API).
- **Grafana**: Para visualizar as métricas coletadas pelo Prometheus em dashboards.
- **Fluentd/Loki**: Para agregação e consulta de logs da aplicação.
