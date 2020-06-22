package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/cobra"
)

// chunkIntervalSetCmd represents the chunk-interval set command
var chunkIntervalSetCmd = &cobra.Command{
	Use:   "set <metric> <duration>",
	Short: "Sets chunk interval in minutes for a specific metric",
	Args:  cobra.ExactArgs(2),
	RunE:  chunkIntervalSet,
}

func init() {
	chunkIntervalCmd.AddCommand(chunkIntervalSetCmd)
}

func chunkIntervalSet(cmd *cobra.Command, args []string) error {
	var err error

	if os.Getenv("PGPASSWORD_POSTGRES") == "" {
		return errors.New("Password for postgres user must be set in environment variable PGPASSWORD_POSTGRES")
	}

	metric := args[0]
	var chunk_interval time.Duration
	chunk_interval, err = time.ParseDuration(args[1])
	if err != nil {
		return err
	}

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

	if chunk_interval.Minutes() < 1.0 {
		return errors.New("Chunk interval must be at least 1 minute")
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

	fmt.Printf("Setting chunk interval of %v to %v\n", metric, chunk_interval)
	_, err = pool.Exec(context.Background(), "SELECT prom_api.set_metric_chunk_interval('"+metric+"', INTERVAL '1 second' * "+strconv.FormatFloat(chunk_interval.Seconds(), 'f', -1, 64)+")")
	if err != nil {
		return err
	}

	return nil
}
