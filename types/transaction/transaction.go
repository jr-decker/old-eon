package transaction

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/Sucks-To-Suck/Eon/core/gas"
	"github.com/Sucks-To-Suck/Eon/eocrypt"
	"github.com/Sucks-To-Suck/Eon/tools/merkle"
)

// A Transaction simply represents a token being sent from one party to another (not even the native token).
type Transaction struct {
	// Token data of the transaction.
	TokenHash eocrypt.Hash `json:"TokenHash"`
	Amount    *big.Int     `json:"Amount"`

	// From who and where to?
	To        []byte `json:"To"`
	From      []byte `json:"From"`
	Signature []byte `json:"Signature"`

	// Where are you getting the tokens from?
	BlockFrom []eocrypt.Hash `json:"BlockFromHash"`
	TxFrom    []eocrypt.Hash `json:"TransactionFromHash"`

	// Basic but otherwise needed data.
	ThreadId *big.Int `json:"Thread"`
	Gas      gas.Gas  `json:"Gas"`
	GasPrice *big.Int `json:"GasPrice"`

	// When did the node get this transaction?
	ReceivedTime time.Time `json:"Time"`
}

// ****
// Transaction basic interaction:

// The following functions all dont use / give pointers.
// This is to allow manipulation of the given / received data that then will not affect the copy in the tx until a function
// is called to set the new value.

func NewTransaction(TokenHash eocrypt.Hash, Amount *big.Int, To []byte, From []byte, Signature []byte,
	BlockFrom []eocrypt.Hash, TxFrom []eocrypt.Hash, ChainId *big.Int, Gas gas.Gas, GasPrice big.Int,
	ReceivedTime time.Time) *Transaction {

	t := new(Transaction)

	t.SetTokenHash(TokenHash)
	t.SetAmount(*Amount)
	t.SetTo(To)
	t.SetFrom(From)
	t.SetSignature(Signature)
	t.SetBlockFrom(BlockFrom)
	t.SetTxFrom(TxFrom)
	t.SetChainId(*ChainId)
	t.SetGas(Gas)
	t.SetGasPrice(GasPrice)
	t.SetReceivedTime(ReceivedTime)

	return t
}

func (t *Transaction) SetTokenHash(h eocrypt.Hash) {

	t.TokenHash = h
}

func (t *Transaction) SetAmount(b big.Int) {

	t.Amount = &b
}

func (t *Transaction) SetTo(p []byte) {

	t.To = p
}

func (t *Transaction) SetFrom(p []byte) {

	t.From = p
}

func (t *Transaction) SetBlockFrom(h []eocrypt.Hash) {

	t.BlockFrom = h
}

func (t *Transaction) SetTxFrom(h []eocrypt.Hash) {

	t.TxFrom = h
}

func (t *Transaction) SetSignature(sig []byte) {

	t.Signature = sig
}

func (t *Transaction) SetChainId(b big.Int) {

	t.ThreadId = &b
}

func (t *Transaction) SetGas(g gas.Gas) {

	t.Gas = g
}

func (t *Transaction) SetGasPrice(g big.Int) {

	t.GasPrice = &g
}

func (t *Transaction) SetReceivedTime(time time.Time) {

	t.ReceivedTime = time
}

func (t *Transaction) GetTokenHash() eocrypt.Hash   { return t.TokenHash }
func (t *Transaction) GetAmount() big.Int           { return *t.Amount }
func (t *Transaction) GetTo() []byte                { return t.To }
func (t *Transaction) GetFrom() []byte              { return t.From }
func (t *Transaction) GetBlockFrom() []eocrypt.Hash { return t.BlockFrom }
func (t *Transaction) GetTxFrom() []eocrypt.Hash    { return t.TxFrom }
func (t *Transaction) GetChainId() big.Int          { return *t.ThreadId }
func (t *Transaction) GetGas() gas.Gas              { return t.Gas }
func (t *Transaction) GetGasPrice() big.Int         { return *t.GasPrice }
func (t *Transaction) GetReceivedTime() time.Time   { return t.ReceivedTime }

