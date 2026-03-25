# 🚀 Omni - Deploy via Coolify (Hostinger VPS)

Guia passo a passo para fazer deploy do Omni usando Coolify no painel da Hostinger.

---

## 📋 Pré-requisitos

- [x] Hostinger VPS com Ubuntu
- [x] Coolify instalado e acessível
- [x] Domínio apontando para o IP da VPS

---

## 🔧 Passo 1: Acessar o Coolify

1. Acesse o Coolify via navegador: `http://IP_DA_VPS:8000`
2. Faça login com sua conta

---

## 🗄️ Passo 2: Criar Recursos do Banco de Dados

### 2.1 Criar PostgreSQL

1. No Coolify, vá em **Resources** → **+ Add New**
2. Selecione **PostgreSQL**
3. Configure:
   - **Name**: `omni-postgres`
   - **Version**: `17` (ou latest)
   - **Database Name**: `omni`
   - **User**: `omni`
   - **Password**: Gere uma senha segura (anote!)
4. Clique em **Create**
5. Anote o **Internal Hostname**: será algo como `fpXXXX-omni-postgres`

### 2.2 Criar Redis

1. **+ Add New** → **Redis**
2. Configure:
   - **Name**: `omni-redis`
   - **Version**: `7`
   - **Password**: Gere uma senha segura (anote!)
3. Clique em **Create**
4. Anote o **Internal Hostname**

---

## 🐳 Passo 3: Criar o Projeto Omni

### 3.1 Criar Projeto

1. **+ Add New** → **Application**
2. Selecione **Docker Compose**

### 3.2 Configurar Repository

- **Source**: GitHub (conecte sua conta se necessário)
- **Repository**: `modestosbussines-source/chat`
- **Branch**: `master`
- **Docker Compose Location**: `docker/coolify.docker-compose.yml`

### 3.3 Configurar Build

- **Build Pack**: Docker Compose
- **Base Directory**: `/` (raiz do repo)

---

## 🔐 Passo 4: Configurar Variáveis de Ambiente

No painel da aplicação, vá em **Environment Variables** e adicione:

### Database (valores do PostgreSQL criado)
```
DATABASE_HOST=<internal_hostname_do_postgres>
DATABASE_PORT=5432
DATABASE_USER=omni
DATABASE_PASSWORD=<sua_senha_postgres>
DATABASE_NAME=omni
DATABASE_SSL_MODE=disable
```

### Redis (valores do Redis criado)
```
REDIS_HOST=<internal_hostname_do_redis>
REDIS_PORT=6379
REDIS_PASSWORD=<sua_senha_redis>
```

### App Core (Gere os secrets!)
```bash
# Execute estes comandos no terminal da VPS para gerar:
# openssl rand -base64 32  (para APP_ENCRYPTION_KEY)
# openssl rand -base64 64  (para JWT_SECRET)
# openssl rand -hex 32     (para tokens)

APP_ENCRYPTION_KEY=<gerado_com_openssl>
JWT_SECRET=<gerado_com_openssl>
```

### Server
```
ALLOWED_ORIGINS=https://seudominio.com
PUBLIC_IP=<ip_publico_da_vps>
TZ=America/Sao_Paulo
```

### Admin
```
ADMIN_EMAIL=admin@seudominio.com
ADMIN_PASSWORD=<senha_forte_12_chars_min>
ADMIN_FULL_NAME=Administrator
```

### Imagem Docker
```
OMNI_IMAGE=omni/omni:latest
```

### Evolution API (opcional - configure depois)
```
EVOLUTION_ENABLED=false
```

---

## 🌐 Passo 5: Configurar Domínio

1. Vá em **Configuration** → **Domains**
2. Adicione: `https://omni.seudominio.com`
3. O Coolify irá gerar o SSL automaticamente via Let's Encrypt

---

## 🚀 Passo 6: Deploy

1. Clique em **Deploy**
2. Aguarde o build e os containers iniciarem
3. Verifique os logs para erros

---

## ✅ Passo 7: Verificar Deploy

1. Acesse: `https://seudominio.com`
2. Login: `admin@seudominio.com` / `<sua_senha>`
3. Verifique o health check: `https://seudominio.com/health`

---

## 🔧 Geração de Secrets (Terminal VPS)

Execute na VPS para gerar todos os secrets de uma vez:

```bash
echo "=== SECRETS GERADOS ==="
echo ""
echo "APP_ENCRYPTION_KEY=$(openssl rand -base64 32)"
echo "JWT_SECRET=$(openssl rand -base64 64)"
echo "DATABASE_PASSWORD=$(openssl rand -base64 32)"
echo "REDIS_PASSWORD=$(openssl rand -base64 32)"
echo "WHATSAPP_VERIFY_TOKEN=$(openssl rand -hex 32)"
echo "EVOLUTION_API_KEY=$(openssl rand -hex 32)"
echo ""
echo "=== Copie e cole no Coolify ==="
```

---

## 🐛 Troubleshooting

### Container não inicia
- Verifique os logs no Coolify
- Confirme que DATABASE_HOST e REDIS_HOST são os internal hostnames do Coolify

### Erro de conexão com banco
- O hostname interno do Coolify é diferente do nome do container
- Use o hostname mostrado no recurso criado (ex: `fp12345-omni-postgres`)

### Frontend não carrega
- Verifique ALLOWED_ORIGINS inclui seu domínio
- O build do frontend precisa ser feito antes (incluído no Dockerfile)

### SSL não funciona
- Aguarde alguns minutos após adicionar domínio
- Verifique se o domínio aponta para o IP da VPS

---

## 📊 Comandos Úteis no Coolify

- **Restart**: Clique no botão de restart nos containers
- **Logs**: Acesse o ícone de logs em cada container
- **Terminal**: Use o terminal integrado do Coolify para debug

---

## 🔄 Atualizações

Para atualizar a aplicação:
1. No Coolify, vá na aplicação
2. Clique em **Redeploy** ou **Force Rebuild**

Ou via Git:
```bash
git pull origin master
# O Coolify detectará as mudanças e fará redeploy automático
```

---

## 📝 Notas Importantes

1. **Primeiro login**: Altere a senha do admin após o primeiro acesso
2. **Backup**: Configure backups regulares no Coolify
3. **Monitoramento**: Use o dashboard do Coolify para monitorar recursos
4. **Logs**: Verifique os logs regularmente nos primeiros dias

---

**Suporte**: Em caso de problemas, verifique os logs no Coolify ou abra uma issue no repositório.
