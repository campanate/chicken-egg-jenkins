package cli

import (
	"fmt"
	"github.com/jessevdk/go-flags"
)

func Parse(osArgs []string, version string) (*CLI, error) {
	var cli CLI
	parser := flags.NewParser(&cli, flags.HelpFlag)
	parser.LongDescription = fmt.Sprintf(`Version %s
		This package creates a Jenkins server for you.`,
		version)
	args, err := parser.ParseArgs(osArgs[1:])
	if err != nil {
		return nil, err
	}
	if len(args) > 0 {
		return nil, fmt.Errorf("Too many argument")
	}
	return &cli, nil
}

type CLI struct {
	S3Bucket      string `long:"s3-bucket" default:"" env:"S3_BUCKET" description:"Bucket for backup"`
	S3File      string `long:"s3-file" default:"" env:"S3_FILE" description:"File in bucket for backup"`
	JenkinsUrl    string    `long:"jenkins-url" default:"jenkins.example.com" env:"JENKINS_URL" description:"DNS address of your jenkins"`
}
