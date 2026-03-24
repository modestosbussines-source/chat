# 📢 COMUNICADO OFICIAL À EQUIPE

**Data:** 23/03/2026  
**De:** Igor (Stakeholder)  
**Para:** Todos os Squads  
**Urgência:** 🔴 ALTA

---

## MENSAGEM

```
NÃO SOMOS WHATOMATE. SOMOS OMNI.

O Whatomate é um projeto open-source de terceiros.
Nós usamos o código como BASE para desenvolver o OMNI.

O OMNI é:
- PRODUTO PRIVADO e PROPRIETÁRIO
- 100% VENDÁVEL
- Para o mercado brasileiro
- Sem vínculo com projeto open-source

REMOVA TODAS as 184 referências a "whatomate" AGORA.
Substitua por "Omni" ou remova.

Isto é OBRIGATÓRIO e IMEDIATO.
```

---

## POR SQUAD

### Terminal 1 — @ux-design-expert (Uma)
```
Sua tarefa: Remover referências a "whatomate" no frontend.

Arquivos para verificar:
- frontend/src/**/*.ts
- frontend/src/components/brand/
- frontend/src/i18n/locales/

Substituir:
- "whatomate" → "omni"
- "Whatomate" → "Omni"

Após: Atualize docs/STATUS.md com ✅
```

### Terminal 2 — @dev (Dex)
```
Sua tarefa: Remover referências a "whatomate" no Go backend.

Arquivos para verificar:
- cmd/whatomate/ → Renomear para cmd/omni/
- go.mod: module path
- internal/**/*.go: imports e strings

Substituir:
- "github.com/shridarpatil/whatomate" → "github.com/omni/platform"
- "Whatomate" → "Omni"
- "whatomate" → "omni"

Testar: go build -o omni ./cmd/omni

Após: Atualize docs/STATUS.md com ✅
```

### Terminal 3 — @devops (Gage)
```
Sua tarefa: Remover referências a "whatomate" na infraestrutura.

Arquivos para verificar:
- docker/Dockerfile
- docker/docker-compose.yml
- docker/.env.example
- .goreleaser.yml
- .github/workflows/*.yml

Substituir:
- "whatomate" → "omni"
- "whatomate_test" → "omni_test"
- Binário: ./whatomate → ./omni

Testar: docker build -t omni .

Após: Atualize docs/STATUS.md com ✅
```

### Terminal 4 — @pm (Morgan)
```
Sua tarefa: Remover referências a "whatomate" na documentação.

Arquivos para verificar:
- docs/src/pages/index.astro
- docs/src/layouts/*.astro
- docs/stories/*.md (5 arquivos)
- docs/prd-documentation-landing.md
- CHANGELOG_OMNI.md
- SECURITY_REPORT.md

Substituir:
- "Whatomate" → "Omni"
- "whatomate" → "omni"
- "shridarpatil/whatomate" → remover
- URLs antigas → omni.com.br

Após: Atualize docs/STATUS.md com ✅
```

### Terminal 5 — @data-engineer (Dara)
```
Sua tarefa: Verificar e remover referências a "whatomate" em handlers de IA.

Arquivos para verificar:
- internal/handlers/chatbot*.go
- internal/handlers/flows*.go
- internal/handlers/custom_actions*.go

Substituir:
- "whatomate" → "omni" (se houver)

Após: Atualize docs/STATUS.md com ✅
```

---

## VALIDAÇÃO FINAL

Após todas as remoções, execute em QUALQUER terminal:

```bash
# Deve retornar 0 resultados
grep -ri "whatomate" --include="*.go" --include="*.vue" --include="*.ts" --include="*.md" --include="*.yml" --include="*.astro" .
```

---

## CHECKLIST DE CONCLUSÃO FASE 1

- [ ] 0 ocorrências de "whatomate" em todo o projeto
- [ ] 0 ocorrências de "shridarpatil" em todo o projeto
- [ ] Build Go funcional
- [ ] Docker build funcional
- [ ] Todos os testes passando
- [ ] Docs mostrando "Omni"
- [ ] Favicon/og-image atualizados

---

**Quando terminar, responda em cada terminal:**

```
✅ [SQUAD NAME] — Whatomate removido. Próximos passos: [listar]
```

---

*Igor*
