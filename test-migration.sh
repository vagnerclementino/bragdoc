#!/bin/bash
set -e

echo "🧪 Testando migração de categories e positions..."

# Criar banco de dados de teste
TEST_DB="test_bragdoc.db"
rm -f "$TEST_DB"

echo "1. Criando banco de dados inicial..."
sqlite3 "$TEST_DB" <<EOF
-- Tabela inicial (001_initial.sql)
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    job_title TEXT,
    company TEXT,
    locale TEXT NOT NULL DEFAULT 'en-US',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE brags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    category INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);

CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    owner_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES users(id),
    UNIQUE(owner_id, name)
);

CREATE TABLE brag_tags (
    brag_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (brag_id, tag_id),
    FOREIGN KEY (brag_id) REFERENCES brags(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);

-- Inserir dados de teste
INSERT INTO users (name, email, job_title, company) VALUES 
('Test User', 'test@example.com', 'Developer', 'Test Corp');

INSERT INTO brags (owner_id, title, description, category) VALUES
(1, 'Test Project', 'A test project brag', 2),  -- CategoryProject (2)
(1, 'Test Achievement', 'A test achievement brag', 3);  -- CategoryAchievement (3)
EOF

echo "✅ Banco inicial criado com 2 brags"

echo "2. Aplicando migração 002_categories_positions.sql..."
sqlite3 "$TEST_DB" < internal/database/migrations/002_categories_positions.sql

echo "3. Verificando estrutura..."
sqlite3 "$TEST_DB" <<EOF
.headers on
.mode column

echo "=== Tabelas existentes ==="
SELECT name FROM sqlite_master WHERE type='table' ORDER BY name;

echo ""
echo "=== Categorias inseridas ==="
SELECT id, name, description FROM categories ORDER BY id;

echo ""
echo "=== Brags migrados ==="
SELECT b.id, b.title, c.name as category_name, p.title as position_title 
FROM brags b 
LEFT JOIN categories c ON b.category_id = c.id 
LEFT JOIN positions p ON b.position_id = p.id;

echo ""
echo "=== Contagem de brags por categoria ==="
SELECT c.name, COUNT(b.id) as brag_count 
FROM categories c 
LEFT JOIN brags b ON c.id = b.category_id 
GROUP BY c.id 
ORDER BY c.name;
EOF

echo ""
echo "🧪 Teste concluído!"
echo "📊 Resultado esperado:"
echo "   - 6 categorias (UNKNOWN, PROJECT, ACHIEVEMENT, SKILL, LEADERSHIP, INNOVATION)"
echo "   - 2 brags migrados (um para PROJECT, outro para ACHIEVEMENT)"
echo "   - Tabela positions vazia (sem dados iniciais)"

rm -f "$TEST_DB"
echo "🧹 Banco de teste removido"