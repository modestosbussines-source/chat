# Omni - Changelog e Migração

## Resumo das Mudanças (2026-03-23)

### 🔒 Vulnerabilidades Corrigidas

| # | Vulnerabilidade | Status | Arquivos |
|---|-----------------|--------|----------|
| 1 | Secrets hardcoded | ✅ Corrigido | `config.toml`, `.env.production` |
| 2 | CORS permissivo em produção | ✅ Corrigido | `middleware.go`, `main.go` |
| 3 | SSL PostgreSQL desabilitado | ✅ Corrigido | `config.toml`, `docker-compose.yml` |
| 4 | Docker rodando como root | ✅ Corrigido | `Dockerfile` |
| 5 | Falta CSP headers | ✅ Corrigido | `middleware.go` |
| 6 | Senha admin fraca | ✅ Corrigido | `config.toml` |

### 🔄 Renomeação Whatomate → Omni

| Arquivo | Mudança |
|---------|---------|
| `go.mod` | Module path: `github.com/omni-platform/omni` |
| `config.toml` | App name: "Omni" |
| `internal/config/config.go` | Env prefix: `OMNI_`, default name: "Omni" |
| `docker-compose.yml` | Container names: `omni_*` |
| `Dockerfile` | Binary name: `omni` |
| `cmd/whatomate/main.go` | CLI output: "Omni" |

### 📁 Arquivos Criados/Modificados

**Novos arquivos:**
- `.env.production` - Template de variáveis de ambiente para VPS
- `DEPLOY_VPS.md` - Guia completo de implantação
- `CHANGELOG_OMNI.md` - Este arquivo

**Arquivos modificados:**
- `config.toml` - Removido secrets, SSL habilitado, Omni branding
- `docker/docker-compose.yml` - Omni containers, SSL PostgreSQL
- `docker/Dockerfile` - Non-root user, Omni binary
- `docker/.env.example` - Omni template
- `internal/config/config.go` - `OMNI_` prefix, Omni default
- `internal/middleware/middleware.go` - CORS seguro, CSP headers
- `internal/handlers/websocket.go` - CORS seguro
- `go.mod` - Module path: `github.com/omni-platform/omni`

---

## Para Fazer na VPS

### 1. Copiar Arquivos
```bash
cd /opt/omni
# Copiar todos os arquivos do projeto
```

### 2. Criar .env
```bash
cp .env.production .env
nano .env  # Preencher valores reais
```

### 3. Gerar Secrets
```bash
# Encryption Key
openssl rand -base64 32

# JWT Secret
openssl rand -base64 64

# Database Password
openssl rand -base64 32

# Redis Password
openssl rand -base64 32
```

### 4. Deploy
```bash
cd docker
docker-compose up -d --build
```

### 5. Verificar
```bash
docker-compose logs -f app
curl http://localhost:8080/health
```

---

## Checklist Pré-Deploy

- [ ] `.env` configurado com secrets fortes
- [ ] Domínio apontando para VPS
- [ ] SSL configurado (Nginx)
- [ ] Firewall configurado (ports 80, 443)
- [ ] Backup script configurado

---

## Comandos Docker

```bash
# Iniciar
docker-compose up -d

# Logs
docker-compose logs -f app

# Parar
docker-compose down

# Rebuild
docker-compose up -d --build

# Status
docker-compose ps
```
