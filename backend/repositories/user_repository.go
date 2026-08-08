package repositories

import (
	"context"
	"errors"

	"iam-security-backend/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrDuplicateEmail = errors.New("email already exists")
	ErrRoleNotFound   = errors.New("role not found")
)

type UserRepository struct{ db *pgxpool.Pool }

func NewUserRepository(db *pgxpool.Pool) *UserRepository    { return &UserRepository{db: db} }
func (repo *UserRepository) Ping(ctx context.Context) error { return repo.db.Ping(ctx) }

func duplicateEmail(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "uq_user_email"
}

func (repo *UserRepository) Register(ctx context.Context, input models.RegisterInput, passwordHash string) (string, error) {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)
	var userID string
	err = tx.QueryRow(ctx, `INSERT INTO identity."user" (email,password_hash,first_name,last_name,actor_type) VALUES ($1,$2,$3,$4,'LEARNER') RETURNING id`, input.Email, passwordHash, input.FirstName, input.LastName).Scan(&userID)
	if duplicateEmail(err) {
		return "", ErrDuplicateEmail
	}
	if err != nil {
		return "", err
	}
	_, err = tx.Exec(ctx, `INSERT INTO rbac.user_role (user_id,role_id,assigned_by) SELECT $1,id,$1 FROM rbac.role WHERE name='LEARNER'`, userID)
	if err != nil {
		return "", err
	}
	if err = tx.Commit(ctx); err != nil {
		return "", err
	}
	return userID, nil
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (models.LoginUser, error) {
	var user models.LoginUser
	err := repo.db.QueryRow(ctx, `SELECT id,email,password_hash,is_active FROM identity."user" WHERE lower(email)=lower($1)`, email).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.IsActive)
	return user, err
}

func (repo *UserRepository) RecordLogin(ctx context.Context, userID *string, email, outcome, ip, userAgent string) {
	_, _ = repo.db.Exec(ctx, `INSERT INTO identity_audit.audit_login (user_id,email_attempted,outcome,ip_address,user_agent) VALUES ($1,$2,$3,$4,$5)`, userID, email, outcome, ip, userAgent)
}

func (repo *UserRepository) CompleteLogin(ctx context.Context, user models.LoginUser, tokenHash, userAgent string) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `UPDATE identity."user" SET last_login_at=now(),failed_attempts=0 WHERE id=$1`, user.ID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO session.refresh_token (user_id,token_hash,device_hint) VALUES ($1,$2,$3)`, user.ID, tokenHash, userAgent)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (repo *UserRepository) CreateResetRequest(ctx context.Context, email, tokenHash, ip string) (bool, error) {
	var userID string
	err := repo.db.QueryRow(ctx, `SELECT id FROM identity."user" WHERE lower(email)=lower($1)`, email).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	_, err = repo.db.Exec(ctx, `INSERT INTO session.password_reset_request (user_id,token_hash,ip_address) VALUES ($1,$2,$3)`, userID, tokenHash, ip)
	return true, err
}

func (repo *UserRepository) ResetPassword(ctx context.Context, tokenHash, passwordHash string) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var userID string
	err = tx.QueryRow(ctx, `SELECT user_id FROM session.password_reset_request WHERE token_hash=$1 AND is_used=false AND expires_at>now() ORDER BY requested_at DESC LIMIT 1 FOR UPDATE`, tokenHash).Scan(&userID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE identity."user" SET password_hash=$1,updated_at=now() WHERE id=$2`, passwordHash, userID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE session.password_reset_request SET is_used=true WHERE token_hash=$1`, tokenHash)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (repo *UserRepository) Profile(ctx context.Context, userID string) (models.Profile, error) {
	var profile models.Profile
	err := repo.db.QueryRow(ctx, `SELECT email,first_name,last_name,actor_type,is_active FROM identity."user" WHERE id=$1`, userID).Scan(&profile.User.Email, &profile.User.FirstName, &profile.User.LastName, &profile.User.ActorType, &profile.User.IsActive)
	if err != nil {
		return profile, err
	}
	rows, err := repo.db.Query(ctx, `SELECT DISTINCT ro.display_name,f.name,rf.scope_type FROM rbac.user_role ur JOIN rbac.role ro ON ro.id=ur.role_id JOIN rbac.role_feature rf ON rf.role_id=ro.id JOIN rbac_catalog.feature f ON f.id=rf.feature_id WHERE ur.user_id=$1 ORDER BY ro.display_name,f.name`, userID)
	if err != nil {
		return profile, err
	}
	defer rows.Close()
	profile.Permissions = []models.Permission{}
	for rows.Next() {
		var permission models.Permission
		if err := rows.Scan(&permission.Role, &permission.Feature, &permission.Scope); err != nil {
			return profile, err
		}
		profile.Permissions = append(profile.Permissions, permission)
	}
	return profile, rows.Err()
}

func (repo *UserRepository) Summary(ctx context.Context) (models.Summary, error) {
	var result models.Summary
	err := repo.db.QueryRow(ctx, `SELECT (SELECT count(*) FROM identity."user"),(SELECT count(*) FROM rbac.role),(SELECT count(*) FROM rbac_catalog.feature),(SELECT count(*) FROM rbac_catalog.module)`).Scan(&result.Users, &result.Roles, &result.Features, &result.Modules)
	return result, err
}

func (repo *UserRepository) HasPermission(ctx context.Context, userID, featureCode string) bool {
	var allowed bool
	err := repo.db.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM rbac.user_role ur JOIN rbac.role_feature rf ON rf.role_id=ur.role_id JOIN rbac_catalog.feature f ON f.id=rf.feature_id WHERE ur.user_id=$1 AND f.code=$2)`, userID, featureCode).Scan(&allowed)
	return err == nil && allowed
}

