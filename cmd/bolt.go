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
	db, err := bolt.Open("/etc/tuber-bolt/db", 0666, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return err
		}
		err = b.Put([]byte("answer"), []byte("42"))
		return err
	})

	if err == nil {
		err = db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte("MyBucket"))
			v := b.Get([]byte("answer"))
			fmt.Printf("The answer is: %s\n", v)
			return nil
		})
	}

	if err != nil {
		fmt.Println(err.Error())
	}
	select {}
	return nil
}

func init() {
	rootCmd.AddCommand(bolterCmd)
}
