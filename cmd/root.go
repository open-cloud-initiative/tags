package cmd

import (
	"context"
	"fmt"
	"net"

	pb "github.com/open-cloud-initiative/specs/gen/go/tags/v1"
	"github.com/open-cloud-initiative/tags/internal/adapters/db"
	config "github.com/open-cloud-initiative/tags/internal/cfg"
	"github.com/open-cloud-initiative/tags/internal/controllers"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/katallaxie/pkg/dbx"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
)

var cfg = config.New()

const versionFmt = "%s (%s %s)"

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func Init() error {
	ctx := context.Background()

	err := cfg.InitDefaultConfig()
	if err != nil {
		return err
	}

	RootCmd.SilenceErrors = true
	RootCmd.SilenceUsage = true

	err = RootCmd.ExecuteContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

var RootCmd = &cobra.Command{
	Use:   "tags",
	Short: "tags",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runRoot(cmd.Context(), args...)
	},
	Version: fmt.Sprintf(versionFmt, version, commit, date),
}

func runRoot(ctx context.Context, _ ...string) error {
	conn, err := gorm.Open(postgres.Open(cfg.Flags.DatabaseURI), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{},
	})
	if err != nil {
		return err
	}

	store, err := dbx.NewDatabase(conn, db.NewReadTx(), db.NewWriteTx())
	if err != nil {
		return err
	}

	err = store.Migrate(ctx)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", cfg.Flags.Addr)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)

	tags := controllers.NewTagsController(store)

	pb.RegisterTagsServiceServer(srv, tags)
	reflection.Register(srv)

	if err := srv.Serve(lis); err != nil {
		return err
	}

	return nil
}
