package internal

import (
	"context"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"

	"github.com/ormanli/requestcounter/internal/app/requestcounter"
	"github.com/ormanli/requestcounter/internal/infra/counter/inmemory"
	"github.com/ormanli/requestcounter/internal/infra/counter/redis"
	"github.com/ormanli/requestcounter/internal/infra/logging"
	"github.com/ormanli/requestcounter/internal/infra/transport/http"
)

// Run starts application with the passed configuration.
func Run(ctx context.Context, cfg requestcounter.Config) error {
	logging.Setup(cfg)

	redisClient := goredis.NewClient(&goredis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort),
	})

	clusterCounter, err := redis.NewCounter(ctx, redisClient, cfg)
	if err != nil {
		return fmt.Errorf("can't initialize Redis counter: %w", err)
	}

	localCounter := &inmemory.Counter{}
	service := requestcounter.NewService(localCounter, clusterCounter)

	httpTransport := http.NewTransport(cfg, service, logging.Middleware)

	g := new(errgroup.Group)

	g.Go(func() error {
		return httpTransport.Run()
	})

	g.Go(func() error {
		<-ctx.Done()

		fiveSecondCtx, fiveSecondCncl := context.WithTimeout(context.Background(), 5*time.Second)
		defer fiveSecondCncl()

		err := httpTransport.Stop(fiveSecondCtx)
		if err != nil {
			return fmt.Errorf("can't stop http transport: %w", err)
		}

		err = redisClient.Close()
		if err != nil {
			return fmt.Errorf("can't close Redis client: %w", err)
		}

		return nil
	})

	return g.Wait()
}
