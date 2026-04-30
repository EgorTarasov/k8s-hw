package crypto

import (
	"github.com/EgorTarasov/shopx/internal/domain"
	"github.com/google/uuid"
)

type UUIDGen struct{}

func NewUUIDGen() UUIDGen { return UUIDGen{} }

func (UUIDGen) New() domain.SessionToken { return domain.SessionToken(uuid.NewString()) }
func (UUIDGen) NewUserID() domain.UserID  { return domain.UserID(uuid.NewString()) }
func (UUIDGen) NewOrderID() domain.OrderID {
	return domain.OrderID(uuid.NewString())
}
