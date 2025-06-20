package crypto

import "os"

func ReadKeys(private_key_path string, public_key_path string) ([]byte, []byte, error) {
	privateKey, err := os.ReadFile(private_key_path)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := os.ReadFile(public_key_path)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}
