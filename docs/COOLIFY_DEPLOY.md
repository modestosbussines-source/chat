# 🚀 Omni + Coolify - Guia Completo de Deploy (Hostinger VPS)

Guia passo a passo para fazer deploy do Omni via Coolify na Hostinger VPS com Ubuntu.

---

## Visão Geral da Arquitetura

```
┌─────────────────────────────────────────────────────────────────┐
│                         COOLIFY                                 │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐             │
│  │   OMNI      │  │     n8n     │  │ Evolution   │             │
│  │   (App)     │◄─┤  (Workflow) │◄─┤    API      │             │
│  │  Port 8080  │  │  Port 5678  │  │  Port 8081  │             │
│  └──────┬──────┘  └─────────────┘  └─────────────┘             │
│         │                                                       │
│  ┌──────┴──────┐  ┌─────────────┐  ┌─────────────┐             │
│  │ PostgreSQL  │  │   Redis     │  │    Nginx    │             │
│  │  Port 5432  │  │  Port 6379  │  │   (Proxy)   │             │
│  └─────────────┘  └─────────────┘  └─────────────┘             │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
                            │
                            ▼
                    ┌───────────────┐
                    │   Internet    │
                    │  (HTTPS:443)  │
                    └───────────────┘
```

---

---

## 📋 Pré-requisitos Rápidos

- Hostinger VPS com Ubuntu 22.04+
- Coolify instalado (comando: `curl -fsSL https://coolify.io/install.sh | bash`)
- Domínio apontando para o IP da VPS

---

## ⚡ Resumo Rápido (Deploy Básico - Omni apenas)

### 1. Acessar Coolify
```
http://IP_DA_VPS:8000
```

### 2. Criar PostgreSQL
- Name: `omni-postgres`
- Database: `omni`, User: `omni`
- **Anote o Internal Hostname** (ex: `fpXXXX-omni-postgres`)

### 3. Criar Redis
- Name: `omni-redis`
- **Anote o Internal Hostname**

### 4. Criar Aplicação (Docker Compose)
- Repo: `modestosbussines-source/chat`
- Branch: `master`
- Docker Compose: `docker/coolify.docker-compose.yml`

### 5. Variáveis Obrigatórias
```
DATABASE_HOST=<hostname_postgres_do_coolify>
DATABASE_PASSWORD=<senha_postgres>
REDIS_HOST=<hostname_redis_do_coolify>
REDIS_PASSWORD=<senha_redis>
APP_ENCRYPTION_KEY=<openssl rand -base64 32>
JWT_SECRET=<openssl rand -base64 64>
ADMIN_EMAIL=admin@seudominio.com
ADMIN_PASSWORD=<senha_forte>
ALLOWED_ORIGINS=https://seudominio.com
TZ=America/Sao_Paulo
```

### 6. Configurar Domínio e Deploy

---

## 📋 Pré-requisitos (Completo)

### VPS Recomendada
- **CPU:** 4+ vCPUs
- **RAM:** 8GB+ (4GB Omni + 2GB n8n + 2GB Evolution)
- **Disco:** 100GB+ SSD
- **SO:** Ubuntu 22.04 LTS
- **Docker:** 24.0+
- **Docker Compose:** v2.0+

### Domínios Necessários
```
omni.seudominio.com      → Omni App
n8n.seudominio.com       → n8n Workflow
evolution.seudominio.com → Evolution API (opcional, pode ser interno)
```

---

## 🔧 Passo 1: Instalar Coolify

```bash
# Conectar na VPS via SSH
ssh root@IP_DA_VPS

# Instalar Coolify (um comando!)
curl -fsSL https://coolify.io/install.sh | bash

# Após a instalação, acessar:
# http://IP_DA_VPS:8000
```

### Configurar Domínio no Coolify
1. Acessar `http://IP_DA_VPS:8000`
2. Criar conta admin
3. Settings → Instance → Configure
4. Adicionar domínio principal: `coolify.seudominio.com`
5. Configurar DNS A record para IP da VPS

---

## 🗄️ Passo 2: Criar Banco de Dados (PostgreSQL)

### No Coolify Dashboard:

1. **Resources → New → PostgreSQL**
2. Configurações:
   ```
   Name: omni-postgres
   Version: 17
   Password: (gerar senha forte)
   ```

3. **Criar Database:**
   ```sql
   -- O Coolify cria automaticamente, mas verificar:
   Database: omni
   User: omni
   ```

4. **Anotar credenciais:**
   - Host: `omni-postgres` (nome do container)
   - Port: `5432`
   - Database: `omni`
   - User: `omni`
   - Password: (a gerada)

---

## 🔴 Passo 3: Criar Redis

### No Coolify Dashboard:

