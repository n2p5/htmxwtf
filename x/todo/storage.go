package main

import (
	"bytes"
	"encoding/gob"
	"log"

	badger "github.com/dgraph-io/badger/v4"
)

type Store interface {
	All() ([]Todo, error)
	Get(id string) (Todo, error)
	New(todo Todo) error
	Update(todo Todo) error
	Delete(id string) error
	Close() error
}

type BadgerStore struct {
	db *badger.DB
}

func NewBadgerStore() *BadgerStore {
	db, err := badger.Open(badger.DefaultOptions("./badgerdb"))
	if err != nil {
		log.Fatal(err)
	}
	return &BadgerStore{db: db}
}

func (s *BadgerStore) Close() error {
	return s.db.Close()
}

func (s *BadgerStore) All() ([]Todo, error) {
	todos := []Todo{}
	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte("todo")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			var todo Todo
			err := item.Value(func(v []byte) error {
				gob.NewDecoder(bytes.NewReader(v)).Decode(&todo)
				return nil
			})
			if err != nil {
				return err
			}
			todos = append(todos, todo)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return todos, nil
}

func (s *BadgerStore) Get(id string) (Todo, error) {
	var todo Todo
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}
		return item.Value(func(v []byte) error {
			gob.NewDecoder(bytes.NewReader(v)).Decode(&todo)
			return nil
		})
	})
	if err != nil {
		return Todo{}, err
	}
	return todo, nil
}

func (s *BadgerStore) New(todo Todo) error {
	return s.db.Update(func(txn *badger.Txn) error {
		var buf bytes.Buffer
		err := gob.NewEncoder(&buf).Encode(todo)
		if err != nil {
			return err
		}
		return txn.Set([]byte(todo.ID), buf.Bytes())
	})
}

func (s *BadgerStore) Update(todo Todo) error {
	return s.New(todo)
}

func (s *BadgerStore) Delete(id string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(id))
	})
}
