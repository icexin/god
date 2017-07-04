package main

import (
	"time"

	"golang.org/x/net/context"

	"github.com/gogo/protobuf/proto"
	"github.com/icexin/god/pb"
	"github.com/jmoiron/sqlx"
)

const (
	agentStatusPending = iota
	agentStatusRunning
	agentStatusFinished
)

const (
	timeLayout = "2006-01-02 15:04:05"
)

func now() string {
	return time.Now().Format(timeLayout)
}

type GodMasterService struct {
	db *sqlx.DB
}

func (s *GodMasterService) newJob(desc *pb.JobDesc) (id int64, err error) {
	str, err := proto.Marshal(desc)
	if err != nil {
		return
	}
	result, err := s.db.Exec("INSERT INTO `job` (desc, start_time) VALUES (?, ?)", str, now())
	if err != nil {
		return
	}
	id, err = result.LastInsertId()

	tx, err := s.db.Begin()
	if err != nil {
		return
	}

	for _, agent := range desc.Agent {
		_, err = tx.Exec("INSERT INTO `task` (jobid, agent, output, status) VALUES (?, ?, ?, ?)", id, agent, "", agentStatusPending)
		if err != nil {
			return
		}
	}
	err = tx.Commit()
	return
}

func (s *GodMasterService) SubmitJob(ctx context.Context, req *pb.SubmitJobRequest) (*pb.SubmitJobResponse, error) {
	id, err := s.newJob(req.GetDesc())
	if err != nil {
		return nil, err
	}
	j := NewJobTracker(id, *req.GetDesc(), s.db)
	go j.Run()
	return &pb.SubmitJobResponse{Id: id}, nil
}

func (s *GodMasterService) StopJob(context.Context, *pb.StopJobRequest) (*pb.StopJobResponse, error) {
	return nil, nil
}

func (s *GodMasterService) ShowJob(context.Context, *pb.ShowJobRequest) (*pb.ShowJobResponse, error) {
	return nil, nil
}
