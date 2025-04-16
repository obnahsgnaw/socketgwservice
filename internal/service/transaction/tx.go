package transaction

// import (
// "github.com/obnahsgnaw/socketgwservice/internal/dal/query"
// "sync"
// )
type Tx int

func Default() Tx {
	return 0
}

//
//var _dbTransIns *DbTransaction
//var _dbTransOnce sync.Once
//
//func Transaction() *DbTransaction {
//	_dbTransOnce.Do(func() {
//		_dbTransIns = newTransaction(query.Q)
//	})
//
//	return _dbTransIns
//}
//
//type DbTransaction struct {
//	sync.Mutex
//	ins   map[Tx]*query.Query
//	index int
//}
//
//func newTransaction(q *query.Query) *DbTransaction {
//	s := &DbTransaction{
//		ins:   make(map[Tx]*query.Query),
//		index: -1,
//	}
//	s.NewTx(q)
//	return s
//}
//
//func (t *DbTransaction) GetTxQuery(tx Tx) *query.Query {
//	t.Lock()
//	defer t.Unlock()
//	if q, ok := t.ins[tx]; ok {
//		return q
//	}
//	return query.Q
//}
//
//func (t *DbTransaction) NewTx(q *query.Query) Tx {
//	t.Lock()
//	defer t.Unlock()
//	t.index++
//	tx := Tx(t.index)
//	t.ins[tx] = q
//	return tx
//}
//
//func (t *DbTransaction) ReleaseTx(tx Tx) {
//	if tx == Default() {
//		return
//	}
//	t.Lock()
//	defer t.Unlock()
//	if _, ok := t.ins[tx]; ok {
//		delete(t.ins, tx)
//	}
//}
//
//func (t *DbTransaction) Trans(tx Tx, cb func(Tx) error) error {
//	return t.GetTxQuery(tx).Transaction(func(q *query.Query) error {
//		tx1 := t.NewTx(q)
//		defer t.ReleaseTx(tx1)
//		return cb(tx1)
//	})
//}
//
//
//type TransRepository struct {
//	trans *DbTransaction
//}
//
//var _dbTransRepoIns *TransRepository
//var _dbTransRepoOnce sync.Once
//
//func TransRepo() *TransRepository {
//	_dbTransRepoOnce.Do(func() {
//		_dbTransRepoIns = &TransRepository{trans: Transaction()}
//	})
//	return _dbTransRepoIns
//}
//
//func (r *TransRepository) Transaction(tx Tx, cb func(tx Tx) error) error {
//	return r.trans.Trans(tx, cb)
//}
