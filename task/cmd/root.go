package cmd

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

var rootCmd = &cobra.Command{
	Use:     "task",
	Short:   "Task is a command line task manager",
	Version: "v0.1.0",
	Run:     list,
}
var path = "tasks.db"

// Execute executes the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type store struct {
	db *bolt.DB
}

type task struct {
	ID   int
	Task string
	Done bool
}

func newStore() (*store, error) {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("tasks"))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &store{db}, nil
}

func (s *store) insertTask(data string) error {
	t := task{Task: data, Done: false}

	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
		id, _ := b.NextSequence()
		t.ID = int(id)

		buf, err := json.Marshal(t)
		if err != nil {
			return err
		}

		return b.Put(itob(t.ID), buf)
	})
}

func (s *store) getAllTasks() ([]task, error) {
	var tasks []task

	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var data task
			json.Unmarshal(v, &data)
			tasks = append(tasks, data)
		}

		return nil
	})

	return tasks, err
}

func (s *store) getActiveTasks() ([]task, error) {
	allTasks, err := s.getAllTasks()
	if err != nil {
		return nil, err
	}

	var tasks []task
	for _, v := range allTasks {
		if !v.Done {
			tasks = append(tasks, v)
		}
	}

	return tasks, nil
}

func (s *store) completeTask(id int) (task, error) {
	var data task
	err := s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("tasks"))
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

// itob returns an 8-byte big endian representation of v.
func itob(v int) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
