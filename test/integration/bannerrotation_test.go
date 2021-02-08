package integration

import (
	"context"
	"github.com/jmoiron/sqlx"
	grpcgenerated "github.com/oleglarionov/otusgolang_finalproject/internal/grpc/generated"
	"github.com/pressly/goose"
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"

	// init driver.
	_ "github.com/lib/pq"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"testing"
)

func init() {
	viper.SetConfigFile(os.Getenv("ENV_FILE"))
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

type BannerRotationSuite struct {
	suite.Suite
	client    grpcgenerated.BannerRotationServiceClient
	ctx       context.Context
	ctxCancel context.CancelFunc
	db        *sqlx.DB
	mu        sync.Mutex
}

func (s *BannerRotationSuite) SetupSuite() {
	s.ctx, s.ctxCancel = context.WithCancel(context.Background())

	host := viper.GetString("APP_HOST")
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		s.FailNow(err.Error())
	}

	s.client = grpcgenerated.NewBannerRotationServiceClient(conn)

	dbDsn := viper.GetString("DB_DSN")
	s.db, err = sqlx.Connect("postgres", dbDsn)
	if err != nil {
		s.FailNow(err.Error())
	}
}

func (s *BannerRotationSuite) TearDownSuite() {
	defer s.ctxCancel()
	defer s.db.Close()
}

func (s *BannerRotationSuite) SetupTest() {
	s.mu.Lock()
	s.cleanupDb()
}

func (s *BannerRotationSuite) TearDownTest() {
	s.mu.Unlock()
}

func (s *BannerRotationSuite) TestAddBanner() {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	_, err := s.client.AddBanner(ctx, &grpcgenerated.AddBannerRequest{
		SlotId:   "slot-1",
		BannerId: "banner-1",
	})
	s.Require().NoError(err)

	response, err := s.client.ChooseBanner(ctx, &grpcgenerated.ChooseBannerRequest{
		SlotId:      "slot-1",
		UserGroupId: "group-1",
	})
	s.Require().NoError(err)
	s.Require().Equal("banner-1", response.BannerId)
}

func (s *BannerRotationSuite) TestRemoveBanner() {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	_, err := s.client.AddBanner(ctx, &grpcgenerated.AddBannerRequest{
		SlotId:   "slot-1",
		BannerId: "banner-1",
	})
	s.Require().NoError(err)

	_, err = s.client.RemoveBanner(ctx, &grpcgenerated.RemoveBannerRequest{
		SlotId:   "slot-1",
		BannerId: "banner-1",
	})
	s.Require().NoError(err)

	response, err := s.client.ChooseBanner(ctx, &grpcgenerated.ChooseBannerRequest{
		SlotId:      "slot-1",
		UserGroupId: "group-1",
	})
	s.Require().NoError(err)
	s.Require().False(response.BannerFound)
}

func (s *BannerRotationSuite) TestRegisterClick() {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	_, err := s.client.AddBanner(ctx, &grpcgenerated.AddBannerRequest{
		SlotId:   "slot-1",
		BannerId: "banner-1",
	})
	s.Require().NoError(err)

	_, err = s.client.ChooseBanner(ctx, &grpcgenerated.ChooseBannerRequest{
		SlotId:      "slot-1",
		UserGroupId: "group-1",
	})
	s.Require().NoError(err)

	_, err = s.client.RegisterClick(ctx, &grpcgenerated.RegisterClickRequest{
		SlotId:      "slot-1",
		BannerId:    "banner-1",
		UserGroupId: "group-1",
	})
	s.Require().NoError(err)
}

func (s *BannerRotationSuite) TestChooseBanner() {
	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	banners := []struct {
		slot   string
		banner string
	}{
		{"slot-1", "banner-1"},
		{"slot-1", "banner-2"},
		{"slot-1", "banner-3"},
		{"slot-2", "banner-1"},
		{"slot-2", "banner-4"},
	}
	for _, banner := range banners {
		_, err := s.client.AddBanner(ctx, &grpcgenerated.AddBannerRequest{
			SlotId:   banner.slot,
			BannerId: banner.banner,
		})
		s.Require().NoError(err)
	}

	viewsByBanner := make(map[string]int)
	n := 100
	for i := 0; i < n; i++ {
		response, err := s.client.ChooseBanner(ctx, &grpcgenerated.ChooseBannerRequest{
			SlotId:      "slot-1",
			UserGroupId: "group-1",
		})
		s.Require().NoError(err)

		banner := response.BannerId
		viewsByBanner[banner] += 1

		if banner == "banner-1" {
			_, err = s.client.RegisterClick(ctx, &grpcgenerated.RegisterClickRequest{
				SlotId:      "slot-1",
				BannerId:    banner,
				UserGroupId: "group-1",
			})
			s.Require().NoError(err)
		}
	}

	for _, views := range viewsByBanner {
		s.Require().GreaterOrEqual(views, 1)
	}
}

func (s *BannerRotationSuite) cleanupDb() {
	_, err := s.db.ExecContext(s.ctx, "DROP SCHEMA public CASCADE; "+
		"CREATE SCHEMA public; "+
		"GRANT ALL ON SCHEMA public TO postgres; "+
		"GRANT ALL ON SCHEMA public TO public;",
	)
	if err != nil {
		s.FailNow(err.Error())
	}

	err = goose.Up(s.db.DB, "../../migrations")
	if err != nil {
		s.FailNow(err.Error())
	}
}

func TestBannerRotationSuite(t *testing.T) {
	suite.Run(t, new(BannerRotationSuite))
}