func (repo *UserRepository) ListUsers(ctx context.Context) ([]models.ManagedUser, error) {
	rows, err := repo.db.Query(ctx, `SELECT u.id,u.email,u.first_name,u.last_name,u.actor_type,u.is_active,COALESCE(string_agg(DISTINCT ro.name,', '),'') FROM identity."user" u LEFT JOIN rbac.user_role ur ON ur.user_id=u.id LEFT JOIN rbac.role ro ON ro.id=ur.role_id GROUP BY u.id ORDER BY u.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []models.ManagedUser{}
	for rows.Next() {
		var user models.ManagedUser
		if err := rows.Scan(&user.ID, &user.Email, &user.FirstName, &user.LastName, &user.ActorType, &user.IsActive, &user.Roles); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (repo *UserRepository) ListRoles(ctx context.Context) ([]models.Role, error) {
	rows, err := repo.db.Query(ctx, `SELECT name,display_name FROM rbac.role ORDER BY display_name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	roles := []models.Role{}
	for rows.Next() {
		var role models.Role
		if err := rows.Scan(&role.Name, &role.DisplayName); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, rows.Err()
}

func (repo *UserRepository) CreateManaged(ctx context.Context, input models.ManagedUserInput, passwordHash, assignedBy string) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var userID string
	err = tx.QueryRow(ctx, `INSERT INTO identity."user" (email,password_hash,first_name,last_name,actor_type) VALUES ($1,$2,$3,$4,$5) RETURNING id`, input.Email, passwordHash, input.FirstName, input.LastName, input.ActorType).Scan(&userID)
	if duplicateEmail(err) {
		return ErrDuplicateEmail
	}
	if err != nil {
		return err
	}
	result, err := tx.Exec(ctx, `INSERT INTO rbac.user_role (user_id,role_id,assigned_by) SELECT $1,id,$2 FROM rbac.role WHERE name=$3`, userID, assignedBy, input.RoleName)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return ErrRoleNotFound
	}
	return tx.Commit(ctx)
}

func (repo *UserRepository) UpdateManaged(ctx context.Context, userID string, input models.ManagedUserInput) error {
	result, err := repo.db.Exec(ctx, `UPDATE identity."user" SET first_name=$1,last_name=$2,actor_type=$3,is_active=$4,updated_at=now() WHERE id=$5`, input.FirstName, input.LastName, input.ActorType, *input.IsActive, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

func (repo *UserRepository) Deactivate(ctx context.Context, userID string) error {
	tx, err := repo.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	result, err := tx.Exec(ctx, `UPDATE identity."user" SET is_active=false,updated_at=now() WHERE id=$1`, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	_, err = tx.Exec(ctx, `UPDATE session.refresh_token SET is_revoked=true,revoked_at=now() WHERE user_id=$1 AND is_revoked=false`, userID)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
