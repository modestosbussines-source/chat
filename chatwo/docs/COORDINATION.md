# OMNICHANNEL PROJECT — COORDENAÇÃO MULTI-SQUAD

**Última atualização:** 23/03/2026  
**Coordenador:** Igor (Stakeholder)  
**PRD:** `/docs/product/PRD-OMNI-PLATFORM.md`

---

## Visão Geral dos Squads Ativos

| Squad | Agente | Terminal | Foco Principal | Status |
|-------|--------|----------|----------------|--------|
| 🎨 UX/UI Design | @ux-design-expert (Uma) | Terminal 1 | Redesign + Design System | 🟢 Ativo |
| 🚀 Core Platform | @dev (Dex) | Terminal 2 | Backend + Integrações | 🟢 Ativo |
| ⚙️ DevOps | @devops (Gage) | Terminal 3 | Infra + CI/CD | 🟢 Ativo |
| 🤖 AI/ML | @data-engineer (Dara) | Terminal 4 | Chatbot + IA | 🟢 Ativo |
| 📈 Growth | @pm (Morgan) | Terminal 5 | Docs + Onboarding | 🟢 Ativo |

---

## Tarefas Focais por Squad

### 🎨 UX/UI Design (Uma) — PRIORIDADE MÁXIMA

**Objetivo:** Reformular completamente a interface do Omni

| # | Tarefa | Comando | Deliverable | Deadline |
|---|--------|---------|-------------|----------|
| 1 | Audit do frontend atual | `*audit ./frontend/src` | Pattern inventory | Hoje |
| 2 | Criar design tokens | `*tokenize` + manual | `/docs/design/tokens.yaml` | Semana 1 |
| 3 | Setup Shadcn/UI | `*bootstrap-shadcn` | Component library base | Semana 1 |
| 4 | Redesign Sidebar | `*build sidebar` | Componente Sidebar | Semana 2 |
| 5 | Redesign Inbox | `*build inbox` | Componente Inbox | Semana 3 |
| 6 | Storybook setup | Configurar docs | `/docs/design/storybook` | Semana 3 |

**Arquivos que pode modificar:**
- `/frontend/src/components/**`
- `/frontend/src/views/**`
- `/docs/design/**`

**NÃO modificar:**
- Backend code
- Docker/deployment files
- Database schema

---

### 🚀 Core Platform (Dex)

**Objetivo:** Solidificar backend e integrações principais

| # | Tarefa | Focus | Deliverable | Deadline |
|---|--------|-------|-------------|----------|
| 1 | Schema DB audit | Revisar schema atual | `/docs/architecture/db-schema.md` | Hoje |
| 2 | WhatsApp integration refactor | API handlers | Código limpo + tests | Semana 1-2 |
| 3 | WebSocket optimization | Real-time messaging | Performance benchmarks | Semana 2 |
| 4 | API documentation | OpenAPI spec | `/docs/api/openapi.yaml` | Semana 2 |
| 5 | Contact/CRM model | Novo model | Migration + model | Semana 3 |

**Arquivos que pode modificar:**
- `/app/models/**`
- `/app/controllers/**`
- `/app/services/**`
- `/spec/**`
- `/db/migrate/**`

**NÃO modificar:**
- Frontend components (deixar para UX/UI)
- Docker files (deixar para DevOps)
- CI/CD pipelines

---

### ⚙️ DevOps (Gage)

**Objetivo:** Infraestrutura production-ready

| # | Tarefa | Focus | Deliverable | Deadline |
|---|--------|-------|-------------|----------|
| 1 | Docker Compose prod | Arquivo de produção | `/docker-compose.prod.yml` | Hoje |
| 2 | CI/CD Pipeline | GitHub Actions | `/.github/workflows/ci.yml` | Semana 1 |
| 3 | Monitoring setup | Prometheus/Grafana | Config files | Semana 1 |
| 4 | SSL/Security | Nginx config | `/docker/nginx/ssl.conf` | Semana 2 |
| 5 | Backup automation | Cron + scripts | `/scripts/backup.sh` | Semana 2 |

**Arquivos que pode modificar:**
- `/docker/**`
- `/.github/**`
- `/scripts/infra/**`
- `/config/deploy/**`

**NÃO modificar:**
- Application code
- Database schema
- Frontend

---

### 🤖 AI/ML Engine (Dara)

**Objetivo:** Chatbot builder e IA conversacional pt-BR

| # | Tarefa | Focus | Deliverable | Deadline |
|---|--------|-------|-------------|----------|
| 1 | Intent classifier pt-BR | ML model | `/ai/models/intent-pt-br.py` | Semana 1-2 |
| 2 | Sentiment analysis | NLP pipeline | `/ai/models/sentiment.py` | Semana 2 |
| 3 | Auto-response engine | RAG + LLM | `/ai/services/auto-responder.py` | Semana 3 |
| 4 | Chatbot flow builder | Visual builder API | `/app/services/chatbot/` | Semana 3-4 |
| 5 | Analytics AI | Métricas bots | `/ai/analytics/` | Semana 4 |

