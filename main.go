package main

import (
	"log"

	pb "github.com/beatitudes/shippy-service-user/proto/user"
	"github.com/micro/go-micro/v2"
)

const schema = `
	create table if not exists users (
		id varchar(36) not null,
		name varchar(125) not null,
		email varchar(225) not null unique,
		password varchar(255) not null,
		company varchar(125),
		primary key(id)
	)
`

func main() {
	db, err := NewConnection()
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	db.MustExec(schema)
	repo := NewPostgresRepository(db)
	tokenService := &TokenService{repo}
	service := micro.NewService(
		micro.Name("shippy.service.user"),
		micro.Version("latest"),
	)
	service.Init()

	if err := pb.RegisterUserServiceHandler(service.Server(), &handler{repo, tokenService}); err != nil {
		log.Panic(err)
	}
	if err := service.Run(); err != nil {
		log.Panic(err)
	}
}
