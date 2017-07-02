package main

import (
	"golang.org/x/net/context"

	"github.com/icexin/god/pb"
	"github.com/jmoiron/sqlx"
)

type GodMasterService struct {
	id int64
	db *sqlx.DB
}

func (s *GodMasterService) SubmitJob(ctx context.Context, req *pb.SubmitJobRequest) (*pb.SubmitJobResponse, error) {
	s.id++
	j := NewJobTracker(s.id, *req.GetDesc())
	go j.Run()
	return new(pb.SubmitJobResponse), nil
}

func (s *GodMasterService) StopJob(context.Context, *pb.StopJobRequest) (*pb.StopJobResponse, error) {
	return nil, nil
}

func (s *GodMasterService) ShowJob(context.Context, *pb.ShowJobRequest) (*pb.ShowJobResponse, error) {
	return nil, nil
}
