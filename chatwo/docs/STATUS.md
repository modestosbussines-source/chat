# STATUS dos Squads — Omni Project

**Última atualização:** 23/03/2026 23:00  
**Coordenador:** Igor (Stakeholder)  
**Prioridade Atual:** 🔴 REMOÇÃO COMPLETA DE "WHATOMATE"

---

## 🚨 ORDEM PRIORITÁRIA

> **TODOS OS SQUADS:** Remover TODAS as 184 referências a "whatomate" do projeto.
> Veja: `/docs/ORDENS/REMOCAO-WHATOMATO.md`

---

## 🎨 UX/UI Design (Uma) — Terminal 1

**Foco Atual:** Redesign + Design System + Remoção Whatomate

### Progresso Fase 1
| # | Tarefa | Status | Notas |
|---|--------|--------|-------|
| 1 | Audit do frontend | ✅ Completo | Pattern inventory gerado |
| 2 | Design tokens | ✅ Completo | `tokens.yaml` criado |
| 3 | Shadcn/UI setup | ✅ Completo | Component library base |
| 4 | Sidebar component | ✅ Completo | Novo componente |
| 5 | Inbox component | ✅ Completo | `InboxView.vue` criado |
| 6 | Brand components | ✅ Completo | Pasta `brand/` criada |
| 7 | Auth components | ✅ Completo | Login/Register redesenhados |
| 8 | i18n pt-BR | ✅ Completo | `pt-BR.json` criado |

### 🚨 NOVA TAREFA - Remoção Whatomate
| # | Tarefa | Status | Prazo |
|---|--------|--------|-------|
| 1 | Verificar `frontend/src/**/*.ts` | ⏳ Pendente | IMEDIATO |
| 2 | Atualizar brand components | ⏳ Pendente | IMEDIATO |
| 3 | Atualizar i18n strings | ⏳ Pendente | IMEDIATO |

**Fase 1 Completa:** 90% → **95% após remoção**

---

## 🚀 Core Platform (Dex) — Terminal 2

**Foco Atual:** Backend + Integrações + Remoção Whatomate

### Progresso Fase 1
| # | Tarefa | Status | Notas |
|---|--------|--------|-------|
| 1 | DB schema audit | ✅ Completo | Schema documentado |
| 2 | Handlers refatorados | ✅ Completo | 700+ linhas em omnis.go |
| 3 | Models atualizados | ✅ Completo | omni.go modificado |
| 4 | Redis integration | ✅ Completo | redis.go atualizado |
| 5 | IVR/Calling | ✅ Completo | 5 arquivos de calling |
| 6 | WebSocket | ✅ Completo | Handlers atualizados |

### 🚨 NOVA TAREFA - Remoção Whatomate
| # | Tarefa | Status | Prazo |
|---|--------|--------|-------|
| 1 | Renomear `cmd/whatomate/` → `cmd/omni/` | ⏳ Pendente | IMEDIATO |
| 2 | Atualizar `go.mod` module path | ⏳ Pendente | IMEDIATO |
| 3 | Atualizar imports Go (~30 arquivos) | ⏳ Pendente | IMEDIATO |
| 4 | Atualizar strings "Whatomate" | ⏳ Pendente | IMEDIATO |
| 5 | Testar build: `go build -o omni ./cmd/omni` | ⏳ Pendente | IMEDIATO |

**Fase 1 Completa:** 75% → **85% após remoção**

---

## ⚙️ DevOps (Gage) — Terminal 3

**Foco Atual:** Infra + CI/CD + Remoção Whatomate (RESPONSÁVEL PRINCIPAL)

### Progresso Fase 1
| # | Tarefa | Status | Notas |
|---|--------|--------|-------|
| 1 | Docker Compose prod | ✅ Completo | Arquivo criado |
| 2 | Dockerfile | ✅ Completo | Atualizado |
| 3 | .env.production | ✅ Completo | Criado |
| 4 | DEPLOY_VPS.md | ✅ Completo | Documentação criada |
| 5 | SECURITY_REPORT.md | ✅ Completo | Security audit |

### 🚨 NOVA TAREFA - Remoção Whatomate (PRIORIDADE MÁXIMA)
| # | Tarefa | Status | Prazo |
|---|--------|--------|-------|
| 1 | Dockerfile: `whatomate` → `omni` | ⏳ Pendente | IMEDIATO |
| 2 | docker-compose: DB names | ⏳ Pendente | IMEDIATO |
| 3 | `.goreleaser.yml` | ⏳ Pendente | IMEDIATO |
| 4 | GitHub Actions (3 workflows) | ⏳ Pendente | IMEDIATO |
| 5 | `.env.example` | ⏳ Pendente | IMEDIATO |
| 6 | Testar: `docker build -t omni .` | ⏳ Pendente | IMEDIATO |

**Fase 1 Completa:** 60% → **80% após remoção**

---

## 📈 Growth (Morgan) — Terminal 4

**Foco Atual:** Docs + Remoção Whatomate (RESPÁVEL DOCS)

