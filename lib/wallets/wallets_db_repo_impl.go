package wallets

import (
	"errors"
	"github.com/reiver/logjam/lib/marshal"
	dbsrv "github.com/reiver/logjam/srv/db"
)

type walletRepo struct {
}

const (
	walletTbl = "wallets"

	ownerIdKey = "ownerId"
)

func GetNewWalletRepository() IWalletRepository {
	return &walletRepo{}
}

func (w *walletRepo) Add(dto AddWalletDTO) (string, error) {
	data, err := marshal.ObjToMap(dto)
	if err != nil {
		return "", err
	}
	return dbsrv.Repository.Insert(walletTbl, data)
}

func (w *walletRepo) UpdateWallet(dto UpdateWalletDTO) error {
	findResult, err := dbsrv.Repository.GetByFilter(walletTbl, map[string]any{
		ownerIdKey: dto.OwnerId,
	})
	if err != nil {
		return err
	}
	if findResult == nil || len(findResult) == 0 {
		return errors.New("couldnt find the wallet, maybe id and ownerId doesnt match")
	}
	data, err := marshal.ObjToMap(dto)
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return errors.New("invalid input")
	}
	return dbsrv.Repository.Update(walletTbl, findResult[0]["id"].(string), data)
}

func (w *walletRepo) GetUserWallet(id string) (wallet *WalletDTO, err error) {
	rows, err := dbsrv.Repository.GetByFilter(walletTbl, map[string]any{
		ownerIdKey: id,
	})
	if err != nil {
		return nil, err
	}
	if rows == nil || len(rows) == 0 {
		return nil, errors.New("couldnt find the user wallet")
	}
	err = marshal.MapToObj(rows[0], &wallet)
	return
}
