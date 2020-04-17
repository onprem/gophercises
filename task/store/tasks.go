package store

import (
	"encoding/binary"
	"encoding/json"
	"fmt"

	bolt "go.etcd.io/bbolt"
)

var path = "tasks.db"
var bucketName = []byte("tasks")

// Store represents a task store with various high level helper methods
type Store struct {
	DB *bolt.DB
}

// Task represents a single task
type Task struct {
	ID    int
	Value string
	Done  bool
}

// NewStore creates and bootstraps a new task store
func NewStore() (*Store, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	s := &Store{DB: db}
	s.bootstrap()

	return s, nil
}

func (s *Store) bootstrap() error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(bucketName)
		return err
	})
}

// InsertTask inserts a new task into DB
func (s *Store) InsertTask(data string) error {
	t := Task{Value: data, Done: false}

	return s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return fmt.Errorf("Bucket doesn't exist")
		}

		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		t.ID = int(id)

		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		return b.Put(itob(t.ID), buf)
	})
}

// GetAllTasks returns all tasks (both active and done)
func (s *Store) GetAllTasks() ([]Task, error) {
	var tasks []Task

	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return fmt.Errorf("Bucket doesn't exist")
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var data Task
			json.Unmarshal(v, &data)
			tasks = append(tasks, data)
		}

		return nil
	})

	return tasks, err
}

// GetActiveTasks return list of currently active tasks
func (s *Store) GetActiveTasks() ([]Task, error) {
	allTasks, err := s.GetAllTasks()
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for _, v := range allTasks {
		if !v.Done {
			tasks = append(tasks, v)
		}
	}

	return tasks, nil
}

// CompleteTask marks a task as done given it's ID
func (s *Store) CompleteTask(id int) (Task, error) {
	var data Task

	err := s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return fmt.Errorf("Bucket doesn't exist")
		}

		v := b.Get(itob(id))
		if v == nil {
			return fmt.Errorf("Key does not exist")
		}

		err := json.Unmarshal(v, &data)
		if err != nil {
			return fmt.Errorf("Error decoding json: %s", err)
		}

		if data.Done {
			return fmt.Errorf("Task is already marked as done")
		}

		data.Done = true
		buf, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("Error encoding json: %s", err)
		}

		return b.Put(itob(id), buf)
	})

	return data, err
}

// Close closes the database backing the store
func (s *Store) Close() {
	s.DB.Close()
}

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