### Progresso Fase 1
| # | Tarefa | Status | Notas |
|---|--------|--------|-------|
| 1 | README principal | ✅ Completo | |
| 2 | CONTRIBUTING.md | ✅ Completo | |
| 3 | Getting started docs | ✅ Completo | 4 arquivos MDX |
| 4 | Layouts docs | ✅ Completo | 2 layouts Astro |
| 5 | Stories docs | ✅ Completo | 5 epics documentados |
| 6 | CHANGELOG_OMNI.md | ✅ Completo | |
| 7 | Landing page PRD | ✅ Completo | |

### 🚨 NOVA TAREFA - Remoção Whatomate (RESPÁVEL PRINCIPAL DOCS)
| # | Tarefa | Status | Prazo |
|---|--------|--------|-------|
| 1 | Reescrever landing page (`index.astro`) | ⏳ Pendente | IMEDIATO |
| 2 | Reescrever `LandingLayout.astro` | ⏳ Pendente | IMEDIATO |
| 3 | Reescrever `OmniLayout.astro` | ⏳ Pendente | IMEDIATO |
| 4 | Atualizar 5 stories (epics 1-5) | ⏳ Pendente | IMEDIATO |
| 5 | Reescrever PRD documentation | ⏳ Pendente | IMEDIATO |
| 6 | Atualizar SECURITY_REPORT | ⏳ Pendente | IMEDIATO |
| 7 | Novo favicon + og-image | ⏳ Pendente | IMEDIATO |

**Fase 1 Completa:** 80% → **100% após remoção**

---

## 🤖 AI/ML (Dara) — Terminal 5

**Foco Atual:** Chatbot + IA + Remoção Whatomate

### Progresso Fase 1
| # | Tarefa | Status | Notas |
|---|--------|--------|-------|
| 1 | Chatbot handlers | ✅ Completo | `chatbot.go` atualizado |
| 2 | Chatbot processor | ✅ Completo | Processor + tests |
| 3 | Flows handlers | ✅ Completo | `flows.go` atualizado |
| 4 | Custom actions | ✅ Completo | `custom_actions.go` |
| 5 | IVR flows | ✅ Completo | `ivr_flows.go` |

### 🚨 NOVA TAREFA - Remoção Whatomate
| # | Tarefa | Status | Prazo |
|---|--------|--------|-------|
| 1 | Verificar imports Go internos | ⏳ Pendente | IMEDIATO |
| 2 | Atualizar model names | ⏳ Pendente | IMEDIATO |
| 3 | Atualizar strings de output | ⏳ Pendente | IMEDIATO |

**Fase 1 Completa:** 50% → **70% após remoção**

---

## 📊 PROGRESSO GERAL

```
┌─────────────────────────────────────────────────────────────────────────┐
│                    PROGRESSO POR SQUAD                                  │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  🎨 UX/UI (Uma)      ████████████████████████████████░░░░  90% → 95%   │
│  📈 Growth (Morgan)  ████████████████████████████░░░░░░░░  80% → 100%  │
│  🚀 Core (Dex)       ████████████████████████████░░░░░░░░  75% → 85%   │
│  ⚙️ DevOps (Gage)    ████████████████████░░░░░░░░░░░░░░░░  60% → 80%   │
│  🤖 AI/ML (Dara)     ████████████████░░░░░░░░░░░░░░░░░░░░  50% → 70%   │
│                                                                         │
├─────────────────────────────────────────────────────────────────────────┤
│  📁 Total: 128+ arquivos modificados                                    │
│  📝 +3,649 linhas de código                                             │
│  🚨 184 referências "whatomate" para remover                            │
└─────────────────────────────────────────────────────────────────────────┘
```

---

## ⏭️ O QUE FALTA PARA FINALIZAR FASE 1

### Tarefas Bloqueantes (DevOps + Core)

| # | Tarefa | Responsável | Bloqueia |
|---|--------|-------------|----------|
| 1 | Remoção whatomate (184 refs) | TODOS | Build, Deploy |
| 2 | Build Go funcional | Core + DevOps | Deploy VPS |
| 3 | Docker build funcional | DevOps | Deploy VPS |
| 4 | Testes passando | QA | Deploy VPS |

### Tarefas de Qualidade

| # | Tarefa | Responsável | Prioridade |
|---|--------|-------------|------------|
| 5 | Storybook documentação | UX/UI | Média |
| 6 | API docs OpenAPI | Core | Média |
| 7 | Load testing | DevOps | Baixa |
| 8 | Security review final | DevOps | Alta |

### Tarefas de Documentação

| # | Tarefa | Responsável | Prioridade |
|---|--------|-------------|------------|
| 9 | User manual pt-BR | Growth | Alta |
| 10 | Video tutorials script | Growth | Média |
| 11 | FAQ | Growth | Média |

---

## PRÓXIMOS PASSOS IMEDIATOS

1. **AGORA:** Todos os squads iniciar remoção de "whatomate"
2. **1H:** DevOps testa build após remoção
3. **2H:** Todos commitam mudanças
4. **4H:** Deploy de staging para validação
5. **FIM DO DIA:** Phase 1 oficialmente completa

---

*Atualizado por Morgan (PM) — 23/03/2026*
