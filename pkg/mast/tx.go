package mast

// Tx holds transaction state to support backtracking.
//
// This is only useful for Capture.
type Tx struct {
	parent  *Tx
	cancels []func()
}

// Begin starts a nested transaction.
func (tx *Tx) Begin() *Tx {
	return &Tx{parent: tx}
}

// End wraps up hte nested transaction.
func (tx *Tx) End() {
	if tx.parent == nil {
		return
	}
	tx.parent.cancels = append(tx.parent.cancels, tx.cancels...)
}

// Cancel cancels all the assignments done until this point.
func (tx *Tx) Cancel() {
	for _, cancel := range tx.cancels {
		cancel()
	}
	tx.cancels = nil
}

// AddCancel registers a cancel function.
func (tx *Tx) AddCancel(cancel func()) {
	tx.cancels = append(tx.cancels, cancel)
}
