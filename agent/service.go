package main

import (
	"io"
	"os/exec"
	"syscall"

	"github.com/icexin/god/pb"
)

type GodAgentService struct {
}

func (g *GodAgentService) RunJob(req *pb.RunJobRequest, stream pb.GodAgent_RunJobServer) error {
	cmd := exec.Command("bash", "-c", req.GetDesc().Cmd)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	cmd.Start()
	r := io.MultiReader(stdout, stderr)
	buf := make([]byte, 1024)
	for {
		n, err := r.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		resp := &pb.RunJobResponse{
			Body: string(buf[:n]),
		}
		err = stream.Send(resp)
		if err != nil {
			return err
		}
	}
	code := int32(0)
	err := cmd.Wait()
	if err != nil {
		code = int32(err.(*exec.ExitError).Sys().(syscall.WaitStatus).ExitStatus())
	}
	resp := &pb.RunJobResponse{Code: code, Stoped: true}
	return stream.Send(resp)
}
