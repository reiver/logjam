package wallets

type IWalletRepository interface {
	Add(AddWalletDTO) (string, error)
	UpdateWallet(UpdateWalletDTO) error
	GetUserWallet(id string) (*WalletDTO, error)
}
