# Relatório de Análise de Segurança - Whatomate

**Data:** 2026-03-23
**Projeto:** Whatomate (WhatsApp Business Platform)
**Analista:** Security Audit (Automated)

---

## 🔴 VULNERABILIDADES CRÍTICAS

### 1. Secrets Hardcoded em Arquivos de Configuração (CRÍTICO)
**Arquivos afetados:**
- `config.toml` (linhas 5, 18, 33, 63)
- `docker/.env.example` (linha 6)

**Secrets expostos:**
```
encryption_key = "4+weOBEEjlNyyqzRiAGmPrrFFZGFWQH29WAdzaLMFfU=" (config.toml:5)
database password = "vyDpsMggRVbKFEWP6Dz/dg==" (config.toml:18)
JWT secret = "JkejFO/3kAbJOJASJipGaIpzzNtkw72BY7Os749fcIAvi3ewlLo9P0T2jpx3X01ZLN4kbxHygl44kW0Q2Ncbp3A==" (config.toml:33)
default_admin password = "WDzR3aztXzCLkdG7" (config.toml:63)
POSTGRES_PASSWORD = "whatomate" (docker/.env.example:6)
```

**Impacto:** Qualquer pessoa com acesso ao repositório pode acessar bancos de dados, assumir contas de admin, e descriptografar dados sensíveis.

**Recomendação:**
- Remover todos os secrets dos arquivos versionados
- Usar variáveis de ambiente ou vault de secrets (HashiCorp Vault, AWS Secrets Manager)
- Adicionar `config.toml` ao `.gitignore`

---

### 2. Senha de Administrador Fraca (ALTO)
**Local:** `config.toml:63` e `config.example.toml:117`

```toml
[default_admin]
email = "admin@seudominio.com"
password = "WDzR3aztXzCLkdG7"  # OU "admin" no exemplo
```

**Impacto:** A senha "admin" no config.example.toml é extremamente fraca e previsível. A senha no config.toml虽然更强，但依然硬编码。

**Recomendação:**
- Usar gerador de senhas forte (32+ caracteres aleatórios)
- Exigir troca de senha no primeiro login
- Remover default_admin do repositório

---

### 3. CORS Permissivo em Produção (ALTO)
**Local:** `internal/middleware/middleware.go:66-68`

```go
func IsOriginAllowed(origin string, allowedOrigins map[string]bool) bool {
    if len(allowedOrigins) == 0 {
        return true // No whitelist configured = allow all (development)
    }
```

```toml
# config.example.toml:19
allowed_origins = ""  # Empty = allow all (dev only)
```

**Impacto:** Se `allowed_origins` estiver vazio em produção, **qualquer origem** pode fazer requisições à API, permitindo ataques CSRF e roubo de dados.

**Recomendação:**
- Validar que `allowed_origins` não está vazio em produção
- Adicionar warning no startup se ambiente for production e allowed_origins estiver vazio

---

### 4. SSL Desabilitado no Banco de Dados (ALTO)
**Local:** `config.toml:20`

```toml
[database]
ssl_mode = "disable"
```

**Impacto:** Tráfego de banco de dados em texto claro, susceptível a interceptação (MITM) em redes não confiáveis.

**Recomendação:**
- Habilitar SSL (`ssl_mode = "require"` ou `"verify-full"`)
- Configurar certificados TLS para PostgreSQL

---

### 5. IP Público Exposto (MÉDIO)
**Local:** `config.toml:59`

```toml
[calling]
public_ip = "76.13.68.93"
```

**Impacto:** Expõe IP real do servidor, facilitando ataques direcionados.

**Recomendação:**
- Usar variável de ambiente para public_ip
- Considerar uso de STUN/TURN servers em vez de IP direto

---

## 🟡 VULNERABILIDADES MÉDIAS

### 6. Validação de Upload Incompleta (MÉDIO)
**Local:** `internal/handlers/media.go`, `internal/handlers/contacts.go`

**Problemas:**
- Não há validação de tamanho máximo de arquivo
- Não há verificação de MIME type real (apenas confia no header)
- Arquivos são escritos com permissão 0644 (muito permissiva)

```go
// media.go:122
if err := os.WriteFile(filePath, data, 0644); err != nil {
```

**Impacto:** Possibilidade de upload de arquivos maliciosos, DoS por arquivos grandes.

**Recomendação:**
- Validar tamanho máximo de arquivo (ex: 25MB para WhatsApp)
- Verificar MIME type real usando biblioteca como `net/http.DetectContentType`
- Usar permissões mais restrictivas (0640)
- Adicionar rate limiting em endpoints de upload

---

### 7. Caminhos de Arquivos Relativos (MÉDIO)
**Local:** `config.toml:39`, `docker/docker-compose.yml:19`

```toml
[storage]
local_path = "./uploads"
```

**Impacto:** Se o working directory mudar, os uploads podem ir para locais errados. Path traversal possível se validação falhar.

**Recomendação:**
- Usar caminhos absolutos
- Validar que paths estão dentro do diretório esperado (já feito parcialmente em media.go:180-189)

