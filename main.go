package main

import (
	"github.com/learninto/go-canal/cmd/job"
	_ "github.com/learninto/go-canal/init"

	"github.com/spf13/cobra"
)

func main() {
	root := cobra.Command{Use: "go-canal"}

	root.AddCommand(
		job.Cmd,
	)

	_ = root.Execute()
}
