# 🚨 ORDEM OFICIAL — REMOÇÃO COMPLETA DE "WHATOMATO"

**Data:** 23/03/2026  
**Emissor:** Igor (Stakeholder)  
**Urgência:** 🔴 ALTA  
**Prazo:** IMEDIATO

---

## DECLARAÇÃO OFICIAL

> **NÃO SOMOS WHATOMATE.**
>
> O Whatomate é um projeto open-source de terceiros (shridarpatil/whatomate).
> Nós usamos o esqueleto/código-fonte como BASE para desenvolver o **OMNI**.
>
> **O Omni é:**
> - Um produto PRIVADO e PROPRIETÁRIO
> - 100% vendável como SaaS ou self-host commercial
> - Rebrandado, melhorado e otimizado para o mercado brasileiro
> - Sem dependência ou associação ao projeto open-source original
>
> **TODA referência a "whatomate" deve ser REMOVIDA completamente.**
> Código, docs, configs, nomes de arquivos, variáveis de ambiente, TUDO.

---

## ESCOPO DA REMOÇÃO

### 184 referências encontradas em:

| Categoria | Arquivos | Ocorrências |
|-----------|----------|-------------|
| **Go Source** | `internal/handlers/*.go`, `internal/models/*.go` | ~30 |
| **Frontend** | `frontend/src/**/*.vue`, `*.ts` | ~15 |
| **Docs** | `docs/**/*.astro`, `*.md`, `*.mdx` | ~80 |
| **CI/CD** | `.github/workflows/*.yml` | ~10 |
| **Docker** | `docker/Dockerfile`, `docker-compose.yml` | ~5 |
| **Config** | `.goreleaser.yml`, `config.toml`, `.env*` | ~10 |
| **Stories** | `docs/stories/*.md` | ~30 |

---

## CHECKLIST POR SQUAD

### ⚙️ DevOps (Gage) — RESPONSÁVEL PRINCIPAL

| # | Tarefa | Arquivos | Status |
|---|--------|----------|--------|
| 1 | Renomear binário no Dockerfile | `docker/Dockerfile` | ⏳ |
| 2 | Atualizar docker-compose | `docker/docker-compose.yml` | ⏳ |
| 3 | Atualizar .env.example | `docker/.env.example` | ⏳ |
| 4 | Atualizar .goreleaser.yml | `.goreleaser.yml` | ⏳ |
| 5 | Atualizar GitHub Actions | `.github/workflows/*.yml` | ⏳ |
| 6 | Renomear bancos de teste | `*.yml` (whatomate_test → omni_test) | ⏳ |

**Substituições:**
```bash
# Dockerfile
FROM whatomate → FROM omni
-o whatomate → -o omni
CMD ["./whatomate" → CMD ["./omni"

# docker-compose
POSTGRES_DB: whatomate → POSTGRES_DB: omni
POSTGRES_DB: whatomate_test → POSTGRES_DB: omni_test
```

---

### 🚀 Core Platform (Dex) — RESPONSÁVEL SECUNDÁRIO

| # | Tarefa | Arquivos | Status |
|---|--------|----------|--------|
| 1 | Atualizar Go imports | `internal/handlers/*.go` | ⏳ |
| 2 | Atualizar module path | `go.mod` | ⏳ |
| 3 | Renomear main.go output | `cmd/whatomate/` → `cmd/omni/` | ⏳ |
| 4 | Atualizar modelos | `internal/models/*.go` | ⏳ |
| 5 | Atualizar handlers | `internal/handlers/*.go` | ⏳ |

**Substituições:**
```go
// go.mod
module github.com/shridarpatil/whatomate → module github.com/omni/platform

// imports
"github.com/shridarpatil/whatomate/internal/models" → "github.com/omni/platform/internal/models"

// output
log.Println("Whatomate") → log.Println("Omni")
```

---

### 📈 Growth (Morgan) — RESPONSÁVEL DOCS

| # | Tarefa | Arquivos | Status |
|---|--------|----------|--------|
| 1 | Reescrever landing page | `docs/src/pages/index.astro` | ⏳ |
| 2 | Reescrever layouts | `docs/src/layouts/*.astro` | ⏳ |
| 3 | Atualizar todas as stories | `docs/stories/*.md` | ⏳ |
| 4 | Reescrever PRD docs | `docs/prd-documentation-landing.md` | ⏳ |
| 5 | Atualizar CHANGELOG | `CHANGELOG_OMNI.md` | ⏳ |
| 6 | Atualizar SECURITY_REPORT | `SECURITY_REPORT.md` | ⏳ |
| 7 | Criar novo favicon/og-image | `docs/public/` | ⏳ |

