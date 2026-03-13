#!/bin/bash
set -e

echo "🚀 Preparando para fazer push das alterações..."

# Verificar se estamos no diretório correto
if [ ! -f "go.mod" ]; then
    echo "❌ Erro: Não está no diretório raiz do projeto bragdoc"
    exit 1
fi

# Verificar status do git
echo "📊 Status do git:"
git status --short

echo ""
echo "📋 Alterações implementadas:"
echo "1. ✅ Database: Migração 002_categories_positions.sql"
echo "2. ✅ Domain: Category e Position como objetos completos"
echo "3. ✅ Repositories: CategoryRepository e PositionRepository"
echo "4. ✅ Services: Validação adaptada para objetos"
echo "5. ✅ CLI: Comandos atualizados para UPPERCASE categories"
echo "6. ✅ Queries: Arquivos SQL gerados manualmente"
echo "7. ✅ Main: Inicialização de todos os repositórios"

echo ""
read -p "📝 Commit message (padrão: 'feat: add categories and positions as domain objects'): " commit_msg
commit_msg=${commit_msg:-"feat: add categories and positions as domain objects"}

echo ""
echo "🧹 Verificando arquivos..."
# Listar arquivos modificados/criados
find . -name "*.go" -type f -newer /tmp/dummy 2>/dev/null || true
find . -name "*.sql" -type f -newer /tmp/dummy 2>/dev/null || true
find . -name "*.md" -type f -newer /tmp/dummy 2>/dev/null || true

echo ""
read -p "✅ Confirmar commit e push? (s/N): " confirm
if [[ ! "$confirm" =~ ^[Ss]$ ]]; then
    echo "❌ Cancelado"
    exit 0
fi

echo ""
echo "📦 Fazendo commit..."
git add .
git commit -m "$commit_msg"

echo ""
echo "🚀 Fazendo push..."
git push

echo ""
echo "🎉 Push concluído com sucesso!"
echo ""
echo "📋 Resumo das alterações:"
echo "   - Database schema atualizado (categories, positions)"
echo "   - Domain models refatorados (objetos completos, não IDs)"
echo "   - Repositories implementados (Category, Position)"
echo "   - CLI atualizado (categorias UPPERCASE)"
echo "   - Migração de dados existentes incluída"
echo ""
echo "⚠️  PRÓXIMOS PASSOS:"
echo "   1. Aplicar migração no banco de produção:"
echo "      sqlite3 ~/.bragdoc/bragdoc.db < internal/database/migrations/002_categories_positions.sql"
echo "   2. Testar com dados reais"
echo "   3. Atualizar outros testes se necessário"