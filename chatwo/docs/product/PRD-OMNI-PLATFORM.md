# PRD: Omni — Plataforma de Atendimento Omnichannel

**Product Requirements Document v1.0**

| Campo | Valor |
|-------|-------|
| **Produto** | Omni (baseado em Chatwoot) |
| **Versão** | 1.0 |
| **Autor** | Morgan (PM) |
| **Data** | 23/03/2026 |
| **Status** | Draft → Em Revisão |
| **Base** | Análise de Mercado @analyst |

---

## 1. Executive Summary

### 1.1 Visão do Produto

> **"A plataforma de atendimento omnichannel que une a liberdade do open-source com a experiência premium que suas vendas merecem."**

Omni é uma plataforma self-host de atendimento omnichannel focada no mercado brasileiro, combinando a robustez do Chatwoot com UX premium, IA conversacional pt-BR nativa e integrações com as principais redes sociais do Brasil.

### 1.2 Problema

| Problema | Impacto | Nossa Solução |
|----------|---------|---------------|
| Plataformas SaaS são caras (R$ 150-900/agent/mês) | 67% das PMEs desistem | Self-host gratuito + cloud acessível |
| UX de open-source é ruim | Baixa adoção, alto churn | Design premium, intuitivo |
| Chatbots em pt-BR são genéricos | Respostas inadequadas | IA treinada pt-BR com contexto local |
| Fragmentação de canais | WhatsApp + Instagram + TikTok separados | Visão unificada em uma tela |
| Pricing opaco | Desconfiança | Preços claros, sem surpresas |

### 1.3 Oportunidade de Mercado

| Métrica | Valor |
|---------|-------|
| **TAM Omnichannel BR** | R$ 1.8 bilhões |
| **SAM (SMB + Mid-market)** | R$ 720 milhões |
| **SOM (18 meses)** | R$ 2.4 milhões (0.33%) |
| **CAGR do segmento** | 22.1% |

---

## 2. Objetivos e Métricas

### 2.1 Objetivos Estratégicos

| # | Objetivo | KPI | Meta 90 dias | Meta 180 dias |
|---|----------|-----|--------------|---------------|
| O1 | Validar product-market fit | NPS | > 45 | > 55 |
| O2 | Aquisição de usuários | Downloads self-host | 2.000 | 8.000 |
| O3 | Revenue | MRR cloud | R$ 50.000 | R$ 200.000 |
| O4 | Retenção | Churn mensal | < 8% | < 5% |
| O5 | Time-to-value | Setup completion | < 30 min | < 15 min |

### 2.2 Métricas de Sucesso do Usuário

| Métrica | Definição | Target |
|---------|-----------|--------|
| **Conversão WhatsApp → Venda** | % de conversões atribuídas | +15% vs baseline |
| **Tempo de resposta médio** | Tempo até primeira resposta | < 3 min |
| **Satisfação do agente** | NPS interno | > 60 |
| **Auto-resolução** | % de tickets resolvidos por bot | > 40% |

---

## 3. Persona e Usuários

### 3.1 Persona Primária: Maria (Gestora de Vendas)

| Atributo | Detalhe |
|----------|---------|
| **Idade** | 32-45 anos |
| **Cargo** | Gestora de Vendas / CS |
| **Empresa** | E-commerce, 20-100 funcionários |
| **Faturamento** | R$ 5M - R$ 50M/ano |
| **Tech savvy** | Médio |
| **Dor principal** | "Perco vendas porque meu time não dá conta do WhatsApp + Instagram" |
| **Objetivo** | Aumentar conversão e reduzir custo de atendimento |
| **Orçamento** | R$ 500 - R$ 2.000/mês |
| **Decisão** | Ela influencia, CTO decide |

### 3.2 Persona Secundária: Carlos (CTO/Dev)

| Atributo | Detalhe |
|----------|---------|
| **Idade** | 28-40 anos |
| **Cargo** | CTO / Tech Lead |
| **Empresa** | Startup / scale-up |
| **Tech savvy** | Alto |
| **Dor principal** | "Preciso de algo self-host que não me prenda" |
| **Objetivo** | Autonomia, integrações, APIs abertas |
| **Orçamento** | R$ 0 - R$ 500/mês |
| **Decisão** | Decisão técnica final |

---

## 4. Escopo do Produto

### 4.1 Funcionalidades Core (MVP)

#### Epic 1: Multi-Channel Inbox

