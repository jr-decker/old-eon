package gas

// Gas (measured in Virtual Gas Units) tracks data transaction size.
// This with gas price is used to determine the fees on the layer 1 threader is applicable.
// 1 Vgu = 1 byte taken up by a data transaction.
type Gas uint

// Returns the Uint value of the gas.
func (g *Gas) Uint() uint {

	return uint(*g)
}