1. **Resources → New → Redis**
2. Configurações:
   ```
   Name: omni-redis
   Version: 7
   Password: (gerar senha forte)
   ```

3. **Anotar credenciais:**
   - Host: `omni-redis`
   - Port: `6379`
   - Password: (a gerada)

---

## 📱 Passo 4: Criar Evolution API

### No Coolify Dashboard:

1. **Resources → New → Application**
2. **Source:** Docker Compose (Custom)
3. **Docker Compose:**

```yaml
version: '3.8'

services:
  evolution-api:
    image: atendai/evolution-api:v2.1.0
    container_name: evolution_api
    restart: unless-stopped
    ports:
      - "8081:8080"
    environment:
      # Server
      SERVER_PORT: 8080
      SERVER_URL: https://evolution.seudominio.com
      
      # Database (PostgreSQL)
      DATABASE_ENABLED: true
      DATABASE_PROVIDER: postgresql
      DATABASE_HOST: ${DATABASE_HOST}
      DATABASE_PORT: 5432
      DATABASE_USER: ${DATABASE_USER}
      DATABASE_PASSWORD: ${DATABASE_PASSWORD}
      DATABASE_NAME: ${DATABASE_NAME}
      
      # Redis
      REDIS_ENABLED: true
      REDIS_HOST: ${REDIS_HOST}
      REDIS_PORT: 6379
      REDIS_PASSWORD: ${REDIS_PASSWORD}
      REDIS_DB: 0
      
      # Auth
      AUTHENTICATION_TYPE: apikey
      AUTHENTICATION_API_KEY: ${EVOLUTION_API_KEY}
      
      # Webhook para Omni
      WEBHOOK_GLOBAL_URL: ${OMNI_WEBHOOK_URL}
      WEBHOOK_GLOBAL_ENABLED: true
      
      # Configurações WhatsApp
      CONFIG_SESSION_PHONE_ENABLED: true
      CONFIG_SESSION_PHONE_NUMBER: ""
      
      # Logs
      LOG_LEVEL: info
      
    volumes:
      - evolution_data:/app/.media
      - evolution_sessions:/app/.sessions

volumes:
  evolution_data:
  evolution_sessions:
```

4. **Variáveis de ambiente no Coolify:**
   ```
   DATABASE_HOST=omni-postgres
   DATABASE_USER=omni
   DATABASE_PASSWORD=[senha_postgres]
   DATABASE_NAME=evolution
   
   REDIS_HOST=omni-redis
   REDIS_PASSWORD=[senha_redis]
   
   EVOLUTION_API_KEY=[gerar_key_32chars]
   OMNI_WEBHOOK_URL=https://omni.seudominio.com/api/webhooks/evolution
   ```

5. **Domínio:** `evolution.seudominio.com`

---

## ⚡ Passo 5: Criar n8n

### No Coolify Dashboard:

1. **Resources → New → Application**
2. **Source:** Docker Compose (Custom)
3. **Docker Compose:**

```yaml
version: '3.8'

services:
  n8n:
    image: n8nio/n8n:latest
    container_name: n8n
    restart: unless-stopped
    ports:
      - "5678:5678"
    environment:
      # General
      TZ: America/Sao_Paulo
      GENERIC_TIMEZONE: America/Sao_Paulo
      
      # Database (PostgreSQL)
      DB_TYPE: postgresdb
      DB_POSTGRESDB_HOST: ${DATABASE_HOST}
      DB_POSTGRESDB_PORT: 5432
      DB_POSTGRESDB_DATABASE: n8n
      DB_POSTGRESDB_USER: ${DATABASE_USER}
      DB_POSTGRESDB_PASSWORD: ${DATABASE_PASSWORD}
      
      # Redis
      QUEUE_BULL_REDIS_HOST: ${REDIS_HOST}
      QUEUE_BULL_REDIS_PORT: 6379
      QUEUE_BULL_REDIS_PASSWORD: ${REDIS_PASSWORD}
      
      # Security
      N8N_BASIC_AUTH_ACTIVE: true
      N8N_BASIC_AUTH_USER: ${N8N_USER}
      N8N_BASIC_AUTH_PASSWORD: ${N8N_PASSWORD}
      
      # Webhook
      WEBHOOK_URL: https://n8n.seudominio.com/
      
      # Execution
      EXECUTIONS_MODE: queue
      GENERIC_TIMEZONE: America/Sao_Paulo
      
      # Encryption
      N8N_ENCRYPTION_KEY: ${N8N_ENCRYPTION_KEY}
      
    volumes:
      - n8n_data:/home/node/.n8n

volumes:
  n8n_data:
```

