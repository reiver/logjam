package wallets

import (
	"github.com/reiver/logjam/lib/marshal"
	dbsrv "github.com/reiver/logjam/srv/db"
)

type walletRepo struct {
}

const (
	walletTbl = "wallets"
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
	//TODO implement me
	panic("implement me")
}

func (w *walletRepo) GetUserWallet(id string) (*WalletDTO, error) {
	//TODO implement me
	panic("implement me")
}