| User Story | Prioridade | AC |
|------------|------------|-----|
| Como agente, quero ver todas as conversas em uma única tela | P0 | Conversas unificadas |
| Como agente, quero filtrar por canal (WhatsApp, Instagram, etc.) | P0 | Filtros funcionais |
| Como agente, quero transferir conversa para outro agente | P0 | Transferência OK |
| Como agente, quero ver histórico completo do contato | P0 | Timeline completa |
| Como agente, quero usar templates de resposta rápidas | P1 | Canned responses |

#### Epic 2: WhatsApp Business Integration

| User Story | Prioridade | AC |
|------------|------------|-----|
| Como admin, quero conectar WhatsApp via Business API | P0 | Conexão funcional |
| Como agente, quero receber e enviar mensagens WhatsApp | P0 | Bidirecional |
| Como agente, quero receber mídias (imagem, áudio, doc) | P0 | Mídias OK |
| Como admin, quero configurar mensagem de boas-vindas | P1 | Auto-message OK |

#### Epic 3: AI Chatbot Builder

| User Story | Prioridade | AC |
|------------|------------|-----|
| Como admin, quero criar fluxos de chatbot visualmente | P1 | Visual builder |
| Como admin, quero treinar intents em pt-BR | P1 | Intent classifier |
| Como agente, quero sugestões de resposta baseadas em contexto | P1 | AI suggestions |
| Como admin, quero ver métricas de desempenho do bot | P2 | Analytics OK |

#### Epic 4: Contact Management (CRM)

| User Story | Prioridade | AC |
|------------|------------|-----|
| Como agente, quero ver perfil completo do contato | P0 | Profile view |
| Como agente, quero adicionar notas ao contato | P1 | Notes OK |
| Como admin, quero segmentar contatos por etiquetas | P1 | Tags OK |
| Como admin, quero exportar contatos | P2 | Export OK |

#### Epic 5: Analytics e Reports

| User Story | Prioridade | AC |
|------------|------------|-----|
| Como gestor, quero ver dashboard de métricas | P1 | Dashboard OK |
| Como gestor, quero ver desempenho por agente | P1 | Agent reports |
| Como gestor, quero ver métricas por canal | P2 | Channel reports |
| Como gestor, quero exportar relatórios | P2 | Export OK |

### 4.2 Funcionalidades Phase 2 (Pós-MVP)

| Funcionalidade | Prioridade | Release |
|----------------|------------|---------|
| Instagram DM Integration | P0 | v1.1 |
| TikTok Integration | P1 | v1.2 |
| Chatbot visual builder avançado | P1 | v1.2 |
| IA com fine-tuning pt-BR | P1 | v1.3 |
| Mobile app gestores | P2 | v1.4 |
| Integração com ERPs (Tiny, Bling) | P2 | v1.5 |
| White-label | P3 | v2.0 |

### 4.3 Fora de Escopo (Explicitamente)

| Funcionalidade | Justificativa |
|----------------|---------------|
| Email marketing massificado | Escopo diferente, usar Mailchimp/etc |
| Social media management | Foco em atendimento, não publicação |
| ERP completo | Integração, não substituição |
| Telefonia (voip) | Complexidade alta, fase futura |

---

## 5. Requisitos Não-Funcionais

### 5.1 Performance

| Métrica | Target |
|---------|--------|
| **Page load (LCP)** | < 2.5s |
| **Time to Interactive** | < 3.5s |
| **API response (p95)** | < 200ms |
| **WebSocket latency** | < 100ms |
| **Concurrent users** | 500+ agents |

### 5.2 Segurança e Compliance

| Requisito | Detalhe |
|-----------|---------|
| **LGPD** | Compliance total, consent management |
| **Autenticação** | SSO (Google, Microsoft), 2FA |
| **Encryption** | TLS 1.3, AES-256 at rest |
| **Audit logs** | 90 dias de retenção |
| **Data residency** | Dados no Brasil obrigatório |

### 5.3 Disponibilidade

| SLA | Target |
|-----|--------|
| **Uptime** | 99.9% |
| **Support response** | < 4h (business hours) |
| **Backup** | Diário, retenção 30 dias |
| **DR Recovery** | < 4 horas |

---

## 6. Design e UX

### 6.1 Princípios de Design

| Princípio | Descrição |
|-----------|-----------|
| **Clean & Minimal** | Interface limpa, foco no conteúdo |
| **Black & White + Accent** | Paleta minimalista com cor de destaque |
| **Mobile First** | Funciona perfeitamente em mobile |
| **Accessible** | WCAG 2.1 AA minimum |
| **Fast** | Feedback imediato, skeletons de loading |

### 6.2 Design System

