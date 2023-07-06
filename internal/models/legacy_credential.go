package models

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/supabase/gotrue/internal/crypto"
	"github.com/supabase/gotrue/internal/storage"
)

// LegacyCredential respresents a registered user with email/password authentication on legacy app
type LegacyCredential struct {
	ID                uuid.UUID          `json:"id" db:"id"`
	Email             storage.NullString `json:"email" db:"email"`
	Phone             storage.NullString `json:"phone" db:"phone"`
	EncryptedPassword string             `json:"-" db:"encrypted_password"`
	CreatedAt         time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt         *time.Time         `json:"updated_at" db:"updated_at"`
}

func findLegacyCredential(tx *storage.Connection, query string, args ...interface{}) (*LegacyCredential, error) {
	obj := &LegacyCredential{}
	if err := tx.Eager().Q().Where(query, args...).First(obj); err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return nil, UserNotFoundError{}
		}
		return nil, errors.Wrap(err, "error finding legacy credential")
	}

	return obj, nil
}

// TableName overrides the table name used by pop
func (LegacyCredential) TableName() string {
	tableName := "legacy_credentials"
	return tableName
}

// FindUserByEmailFromLegacy finds a user with the matching email from legacy credentials datas.
func FindUserByEmailFromLegacy(tx *storage.Connection, email string) (*LegacyCredential, error) {
	return findLegacyCredential(tx, "LOWER(email) = ?", strings.ToLower(email))
}

// Authenticate a user from a password
func (u *LegacyCredential) Authenticate(password string) bool {
	if u.EncryptedPassword == "" {
		return false
	}
	err := crypto.CompareLegacyHashAndPassword(context.Background(), u.EncryptedPassword, password)
	return err == nil
}
