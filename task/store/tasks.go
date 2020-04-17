package store

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"time"

	bolt "go.etcd.io/bbolt"
)

var dbPath = "tasks.db"
var bucketName = []byte("tasks")

// Store represents a task store with various high level helper methods
type Store struct {
	DB *bolt.DB
}

// Task represents a single task
type Task struct {
	ID        int
	Value     string
	Done      bool
	TimeStamp string
}

// NewStore creates and bootstraps a new task store
func NewStore() (*Store, error) {
	confDir, err := os.UserConfigDir()
	if err != nil {
		return nil, err
	}
	baseDir := path.Join(confDir, "task")
	err = os.MkdirAll(baseDir, 0755)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(path.Join(baseDir, dbPath), 0600, nil)
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
	t := Task{Value: data, Done: false, TimeStamp: time.Now().Format(time.RFC3339)}

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

// GetAllCompletedTasks return list of completed tasks
func (s *Store) GetAllCompletedTasks() ([]Task, error) {
	allTasks, err := s.GetAllTasks()
	if err != nil {
		return nil, err
	}

	var tasks []Task
	for _, v := range allTasks {
		if v.Done {
			tasks = append(tasks, v)
		}
	}

	return tasks, nil
}

// GetTasksDoneToday returns task which are marked as done today
func (s *Store) GetTasksDoneToday() ([]Task, error) {
	var tasks []Task
	t := time.Now()
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)

	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return fmt.Errorf("Bucket doesn't exist")
		}

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var data Task
			json.Unmarshal(v, &data)

			tst, _ := time.Parse(time.RFC3339, data.TimeStamp)

			if !data.Done {
				continue
			}
			if tst.Before(midnight) {
				continue
			}

			tasks = append(tasks, data)
		}

		return nil
	})

	return tasks, err
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
		data.TimeStamp = time.Now().Format(time.RFC3339)
		buf, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("Error encoding json: %s", err)
		}

		return b.Put(itob(id), buf)
	})

	return data, err
}

// DeleteTask deletes a specific task given it's ID
func (s *Store) DeleteTask(id int) (Task, error) {
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

		return b.Delete(itob(id))
	})

	return data, err
}

// DeleteAllTasks deletes all tasks by deleting thw whole bucket
func (s *Store) DeleteAllTasks() error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		return tx.DeleteBucket(bucketName)
	})
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
