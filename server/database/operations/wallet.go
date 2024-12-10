package operations

import (
	"database/sql"
	"fmt"
	"server/database/models"

	"github.com/btcsuite/btcd/rpcclient"
)

// UpdateWalletAddress updates the only record in the WalletInfo table.
func UpdateWalletAddress(db *sql.DB, address string) error {
	query := `UPDATE WalletInfo SET address = ?`
	_, err := db.Exec(query, address)
	if err != nil {
		return fmt.Errorf("error updating record from WalletInfo: %v", err)
	}

	fmt.Printf("Wallet address updated successfully in WalletInfo.\n")
	return nil
}

// UpdateWalletPassphrases updates the only record in the WalletInfo table.
func UpdateWalletPassphrases(db *sql.DB, pubPassphrase, privPassphrase string) error {
	query := `UPDATE WalletInfo SET pubPassphrase = ?, privPassphrase = ?`
	_, err := db.Exec(query, pubPassphrase, privPassphrase)
	if err != nil {
		return fmt.Errorf("error updating record from WalletInfo: %v", err)
	}

	fmt.Printf("Wallet passphrases updated successfully in WalletInfo.\n")
	return nil
}

// Store the wallet address.
func StoreAddress(btcwallet *rpcclient.Client, db *sql.DB) (string, error) {
	// Query the RPC server for the wallet address.
	address, err := btcwallet.GetAccountAddress("default")
	if err != nil {
		return "", err
	}

	// Store wallet address for transactions.
	err = UpdateWalletAddress(db, address.String())
	if err != nil {
		return "", err
	}

	return address.String(), nil
}

// GetWalletInfo retrieves the only record from the WalletInfo table.
func GetWalletInfo(db *sql.DB) (*models.WalletInfo, error) {
	var walletInfo models.WalletInfo
	query := `SELECT address, pubPassphrase, privPassphrase FROM WalletInfo`
	err := db.QueryRow(query).Scan(&walletInfo.Address, &walletInfo.PubPassphrase, &walletInfo.PrivPassphrase)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No record found
		}
		return nil, fmt.Errorf("error finding record in WalletInfo: %v", err)
	}

	return &walletInfo, nil
}
