package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"
)

// chunkIntervalResetCmd represents the chunk-interval reset command
var chunkIntervalResetCmd = &cobra.Command{
	Use:   "reset <metric>",
	Short: "Resets the chunk interval for a specific metric back to the default",
	Args:  cobra.ExactArgs(1),
	RunE:  chunkIntervalReset,
}

func init() {
	chunkIntervalCmd.AddCommand(chunkIntervalResetCmd)
}

func chunkIntervalReset(cmd *cobra.Command, args []string) error {
	var err error

	if os.Getenv("PGPASSWORD_POSTGRES") == "" {
		return errors.New("Password for postgres user must be set in environment variable PGPASSWORD_POSTGRES")
	}

	metric := args[0]

	var name string
	name, err = cmd.Flags().GetString("name")
	if err != nil {
		return err
	}

	var namespace string
	namespace, err = cmd.Flags().GetString("namespace")
	if err != nil {
		return err
	}

	podName, err := KubeGetPodName(namespace, map[string]string{"release": name, "role": "master"})
	if err != nil {
		return err
	}

	err = KubePortForwardPod(namespace, podName, LISTEN_PORT_TSDB, FORWARD_PORT_TSDB)
	if err != nil {
		return err
	}

	var pool *pgxpool.Pool
	pool, err = pgxpool.Connect(context.Background(), "postgres://postgres:"+os.Getenv("PGPASSWORD_POSTGRES")+"@localhost:"+strconv.Itoa(LISTEN_PORT_TSDB)+"/postgres")
	if err != nil {
		return err
	}
	defer pool.Close()

	fmt.Printf("Resetting chunk interval for %v back to default\n", metric)
	_, err = pool.Exec(context.Background(), "SELECT prom_api.reset_metric_chunk_interval('"+metric+"')")
	if err != nil {
		return err
	}

	return nil
}
