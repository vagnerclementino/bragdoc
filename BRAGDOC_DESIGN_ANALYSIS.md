# Bragdoc - Design Document & Data Model Analysis

## Visão Geral
**Bragdoc** é um sistema para registro e documentação de conquistas profissionais ("brags"). Permite que usuários registrem projetos, habilidades, liderança, inovações e outras conquistas, organizando-as com tags e categorias para posterior geração de documentos profissionais.

**Branch analisada:** `beta_version`
**Data da análise:** 23 de fevereiro de 2026
**Tecnologia principal:** Go + SQLite

---

## 1. Modelo de Dados Atual

### 1.1 Diagrama ER (Entidade-Relacionamento)

```mermaid
erDiagram
    USERS ||--o{ BRAGS : "owns (1:N)"
    USERS ||--o{ TAGS : "creates (1:N)"
    BRAGS ||--o{ BRAG_TAGS : "has (1:N)"
    TAGS ||--o{ BRAG_TAGS : "assigned_to (1:N)"
    
    USERS {
        int64 id PK "AUTOINCREMENT"
        string name "NOT NULL"
        string email "UNIQUE, NOT NULL"
        string job_title "NULLABLE"
        string company "NULLABLE"
        string locale "DEFAULT 'en-US'"
        datetime created_at "DEFAULT CURRENT_TIMESTAMP"
        datetime updated_at "NULLABLE"
    }
    
    BRAGS {
        int64 id PK "AUTOINCREMENT"
        int64 owner_id FK "NOT NULL"
        string title "NOT NULL"
        string description "NOT NULL"
        int64 category "NOT NULL"
        datetime created_at "DEFAULT CURRENT_TIMESTAMP"
        datetime updated_at "NULLABLE"
    }
    
    TAGS {
        int64 id PK "AUTOINCREMENT"
        string name "NOT NULL, CHECK(length 2-20)"
        int64 owner_id FK "NOT NULL"
        datetime created_at "DEFAULT CURRENT_TIMESTAMP"
        UNIQUE(name, owner_id)
    }
    
    BRAG_TAGS {
        int64 brag_id PK,FK "NOT NULL"
        int64 tag_id PK,FK "NOT NULL"
        PRIMARY KEY(brag_id, tag_id)
        FOREIGN KEY (brag_id) REFERENCES brags(id) ON DELETE CASCADE
        FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
    }
```

### 1.2 Esquema de Banco de Dados

#### **Tabela: users**
```sql
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    job_title TEXT,
    company TEXT,
    locale TEXT NOT NULL DEFAULT 'en-US',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
);
```

#### **Tabela: brags**
```sql
CREATE TABLE brags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    category INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    FOREIGN KEY (owner_id) REFERENCES users(id)
);
```

#### **Tabela: tags**
```sql
CREATE TABLE tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL CHECK(length(name) >= 2 AND length(name) <= 20),
    owner_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(name, owner_id),
    FOREIGN KEY (owner_id) REFERENCES users(id)
);
```

#### **Tabela: brag_tags** (tabela de junção)
```sql
CREATE TABLE brag_tags (
    brag_id INTEGER NOT NULL,
    tag_id INTEGER NOT NULL,
    PRIMARY KEY (brag_id, tag_id),
    FOREIGN KEY (brag_id) REFERENCES brags(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
);
```

### 1.3 Índices
```sql
CREATE INDEX idx_brag_tags_brag_id ON brag_tags(brag_id);
CREATE INDEX idx_brag_tags_tag_id ON brag_tags(tag_id);
CREATE INDEX idx_tags_owner_id ON tags(owner_id);
CREATE INDEX idx_tags_name ON tags(name);
CREATE INDEX idx_brags_owner_id ON brags(owner_id);
CREATE INDEX idx_brags_category ON brags(category);
```

### 1.4 Modelo de Domínio (Go)