// Transaction basic interaction:
// ****

// ****
// Transaction advanced interaction:

// Uses the gob encoder to encode and return the transaction as a Bytes Buffer.
func (t *Transaction) EncodeToBuffer() (*bytes.Buffer, error) {

	buff := new(bytes.Buffer)

	// Encode the transaction into the Bytes Buffer.
	encodeErr := gob.NewEncoder(buff).Encode(t)

	return buff, encodeErr
}

// Uses the gob encoder with a provided Bytes Buffer to encode the transaction into that Buffer.
func (t *Transaction) EncodeWithBuffer(b *bytes.Buffer) error {

	// Encode the transaction into the Bytes Buffer.
	return gob.NewEncoder(b).Encode(t)
}

// Encode the transaction into JSON format. Returns the encoded bytes in a Bytes Buffer.
func (t *Transaction) EncodeJSON() (*bytes.Buffer, error) {

	buff := new(bytes.Buffer)

	// Encode the transaction into the Bytes Buffer.
	encodeErr := json.NewEncoder(buff).Encode(t)

	return buff, encodeErr
}

// Encode the transaction with the provided Bytes Buffer. The encoded bytes will reside there.
func (t *Transaction) EncodeJSONwithBuff(b *bytes.Buffer) error {

	// Encode the transaction into the Bytes Buffer.
	return json.NewEncoder(b).Encode(t)
}

// Trys to decode the given Bytes Buffer (suppose to be the encoded form by gob). Returns the decoded transaction and any errors.
func Decode(b *bytes.Buffer) (*Transaction, error) {

	t := new(Transaction)

	// Try to decode the Bytes Buffer into a transaction.
	decodeErr := gob.NewDecoder(b).Decode(t)

	return t, decodeErr
}

// Trys to decode the given Bytes Buffer (suppose to be encoded JSON form). Returns the decoded transaction and any errors.
func DecodeJSON(b *bytes.Buffer) (*Transaction, error) {

	t := new(Transaction)

	// Try to decode the Bytes Buffer into a transaction.
	decodeErr := json.NewDecoder(b).Decode(t)

	return t, decodeErr
}

// Get the hash of the transaction. Will be whats signed by the transaction sender.
func (t *Transaction) Hash() *eocrypt.Hash {

	return eocrypt.HashInterface(
		[]interface{}{
			t.TokenHash.Bytes(),
			t.Amount.Bytes(),
			t.To,
			t.From,
			t.BlockFromRoot().Bytes(),
			t.TxFromRoot().Bytes(),
			t.ThreadId.Bytes(),
			t.Gas.Uint(),
			t.GasPrice.Bytes(),
		},
	)
}

// Prints the main important information about the transaction.
func (t *Transaction) Print() {

	fmt.Printf(`
	[
		Transaction: %x
		
		Token Hash: %x
		Amount: %d
		To: %x
		ChainId: %x
		Gas: %v
		GasPrice: %v
	]`, t.Hash().Bytes(), t.TokenHash.Bytes(), t.Amount.Bytes()[:], t.To, t.ThreadId, t.Gas, t.GasPrice)
}

// Gets the merkle root of the block from hashes, aka the block hashes of the blocks you are pointing
// to in the transaction to show where the token being spent is coming from.
func (t *Transaction) BlockFromRoot() *eocrypt.Hash {

	// Create the merkle tree.
	m := merkle.NewMerkleTree(t.GetBlockFrom())

	// Return the merkle root.
	return m.FindRoot()
}

// Get the merkle root from the tx hashes included in the tx, aka which transaction are you getting
// the token from.
func (t *Transaction) TxFromRoot() *eocrypt.Hash {

	// Create a merkle tree.
	m := merkle.NewMerkleTree(t.GetTxFrom())

	// Return the merkle root.
	return m.FindRoot()
}

// Transaction advanced interaction:
// ****
