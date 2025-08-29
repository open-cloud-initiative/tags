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
	"github.com/katallaxie/pkg/server"
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

var _ server.Listener = (*WebSrv)(nil)

// WebSrv is the server that implements the Noop interface.
type WebSrv struct {
	cfg *config.Config
}

// NewWebSrv returns a new instance of NoopSrv.
func NewWebSrv(cfg *config.Config) *WebSrv {
	return &WebSrv{cfg}
}

// Start starts the server.
func (s *WebSrv) Start(ctx context.Context, ready server.ReadyFunc, run server.RunFunc) func() error {
	return func() error {
		// Start the server
		return nil
	}
}

func Init() error {
	ctx := context.Background()

	err := cfg.InitDefaultConfig()
	if err != nil {
		return err
	}

	RootCmd.PersistentFlags().BoolVarP(&cfg.Flags.Verbose, "verbose", "v", cfg.Flags.Verbose, "verbose output")
	RootCmd.PersistentFlags().BoolVarP(&cfg.Flags.Dry, "dry", "d", cfg.Flags.Dry, "dry run")
	RootCmd.PersistentFlags().BoolVarP(&cfg.Flags.Root, "root", "r", cfg.Flags.Root, "run as root")
	RootCmd.PersistentFlags().BoolVarP(&cfg.Flags.Force, "force", "f", cfg.Flags.Force, "force init")

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
	s, _ := server.WithContext(ctx)
	s.SetLimit(3)

	conn, err := gorm.Open(postgres.Open(""), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{},
	})
	if err != nil {
		return err
	}

	store, err := dbx.NewDatabase(conn, db.NewReadTx(), db.NewWriteTx())
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", ":4040")
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