#### **User** (`internal/domain/user.go`)
```go
type User struct {
    ID        int64
    Name      string
    Email     string
    JobTitle  string
    Company   string
    Locale    Locale  // language-COUNTRY format
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

#### **Brag** (`internal/domain/brag.go`)
```go
type Brag struct {
    ID          int64
    Owner       User
    Title       string
    Description string
    Category    Category  // Enum: 1-6
    Tags        []*Tag
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

// Categorias disponíveis (enum)
const (
    CategoryUnknown Category = iota + 1  // 1
    CategoryProject                      // 2
    CategoryAchievement                  // 3
    CategorySkill                        // 4
    CategoryLeadership                   // 5
    CategoryInnovation                   // 6
)
```

#### **Tag** (`internal/domain/tag.go`)
```go
type Tag struct {
    ID        int64
    Name      string
    OwnerID   int64
    CreatedAt time.Time
}
```

#### **Document** (`internal/domain/document.go`)
```go
type Document struct {
    Format   DocumentFormat  // markdown, pdf, docx
    Template string
    Content  []byte
    Metadata DocumentMetadata
}

type DocumentMetadata struct {
    Title       string
    Author      string
    BragCount   int
    Categories  []string
    Tags        []string
    GeneratedAt string
}
```

---

## 2. Arquitetura do Sistema

### 2.1 Estrutura de Diretórios
```
bragdoc/
├── cmd/                    # Ponto de entrada da aplicação
├── config/                 # Configurações
├── docs/                   # Documentação
├── internal/
│   ├── database/          # Camada de persistência
│   │   ├── migrations/    # Migrações SQL
│   │   ├── queries/       # SQLc generated code
│   │   └── *.go           # Database implementations
│   ├── domain/            # Entidades de domínio
│   └── service/           # Lógica de negócio
├── testdata/              # Dados de teste
└── *.md                   # Documentação
```

### 2.2 Padrões de Design Identificados

#### **Clean Architecture**
- **Domain:** Entidades puras sem dependências externas
- **Database:** Implementação de persistência
- **Service:** Lógica de negócio

#### **Repository Pattern**
- Separação entre domínio e persistência
- SQLc para geração de código SQL type-safe

#### **Domain-Driven Design (DDD)**
- Entidades com identidade (User, Brag, Tag)
- Value Objects (Category, Locale, DocumentFormat)
- Agregados (Brag com Tags)

### 2.3 Fluxo de Dados
```
1. Usuário cria conta → tabela `users`
2. Usuário cria brag → tabela `brags`
3. Usuário cria tags → tabela `tags`
4. Associa tags a brags → tabela `brag_tags`
5. Gera documento → domínio `Document`
```

---

## 3. Análise de Pontos Fortes ✅

### 3.1 Design de Banco de Dados
- **Normalização adequada:** 3NF com tabelas bem separadas
- **Constraints robustas:** UNIQUE, CHECK, FOREIGN KEY
- **Cascading delete:** ON DELETE CASCADE em relações N:M
- **Índices otimizados:** Para queries frequentes

### 3.2 Arquitetura de Software
- **Separação clara:** Domínio × Persistência × Serviço
- **Imutabilidade:** Estruturas de domínio sem side effects
- **Type safety:** SQLc gera código type-safe do SQL
- **Testabilidade:** Domínio puro facilita testes unitários

### 3.3 Modelo de Domínio
- **Enums type-safe:** Category, DocumentFormat
- **Validações no banco:** CHECK constraints
- **Relacionamentos explícitos:** 1:N e N:M bem definidos
- **Internacionalização:** Suporte a locale desde o início

### 3.4 Performance
- **Índices estratégicos:** Otimizados para buscas comuns
- **SQLite leve:** Bom para MVP e small-scale
- **Caching natural:** Tags por usuário podem ser cacheadas

---

## 4. Análise de Pontos de Melhoria ⚠️

### 4.1 Limitações do Modelo Atual

#### **1. Categorias Fixas (Hardcoded)**
**Problema:** Categorias são enum hardcoded (1-6)
**Impacto:** Usuários não podem criar categorias personalizadas
**Exemplo:** Não pode adicionar "Publicação", "Certificação", etc.

#### **2. Falta de Soft Delete**
**Problema:** Exclusões são físicas (DELETE)
**Impacto:** 
- Perda irreversível de dados históricos
- Impossibilidade de recuperação acidental
- Dificuldade em analytics históricos

#### **3. Ausência de Versionamento**
**Problema:** Não há histórico de alterações em brags
**Impacto:**
- Não sabe o que mudou em uma brag ao longo do tempo
- Impossibilidade de reverter alterações
- Sem audit trail

#### **4. Metadados Limitados**
**Problema:** Campos fixos (title, description, category)
**Impacto:** Dificuldade em adicionar campos específicos
**Exemplo:** 
- Projetos: datas, orçamento, equipe
- Conquistas: métricas, impacto
- Habilidades: nível, anos de experiência

#### **5. Sistema de Permissões Básico**
**Problema:** Todos os brags são privados ao owner
**Impacto:** 
- Não permite compartilhamento seletivo
- Não suporta colaboração em equipe
- Limita casos de uso corporativos

#### **6. Internacionalização Parcial**
**Problema:** Apenas locale do usuário, não do conteúdo
**Impacto:** 
- Não suporta brags multilíngues
- Usuário precisa escolher um idioma principal

#### **7. Relacionamentos entre Brags**
**Problema:** Brags são entidades isoladas
**Impacto:** 
- Não pode criar sequências de conquistas
- Não pode relacionar projetos dependentes
- Perde contexto entre conquistas

### 4.2 Problemas Técnicos

#### **1. SQLite Limitations**
- Escalabilidade limitada para muitos usuários
- Concorrência em escrita limitada
- Migrações mais complexas que em outros bancos

#### **2. Ausência de Full-Text Search**
- Busca apenas por tags e título básico
- Não busca em description
- Performance ruim para muitas brags

#### **3. Cache Strategy**
- Nenhuma estratégia de cache implementada
- Queries repetidas para mesmo usuário
- Performance degrada com crescimento

#### **4. Monitoring e Logging**
- Não há logs estruturados
- Não há métricas de performance
- Dificuldade em debug production issues

---

## 5. Propostas de Melhoria

### 5.1 Modelo de Dados Aprimorado

#### **Tabelas Novas Propostas**
```sql
-- 1. Categorias personalizáveis
CREATE TABLE categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    owner_id INTEGER,  -- NULL para categorias do sistema
    name TEXT NOT NULL,
    description TEXT,
    color TEXT,
    icon TEXT,
    is_system BOOLEAN DEFAULT FALSE,
    sort_order INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME,
    UNIQUE(owner_id, name),
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

-- 2. Versionamento de brags
CREATE TABLE brag_versions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    brag_id INTEGER NOT NULL,
    version INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    metadata JSON,
    changed_by INTEGER,
    change_reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (brag_id) REFERENCES brags(id) ON DELETE CASCADE,
    FOREIGN KEY (changed_by) REFERENCES users(id),
    UNIQUE(brag_id, version)
);

