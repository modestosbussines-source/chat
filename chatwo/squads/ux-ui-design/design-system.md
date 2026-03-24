# OMNI Design System

> Visual moderno, minimalista para o público brasileiro. Autoridade OMNICall.

## Visão Geral

O OMNI Design System é um sistema de design completo, moderno e minimalista, criado especificamente para o mercado brasileiro. Ele combina elegância funcional com a vibrância da estética brasileira contemporânea.

### Princípios Fundamentais

1. **Minimalismo Funcional** - Cada elemento tem um propósito claro
2. **Acessibilidade Primeiro** - WCAG 2.1 AA compliance
3. **Mobile-First** - Design responsivo com foco em dispositivos móveis
4. **Performance** - Carregamento rápido e animações suaves
5. **Consistência** - Linguagem visual unificada em todo o produto

---

## Sistema de Cores

### Paleta Principal

#### Primary - Indigo/Violet
A cor principal do OMNI, transmitindo profissionalismo e inovação.

| Token | Light Mode | Dark Mode |
|-------|------------|-----------|
| `--primary` | `hsl(252, 70%, 55%)` | `hsl(252, 80%, 72%)` |
| `--primary-foreground` | `hsl(0, 0%, 100%)` | `hsl(220, 20%, 8%)` |

#### Secondary - Teal/Cyan
Cor secundária para acentos e elementos interativos.

| Token | Light Mode | Dark Mode |
|-------|------------|-----------|
| `--secondary` | `hsl(199, 80%, 50%)` | `hsl(199, 100%, 65%)` |
| `--secondary-foreground` | `hsl(0, 0%, 100%)` | `hsl(220, 20%, 8%)` |

### Acentos Brasileiros

Cores vibrantes inspiradas na estética brasileira:

| Nome | Token | Cor |
|------|-------|-----|
| Coral | `--coral` | `hsl(12, 100%, 62%)` |
| Esmeralda | `--teal` | `hsl(168, 80%, 50%)` |
| Âmbar | `--amber` | `hsl(38, 92%, 58%)` |
| Rosa | `--pink` | `hsl(330, 100%, 68%)` |

### Cores Semânticas

| Estado | Light | Dark | Uso |
|--------|-------|------|-----|
| Success | `hsl(142, 76%, 36%)` | `hsl(142, 76%, 50%)` | Sucesso, conclusão |
| Warning | `hsl(38, 92%, 50%)` | `hsl(38, 92%, 60%)` | Avisos, atenção |
| Error | `hsl(0, 84%, 60%)` | `hsl(0, 84%, 65%)` | Erros, destructives |
| Info | `hsl(199, 100%, 50%)` | `hsl(199, 100%, 65%)` | Informações |

---

## Tipografia

### Famílias

```css
font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
font-family-mono: 'JetBrains Mono', monospace;
```

### Escala

| Name | Size | Weight | Line Height |
|------|------|--------|-------------|
| `display-large` | 3rem / 48px | 700 | 1.1 |
| `display-medium` | 2.25rem / 36px | 700 | 1.15 |
| `headline-large` | 1.875rem / 30px | 600 | 1.2 |
| `headline-medium` | 1.5rem / 24px | 600 | 1.25 |
| `title-large` | 1.25rem / 20px | 600 | 1.3 |
| `body-large` | 1rem / 16px | 400 | 1.5 |
| `body-medium` | 0.875rem / 14px | 400 | 1.5 |
| `label-large` | 0.875rem / 14px | 500 | 1.4 |

---

## Espaçamento

Sistema baseado em 4px:

```css
/* Espaçamentos comuns */
--spacing-1: 0.25rem;  /* 4px */
--spacing-2: 0.5rem;   /* 8px */
--spacing-3: 0.75rem;  /* 12px */
--spacing-4: 1rem;     /* 16px */
--spacing-6: 1.5rem;   /* 24px */
--spacing-8: 2rem;     /* 32px */
```

---

## Bordas Arredondadas

```css
--radius-none: 0;
--radius-xs: 0.25rem;   /* 4px */
--radius-sm: 0.375rem;  /* 6px */
--radius-md: 0.5rem;    /* 8px */
--radius-lg: 0.75rem;   /* 12px */
--radius-xl: 1rem;      /* 16px */
--radius-2xl: 1.25rem;  /* 20px */
--radius-full: 9999px;
```

---

## Sombras

Sistema de elevação moderno:

```css
--shadow-sm: 0 1px 3px 0 hsl(220 20% 8% / 0.06);
--shadow-md: 0 4px 6px -1px hsl(220 20% 8% / 0.08);
--shadow-lg: 0 10px 15px -3px hsl(220 20% 8% / 0.1);
--shadow-xl: 0 20px 25px -5px hsl(220 20% 8% / 0.12);
```

---

## Glassmorphism

Efeito de vidro sutil para cards e overlays:

```css
.glass {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(12px);
  border: 1px solid rgba(0, 0, 0, 0.08);
}
```

---

## Componentes do Inbox

### Estrutura

```
inbox/
├── ConversationList.vue     # Lista de conversas
├── ConversationItem.vue     # Item individual
├── ContactCard.vue          # Card de contato
├── InboxHeader.vue          # Header do inbox
├── MessageThread.vue        # Thread de mensagens
├── MessageInput.vue         # Input de mensagens
└── index.ts                 # Exportações
```

