package transactions

// The transaction that gives block reward. Only 1 of these are in each block. Has a pre-defined lock time.
type PLuX struct {
	blockReward uint64

	luckyMiner []byte

	weight Weight
}

// A basic transaction from 1 person to another. Has a customisable lock time.
type BLuX struct {
	txInput []Input
	pubKey  []byte
}

// A more advanced transaction that will have very basic scripting functionality.
type ALuX struct{}

type Input struct {
	blockNumber uint32
	txHash      []byte
	amount      uint64
}
