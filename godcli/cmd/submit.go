// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/icexin/god/pb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "submit job",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("usage: godcli submit [option] job.json")
		}
		buf, err := ioutil.ReadFile(args[0])
		if err != nil {
			log.Fatal(err)
		}

		var desc pb.JobDesc
		err = jsonpb.UnmarshalString(string(buf), &desc)
		if err != nil {
			log.Fatal(err)
		}

		host := viper.GetString("master")
		conn, err := grpc.Dial(host, grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		client := pb.NewGodMasterClient(conn)
		req := &pb.SubmitJobRequest{
			Desc: &desc,
		}
		resp, err := client.SubmitJob(context.TODO(), req)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp.Id)
	},
}

func init() {
	RootCmd.AddCommand(submitCmd)
}
