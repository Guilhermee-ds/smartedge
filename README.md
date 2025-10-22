# 🚀 SmartEdge — Load Balancer Inteligente em Go

O **SmartEdge** é um balanceador de carga leve e eficiente desenvolvido em **Go**, projetado para distribuir requisições entre múltiplos servidores backend, monitorar a saúde de cada instância e permitir atualização dinâmica de configuração sem interrupção do serviço.

---

## 🧠 Visão Geral da Arquitetura

```
                                   ┌──────────────────────┐
                                   │   Cliente / Usuário  │
                                   └──────────┬───────────┘
                                              │
                                              ▼
                                    ┌──────────────────┐
                                    │   SmartEdge LB   │  (porta :8080)
                                    │  ─────────────── │
                                    │  ✅ Health Check  │
                                    │  🔄 /api/reload   │
                                    │  ⚙️ Round-Robin   │
                                    │  📊 /metrics      │
                                    └────────┬─────────┘
                                             │
                        ┌────────────────────┴────────────────────┐
                        │                                         │
                        ▼                                         ▼
                ┌────────────────────┐                 ┌────────────────────┐
                │ Backend 1 (Python) │                 │ Backend 2 (Python) │
                │ Porta: :8081       │                 │ Porta: :8082       │
                │ /health            │                 │ /health            │
                └────────────────────┘                 └────────────────────┘
```


O balanceador recebe requisições HTTP na porta `8080` e as distribui entre os backends ativos usando o algoritmo **Round-Robin**.  
Cada backend é monitorado por requisições periódicas em `/health` para garantir alta disponibilidade.

---

## ⚙️ Fluxo de Funcionamento

1. **Inicialização**
   - O SmartEdge inicia e carrega os backends configurados.
   - Realiza checagem de saúde em cada servidor.
   - Ativa somente os backends disponíveis.

2. **Distribuição de Requisições**
   - Cada requisição recebida é enviada ao próximo backend ativo em ordem circular.
   - Caso um backend falhe, é removido temporariamente da rotação.

3. **Reload Dinâmico**
   - Endpoint `/api/reload` permite atualizar a lista de backends **sem reiniciar o servidor Go**.
   - Persistência garante que os backends ativos sejam mantidos entre reloads.

4. **Monitoramento**
   - Logs detalhados informam o estado dos backends e requisições.
   - Mensagens como `✅ OK` ou `❌ Offline` indicam status em tempo real.

---

## 🧩 Estrutura do Projeto

```
smartedge/
├── cmd/
│   └── main.go          # Entrada principal da aplicação Go
├── internal/
│   ├── backend/         # Gerenciamento e descoberta de backends
│   ├── balancer/        # Algoritmos de balanceamento (Round-Robin, EWMA)
│   ├── proxy/           # Proxy reverso para rotear requisições
│   └── metrics/         # Exposição de métricas Prometheus
├── server1.py           # Backend de exemplo 1 (porta 8081)
├── server2.py           # Backend de exemplo 2 (porta 8082)
├── backends.json        # Persistência dos backends ativos (opcional)
└── README.md
```

---

## 🧪 Executando o Projeto

### 1. Clonar o repositório
```bash
git clone https://github.com/seu-usuario/smartedge.git
cd smartedge
```

### 2. Iniciar os servidores backend
Em terminais separados (ou em background):
```bash
python3 server1.py &
python3 server2.py &
```

Esses servidores de exemplo rodam nas portas `8081` e `8082` e respondem em `/health`.

### 3. Iniciar o SmartEdge
```bash
go run ./cmd/main.go
```

Saída esperada:
```
🚀 SmartEdge iniciado na porta 8080
🔄 Backends atualizados via /api/reload
✅ http://localhost:8081 OK
✅ http://localhost:8082 OK
```

---

## 🧠 Endpoints Principais

| Método | Endpoint        | Descrição |
|--------|-----------------|------------|
| `GET`  | `/`             | Redireciona requisição para o backend ativo |
| `GET`  | `/health`       | Health check interno do balanceador |
| `POST` | `/api/reload`   | Atualiza lista de backends dinamicamente |
| `GET`  | `/api/status`   | Retorna status dos backends e distribuição |

---

## ⚡ Teste de Carga

Use o **Apache Benchmark** para validar o balanceamento:

```bash
ab -n 100 -c 10 http://localhost:8080/
```

Exemplo de saída:
```
Concurrency Level:      10
Time taken for tests:   1.245 seconds
Complete requests:      100
Requests per second:    80.32 [#/sec]
```

Durante o teste, o SmartEdge registra a distribuição das requisições entre os backends.

---

## 💾 Persistência de Backends (opcional)

Para manter os backends mesmo após reiniciar o balanceador, basta salvar em `backends.json`:

```json
{
  "backends": [
    "http://localhost:8081",
    "http://localhost:8082"
  ]
}
```

O SmartEdge lê esse arquivo na inicialização e restaura automaticamente as instâncias conhecidas.

---

## 🧱 Tecnologias Utilizadas

| Tecnologia | Uso |
|-------------|-----|
| **Go** | Backend e balanceamento de carga |
| **Python** | Servidores de teste simulando microserviços |
| **net/http** | Servidor HTTP em Go |
| **encoding/json** | Persistência e reload dinâmico |
| **log** | Logging estruturado |
| **Apache Benchmark (ab)** | Testes de carga |

---

## 🔍 Próximos Passos

- 🧱 Circuit Breaker → detectar falhas frequentes e desativar backends temporariamente  
- 📊 Dashboard Grafana → visualizar métricas de latência e falhas  
- 🔁 Hot Reload avançado → atualizar backends e estratégias sem reiniciar  
- 🌍 GeoAffinity → priorizar servidores próximos geograficamente  
- 💾 Descoberta automática via Consul / etcd

---

## 👨‍💻 Autor

**Guilherme Oliveira**  
Desenvolvedor Backend • Foco em Go, Docker, Python e Arquitetura de Sistemas Distribuídos

---

## 🏁 Licença

Distribuído sob a licença MIT. Consulte o arquivo `LICENSE` para mais detalhes.
