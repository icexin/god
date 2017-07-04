package main

import (
	"html/template"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

var gdb *sqlx.DB

var (
	jobstpl  = template.Must(template.New("jobs").Parse(string(MustAsset("tpl/jobs.tpl"))))
	jobtpl   = template.Must(template.New("job").Parse(string(MustAsset("tpl/job.tpl"))))
	agenttpl = template.Must(template.New("agent").Parse(string(MustAsset("tpl/agent.tpl"))))
)

type job struct {
	Id        string `db:"id"`
	StartTime string `db:"start_time"`
	EndTime   string `db:"end_time"`
	Status    int
}

// url: /jobs
func ListJobs(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var jobs []job
	err := gdb.Select(&jobs, "SELECT id, start_time, end_time FROM `job`")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = jobstpl.Execute(w, jobs)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

type task struct {
	Id        string `db:"jobid"`
	Agent     string `db:"agent"`
	StartTime string `db:"start_time"`
	EndTime   string `db:"end_time"`
	ExitCode  int    `db:"exit_code"`
	Status    int    `db:"status"`
}

// url: /job/$id
func ShowJob(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	jobid := ps[0].Value
	var tasks []task
	err := gdb.Select(&tasks, "SELECT jobid, agent, start_time, end_time, exit_code, status FROM `task` WHERE jobid=?", jobid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = jobtpl.Execute(w, tasks)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

// url: /job/$id/$agent
func ShowAgent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var output string
	jobid, agent := ps[0].Value, ps[1].Value
	err := gdb.Get(&output, "SELECT output FROM `task` WHERE jobid=? and agent=?", jobid, agent)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	err = agenttpl.Execute(w, output)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

}

func init() {
	r := httprouter.New()
	r.GET("/jobs", ListJobs)
	r.GET("/job/:id", ShowJob)
	r.GET("/job/:id/:agent/output", ShowAgent)
	http.Handle("/", r)
}