### ConversationList

Lista principal de conversas com filtros e busca.

```vue
<ConversationList
  :conversations="conversations"
  :active-conversation-id="activeId"
  filter="unread"
  @select="handleSelect"
  @filter-change="handleFilterChange"
/>
```

**Props:**
- `conversations: Conversation[]` - Lista de conversas
- `activeConversationId?: string` - ID da conversa ativa
- `loading?: boolean` - Estado de carregamento
- `filter?: 'all' | 'unread' | 'assigned' | 'unassigned'` - Filtro ativo

### ConversationItem

Item individual na lista de conversas.

```vue
<ConversationItem
  :conversation="conversation"
  :active="true"
  @click="handleClick"
/>
```

### ContactCard

Card para exibir informações de contato.

```vue
<ContactCard
  :contact="contact"
  :compact="false"
  @edit="handleEdit"
  @message="handleMessage"
  @call="handleCall"
/>
```

### InboxHeader

Header do inbox com ações.

```vue
<InboxHeader
  :conversation="activeConversation"
  @call="handleCall"
  @video-call="handleVideoCall"
  @assign="handleAssign"
  @star="handleStar"
  @archive="handleArchive"
  @delete="handleDelete"
/>
```

### MessageThread

Thread de mensagens com agrupamento por data.

```vue
<MessageThread
  :messages="messages"
  :contact-name="contact.name"
  :contact-avatar="contact.avatar"
  @reply="handleReply"
  @react="handleReact"
  @delete="handleDelete"
/>
```

### MessageInput

Campo de entrada de mensagens com attachments.

```vue
<MessageInput
  placeholder="Digite sua mensagem..."
  :reply-to="replyToMessage"
  @send="handleSend"
  @typing="handleTyping"
  @attach="handleAttach"
  @cancel-reply="clearReply"
/>
```

---

## Animações

### Transições

```css
--transition-all: all 200ms cubic-bezier(0.4, 0, 0.2, 1);
--transition-colors: color, background-color, border-color 200ms;
--transition-opacity: opacity 200ms;
--transition-shadow: box-shadow 200ms;
--transition-transform: transform 200ms;
```

### Keyframes

```css
@keyframes fade-in {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

@keyframes slide-in-right {
  from { opacity: 0; transform: translateX(10px); }
  to { opacity: 1; transform: translateX(0); }
}

@keyframes pulse-ring {
  0%, 100% { transform: scale(0.95); opacity: 0.7; }
  50% { transform: scale(1.05); opacity: 1; }
}
```

---

## Acessibilidade

### Focus Ring

```css
.focus-ring {
  outline: none;
  ring: 2px solid hsl(var(--ring));
  ring-offset: 2px;
  ring-offset-color: hsl(var(--background));
}
```

### Touch Targets

Mínimo de 44px para dispositivos touch:

```css
@media (pointer: coarse) {
  button, [role="button"], a {
    min-height: 44px;
    min-width: 44px;
  }
}
```

### Reduced Motion

Respeita preferências do usuário:

```css
@media (prefers-reduced-motion: reduce) {
  *, *::before, *::after {
    animation-duration: 0.01ms !important;
    transition-duration: 0.01ms !important;
  }
}
```

---

## Tema Escuro

O tema escuro é ativado adicionando a classe `.dark` ao elemento `<html>`:

```html
<html class="dark">
  <!-- Conteúdo -->
</html>
```

Ou usando o toggle de tema:

```vue
<script setup>
import { useTheme } from '@/composables/useTheme'
const { toggleTheme, isDark } = useTheme()
</script>
```

---

## Arquivos do Sistema

| Arquivo | Descrição |
|---------|-----------|
| `tokens.yaml` | Design tokens em YAML |
| `index.css` | CSS custom properties |
| `tailwind.config.cjs` | Configuração do Tailwind |
| `design-system.md` | Esta documentação |

---

## Resumo das Mudanças

### Antes (Sistema Antigo)
- Paleta monocromática (preto e branco)
- Glassmorphism intenso
- Estilo Linear/Vercel
- Foco em developers

### Agora (OMNI Design System)
- Paleta colorida com personalidade brasileira
- Glassmorphism sutil
- Visual moderno e minimalista
- Foco em UX para mercado brasileiro
- Maior acessibilidade
- Tema claro como padrão

---

## Mantendo o Design System

### Ao Adicionar Novos Cores

1. Adicione ao `tokens.yaml`
2. Atualize `index.css` com as variáveis CSS
3. Adicione ao `tailwind.config.cjs`
4. Documente nesta página

### Ao Criar Novos Componentes

1. Use os tokens existentes
2. Siga as convenções de nomenclatura
3. Inclua suporte a tema claro/escuro
4. Adicione estados de focus e hover
5. Teste acessibilidade

---

## Recursos

- [Design Tokens YAML](./tokens.yaml)
- [Componentes de Inbox](../frontend/src/components/inbox/)
- [Configuração do Tailwind](../frontend/tailwind.config.cjs)
- [CSS Base](../frontend/src/assets/index.css)

---

## Contato

**UX/UI Design Squad**  
Responsável pelo design system OMNI

- Lead: ux-design-expert (Uma)
- Sprint: 2 semanas
- Quality Gates: Design Review, A11y Audit (AA), Performance