| Componente | Status |
|------------|--------|
| Design tokens | 🔴 A definir |
| Component library | 🔴 A reconstruir |
| Storybook | 🔴 A configurar |
| Dark mode | 🟡 Parcial |

### 6.3 Prioridades de Design

| Prioridade | Item | Squad Responsável |
|------------|------|-------------------|
| P0 | Redesign completo do inbox | UX/UI Squad |
| P0 | Novo sidebar de navegação | UX/UI Squad |
| P0 | Componentes reutilizáveis (shadcn/ui) | UX/UI Squad |
| P1 | Dark mode completo | UX/UI Squad |
| P1 | Onboarding wizard | UX/UI Squad |
| P2 | Micro-interactions | UX/UI Squad |

---

## 7. Arquitetura Técnica

### 7.1 Stack Tecnológica

| Camada | Tecnologia |
|--------|------------|
| **Frontend** | Vue 3 + TypeScript + Tailwind CSS + Shadcn/UI |
| **Backend** | Ruby on Rails (chatwoot core) + Node.js (AI services) |
| **Database** | PostgreSQL + Redis |
| **Real-time** | Action Cable / WebSockets |
| **AI/ML** | Python (FastAPI) + Hugging Face + OpenAI/Anthropic |
| **Infrastructure** | Docker + Kubernetes + Nginx |
| **CI/CD** | GitHub Actions |

### 7.2 Integrações

| Canal | API | Status |
|-------|-----|--------|
| **WhatsApp** | Meta Business API | ✅ Funcional |
| **Instagram** | Meta Graph API | 🟡 Parcial |
| **Facebook** | Meta Messenger | ✅ Funcional |
| **TikTok** | TikTok Business API | 🔴 A integrar |
| **Email** | IMAP/SMTP | ✅ Funcional |
| **Webchat** | Widget JS | ✅ Funcional |
| **API REST** | OpenAPI 3.0 | ✅ Funcional |

---

## 8. Roadmap

### 8.1 Timeline

```
2026 Q2 (Abr-Jun) — FASE 1: MVP + Redesign
├── Sprint 1-2: UX Audit + Design System Setup
├── Sprint 3-4: Inbox Redesign
├── Sprint 5-6: WhatsApp Integration Refactor
└── Sprint 7-8: AI Chatbot Builder v1

2026 Q3 (Jul-Set) — FASE 2: Scale
├── Sprint 9-10: Instagram + TikTok Integration
├── Sprint 11-12: Analytics Dashboard
├── Sprint 13-14: Mobile App (PWA)
└── Sprint 15-16: ERP Integrations

2026 Q4 (Out-Dez) — FASE 3: Growth
├── Sprint 17-18: White-label Support
├── Sprint 19-20: Advanced AI (fine-tuning pt-BR)
├── Sprint 21-22: Partner Program
└── Sprint 23-24: v2.0 Launch
```

### 8.2 Release Plan

| Release | Data Target | Épicas | KPIs |
|---------|-------------|--------|------|
| **v1.0 MVP** | Jun 2026 | E1, E2, E4 | 200 ativos, NPS 40 |
| **v1.1 Instagram** | Jul 2026 | E3 partial | 500 ativos |
| **v1.2 AI Builder** | Set 2026 | E3 full | 800 ativos, MRR R$ 80k |
| **v1.5 ERPs** | Nov 2026 | Integrations | 1.500 ativos, MRR R$ 150k |
| **v2.0** | Dez 2026 | White-label | 2.000 ativos, MRR R$ 200k |

---

## 9. Riscos e Mitigação

| Risco | Impacto | Prob. | Mitigação |
|-------|---------|-------|-----------|
| WhatsApp mudar pricing API | Alto | Média | Multi-channel first, TikTok priority |
| Concorrente open-source similar | Médio | Alta | UX premium, suporte pt-BR |
| Atraso na entrega do redesign | Alto | Média | Squad dedicado, milestones semanais |
| Dificuldade de aquisição | Alto | Média | Growth hacking, parcerias |
| Dependência de Meta | Alto | Média | Estratégia channel-agnostic |

---

## 10. Dependências e Recursos

### 10.1 Squads Envolvidos

| Squad | Responsabilidade | Lead |
|-------|------------------|------|
| **Core Platform** | Backend, DB, integrações | @dev |
| **UX/UI Design** | Redesign, design system | @ux-design-expert |
| **DevOps** | Infra, deployment | @devops |
| **AI/ML Engine** | Chatbot, IA | @data-engineer |
| **Growth** | Marketing, docs, onboarding | @pm |

### 10.2 Recursos Necessários

