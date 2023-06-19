package testing

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/oslokommune/common-lib-go/db"
	"github.com/testcontainers/testcontainers-go"
	"testing"
)

const (
	image = "postgres:14.2-alpine"
)

//EmbeddedPostgres spins up a postgres container.
func ContainerizedPostgres(t *testing.T, conf *db.DbConf) {
	t.Helper()
	ctx := context.Background()
	natPort := fmt.Sprintf("%d/tcp", conf.Port)

	// Setup container
	req := testcontainers.ContainerRequest{
		Image:        image,
		ExposedPorts: []string{natPort},
		Env: map[string]string{
			"POSTGRES_PASSWORD": conf.Password,
			"POSTGRES_USER":     conf.Username,
			"POSTGRES_DATABASE": conf.Database,
		},
		WaitingFor: NewPostgresStrategy(nat.Port(natPort), conf),
	}

	// Start container
	pg, err := testcontainers.GenericContainer(
		ctx,
		testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		},
	)
	if err != nil {
		t.Error(err)
	}

	// When test is done terminate container
	t.Cleanup(func() {
		_ = pg.Terminate(ctx)
	})

	// Get the container info needed
	mp, err := pg.MappedPort(ctx, nat.Port(natPort))
	if err != nil {
		t.Error(err)
	}
	ma, err := pg.Host(ctx)
	if err != nil {
		t.Error(err)
	}

	// Update the config with the containers host and portnumber
	conf.UpdateHostAndPort(ma, mp.Int())
}
