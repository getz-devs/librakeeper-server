package suite

import (
	"context"
	searcherv1 "github.com/getz-devs/librakeeper-protos/gen/go/searcher"
	"github.com/getz-devs/librakeeper-server/internal/searcher/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"

	"testing"
)

type Suite struct {
	*testing.T
	Cfg          *config.Config
	SearcherServ searcherv1.SearcherClient
}

const (
	grpcHost = "localhost"
)

func New(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	// TODO: Read test config from env

	cfgManager := config.MustLoadByPath("../config/searcher/local.yaml")
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfgManager.GRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	//cc, err := grpc.DialContext(ctx,
	//	grpcAdress(&cfgManager.GRPC),
	//	grpc.WithTransportCredentials(insecure.NewCredentials()),
	//)
	cc, err := grpc.NewClient(grpcAdress(&cfgManager.GRPC), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		t.Fatalf("grpc server connection failed %v", err)
	}

	return ctx, &Suite{
		T:            t,
		Cfg:          cfgManager,
		SearcherServ: searcherv1.NewSearcherClient(cc),
	}

}

func grpcAdress(cfg *config.GRPCConfig) string {
	return net.JoinHostPort(grpcHost, strconv.Itoa(cfg.Port))
}
