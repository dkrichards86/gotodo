package gotodo

import (
	"errors"
	"os"
	"strconv"

	"github.com/boltdb/bolt"
	homedir "github.com/mitchellh/go-homedir"
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
type BoltStorage struct{}

// TodoDir is the base directory for gotodo
const todoDir = ".gotodo"

// TodoDBFile is the name of the todo file in the user's home directory
const todoDBFile = "todo.db"

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

// makeTodoDir creates a directory at dirPath
func makeTodoDir(dirPath string) (string, error) {
	err := os.Mkdir(dirPath, 0755)
	if err != nil {
		return "", nil
	}
	return dirPath, nil
}

// getTodoDir returns the absolute path to the gotodo directory
func getTodoDir() (string, error) {
	homePath, err := getHomeDir()
	if err != nil {
		return "", err
	}

	dirPath := homePath + "/" + todoDir

	_, err = os.Stat(dirPath)
	if os.IsNotExist(err) {
		return makeTodoDir(dirPath)
	}

	return dirPath, nil
}

// getDB returns an instance of *bolt.DB
func (me *BoltStorage) getDB() (*bolt.DB, error) {
	absPath, err := getTodoDir()
	if err != nil {
		return nil, err
	}
	dbPath := absPath + "/" + todoDBFile

	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("Todos"))
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

// Create inserts a new *Todo
func (me *BoltStorage) Create(todo *Todo) error {
	db, err := me.getDB()
	if err != nil {
		return err
	}
	defer db.Close()

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Todos"))
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

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Todos"))

		key := strconv.Itoa(todoID)
		v := b.Get([]byte(key))

		todo = FromString(string(v))
		todo.TodoID = todoID
		return nil
	})

	if err != nil {
		return todo, err
	}

	// Check to see if the todo exists
	if todo.Description == "" {
		return todo, errors.New("Todo ID does not exist")
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
		b := tx.Bucket([]byte("Todos"))
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

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Todos"))

		key := strconv.Itoa(todoID)
		value := todo.String()
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

	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Todos"))

		key := strconv.Itoa(todoID)
		return b.Delete([]byte(key))
	})
}