**Substituições:**
```
Whatomate → Omni
whatomate → omni
shridarpatil/whatomate → (remover referência)
https://shridarpatil.github.io/whatomate → https://omni.com.br
"Modern WhatsApp Business Platform" → "Plataforma Omnichannel Brasileira"
```

---

### 🎨 UX/UI (Uma) — VERIFICAÇÃO

| # | Tarefa | Arquivos | Status |
|---|--------|----------|--------|
| 1 | Verificar frontend imports | `frontend/src/**/*.ts` | ⏳ |
| 2 | Atualizar brand components | `frontend/src/components/brand/*` | ⏳ |
| 3 | Atualizar i18n strings | `frontend/src/i18n/locales/*.json` | ⏳ |

---

### 🤖 AI/ML (Dara) — VERIFICAÇÃO

| # | Tarefa | Arquivos | Status |
|---|--------|----------|--------|
| 1 | Verificar AI handlers | `internal/handlers/chatbot*.go` | ⏳ |
| 2 | Atualizar model names | Models internos | ⏳ |

---

## SUBSTITUIÇÕES GLOBAIS

### Arquivos Go
```bash
# Imports
find . -name "*.go" -exec sed -i 's|github.com/shridarpatil/whatomate|github.com/omni/platform|g' {} +

# Strings
find . -name "*.go" -exec sed -i 's|Whatomate|Omni|g' {} +
find . -name "*.go" -exec sed -i 's|whatomate|omni|g' {} +
```

### Arquivos Markdown/MDX
```bash
find . -name "*.md" -o -name "*.mdx" | xargs sed -i 's|Whatomate|Omni|g'
find . -name "*.md" -o -name "*.mdx" | xargs sed -i 's|whatomate|omni|g'
find . -name "*.md" -o -name "*.mdx" | xargs sed -i 's|shridarpatil/whatomate||g'
```

### Arquivos Astro
```bash
find . -name "*.astro" | xargs sed -i 's|Whatomate|Omni|g'
find . -name "*.astro" | xargs sed -i 's|whatomate|omni|g'
```

### YAML/YML
```bash
find . -name "*.yml" -o -name "*.yaml" | xargs sed -i 's|whatomate|omni|g'
```

---

## VALIDAÇÃO FINAL

Após todas as substituições, verificar:

```bash
# Deve retornar 0 resultados
grep -ri "whatomate" --include="*.go" --include="*.vue" --include="*.ts" --include="*.md" --include="*.yml" --include="*.astro" .

# Deve retornar 0 resultados
grep -ri "shridarpatil" --include="*.go" --include="*.vue" --include="*.ts" --include="*.md" --include="*.yml" --include="*.astro" .

# Deve retornar 0 resultados
grep -ri "Whatomate" .
```

---

## RENOMEAÇÃO DE DIRETÓRIOS

| De | Para |
|----|------|
| `cmd/whatomate/` | `cmd/omni/` |

```bash
mv cmd/whatomate cmd/omni
```

---

## COMUNICAÇÃO À EQUIPE

### Para cada agente nos terminais:

```
📋 ORDEM DO STAKEHOLDER (Igor):

O projeto NÃO é Whatomate. O Whatomate é um projeto open-source de terceiros
que usamos como BASE/CODE SKELETON.

Nós desenvolvemos o OMNI — um produto PRIVADO, PROPRIETÁRIO e 100% VENDÁVEL.

TODAS as referências a "whatomate", "Whatomate", "shridarpatil/whatomate"
devem ser REMOVIDAS completamente do seu escopo.

Substitua por "Omni", "omni", ou remova a referência.

Isso é OBRIGATÓRIO e IMEDIATO.

Arquivos que você precisa verificar:
- [listar arquivos relevantes para cada squad]

Execute as substituições e atualize o STATUS.md quando concluído.
```

---

## PRIORIDADE

1. 🔴 **IMEDIATO:** CI/CD, Docker, Go imports (bloqueia build)
2. 🟡 **ALTO:** Docs, stories, landing page
3. 🟢 **MÉDIO:** Testes, validação final

---

## DEFINIÇÃO DE SUCESSO

- [ ] 0 ocorrências de "whatomate" em todo o código
- [ ] 0 ocorrências de "shridarpatil" em todo o código
- [ ] Build funciona: `go build -o omni ./cmd/omni`
- [ ] Docker build funciona: `docker build -t omni .`
- [ ] Todos os testes passam
- [ ] Docs mostram "Omni" em vez de "Whatomate"

---

*Esta ordem é OBRIGATÓRIA e deve ser executada IMEDIATAMENTE por todos os squads.*

**— Igor, definindo a direção 🎯**
