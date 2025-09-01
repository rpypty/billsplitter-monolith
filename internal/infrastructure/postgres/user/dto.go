package user

import (
	domain "billsplitter-monolith/internal/domain/user"
	vo "billsplitter-monolith/internal/domain/valueobject"
	"billsplitter-monolith/internal/utils/pg"
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)

type userEntity struct {
	gorm.Model
	ID        int64         `gorm:"column:id"`
	Username  string        `gorm:"column:username"`
	FirstName string        `gorm:"column:first_name"`
	LastName  string        `gorm:"column:last_name"`
	Extra     *pgtype.JSONB `gorm:"column:extra"`
}

func (userEntity) TableName() string {
	return "users"
}

func fromDomain(d *domain.User) *userEntity {
	if d == nil {
		return nil
	}

	extra, _ := pg.ToJsonb(&d.Extra)

	return &userEntity{
		ID:        int64(d.ID),
		Username:  d.Username,
		FirstName: d.FirstName,
		LastName:  d.LastName,
		Extra:     extra,
	}
}

func toDomain(e *userEntity) *domain.User {
	if e == nil {
		return nil
	}

	out := domain.User{
		ID:        vo.UserID(e.ID),
		Username:  e.Username,
		FirstName: e.FirstName,
		LastName:  e.LastName,
		Extra:     domain.ExtraInfo{},
	}

	v, _ := pg.FromJsonb[domain.ExtraInfo](e.Extra)
	if v != nil {
		out.Extra = *v
	}

	return &out
}
