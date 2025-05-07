package hdkey

type GenerateKey interface {
	Generate(seed []byte) (privateKey string, address string, err error)
}
