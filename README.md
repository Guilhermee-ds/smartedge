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



⚙️ Componentes Principais

1️⃣ Backend Manager ("backend/")

Gerencia a lista de backends ativos:

- Verifica a saúde de cada instância ("/health")
- Marca como “offline” backends indisponíveis
- Suporta reload dinâmico via API ("/api/reload")

---

2️⃣ Balanceadores ("balancer/")

Atualmente disponíveis:

- Round-Robin: distribuição simples e uniforme.
- EWMA: pondera os tempos de resposta dos backends e ajusta o tráfego dinamicamente.

---

3️⃣ Proxy ("proxy/")

Camada reversa que recebe as requisições dos clientes e as redireciona para o backend ativo.

- Implementa "httputil.ReverseProxy"
- Mede tempo e sucesso das requisições
- Reporta métricas via Prometheus

---

4️⃣ Métricas ("metrics/")

Expostas em "/metrics", prontas para integração com Prometheus e Grafana.

---

5️⃣ Descoberta Automática

Suporte a descoberta de backends via Consul (mockado, mas preparado para integração real).

---

🚀 Como Executar o Projeto

✅ Pré-requisitos

- Go 1.21+
- Python 3.10+
- ab (ApacheBench) opcional, para testes de carga

---

🧩 Passo a Passo

1. Clonar o repositório:

git clone https://github.com/seuusuario/smartedge.git
cd smartedge

2. Iniciar os backends simulados (Python):

python3 server1.py &
python3 server2.py &

3. Executar o Load Balancer:

go run ./cmd/main.go

4. Testar o balanceamento:

ab -n 100 -c 10 http://localhost:8080/

5. Ver métricas Prometheus:

http://localhost:8080/metrics

---

🔄 APIs Principais

Endpoint| Método| Descrição
"/api/reload"| "POST"| Recarrega manualmente os backends
"/metrics"| "GET"| Exibe métricas Prometheus
"/"| "GET"| Endpoint balanceado (proxy reverso)

---

🧪 Teste de Performance (exemplo)

ab -n 100 -c 10 http://localhost:8080/

Saída esperada:

Requests per second:    ~14000 [#/sec]
Failed requests:        0

---

🧱 Estrutura do Projeto

smartedge/
├── cmd/
│   └── main.go                # Entry point
├── internal/
│   ├── backend/               # Gerenciamento de backends
│   ├── balancer/              # Estratégias de balanceamento
│   ├── metrics/               # Prometheus integration
│   └── proxy/                 # Reverse proxy
├── server1.py                 # Backend simulado 1
├── server2.py                 # Backend simulado 2
└── README.md

---

📈 Exemplos de Log

2025/10/21 23:25:39 🚀 SmartEdge iniciado na porta 8080
2025/10/21 23:25:39 ✅ http://localhost:8081 OK
2025/10/21 23:25:39 ✅ http://localhost:8082 OK
2025/10/21 23:25:49 🔄 Backends atualizados via /api/reload

---

💡 Diferenciais Técnicos

- 🔁 Reload dinâmico de backends sem reiniciar o servidor
- ⚙️ Health Check automático com detecção de falhas
- 🧠 EWMA adaptativo com priorização inteligente
- 📊 Métricas Prometheus nativas
- 🧩 Arquitetura modular e extensível

---

🧑‍💻 Autor

Desenvolvido por Guilherme Oliveira
💼 Projeto técnico de demonstração — Engenharia de Software / Sistemas Distribuídos