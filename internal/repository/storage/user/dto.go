package user

import (
	"billsplitter-monolith/internal/domain/auth"
	"billsplitter-monolith/internal/utils/pg"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type userEntity struct {
	gorm.Model
	ID        string        `gorm:"column:id"`
	Username  string        `gorm:"column:username"`
	FirstName string        `gorm:"column:first_name"`
	LastName  string        `gorm:"column:last_name"`
	Extra     *pgtype.JSONB `gorm:"column:extra"`
}

func (userEntity) TableName() string {
	return "users"
}

func fromDomain(d *auth.User) *userEntity {
	if d == nil {
		return nil
	}

	extra, _ := pg.ToJsonb(&d.Extra)

	return &userEntity{
		ID:        d.ID,
		Username:  d.Username,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Extra:     extra,
	}
}

func toDomain(e *userEntity) *auth.User {
	if e == nil {
		return nil
	}

	out := auth.User{
		ID:        e.ID,
		Username:  e.Username,
		FirstName: e.FirstName,
		LastName:  e.LastName,
		Extra:     auth.UserExtra{},
	}

	v, _ := pg.FromJsonb[auth.UserExtra](e.Extra)
	if v != nil {
		out.Extra = *v
	}

	return &out
}