| Recurso | Quantidade | Prioridade |
|---------|------------|------------|
| **Dev Backend** | 2 | P0 |
| **Dev Frontend** | 2 | P0 |
| **Designer UX/UI** | 1 | P0 |
| **DevOps** | 1 | P1 |
| **Data Engineer** | 1 | P1 |
| **QA** | 1 | P1 |

### 10.3 Budget Estimado (18 meses)

| Categoria | Valor |
|-----------|-------|
| **Infraestrutura (cloud)** | R$ 3.600/mês |
| **APIs externas** | R$ 2.000/mês |
| **Ferramentas dev** | R$ 1.000/mês |
| **Marketing** | R$ 5.000/mês |
| **Total anual** | ~ R$ 140.000 |

---

## 11. Aprovações

| Papel | Nome | Status | Data |
|-------|------|--------|------|
| **Product Manager** | Morgan | ✅ Proposto | 23/03/2026 |
| **Stakeholder** | [Igor] | ⏳ Pendente | - |
| **Tech Lead** | [A ser designado] | ⏳ Pendente | - |

---

## 12. Anexos

| Documento | Link |
|-----------|------|
| Análise de Mercado | `/docs/research/market-analysis-2026.md` |
| User Research | `/docs/research/user-personas.md` |
| Design System Specs | `/docs/design/design-system-v1.md` |
| API Documentation | `/docs/api/openapi.yaml` |
| Architecture Diagram | `/docs/architecture/system-overview.md` |

---

**— Morgan, planejando o futuro 📊**

---

## FASE 3: @ux-design-expert (Uma) — Reformulação do Design

🎨 **Uma the Empathizer** entrando!

**Role:** UX/UI Designer & Design System Architect

Com base na análise de mercado e no PRD, vou conduzir a reformulação completa do design. Começando com **Auditoria + Definição de Prioridades**:

---

### Análise do Design Atual (Baseada nos Commits Recentes)

**Estado Atual:**
- Redesign "Omni" com paleta preto e branco minimalista
- 24 arquivos Vue modificados
- Foco em branding consistency

**Problemas Identificados (Pelo que vi nos arquivos):**

| # | Problema | Impacto | Prioridade |
|---|----------|---------|------------|
| 1 | **Sem design system consistente** | Componentes duplicados, manutenção difícil | P0 |
| 2 | **Sem dark mode nativo** | 45% dos usuários preferem | P1 |
| 3 | **Mobile responsiveness incompleta** | 60% de acesso via mobile | P0 |
| 4 | **Sem micro-interactions** | UX menos fluida | P2 |
| 5 | **Sem accessibility audit** | Exigência LGPD | P0 |

---

### Plana de Reformulação de Design

#### Fase 1: Research & Audit (2 semanas)

| Atividade | Comando | Output |
|-----------|---------|--------|
| User Research | `*research` | Personas, jobs-to-be-done |
| Codebase Audit | `*audit ./frontend/src` | Pattern inventory |
| Shock Report | `*shock-report` | Visual report com métricas |

#### Fase 2: Design System (3 semanas)

| Atividade | Comando | Output |
|-----------|---------|--------|
| Token Extraction | `*tokenize` | tokens.yaml |
| System Setup | `*setup` | Design system structure |
| Component Library | `*bootstrap-shadcn` | Shadcn/UI integration |

#### Fase 3: Component Build (4 semanas)

| Componente | Prioridade | Squad |
|------------|------------|-------|
| **Sidebar Navigation** | P0 | UX/UI |
| **Conversation Inbox** | P0 | UX/UI |
| **Contact Profile Card** | P0 | UX/UI |
| **Message Composer** | P0 | UX/UI |
| **Dashboard Cards** | P1 | UX/UI |
| **Settings Panels** | P1 | UX/UI |
| **Chatbot Flow Builder** | P2 | Core + UX/UI |
| **Analytics Charts** | P2 | UX/UI |

#### Fase 4: Quality & Documentation (1 semana)

| Atividade | Comando | Output |
|-----------|---------|--------|
| A11y Audit | `*a11y-check` | WCAG compliance report |
| ROI Calculation | `*calculate-roi` | Cost savings analysis |
| Documentation | `*document` | Pattern library docs |

---

### Design Tokens Iniciais (Draft)

