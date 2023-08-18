package nonce

import (
	"context"
	"time"
	"web3Tarot-backend/log"
	"web3Tarot-backend/models"
	"web3Tarot-backend/util"
)

const TTL = 300 * time.Second

//type Service struct {
//	nonceDB *models
//}

//type Option func(svc *Service) error

//func New(nonceDB *models.NonceDB, opts ...Option) (*Service, error) {
//	srv := &Service{
//		nonceDB: nonceDB,
//	}
//	for _, opt := range opts {
//		if err := opt(srv); err != nil {
//			return nil, err
//		}
//	}
//	return srv, nil
//}

// GetNonce return a random string as nonce
func GetNonce(ctx context.Context) (*GetNonceData, error) {
	ulid := util.GenerateULID()
	exp := time.Now().Add(TTL).Unix()

	nonce := models.Nonce{
		Nonce:      ulid,
		Expiration: uint(exp),
	}
	err := models.CreateNonce(ctx, &nonce)
	if err != nil {
		log.Errorf("create nonce error: %v", err)
		return nil, err
	}
	data := GetNonceData(nonce.Nonce)
	return &data, nil
}
