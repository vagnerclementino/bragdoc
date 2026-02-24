package database

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "github.com/vagnerclementino/bragdoc/internal/database/queries"
    "github.com/vagnerclementino/bragdoc/internal/domain"
    "github.com/vagnerclementino/bragdoc/internal/repository"
)

type sqlitePositionRepository struct {
    db       *SQLiteDB
    userRepo repository.UserRepository
}

func NewPositionRepository(db *SQLiteDB, userRepo repository.UserRepository) repository.PositionRepository {
    return &sqlitePositionRepository{db: db, userRepo: userRepo}
}

func (r *sqlitePositionRepository) Get(ctx context.Context, id int64) (*domain.Position, error) {
    dbPosition, err := r.db.Queries().GetPosition(ctx, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("position not found: %d", id)
        }
        return nil, fmt.Errorf("failed to get position: %w", err)
    }
    return r.toDomainPosition(ctx, &dbPosition)
}

func (r *sqlitePositionRepository) ListByUser(ctx context.Context, userID int64) ([]*domain.Position, error) {
    dbPositions, err := r.db.Queries().ListPositionsByUser(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to list positions: %w", err)
    }
    
    positions := make([]*domain.Position, 0, len(dbPositions))
    for _, dbPosition := range dbPositions {
        position, err := r.toDomainPosition(ctx, &dbPosition)
        if err != nil {
            return nil, err
        }
        positions = append(positions, position)
    }
    return positions, nil
}

func (r *sqlitePositionRepository) Create(ctx context.Context, position *domain.Position) (*domain.Position, error) {
    if err := r.validatePosition(position); err != nil {
        return nil, fmt.Errorf("invalid position: %w", err)
    }
    
    var startDate, endDate sql.NullTime
    if position.StartDate != nil {
        startDate = sql.NullTime{Time: *position.StartDate, Valid: true}
    }
    if position.EndDate != nil {
        endDate = sql.NullTime{Time: *position.EndDate, Valid: true}
    }
    
    dbPosition, err := r.db.Queries().CreatePosition(ctx, queries.CreatePositionParams{
        UserID:    position.User.ID,
        Title:     position.Title,
        Company:   sql.NullString{String: position.Company, Valid: position.Company != ""},
        StartDate: startDate,
        EndDate:   endDate,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create position: %w", err)
    }
    return r.toDomainPosition(ctx, &dbPosition)
}

func (r *sqlitePositionRepository) Update(ctx context.Context, position *domain.Position) (*domain.Position, error) {
    if err := r.validatePosition(position); err != nil {
        return nil, fmt.Errorf("invalid position: %w", err)
    }
    
    var startDate, endDate sql.NullTime
    if position.StartDate != nil {
        startDate = sql.NullTime{Time: *position.StartDate, Valid: true}
    }
    if position.EndDate != nil {
        endDate = sql.NullTime{Time: *position.EndDate, Valid: true}
    }
    
    dbPosition, err := r.db.Queries().UpdatePosition(ctx, queries.UpdatePositionParams{
        Title:     position.Title,
        Company:   sql.NullString{String: position.Company, Valid: position.Company != ""},
        StartDate: startDate,
        EndDate:   endDate,
        ID:        position.ID,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to update position: %w", err)
    }
    return r.toDomainPosition(ctx, &dbPosition)
}

func (r *sqlitePositionRepository) Delete(ctx context.Context, id int64) error {
    // Check if position is in use
    count, err := r.db.Queries().CountBragsByPosition(ctx, id)
    if err != nil {
        return fmt.Errorf("failed to check position usage: %w", err)
    }
    if count > 0 {
        return fmt.Errorf("cannot delete position %d: %d brags use it", id, count)
    }
    
    if err := r.db.Queries().DeletePosition(ctx, id); err != nil {
        return fmt.Errorf("failed to delete position: %w", err)
    }
    return nil
}

func (r *sqlitePositionRepository) validatePosition(position *domain.Position) error {
    if position == nil {
        return errors.New("position cannot be nil")
    }
    
    if position.User.ID == 0 {
        return errors.New("position user ID cannot be empty")
    }
    
    if position.Title == "" {
        return errors.New("position title cannot be empty")
    }
    
    if position.Company == "" {
        return errors.New("position company cannot be empty")
    }
    
    // Validate date range if both dates are set
    if position.StartDate != nil && position.EndDate != nil {
        if position.EndDate.Before(*position.StartDate) {
            return errors.New("end date cannot be before start date")
        }
    }
    
    return nil
}

func (r *sqlitePositionRepository) toDomainPosition(ctx context.Context, dbPosition *queries.Position) (*domain.Position, error) {
    user, err := r.userRepo.Select(ctx, dbPosition.UserID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, fmt.Errorf("user not found for position %d: %w", dbPosition.ID, err)
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    position := &domain.Position{
        ID:      dbPosition.ID,
        User:    *user,
        Title:   dbPosition.Title,
    }
    
    if dbPosition.Company.Valid {
        position.Company = dbPosition.Company.String
    }
    
    if dbPosition.StartDate.Valid {
        startDate := dbPosition.StartDate.Time
        position.StartDate = &startDate
    }
    
    if dbPosition.EndDate.Valid {
        endDate := dbPosition.EndDate.Time
        position.EndDate = &endDate
    }
    
    if dbPosition.CreatedAt.Valid {
        position.CreatedAt = dbPosition.CreatedAt.Time
    }
    if dbPosition.UpdatedAt.Valid {
        position.UpdatedAt = dbPosition.UpdatedAt.Time
    }
    
    return position, nil
}