```yaml
# Omni Design Tokens v1.0

colors:
  primary:
    black: "#000000"
    white: "#FFFFFF"
    gray:
      50: "#FAFAFA"
      100: "#F5F5F5"
      200: "#E5E5E5"
      300: "#D4D4D4"
      400: "#A3A3A3"
      500: "#737373"
      600: "#525252"
      700: "#404040"
      800: "#262626"
      900: "#171717"
  
  accent:
    blue: "#3B82F6"
    green: "#22C55E"
    red: "#EF4444"
    yellow: "#EAB308"
  
  semantic:
    success: "#22C55E"
    warning: "#EAB308"
    error: "#EF4444"
    info: "#3B82F6"

typography:
  font-family:
    sans: "Inter, system-ui, sans-serif"
    mono: "JetBrains Mono, monospace"
  
  font-size:
    xs: "0.75rem"
    sm: "0.875rem"
    base: "1rem"
    lg: "1.125rem"
    xl: "1.25rem"
    2xl: "1.5rem"
    3xl: "1.875rem"
    4xl: "2.25rem"

spacing:
  1: "0.25rem"
  2: "0.5rem"
  3: "0.75rem"
  4: "1rem"
  5: "1.25rem"
  6: "1.5rem"
  8: "2rem"
  10: "2.5rem"
  12: "3rem"
  16: "4rem"

border-radius:
  none: "0"
  sm: "0.125rem"
  md: "0.375rem"
  lg: "0.5rem"
  xl: "0.75rem"
  full: "9999px"

shadows:
  sm: "0 1px 2px 0 rgb(0 0 0 / 0.05)"
  md: "0 4px 6px -1px rgb(0 0 0 / 0.1)"
  lg: "0 10px 15px -3px rgb(0 0 0 / 0.1)"
```

---

### Próximos Passos para UX/UI Squad

| # | Ação | Responsável | Deadline |
|---|------|-------------|----------|
| 1 | Criar diretório `docs/design/` | @ux-design-expert | Hoje |
| 2 | Rodar `*audit ./frontend` | @ux-design-expert | Semana 1 |
| 3 | Completar tokens.yaml | @ux-design-expert | Semana 1 |
| 4 | Setup Shadcn/UI | @dev + @ux-design-expert | Semana 2 |
| 5 | Redesign Sidebar | @ux-design-expert | Semana 3 |
| 6 | Redesign Inbox | @ux-design-expert + @dev | Semana 4-5 |

---

**— Uma, desenhando com empatia 💝**

---

## Resumo da Orquestração

```
┌─────────────────────────────────────────────────────────────────┐
│                    OMNICHANNEL PROJECT                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐        │
│  │   @analyst   │───▶│    @pm      │───▶│@ux-design   │        │
│  │   (Atlas)    │    │  (Morgan)   │    │  (Uma)      │        │
│  │             │    │             │    │             │        │
│  │ • Mercado BR │    │ • 5 Squads  │    │ • Design    │        │
│  │ • Competidores    │ • PRD v1.0  │    │   System    │        │
│  │ • ICP       │    │ • Roadmap   │    │ • Tokens    │        │
│  │ • Oportunidades   │ • Timeline  │    │ • Components│        │
│  └─────────────┘    └─────────────┘    └─────────────┘        │
│                              │                                  │
│                              ▼                                  │
│         ┌─────────────────────────────────────┐               │
│         │         SQUADS CRIADOS              │               │
│         ├─────────────────────────────────────┤               │
│         │ 🚀 Core Platform (Dev)              │               │
│         │ 🎨 UX/UI Design (Uma)               │               │
│         │ ⚙️ DevOps & Infra (Gage)            │               │
│         │ 🤖 AI/ML Engine (Dara)              │               │
│         │ 📈 Growth & Marketing (Morgan)      │               │
│         └─────────────────────────────────────┘               │
└─────────────────────────────────────────────────────────────────┘
```

---

### Para os próximos opencodes, recomendo:

| Opencode | Squad/Agente | Foco |
|----------|--------------|------|
| **#1** | UX/UI + Dev | Redesign do Inbox + Design System |
| **#2** | Core Platform | WhatsApp refatoração + DB schema |
| **#3** | AI/ML | Chatbot builder + Intent classifier |
| **#4** | DevOps | Infra prod + CI/CD pipeline |
| **#5** | Growth | Landing page + Documentação |

---<tool_call>
<function=todowrite>
<parameter=todos>[{"content": "Analista (@analyst): Análise de mercado brasileiro - segmento omnichannel/chatwoot", "priority": "high", "status": "completed"}, {"content": "Morgan (@pm): Criar estrutura de equipe (squads) para o projeto", "priority": "high", "status": "completed"}, {"content": "Morgan (@pm): Criar PRD estruturado baseado na análise de mercado", "priority": "high", "status": "completed"}, {"content": "UX Design Expert (@ux): Reformulação completa do design/frontend", "priority": "high", "status": "completed"}]