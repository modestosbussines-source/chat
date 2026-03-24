# 🚀 Omni - Guia de Implantação na VPS

## Prerequisitos

- VPS com Ubuntu 22.04+ (mínimo 2GB RAM, 2 vCPUs)
- Domínio apontando para o IP da VPS
- Docker e Docker Compose instalados
- Certificado SSL (Let's Encrypt)

---

## 1. Preparar a VPS

```bash
# Atualizar sistema
sudo apt update && sudo apt upgrade -y

# Instalar Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Instalar Docker Compose
sudo apt install docker-compose -y

# Adicionar usuário ao grupo docker
sudo usermod -aG docker $USER
newgrp docker
```

## 2. Clonar e Configurar

```bash
# Criar diretório
sudo mkdir -p /opt/omni
sudo chown $USER:$USER /opt/omni
cd /opt/omni

# Clonar repositório (ou copiar arquivos via SCP/RSYNC)
git clone https://github.com/seu-usuario/omni.git .

# Criar arquivo .env
cp .env.production .env
nano .env  # Editar com valores reais
```

## 3. Gerar Secrets Seguros

```bash
# Gerar encryption key
openssl rand -base64 32

# Gerar JWT secret
openssl rand -base64 64

# Gerar senha PostgreSQL
openssl rand -base64 32

# Gerar senha Redis
openssl rand -base64 32

# Gerar webhook verify token
openssl rand -hex 32
```

## 4. Configurar .env

```bash
# Editar .env com os valores gerados
nano /opt/omni/.env

# Valores obrigatórios:
OMNI_APP_ENCRYPTION_KEY=<gerado acima>
OMNI_DATABASE_PASSWORD=<gerado acima>
OMNI_JWT_SECRET=<gerado acima>
OMNI_DEFAULT_ADMIN_EMAIL=admin@seudominio.com
OMNI_DEFAULT_ADMIN_PASSWORD=<senha forte!>
OMNI_SERVER_ALLOWED_ORIGINS=https://omni.seudominio.com
```

## 5. Configurar Nginx (Reverse Proxy)

```bash
# Instalar Nginx
sudo apt install nginx -y

# Criar configuração
sudo nano /etc/nginx/sites-available/omni
```

```nginx
server {
    listen 80;
    server_name omni.seudominio.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name omni.seudominio.com;

    # SSL (Let's Encrypt)
    ssl_certificate /etc/letsencrypt/live/omni.seudominio.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/omni.seudominio.com/privkey.pem;
    
    # SSL Config
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;
    ssl_session_cache shared:SSL:10m;

    # Security Headers
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Proxy to Omni
    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
        proxy_read_timeout 90;
    }

    # Upload limit
    client_max_body_size 25M;
}
```

```bash
# Ativar site
sudo ln -s /etc/nginx/sites-available/omni /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

## 6. Configurar SSL (Let's Encrypt)

```bash
# Instalar Certbot
sudo apt install certbot python3-certbot-nginx -y

# Obter certificado
sudo certbot --nginx -d omni.seudominio.com

# Auto-renewal (já configurado pelo certbot)
sudo certbot renew --dry-run
```

## 7. Deploy com Docker Compose

```bash
cd /opt/omni/docker

# Build e iniciar
docker-compose up -d --build

# Verificar logs
docker-compose logs -f app

# Verificar status
docker-compose ps
```

## 8. Criar Script de Backup

```bash
sudo nano /opt/omni/backup.sh
```

```bash
#!/bin/bash
set -e

BACKUP_DIR="/opt/omni/backups"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
mkdir -p "$BACKUP_DIR"

# Backup PostgreSQL
docker-compose exec -T db pg_dump -U omni omni > "$BACKUP_DIR/omni_$TIMESTAMP.sql"
gzip "$BACKUP_DIR/omni_$TIMESTAMP.sql"

# Backup uploads
tar -czf "$BACKUP_DIR/uploads_$TIMESTAMP.tar.gz" /opt/omni/uploads

# Manter apenas últimos 7 dias
find "$BACKUP_DIR" -mtime +7 -delete

echo "Backup completed: $TIMESTAMP"
```

```bash
sudo chmod +x /opt/omni/backup.sh

# Agendar backup diário às 2h
sudo crontab -e
# Adicionar: 0 2 * * * /opt/omni/backup.sh
```

## 9. Criar Systemd Service (opcional)

```bash
sudo nano /etc/systemd/system/omni.service
```

```ini
[Unit]
Description=Omni WhatsApp Business Platform
Requires=docker.service
After=docker.service

[Service]
Type=oneshot
RemainAfterExit=yes
WorkingDirectory=/opt/omni/docker
ExecStart=/usr/bin/docker-compose up -d
ExecStop=/usr/bin/docker-compose down
TimeoutStartSec=300

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable omni
sudo systemctl start omni
```

## 10. Monitoramento

```bash
# Ver logs em tempo real
docker-compose logs -f app

# Ver uso de recursos
docker stats

# Health check
curl -f http://localhost:8080/health || echo "DOWN"
```

---

## Checklist de Segurança

- [ ] Senhas fortes configuradas no `.env`
- [ ] SSL habilitado (Nginx + PostgreSQL)
- [ ] CORS configurado com domínios específicos
- [ ] Firewall configurado (apenas ports 80, 443)
- [ ] Backup automático configurado
- [ ] Fail2Ban instalado
- [ ] Log rotation configurado
- [ ] Docker rodando como non-root

## Comandos Úteis

```bash
# Reiniciar serviços
docker-compose restart

# Atualizar
git pull
docker-compose up -d --build

# Verificar vulnerabilidades
docker scout cves omni/omni:latest

# Limpar imagens antigas
docker image prune -a
```
