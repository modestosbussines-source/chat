#!/bin/bash
# ============================================================
# SCRIPT: Remoção Completa de "Whatomate" → "Omni"
# Autor: Igor (Stakeholder)
# Data: 23/03/2026
# Uso: bash scripts/remove-whatomate.sh [frontend|backend|docs|infra|all]
# ============================================================

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

PROJECT_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_ROOT"

echo -e "${RED}🚨 INICIANDO REMOÇÃO DE WHATOMATE → OMNI${NC}"
echo "============================================"

count_before=$(grep -ri "whatomate" --include="*.go" --include="*.vue" --include="*.ts" --include="*.md" --include="*.yml" --include="*.astro" --include="*.json" --include="*.toml" --include="*.html" --include="*.css" . 2>/dev/null | wc -l || echo "0")
echo -e "${YELLOW}Referências encontradas ANTES: ${count_before}${NC}"
echo ""

TARGET="${1:-all}"

# ============================================================
# FUNÇÕES
# ============================================================

replace_in_files() {
    local pattern="$1"
    local replacement="$2"
    local filetypes="$3"
    
    for ext in $filetypes; do
        find . -type f -name "$ext" -not -path "./.git/*" -not -path "./node_modules/*" 2>/dev/null | while read file; do
            if grep -q "$pattern" "$file" 2>/dev/null; then
                sed -i "s/${pattern}/${replacement}/g" "$file" 2>/dev/null || true
                echo "  ✓ $file"
            fi
        done
    done
}

replace_case_insensitive() {
    local pattern="$1"
    local replacement="$2"
    local filetypes="$3"
    
    for ext in $filetypes; do
        find . -type f -name "$ext" -not -path "./.git/*" -not -path "./node_modules/*" 2>/dev/null | while read file; do
            if grep -qi "$pattern" "$file" 2>/dev/null; then
                sed -i "s/${pattern}/${replacement}/gi" "$file" 2>/dev/null || true
                echo "  ✓ $file"
            fi
        done
    done
}

# ============================================================
# FRONTEND (UX/UI Squad)
# ============================================================

remove_frontend() {
    echo -e "\n${GREEN}🎨 FRONTEND - UX/UI Squad${NC}"
    echo "----------------------------"
    
    echo "Substituindo 'whatomate' → 'omni'..."
    replace_in_files "whatomate" "omni" "*.ts *.vue *.json *.css *.html"
    
    echo "Substituindo 'Whatomate' → 'Omni'..."
    replace_in_files "Whatomate" "Omni" "*.ts *.vue *.json *.css *.html"
    
    echo "Substituindo 'WHATOMATE' → 'OMNI'..."
    replace_in_files "WHATOMATE" "OMNI" "*.ts *.vue *.json *.css *.html"
    
    echo -e "${GREEN}✓ Frontend processado${NC}"
}

# ============================================================
# BACKEND (Core Squad)
# ============================================================

remove_backend() {
    echo -e "\n${GREEN}🚀 BACKEND - Core Squad${NC}"
    echo "----------------------------"
    
    # Renomear diretório
    if [ -d "cmd/whatomate" ]; then
        echo "Renomeando cmd/whatomate → cmd/omni..."
        mv cmd/whatomate cmd/omni
        echo "  ✓ Diretório renomeado"
    fi
    
    # Atualizar go.mod
    if [ -f "go.mod" ]; then
        echo "Atualizando go.mod..."
        sed -i 's|github.com/shridarpatil/whatomate|github.com/omni/platform|g' go.mod 2>/dev/null || true
        echo "  ✓ go.mod atualizado"
    fi
    
    # Substituir em arquivos Go
    echo "Substituindo imports Go..."
    find . -name "*.go" -not -path "./.git/*" -not -path "./node_modules/*" 2>/dev/null | while read file; do
        if grep -q "shridarpatil/whatomate" "$file" 2>/dev/null; then
            sed -i 's|github.com/shridarpatil/whatomate|github.com/omni/platform|g' "$file" 2>/dev/null || true
            echo "  ✓ $file (import)"
        fi
        if grep -q "whatomate" "$file" 2>/dev/null; then
            sed -i 's/whatomate/omni/g' "$file" 2>/dev/null || true
            echo "  ✓ $file (string)"
        fi
        if grep -q "Whatomate" "$file" 2>/dev/null; then
            sed -i 's/Whatomate/Omni/g' "$file" 2>/dev/null || true
            echo "  ✓ $file (Whatomate)"
        fi
    done
    
    echo -e "${GREEN}✓ Backend processado${NC}"
}

# ============================================================
# DOCS (Growth Squad)
# ============================================================