4. **Variáveis de ambiente:**
   ```
   DATABASE_HOST=omni-postgres
   DATABASE_USER=omni
   DATABASE_PASSWORD=[senha_postgres]
   
   REDIS_HOST=omni-redis
   REDIS_PASSWORD=[senha_redis]
   
   N8N_USER=admin
   N8N_PASSWORD=[senha_forte]
   N8N_ENCRYPTION_KEY=[gerar_32chars]
   ```

5. **Domínio:** `n8n.seudominio.com`

---

## 🎯 Passo 6: Criar Omni App

### No Coolify Dashboard:

1. **Resources → New → Application**
2. **Source:** Docker Compose (Custom)
3. **Docker Compose:**

```yaml
version: '3.8'

services:
  omni:
    image: omni/omni:latest
    container_name: omni_app
    restart: unless-stopped
    ports:
      - "8080:8080"
    environment:
      # App
      OMNI_APP_NAME=Omni
      OMNI_APP_ENVIRONMENT=production
      OMNI_APP_ENCRYPTION_KEY=${APP_ENCRYPTION_KEY}
      
      # Server
      OMNI_SERVER_HOST=0.0.0.0
      OMNI_SERVER_PORT=8080
      OMNI_SERVER_ALLOWED_ORIGINS=${ALLOWED_ORIGINS}
      
      # Database
      OMNI_DATABASE_HOST=${DATABASE_HOST}
      OMNI_DATABASE_PORT=5432
      OMNI_DATABASE_USER=${DATABASE_USER}
      OMNI_DATABASE_PASSWORD=${DATABASE_PASSWORD}
      OMNI_DATABASE_NAME=omni
      OMNI_DATABASE_SSL_MODE=require
      
      # Redis
      OMNI_REDIS_HOST=${REDIS_HOST}
      OMNI_REDIS_PORT=6379
      OMNI_REDIS_PASSWORD=${REDIS_PASSWORD}
      
      # JWT
      OMNI_JWT_SECRET=${JWT_SECRET}
      OMNI_JWT_ACCESS_EXPIRY_MINS=15
      OMNI_JWT_REFRESH_EXPIRY_DAYS=7
      
      # WhatsApp (Meta) - Opcional se usar Evolution
      OMNI_WHATSAPP_WEBHOOK_VERIFY_TOKEN=${WHATSAPP_VERIFY_TOKEN}
      OMNI_WHATSAPP_API_VERSION=v21.0
      
      # Evolution API Integration
      OMNI_EVOLUTION_ENABLED=true
      OMNI_EVOLUTION_API_URL=${EVOLUTION_API_URL}
      OMNI_EVOLUTION_API_KEY=${EVOLUTION_API_KEY}
      
      # n8n Integration
      OMNI_N8N_ENABLED=true
      OMNI_N8N_WEBHOOK_URL=${N8N_WEBHOOK_URL}
      OMNI_N8N_API_KEY=${N8N_API_KEY}
      
      # Calling
      OMNI_CALLING_PUBLIC_IP=${PUBLIC_IP}
      OMNI_CALLING_RECORDING_ENABLED=true
      
      # Cookie
      OMNI_COOKIE_DOMAIN=${COOKIE_DOMAIN}
      OMNI_COOKIE_SECURE=true
      
      # Rate Limit
      OMNI_RATE_LIMIT_ENABLED=true
      OMNI_RATE_LIMIT_TRUST_PROXY=true
      
      # Default Admin
      OMNI_DEFAULT_ADMIN_EMAIL=${ADMIN_EMAIL}
      OMNI_DEFAULT_ADMIN_PASSWORD=${ADMIN_PASSWORD}
      OMNI_DEFAULT_ADMIN_FULL_NAME=Administrator
      
    volumes:
      - omni_uploads:/app/uploads
      - omni_audio:/app/audio
      
    depends_on:
      - omni-postgres
      - omni-redis
      
    command: ["./omni", "server", "-config", "config.toml", "-migrate"]

  # PostgreSQL
  omni-postgres:
    image: postgres:17-alpine
    container_name: omni_postgres
    restart: unless-stopped
    environment:
      POSTGRES_USER: ${DATABASE_USER:-omni}
      POSTGRES_PASSWORD: ${DATABASE_PASSWORD}
      POSTGRES_DB: ${DATABASE_NAME:-omni}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    command: >
      postgres
        -c ssl=on
        -c shared_preload_libraries=pg_stat_statements
        -c max_connections=200

  # Redis
  omni-redis:
    image: redis:7-alpine
    container_name: omni_redis
    restart: unless-stopped
    command: redis-server --requirepass ${REDIS_PASSWORD} --maxmemory 256mb --maxmemory-policy allkeys-lru
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
  omni_uploads:
  omni_audio:
```

