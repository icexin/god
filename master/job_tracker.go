package main

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	"github.com/icexin/god/pb"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type JobTracker struct {
	id   int64
	desc pb.JobDesc
	ch   chan string
	wait sync.WaitGroup
	idx  int
	db   *sqlx.DB
}

func NewJobTracker(id int64, desc pb.JobDesc, db *sqlx.DB) *JobTracker {
	return &JobTracker{
		id:   id,
		desc: desc,
		ch:   make(chan string),
		db:   db,
	}
}

func (j *JobTracker) pickAgent() string {
	if j.idx >= len(j.desc.Agent) {
		return ""
	}
	i := j.idx
	j.idx++
	return j.desc.Agent[i]
}

func (j *JobTracker) updateAgentOutput(agent, content string) error {
	_, err := j.db.Exec("UPDATE `task` SET output=output||? WHERE jobid=? and agent=?",
		content, j.id, agent)
	if err != nil {
		log.Printf("update error:%s", err)
	}
	return err
}

func (j *JobTracker) recordJobFinished() error {
	_, err := j.db.Exec("UPDATE `job` SET end_time=? WHERE id=?",
		now(), j.id)
	if err != nil {
		return err
	}
	log.Printf("job %d finished", j.id)
	return nil
}

func (j *JobTracker) Run() {
	for i := int32(0); i < j.desc.Concurrent; i++ {
		j.wait.Add(1)
		go j.workloop()
	}

	for {
		agent := j.pickAgent()
		if agent == "" {
			break
		}
		j.ch <- agent
		time.Sleep(time.Duration(j.desc.Interval))
	}
	close(j.ch)
	j.wait.Wait()
	j.recordJobFinished()
}

func (j *JobTracker) runJob(agent string) error {
	log.Printf("%s started.", agent)
	conn, err := grpc.Dial(agent, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewGodAgentClient(conn)
	req := &pb.RunJobRequest{
		Id:   j.id,
		Desc: &j.desc,
	}
	stream, err := client.RunJob(context.TODO(), req)
	if err != nil {
		return err
	}
	var resp *pb.RunJobResponse
	for {
		resp, err = stream.Recv()
		if err == io.EOF {
			break
		}
		j.updateAgentOutput(agent, resp.Body)
		if resp.Stoped {
			_, err := j.db.Exec("UPDATE `task` SET status=?, end_time=?, exit_code=? WHERE jobid=? and agent=?",
				agentStatusFinished, now(), resp.Code, j.id, agent)
			if err != nil {
				log.Printf("error when update agent:%s", err)
			}
			log.Printf("%s finished, code:%d", agent, resp.Code)
		}
	}

	return nil

}

func (j *JobTracker) workloop() {
	defer j.wait.Done()
	for agent := range j.ch {
		_, err := j.db.Exec("UPDATE `task` SET status=?, start_time=? WHERE jobid=? and agent=?",
			agentStatusRunning, now(), j.id, agent)
		if err != nil {
			log.Print(err)
			continue
		}
		j.runJob(agent)
	}
}