-- 3. Sistema de compartilhamento
CREATE TABLE brag_shares (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    brag_id INTEGER NOT NULL,
    shared_by INTEGER NOT NULL,
    shared_with INTEGER NOT NULL,
    permission_level TEXT NOT NULL DEFAULT 'view',
    expires_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (brag_id) REFERENCES brags(id) ON DELETE CASCADE,
    FOREIGN KEY (shared_by) REFERENCES users(id),
    FOREIGN KEY (shared_with) REFERENCES users(id),
    UNIQUE(brag_id, shared_with)
);

-- 4. Relacionamentos entre brags
CREATE TABLE brag_relationships (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    source_brag_id INTEGER NOT NULL,
    target_brag_id INTEGER NOT NULL,
    relationship_type TEXT NOT NULL,
    description TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (source_brag_id) REFERENCES brags(id) ON DELETE CASCADE,
    FOREIGN KEY (target_brag_id) REFERENCES brags(id) ON DELETE CASCADE,
    CHECK(source_brag_id != target_brag_id)
);
```

#### **Modificações em Tabelas Existentes**
```sql
-- Adicionar soft delete
ALTER TABLE users ADD COLUMN deleted_at DATETIME;
ALTER TABLE brags ADD COLUMN deleted_at DATETIME;
ALTER TABLE tags ADD COLUMN deleted_at DATETIME;

-- Adicionar metadados flexíveis
ALTER TABLE brags ADD COLUMN metadata JSON;

-- Substituir category integer por foreign key
ALTER TABLE brags ADD COLUMN category_id INTEGER;
ALTER TABLE brags ADD FOREIGN KEY (category_id) REFERENCES categories(id);