4. **Variáveis de ambiente (TODAS as necessárias):**

```bash
# === DATABASE ===
DATABASE_HOST=omni-postgres
DATABASE_USER=omni
DATABASE_PASSWORD=[gerar_senha_forte]
DATABASE_NAME=omni

# === REDIS ===
REDIS_HOST=omni-redis
REDIS_PASSWORD=[gerar_senha_forte]

# === APP ===
APP_ENCRYPTION_KEY=[openssl rand -base64 32]
JWT_SECRET=[openssl rand -base64 64]
ALLOWED_ORIGINS=https://omni.seudominio.com,https://n8n.seudominio.com

# === WHATSAPP (Meta - opcional) ===
WHATSAPP_VERIFY_TOKEN=[gerar_token]

# === EVOLUTION API ===
EVOLUTION_API_URL=https://evolution.seudominio.com
EVOLUTION_API_KEY=[gerar_key_32chars]

# === n8n ===
N8N_WEBHOOK_URL=https://n8n.seudominio.com/webhook
N8N_API_KEY=[gerar_key]

# === CALLING ===
PUBLIC_IP=[ip_publico_vps]
COOKIE_DOMAIN=.seudominio.com

# === ADMIN ===
ADMIN_EMAIL=admin@seudominio.com
ADMIN_PASSWORD=[senha_forte]
```

5. **Domínio:** `omni.seudominio.com`

---

## 🔐 Passo 7: Gerar Secrets

### Na VPS ou local:

```bash
# App Encryption Key
openssl rand -base64 32
# Resultado: Ex: 4+weOBEEjlNyyqzRiAGmPrrFFZGFWQH29WAdzaLMFfU=

# JWT Secret (mais longo)
openssl rand -base64 64

# Database Password
openssl rand -base64 32

# Redis Password
openssl rand -base64 32

# API Keys
openssl rand -hex 32  # Para Evolution API Key
openssl rand -hex 16  # Para n8n API Key
```

---

## 🔄 Passo 8: Configurar Integrações

### 8.1 Conectar Omni com Evolution API

No Omni Dashboard (após deploy):
1. Settings → Integrations → Evolution API
2. Configurar:
   ```
   URL: https://evolution.seudominio.com
   API Key: [key que você gerou]
   ```
3. Salvar e testar conexão

### 8.2 Conectar Omni com n8n

No Omni Dashboard:
1. Settings → Integrations → n8n
2. Configurar:
   ```
   Webhook URL: https://n8n.seudominio.com/webhook/omni
   API Key: [key que você gerou]
   ```

No n8n Dashboard:
1. Criar workflow novo
2. Adicionar trigger "Webhook"
3. URL será: `https://n8n.seudominio.com/webhook/omni`
4. Configurar ações (enviar mensagem, criar contato, etc.)

---

## 📊 Passo 9: Monitoramento

### Coolify Dashboard
- Ver logs de cada serviço
- Monitorar CPU/Memória
- Configurar alertas

### Health Checks
```bash
# Omni
curl https://omni.seudominio.com/health

# n8n
curl https://n8n.seudominio.com/healthz

# Evolution API
curl https://evolution.seudominio.com/health
```

---

## 🔧 Troubleshooting

### Problema: Coolify não acessível
```bash
# Verificar se Coolify está rodando
docker ps | grep coolify

# Restart Coolify
docker restart coolify-proxy coolify
```

### Problema: Domínios não resolvem
```bash
# Verificar DNS
dig omni.seudominio.com

# Verificar Certbot
certbot certificates
```

### Problema: Banco de dados
```bash
# Logs PostgreSQL
docker logs omni_postgres

# Conectar no psql
docker exec -it omni_postgres psql -U omni
```

---

## 📁 Checklist Final

- [ ] Coolify instalado e configurado
- [ ] Domínios apontando para VPS
- [ ] PostgreSQL criado e funcionando
- [ ] Redis criado e funcionando
- [ ] Evolution API deployed
- [ ] n8n deployed
- [ ] Omni App deployed
- [ ] SSL configurado para todos
- [ ] Integração Omni ↔ Evolution testada
- [ ] Integração Omni ↔ n8n testada
- [ ] Admin account criado
- [ ] Backup configurado

---

## 🔗 URLs Finais

| Serviço | URL |
|---------|-----|
| **Omni App** | https://omni.seudominio.com |
| **n8n** | https://n8n.seudominio.com |
| **Evolution API** | https://evolution.seudominio.com |
| **Coolify** | http://IP:8000 |

---

## 📞 Suporte

Em caso de problemas:
1. Verificar logs no Coolify Dashboard
2. Consultar `coolify/docs` no GitHub
3. Community Discord: https://coolify.io/discord
