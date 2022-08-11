package svc

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/bnb-chain/zkbas/common/model/proof"

	"github.com/bnb-chain/zkbas/common/model/blockForProof"
	"github.com/bnb-chain/zkbas/service/cronjob/prover/internal/config"
)

type ServiceContext struct {
	Config config.Config

	RedisConn *redis.Redis

	ProofSenderModel   proof.ProofModel
	BlockForProofModel blockForProof.BlockForProofModel
}

func WithRedis(redisType string, redisPass string) redis.Option {
	return func(p *redis.Redis) {
		p.Type = redisType
		p.Pass = redisPass
	}
}
func NewServiceContext(c config.Config) *ServiceContext {
	gormPointer, err := gorm.Open(postgres.Open(c.Postgres.DataSource))
	if err != nil {
		logx.Errorf("gorm connect db error, err = %s", err.Error())
	}
	conn := sqlx.NewSqlConn("postgres", c.Postgres.DataSource)
	redisConn := redis.New(c.CacheRedis[0].Host, WithRedis(c.CacheRedis[0].Type, c.CacheRedis[0].Pass))
	return &ServiceContext{
		Config:             c,
		RedisConn:          redisConn,
		BlockForProofModel: blockForProof.NewBlockForProofModel(conn, c.CacheRedis, gormPointer),
		ProofSenderModel:   proof.NewProofModel(gormPointer),
	}
}
