package main

import (
	"github.com/youngshawn/go-project-demo/course/cmd"
	_ "github.com/youngshawn/go-project-demo/course/config"
	_ "github.com/youngshawn/go-project-demo/course/models"
)

func main() {

	cmd.Execute()

}
