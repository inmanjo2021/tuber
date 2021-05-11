package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	bolt "go.etcd.io/bbolt"
)

var bolterCmd = &cobra.Command{
	SilenceUsage: true,
	Use:          "bolter",
	Short:        "boltss",
	RunE:         bolter,
}

func bolter(cmd *cobra.Command, args []string) error {
	db, err := bolt.Open("/etc/tuber-bolt", 0666, nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		err := b.Put([]byte("answer"), []byte("42"))
		return err
	})
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("MyBucket"))
		v := b.Get([]byte("answer"))
		fmt.Printf("The answer is: %s\n", v)
		return nil
	})
	return nil
}

func init() {
	rootCmd.AddCommand(bolterCmd)
}
