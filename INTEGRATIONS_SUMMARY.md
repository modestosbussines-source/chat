# 🔗 Omni - Resumo das Integrações

## Arquitetura Final

```
┌─────────────────────────────────────────────────────────────────────┐
│                          COOLIFY                                    │
├─────────────────────────────────────────────────────────────────────┤
│                                                                     │
│   ┌─────────────────────────────────────────────────────────────┐   │
│   │                    OMNI PLATFORM                             │   │
│   │  ┌─────────────────────────────────────────────────────┐    │   │
│   │  │                    Main App                          │    │   │
│   │  │  Port 8080                                          │    │   │
│   │  └─────────────┬───────────────────────┬───────────────┘    │   │
│   │                │                       │                    │   │
│   │                ▼                       ▼                    │   │
│   │  ┌──────────────────────┐  ┌──────────────────────┐        │   │
│   │  │   Evolution API      │  │        n8n           │        │   │
│   │  │   Port 8081          │  │   Port 5678          │        │   │
│   │  │                      │  │                      │        │   │
│   │  │  • WhatsApp Multi    │  │  • Workflows         │        │   │
│   │  │  • Webhooks          │  │  • Automações        │        │   │
│   │  │  • QR Code Connect   │  │  • Integrações       │        │   │
│   │  └──────────────────────┘  └──────────────────────┘        │   │
│   └─────────────────────────────────────────────────────────────┘   │
│                                                                     │
│   ┌─────────────────┐  ┌─────────────────┐                         │
│   │   PostgreSQL    │  │     Redis       │                         │
│   │   Port 5432     │  │   Port 6379     │                         │
│   └─────────────────┘  └─────────────────┘                         │
│                                                                     │
└─────────────────────────────────────────────────────────────────────┘
```

---

## 📁 Arquivos Criados

| Arquivo | Descrição |
|---------|-----------|
| `docs/COOLIFY_DEPLOY.md` | Guia completo de deploy no Coolify |
| `docker/coolify.docker-compose.yml` | Docker Compose otimizado para Coolify |
| `docker/.env.coolify` | Template de variáveis de ambiente |
| `internal/integrations/evolution.go` | Cliente Evolution API |
| `internal/integrations/n8n.go` | Cliente n8n |
| `internal/integrations/integrations.go` | Manager de integrações |

---

## 🔄 Fluxo de Integrações

### Evolution API → Omni
```
1. Usuário conecta WhatsApp via QR Code
2. Evolution API recebe mensagens
3. Evolution envia webhook para Omni
4. Omni processa e salva no banco
5. Omni notifica n8n (se habilitado)
```

### Omni → n8n
```
1. Evento acontece no Omni (mensagem, contato, etc.)
2. Omni envia webhook para n8n
3. n8n executa workflow configurado
4. n8n pode: enviar email, criar tarefa, notificar Slack, etc.
```

### Omni → Evolution API
```
1. Omni precisa enviar mensagem
2. Omni chama Evolution API
3. Evolution envia via WhatsApp
4. Evolution retorna status
```

---

## ⚙️ Variáveis de Ambiente Principais

### Evolution API
```bash
OMNI_EVOLUTION_ENABLED=true
OMNI_EVOLUTION_API_URL=http://evolution-api:8080
OMNI_EVOLUTION_API_KEY=[gerar_key_32chars]
```

### n8n
```bash
OMNI_N8N_ENABLED=true
OMNI_N8N_WEBHOOK_URL=http://n8n:5678/webhook
OMNI_N8N_API_KEY=[gerar_key]
OMNI_N8N_WEBHOOK_SECRET=[gerar_secret]
```

---

## 🚀 Passo a Passo Rápido Coolify

### 1. Instalar Coolify
```bash
curl -fsSL https://coolify.io/install.sh | bash
```

### 2. Criar Recursos no Painel
1. PostgreSQL (nome: omni-postgres)
2. Redis (nome: omni-redis)
3. Evolution API (Docker Compose)
4. n8n (Docker Compose)
5. Omni App (Docker Compose)

### 3. Configurar Domínios
- `omni.seudominio.com`
- `n8n.seudominio.com`
- `evolution.seudominio.com`

### 4. Adicionar Variáveis
Copiar de `docker/.env.coolify` e preencher valores

### 5. Deploy!
Coolify faz deploy automático após configuração

---

## 📱 Configuração Evolution API

### Conectar WhatsApp
1. Acessar `https://evolution.seudominio.com`
2. Criar instância via API ou painel
3. Escanear QR Code
4. WhatsApp conectado!

### API Endpoints
```
GET  /instance/{name}/qrcode    → Pegar QR Code
POST /instance/{name}/message   → Enviar mensagem
GET  /instance/{name}/connection → Status conexão
```

---

## ⚡ Configuração n8n

### Acessar n8n
1. `https://n8n.seudominio.com`
2. Login: admin / [senha configurada]

### Criar Workflow
1. New Workflow
2. Add Trigger: Webhook
3. URL: `/webhook/omni`
4. Adicionar ações (enviar email, Slack, etc.)
5. Ativar workflow

### Eventos Disponíveis
- `message.sent` - Mensagem enviada
- `message.received` - Mensagem recebida
- `contact.created` - Contato criado
- `campaign.started` - Campanha iniciada
- `campaign.completed` - Campanha finalizada

---

## 🔧 Integração no Código Go

### Usar Evolution API
```go
import "github.com/omni-platform/omni/internal/integrations"

// Enviar mensagem
err := app.Integrations.Evolution().SendMessage(ctx, "instance1", "5511999999999", "Olá!")
```

### Usar n8n
```go
// Notificar sobre evento
err := app.Integrations.N8n().TriggerEvent(ctx, "custom.event", data)
```

### Notificações Automáticas
```go
// Omni notifica automaticamente:
app.Integrations.NotifyMessageSent(ctx, msgID, contactID, content, metadata)
app.Integrations.NotifyMessageReceived(ctx, msgID, contactID, content, metadata)
app.Integrations.NotifyContactCreated(ctx, contactID, name, phone)
```

---

## 📋 Checklist Deploy

### Prerrequisitos
- [ ] VPS com 8GB+ RAM, 4+ vCPUs
- [ ] Coolify instalado
- [ ] Domínios configurados
- [ ] DNS apontando para VPS

### Deploy Omni
- [ ] PostgreSQL criado
- [ ] Redis criado
- [ ] Omni App deployed
- [ ] Health check OK

### Deploy Evolution
- [ ] Evolution API deployed
- [ ] Instância criada
- [ ] WhatsApp conectado

### Deploy n8n
- [ ] n8n deployed
- [ ] Login configurado
- [ ] Workflows criados

### Integrações
- [ ] Omni ↔ Evolution testada
- [ ] Omni ↔ n8n testada
- [ ] Webhooks funcionando

---

## 📞 URLs de Acesso

| Serviço | URL |
|---------|-----|
| Omni | https://omni.seudominio.com |
| n8n | https://n8n.seudominio.com |
| Evolution API | https://evolution.seudominio.com |
| Coolify | http://IP:8000 |

---

## 🛠️ Comandos Úteis

```bash
# Ver logs Omni
docker logs omni_app -f

# Ver logs Evolution
docker logs evolution_api -f

# Ver logs n8n
docker logs n8n -f

# Restart Omni
docker restart omni_app

# Backup banco
docker exec omni_postgres pg_dump -U omni omni > backup.sql
```