-- Adicionar campos de auditoria
ALTER TABLE brags ADD COLUMN created_by INTEGER REFERENCES users(id);
ALTER TABLE brags ADD COLUMN updated_by INTEGER REFERENCES users(id);
```

### 5.2 Novos Índices Recomendados
```sql
CREATE INDEX idx_brags_deleted_at ON brags(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_categories_owner_id ON categories(owner_id);
CREATE INDEX idx_brag_versions_brag_id ON brag_versions(brag_id);
CREATE INDEX idx_brag_shares_brag_id ON brag_shares(brag_id);
CREATE INDEX idx_brag_shares_shared_with ON brag_shares(shared_with);
CREATE INDEX idx_brag_relationships_source ON brag_relationships(source_brag_id);
CREATE INDEX idx_brag_relationships_target ON brag_relationships(target_brag_id);
```

### 5.3 Modelo de Domínio Aprimorado

```go
// Estruturas propostas
type Brag struct {
    ID          int64
    Owner       User
    Title       string
    Description string
    Category    *Category  // Referência em vez de enum
    Tags        []*Tag
    Metadata    map[string]interface{}
    Versions    []BragVersion
    Shares      []BragShare
    Relationships []BragRelationship
    Visibility  VisibilityLevel
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   *time.Time
    CreatedBy   *User
    UpdatedBy   *User
}

type Category struct {
    ID          int64
    Owner       *User  // NULL para system categories
    Name        string
    Description string
    Color       string
    Icon        string
    IsSystem    bool
    SortOrder   int
    CreatedAt   time.Time
    UpdatedAt   time.Time
}

type BragVersion struct {
    ID           int64
    BragID       int64
    Version      int
    Title        string
    Description  string
    Metadata     map[string]interface{}
    ChangedBy    *User
    ChangeReason string
    CreatedAt    time.Time
}

type BragShare struct {
    ID             int64
    Brag           *Brag
    SharedBy       *User
    SharedWith     *User
    Permission     PermissionLevel
    ExpiresAt      *time.Time
    CreatedAt      time.Time
}

type BragRelationship struct {
    ID               int64
    SourceBrag       *Brag
    TargetBrag       *Brag
    RelationshipType RelationshipType
    Description      string
    CreatedAt        time.Time
}
```

---

## 6. Roadmap de Implementação

### Fase 1: Correções Críticas (Sprint 1-2)
1. **Soft Delete** - Adicionar `deleted_at` em todas as tabelas
2. **Metadados Flexíveis** - Campo JSON `metadata` em brags
3. **Auditoria Básica** - `created_by`, `updated_by`

### Fase 2: Flexibilidade (Sprint 3-4)
4. **Categorias Personalizáveis** - Nova tabela `categories`
5. **Migração de Dados** - Converter enum para foreign key
6. **Versionamento Básico** - Tabela `brag_versions`

### Fase 3: Colaboração (Sprint 5-6)
7. **Sistema de Permissões** - Tabela `brag_shares`
8. **Relacionamentos** - Tabela `brag_relationships`
9. **APIs de Compartilhamento** - Endpoints REST/GraphQL

### Fase 4: Otimizações (Sprint 7-8)
10. **Full-Text Search** - Índices de busca em texto
11. **Cache Strategy** - Redis/Memcached para queries frequentes
12. **Monitoring** - Métricas, logs estruturados, alertas

### Fase 5: Escalabilidade (Sprint 9-10)
13. **Migração para PostgreSQL** - Para maior escala
14. **Microservices** - Separação por domínio
15. **Event Sourcing** - Para audit trail completo

---

## 7. Impacto Esperado

### Para Usuários
- ✅ **Mais flexibilidade:** Categorias e metadados customizados
- ✅ **Mais segurança:** Recuperação de dados excluídos
- ✅ **Colaboração:** Compartilhamento controlado
- ✅ **Contexto:** Relacionamentos entre conquistas

### Para Desenvolvedores
- ✅ **Manutenibilidade:** Modelo mais extensível
- ✅ **Testabilidade:** Soft delete facilita testes
- ✅ **Performance:** Índices otimizados, cache
- ✅ **Observability:** Métricas e logs para debug

### Para Negócio
- ✅ **Retenção:** Recursos avançados mantêm usuários
- ✅ **Monetização:** Funcionalidades premium possíveis
- ✅ **Escalabilidade:** Pronto para crescimento
- ✅ **Competitividade:** Diferenciação no mercado

---

## 8. Riscos e Mitigações

| Risco | Impacto | Mitigação |
|-------|---------|-----------|
| Breaking changes | Alto | Migração gradual, compatibilidade retroativa |
| Performance degradation | Médio | Índices, paginação, otimizações incrementais |
| Data migration errors | Alto | Backups, scripts testados, rollback plan |
| Increased complexity | Médio | Documentação, exemplos, onboarding |
| User adoption | Médio | UI/UX cuidadosa, tutoriais, feedback |

---

## 9. Recomendações Imediatas

### Prioridade 1 (Fácil, Alto Impacto)
1. **Implementar soft delete** - `deleted_at` em todas as tabelas
2. **Adicionar metadados JSON** - Campo `metadata` em brags
3. **Criar plano de migração** - Para futuras breaking changes

### Prioridade 2 (Médio Esforço, Médio Impacto)
4. **Categorias personalizáveis** - Tabela `categories`
5. **Versionamento básico** - Tabela `brag_versions`
6. **Sistema de permissões** - Tabela `brag_shares`

### Prioridade 3 (Alto Esforço, Alto Impacto)
7. **Migração para PostgreSQL** - Para escala
8. **Full-text search** - Para buscas avançadas
9. **APIs GraphQL** - Para flexibilidade de queries

---

## 10. Conclusão

### Estado Atual
O Bragdoc possui uma **base sólida** com:
- Modelo relacional bem estruturado
- Arquitetura limpa (Clean Architecture)
- Separação clara de responsabilidades
- Boas práticas de desenvolvimento Go

### Limitações Principais
1. **Rigidez:** Categorias fixas, metadados limitados
2. **Fragilidade:** Exclusões permanentes, sem versionamento
3. **Isolamento:** Sem colaboração ou relacionamentos

### Visão Futura
Com as melhorias propostas, o Bragdoc pode evoluir para:
- **Sistema flexível:** Adaptável a diferentes casos de uso
- **Plataforma colaborativa:** Compartilhamento e trabalho em equipe
- **Ferramenta profissional:** Para carreira e desenvolvimento pessoal
- **Produto escalável:** Pronto para crescimento de usuários

### Próximos Passos
1. **Validar** propostas com stakeholders
2. **Priorizar** baseado em recursos disponíveis
3. **Criar** issues detalhadas no GitHub
4. **Implementar** em sprints incrementais
5. **Comunicar** mudanças aos usuários

---

## Apêndices

### A. Códigos de Categoria Atuais
```go
CategoryUnknown     = 1  // "unknown"
CategoryProject     = 2  // "project"
CategoryAchievement = 3  // "achievement"
CategorySkill       = 4  // "skill"
CategoryLeadership  = 5  // "leadership"
CategoryInnovation  = 6  // "innovation"
```

### B. Exemplo de Metadata JSON
```json
{
  "project": {
    "start_date": "2024-01-15",
    "end_date": "2024-06-30",
    "budget": 50000,
    "team_size": 5,
    "technologies": ["Go", "React", "PostgreSQL"],
    "url": "https://github.com/example/project"
  },
  "achievement": {
    "metric": "revenue_increase",
    "value": 25,
    "unit": "percent",
    "timeframe": "quarter",
    "evidence": "Q3-2024-report.pdf"
  }
}
```

### C. Script de Migração Exemplo
```sql
BEGIN TRANSACTION;

-- 1. Criar tabela de categorias
CREATE TABLE categories (...);

-- 2. Migrar categorias existentes
INSERT INTO categories (owner_id, name, is_system, sort_order)
SELECT NULL, 'project', TRUE, 1
UNION ALL SELECT NULL, 'achievement', TRUE, 2
UNION ALL SELECT NULL, 'skill', TRUE, 3
UNION ALL SELECT NULL, 'leadership', TRUE, 4
UNION ALL SELECT NULL, 'innovation', TRUE, 5;

-- 3. Adicionar category_id às brags existentes
UPDATE brags SET category_id = 
  CASE category
    WHEN 2 THEN (SELECT id FROM categories WHERE name = 'project')
    WHEN 3 THEN (SELECT id FROM categories WHERE name = 'achievement')
    WHEN 4 THEN (SELECT id FROM categories WHERE name = 'skill')
    WHEN 5 THEN (SELECT id FROM categories WHERE name = 'leadership')
    WHEN 6 THEN (SELECT id FROM categories WHERE name = 'innovation')
    ELSE (SELECT id FROM categories WHERE name = 'unknown')
  END;

-- 4. Adicionar soft delete
ALTER TABLE users ADD COLUMN deleted_at DATETIME;
ALTER TABLE brags ADD COLUMN deleted_at DATETIME;
ALTER TABLE tags ADD COLUMN deleted_at DATETIME;

-- 5. Adicionar metadados
ALTER TABLE brags ADD COLUMN metadata JSON;

COMMIT;
```

### D. Referências
- [SQLite Documentation](https://www.sqlite.org/docs.html)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain-Driven Design](https://domainlanguage.com/ddd/)
- [12 Factor App](https://12factor.net/)

---

**Documento gerado em:** 2026-02-23  
**Analista:** OpenClaw AI Assistant  
**Repositório:** https://github.com/vagnerclementino/bragdoc  
**Branch:** beta_version  
**Commit analisado:** Último commit da branch beta_version