**Arquivos que pode modificar:**
- `/ai/**`
- `/app/services/chatbot/**`
- `/ml/**`
- `/data/models/**`

**NÃO modificar:**
- Core backend (deixar para Core Platform)
- Frontend UI (deixar para UX/UI)
- Infrastructure

---

### 📈 Growth (Morgan)

**Objetivo:** Documentação, onboarding e aquisição

| # | Tarefa | Focus | Deliverable | Deadline |
|---|--------|-------|-------------|----------|
| 1 | README principal | Projeto overview | `/README.md` | Hoje |
| 2 | Contributing guide | Guia de contribuição | `/CONTRIBUTING.md` | Hoje |
| 3 | Setup guide pt-BR | Instalação | `/docs/setup/INSTALL.md` | Semana 1 |
| 4 | User manual | Uso da plataforma | `/docs/user/MANUAL.md` | Semana 2 |
| 5 | API documentation | Integrações | `/docs/api/README.md` | Semana 2 |
| 6 | Changelog template | Versões | `/CHANGELOG.md` | Semana 1 |

**Arquivos que pode modificar:**
- `/README.md`
- `/docs/**`
- `/CONTRIBUTING.md`
- `/CHANGELOG.md`
- `/docs/user/**`
- `/docs/setup/**`

**NÃO modificar:**
- Application code
- Infrastructure config
- Design files

---

## Protocolo de Sync entre Squads

### Como Coordenar (Para Você, Igor)

#### 1. Verificar Progresso de Cada Squad

Pergunte em cada terminal:
```
Qual é o status atual? O que foi completado?
```

#### 2. Identificar Blockers

Se um squad estiver travado:
```
Qual é o blocker? Precisa de input de qual outro squad?
```

#### 3. Sincronizar Dependências

| Dependência | De | Para | Ação |
|-------------|-----|------|------|
| Design tokens | UX/UI | Core, AI | UX/UI cria → outros consomem |
| API endpoints | Core | AI, Frontend | Core define → outros integram |
| Infra config | DevOps | Todos | DevOps setup → outros deploy |
| Docs template | Growth | Todos | Growth cria → outros contribuem |

### Arquivo de Status compartilhado

Cada squad atualiza este arquivo ao completar tarefas:

**Arquivo:** `/docs/STATUS.md`

```markdown
# Status dos Squads — Atualizado: [DATA]

## UX/UI (Uma)
- [x] Audit do frontend
- [ ] Design tokens
- [ ] Shadcn/UI setup
- Blockers: aguardando definição de palette

## Core Platform (Dex)
- [x] DB schema audit
- [ ] WhatsApp refactor
- [ ] API docs
- Blockers: nenhum

## DevOps (Gage)
- [x] Docker Compose prod
- [ ] CI/CD pipeline
- Blockers: aguardando definição de env vars

## AI/ML (Dara)
- [ ] Intent classifier
- [ ] Sentiment analysis
- Blockers: aguardando API endpoints do Core

## Growth (Morgan)
- [x] README
- [ ] Setup guide
- Blockers: aguardando output dos outros squads
```

---

## Comandos Rápidos por Terminal

### Para você passar para cada agente:

**UX/UI (Uma):**
```
 Leia o PRD em docs/product/PRD-OMNI-PLATFORM.md
 Foque na seção 6 (Design) e 12 (Anexos)
 Comece com *audit ./frontend/src
```

**Core (Dex):**
```
 Leia o PRD em docs/product/PRD-OMNI-PLATFORM.md  
 Foque na seção 4.1 (Epic 1 e 2) e 7 (Arquitetura)
 Comece auditando o schema DB atual
```

**DevOps (Gage):**
```
 Leia o PRD em docs/product/PRD-OMNI-PLATFORM.md
 Foque na seção 5 (Non-functional) e 7.1 (Stack)
 Comece criando docker-compose.prod.yml
```

**AI/ML (Dara):**
```
 Leia o PRD em docs/product/PRD-OMNI-PLATFORM.md
 Foque na seção 4.1 (Epic 3 - Chatbot)
 Comece definindo estrutura do modelo de intents pt-BR
```

**Growth (Morgan):**
```
 Leia o PRD em docs/product/PRD-OMNI-PLATFORM.md
 Foque na seção 8 (Roadmap) e 10 (Recursos)
 Comece criando README.md principal
```

---

## Reunião de Sync (Opcional)

Quando quiser alinhar todos, pergunte em cada terminal:

```
1. O que você completou?
2. O que está trabalhando agora?
3. Precisa de algo de outro squad?
4. Quando estará pronto o próximo deliverable?
```

---

## Métricas de Progresso Semanal

| Squad | KPI Semanal |
|-------|-------------|
| UX/UI | # componentes criados, % cobertura tokens |
| Core | # endpoints documentados, % test coverage |
| DevOps | # pipelines funcionando, uptime |
| AI/ML | # modelos treinados, accuracy |
| Growth | # páginas docs, completeness |

---

*Coordenação criada por Morgan (PM) — 23/03/2026*