remove_docs() {
    echo -e "\n${GREEN}📈 DOCS - Growth Squad${NC}"
    echo "----------------------------"
    
    # Docs markdown/mdx
    echo "Processando markdown/mdx..."
    find docs -type f \( -name "*.md" -o -name "*.mdx" \) 2>/dev/null | while read file; do
        if grep -qi "whatomate\|shridarpatil" "$file" 2>/dev/null; then
            sed -i 's/Whatomate/Omni/g' "$file" 2>/dev/null || true
            sed -i 's/whatomate/omni/g' "$file" 2>/dev/null || true
            sed -i 's|shridarpatil/whatomate||g' "$file" 2>/dev/null || true
            sed -i 's|https://shridarpatil.github.io/whatomate|https://omni.com.br|g' "$file" 2>/dev/null || true
            sed -i 's|https://github.com/shridarpatil/whatomate|https://github.com/omni/platform|g' "$file" 2>/dev/null || true
            echo "  ✓ $file"
        fi
    done
    
    # Docs Astro
    echo "Processando arquivos Astro..."
    find docs -type f -name "*.astro" 2>/dev/null | while read file; do
        if grep -qi "whatomate\|shridarpatil" "$file" 2>/dev/null; then
            sed -i 's/Whatomate/Omni/g' "$file" 2>/dev/null || true
            sed -i 's/whatomate/omni/g' "$file" 2>/dev/null || true
            sed -i 's|shridarpatil/whatomate||g' "$file" 2>/dev/null || true
            sed -i 's|/whatomate/|/omni/|g' "$file" 2>/dev/null || true
            echo "  ✓ $file"
        fi
    done
    
    # Outros docs
    for file in CHANGELOG_OMNI.md SECURITY_REPORT.md docs/prd-documentation-landing.md; do
        if [ -f "$file" ] && grep -qi "whatomate" "$file" 2>/dev/null; then
            sed -i 's/Whatomate/Omni/g' "$file" 2>/dev/null || true
            sed -i 's/whatomate/omni/g' "$file" 2>/dev/null || true
            sed -i 's/Whatomate/Omni/g' "$file" 2>/dev/null || true
            echo "  ✓ $file"
        fi
    done
    
    echo -e "${GREEN}✓ Docs processados${NC}"
}

# ============================================================
# INFRA (DevOps Squad)
# ============================================================

remove_infra() {
    echo -e "\n${GREEN}⚙️ INFRA - DevOps Squad${NC}"
    echo "----------------------------"
    
    # Docker
    echo "Processando Docker..."
    for file in docker/Dockerfile docker/docker-compose.yml docker/.env.example docker-compose.yml; do
        if [ -f "$file" ] && grep -q "whatomate" "$file" 2>/dev/null; then
            sed -i 's/whatomate/omni/g' "$file" 2>/dev/null || true
            sed -i 's/Whatomate/Omni/g' "$file" 2>/dev/null || true
            echo "  ✓ $file"
        fi
    done
    
    # GitHub Actions
    echo "Processando GitHub Actions..."
    find .github -type f -name "*.yml" 2>/dev/null | while read file; do
        if grep -q "whatomate" "$file" 2>/dev/null; then
            sed -i 's/whatomate/omni/g' "$file" 2>/dev/null || true
            sed -i 's/Whatomate/Omni/g' "$file" 2>/dev/null || true
            echo "  ✓ $file"
        fi
    done
    
    # goreleaser
    if [ -f ".goreleaser.yml" ] && grep -q "whatomate" ".goreleaser.yml" 2>/dev/null; then
        sed -i 's/whatomate/omni/g' .goreleaser.yml 2>/dev/null || true
        sed -i 's|cmd/whatomate|cmd/omni|g' .goreleaser.yml 2>/dev/null || true
        echo "  ✓ .goreleaser.yml"
    fi
    
    # gitignore
    if [ -f ".gitignore" ] && grep -q "whatomate" ".gitignore" 2>/dev/null; then
        sed -i 's/whatomate/omni/g' .gitignore 2>/dev/null || true
        echo "  ✓ .gitignore"
    fi
    
    echo -e "${GREEN}✓ Infra processada${NC}"
}

# ============================================================
# EXECUÇÃO
# ============================================================

case "$TARGET" in
    frontend)
        remove_frontend
        ;;
    backend)
        remove_backend
        ;;
    docs)
        remove_docs
        ;;
    infra)
        remove_infra
        ;;
    all)
        remove_frontend
        remove_backend
        remove_docs
        remove_infra
        ;;
    *)
        echo "Uso: $0 [frontend|backend|docs|infra|all]"
        exit 1
        ;;
esac

# ============================================================
# VALIDAÇÃO
# ============================================================

echo ""
echo -e "${YELLOW}============================================${NC}"
echo -e "${YELLOW}🔍 VALIDAÇÃO FINAL${NC}"
echo -e "${YELLOW}============================================${NC}"

count_after=$(grep -ri "whatomate" --include="*.go" --include="*.vue" --include="*.ts" --include="*.md" --include="*.yml" --include="*.astro" --include="*.json" --include="*.toml" . 2>/dev/null | grep -v ".git/" | grep -v "node_modules/" | wc -l || echo "0")

echo ""
echo "Referências ANTES: ${count_before}"
echo "Referências DEPOIS: ${count_after}"

if [ "$count_after" -eq 0 ]; then
    echo -e "\n${GREEN}✅ SUCESSO! Todas as referências a 'whatomate' foram removidas.${NC}"
else
    echo -e "\n${RED}⚠️  ATENÇÃO! Ainda restam ${count_after} referências.${NC}"
    echo "Execute: grep -ri \"whatomate\" --include=\"*.go\" --include=\"*.vue\" --include=\"*.ts\" --include=\"*.md\" --include=\"*.yml\" --include=\"*.astro\" ."
fi

echo ""
echo -e "${GREEN}🎉 Processo concluído!${NC}"
