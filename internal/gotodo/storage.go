package gotodo

import (
	"errors"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/mitchellh/go-homedir"
)

// Storage provides an interface for various todo storage mechanisms (in-memory, filesystem)
type Storage interface {
	Create(todo *Todo) error
	List() (TodoList, error)
	Get(todoID int) (*Todo, error)
	Update(todoID int, todo *Todo) error
	Delete(todoID int) error
}

// BoltStorage implements Storage, saving items to a file in the filesystem
type BoltStorage struct {
	Bucket []byte
}

// TodoDBFile is the name of the todo file in the user's home directory
const todoDBFile = ".gotodo.db"

// getHomeDir returns the user's home directory
func getHomeDir() (string, error) {
	relPath, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	absPath, err := homedir.Expand(relPath)
	if err != nil {
		return "", err
	}

	return absPath, nil
}

// getDB returns an instance of *bolt.DB
func (me *BoltStorage) getDB() (*bolt.DB, error) {
	absPath, err := getHomeDir()
	if err != nil {
		return nil, err
	}
	dbPath := absPath + "/" + todoDBFile

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(me.Bucket)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// checkKey checks for the existence of a key in Bolt.
func (me *BoltStorage) checkKey(key string, db *bolt.DB) error {
	return db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(me.Bucket)

		v := b.Get([]byte(key))
		if v == nil {
			return errors.New("Todo ID does not exist")
		}

		return nil
	})
}

// Create inserts a new *Todo
func (me *BoltStorage) Create(todo *Todo) error {
	db, err := me.getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(me.Bucket)
		id, _ := b.NextSequence()
		todo.TodoID = int(id)
		key := strconv.Itoa(todo.TodoID)
		value := todo.String()
		return b.Put([]byte(key), []byte(value))
	})
}

// Get retrieves the *Todo identified by todoID
func (me *BoltStorage) Get(todoID int) (*Todo, error) {
	var todo *Todo

	db, err := me.getDB()
	if err != nil {
		return todo, err
	}
	defer db.Close()

	key := strconv.Itoa(todoID)
	err = me.checkKey(key, db)
	if err != nil {
		return nil, err
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(me.Bucket)
		v := b.Get([]byte(key))
		todo = FromString(string(v))
		todo.TodoID = todoID
		return nil
	})

	if err != nil {
		return todo, err
	}

	return todo, nil
}

// List reads all Todos
func (me *BoltStorage) List() (TodoList, error) {
	items := make(TodoList, 0)

	db, err := me.getDB()
	if err != nil {
		return items, err
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(me.Bucket)
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			todo := FromString(string(v))
			todo.TodoID, _ = strconv.Atoi(string(k))
			items = append(items, todo)
		}

		return nil
	})

	return items, err
}

// Update modifies the *Todo identified by todoID
func (me *BoltStorage) Update(todoID int, todo *Todo) error {
	db, err := me.getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	key := strconv.Itoa(todoID)
	value := todo.String()

	// make sure the key exists before working with it
	err = me.checkKey(key, db)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(me.Bucket)
		return b.Put([]byte(key), []byte(value))
	})
}

// Delete removes the *Todo identified by todoID
func (me *BoltStorage) Delete(todoID int) error {
	db, err := me.getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	key := strconv.Itoa(todoID)

	// make sure the key exists before working with it
	err = me.checkKey(key, db)
	if err != nil {
		return err
	}

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(me.Bucket)
		return b.Delete([]byte(key))
	})
}