---

### 8. Falta de Rate Limiting em Alguns Endpoints (MÉDIO)
**Local:** `cmd/whatomate/main.go`

Rate limiting configurado apenas para auth endpoints:
```toml
[rate_limit]
login_max_attempts = 10
register_max_attempts = 5
```

**Impacto:** Outros endpoints sensíveis (API, webhooks) podem ser abusados.

**Recomendação:**
- Aplicar rate limiting global em todos os endpoints
- Considerar rate limiting por usuário em endpoints de API

---

### 9. Dockerfile Rodando Como Root (MÉDIO)
**Local:** `docker/Dockerfile`

Não há usuário não-root configurado:
```dockerfile
WORKDIR /app
# ... sem USER指令
CMD ["./whatomate", "server", "-config", "config.toml", "-migrate"]
```

**Impacto:** Se container for comprometido, atacante tem acesso root.

**Recomendação:**
- Criar usuário não-root
- Usar `USER nonroot` antes do CMD
- Usar read-only filesystem onde possível

---

### 10. Falta de CSP Headers (MÉDIO)
**Local:** `internal/middleware/middleware.go:97-108`

Security headers existem, mas falta Content-Security-Policy:
```go
func SecurityHeaders() fastglue.FastMiddleware {
    h.Set("X-Content-Type-Options", "nosniff")
    h.Set("X-Frame-Options", "DENY")
    // Mas falta: Content-Security-Policy
```

**Impacto:** Vulnerabilidade a XSS se houver slots para injeção de scripts.

**Recomendação:**
- Adicionar CSP header adequado
- `Content-Security-Policy: default-src 'self'; script-src 'self' 'unsafe-inline'`

---

## 🟢 VULNERABILIDADES BAIXAS

### 11. JWT Secret Validation Fraca (BAIXO)
**Local:** `cmd/whatomate/main.go:123-128`

```go
// Validate JWT secret
if a.Config.App.Environment == "production" && len(a.Config.JWT.Secret) < 32 {
    lo.Fatal("JWT secret must be at least 32 characters in production")
}
lo.Warn("JWT secret is empty, using a random secret")
```

**Problema:** Em desenvolvimento, secret vazio gera secret aleatório que muda a cada restart, causando logouts.

**Recomendação:**
- Advertência mais clara sobre riscos de secret vazio

---

### 12. Logs Podem Conter Dados Sensíveis (BAIXO)
**Local:** Vários handlers

```go
a.Log.Info("Registration completed", "user_id", user.ID, "org_id", req.OrganizationID)
a.Log.Info("Media saved", "path", relativePath, "size", len(data))
```

**Problema:** Alguns logs podem conter PII ou paths de arquivos sensíveis.

**Recomendação:**
- Revisar logs para garantir que não logam dados sensíveis
- Usar structured logging com redaction

---

### 13. Dependências Potencialmente Desatualizadas (BAIXO)
**Local:** `go.mod`

Algumas dependências estão em versões específicas que podem ter vulnerabilidades conhecidas:
- `golang.org/x/crypto v0.48.0` - Verificar CVEs recentes
- `github.com/valyala/fasthttp v1.58.0` - Atualizar para última versão

**Recomendação:**
- Rodar `govulncheck` regularmente
- Configurar Dependabot/Renovate para atualizações automáticas

---

## 📋 RECOMENDAÇÕES GERAIS

### Prioridade 1 (Imediato)
1. **Remover todos os secrets do repositório**
2. **Configurar `.gitignore` para `config.toml` e `.env`**
3. **Habilitar SSL no PostgreSQL**
4. **Configurar `allowed_origins` em produção**

### Prioridade 2 (Curto prazo)
1. **Implementar validação completa de uploads**
2. **Adicionar CSP headers**
3. **Configurar usuário não-root no Docker**
4. **Implementar rate limiting global**

### Prioridade 3 (Médio prazo)
1. **Configurar secret management (Vault/AWS)**
2. **Implementar auditoria de logs**
3. **Configurar monitoramento de vulnerabilidades**
4. **Realizar pentest completo**

---

## 🔍 TESTES REALIZADOS

| Categoria | Status | Observações |
|-----------|--------|-------------|
| SQL Injection | ✅ OK | Usa GORM com prepared statements |
| XSS | ⚠️ Parcial | Falta CSP header |
| CSRF | ⚠️ Parcial | CORS pode ser permissivo |
| Auth | ✅ OK | JWT + bcrypt implementados corretamente |
| File Upload | ⚠️ Parcial | Falta validação de tamanho e MIME real |
| Docker Security | ⚠️ Parcial | Roda como root |
| Secrets Management | ❌ CRÍTICO | Secrets hardcoded |

---

## 📞 CONTATO

Para reportar vulnerabilidades de segurança, contate:
- Email: security@seudominio.com
- GitHub Security Advisories

---

*Este relatório foi gerado automaticamente. Recomenda-se validação manual e testes adicionais.*
