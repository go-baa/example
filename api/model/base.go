package model

import (
	"sync"
)

// db simple database
type db struct {
	store map[string]*table
	lock  sync.RWMutex
}

// table table has rows of same type
type table struct {
	rows []*row
	lock sync.RWMutex
}

// row base component of database
type row struct {
	data    interface{}
	deleted bool
	lock    sync.RWMutex
}

// Table get table store, return or create it
func (t *db) Table(name string) (*table, error) {
	t.lock.Lock()
	tab, ok := t.store[name]
	if !ok {
		tab = new(table)
		t.store[name] = tab
	}
	t.lock.Unlock()
	return tab, nil
}

// Rows returns all of data
func (t *table) Rows(v interface{}) error {

}

// Row return a row of data
func (t *table) Row(id int, v interface{}) error {

}
