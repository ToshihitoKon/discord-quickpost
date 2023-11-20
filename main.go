package main

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"runtime/debug"

	flag "github.com/spf13/pflag"
)

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		// Goモジュールが無効など
		return "(devel)"
	}
	return info.Main.Version
}

type CliOutput struct {
	// Channel   string `json:"channel"`
	// Timestamp string `json:"timestamp"`
}

type Options struct {
	token    string
	text     string
	filepath string
	postOpts *PostOptions

	// mode        string
	// blocks   string
	// snippetMode bool
	// slackClient SlackClient
}

type PostOptions struct {
	username string
	channel  string
	// iconEmoji string
	// iconUrl   string
	// threadTs  string
}

func strGetFirstOne(vars ...string) string {
	for _, v := range vars {
		if v != "" {
			return v
		}
	}
	return ""
}

func main() {
	var (
		// mode: print version
		printVersion = flag.Bool("version", false, "print version")

		// mode: post text
		optText     = flag.String("text", "", "post text")
		optTextFile = flag.String("textfile", "", "post text file path")
		// snippetMode = flag.Bool("snippet", false, "post text as snippet")

		// mode: post file
		filepath = flag.String("file", "", "post file path")

		// mode: post blocks json
		// optBlocks = flag.String("blocks", "", "post BlockKit json")

		// must options
		envToken   = os.Getenv("SLACK_TOKEN")
		optToken   = flag.String("token", "", "slack app OAuth token")
		optChannel = flag.String("channel", "", "post slack channel id")

		// optional
		envProfile = os.Getenv("SLACK_QUICKPOST_PROFILE")
		optProfile = flag.String("profile", "", "slack quickpost profile name")
		// threadTs   = flag.String("thread-ts", "", "post under thread")
		// iconEmoji  = flag.String("icon", "", "icon emoji")
		// iconUrl    = flag.String("icon-url", "", "icon image url")
		username = flag.String("username", "", "user name")

		noFail = flag.Bool("nofail", false, "always return success code(0)")

		errText []string
	)
	flag.Parse()

	if *printVersion {
		fmt.Println(version())
		os.Exit(0)
	}

	opts := &Options{
		// snippetMode: *snippetMode,
		filepath: *filepath,
		postOpts: &PostOptions{
			username: *username,
			// iconEmoji: *iconEmoji,
			// iconUrl:   *iconUrl,
			// threadTs:  *threadTs,
		},
	}
	_ = optText
	_ = optTextFile
	usr, err := user.Current()
	if err != nil {
		log.Printf("error: user.Current(). %s", err)
		os.Exit(1)
	}

	profileName := strGetFirstOne(*optProfile, envProfile)

	var profile = &Profile{}
	if profileName != "" {
		profPath := path.Join(usr.HomeDir, ".config", "discord-quickpost", profileName+".yaml")
		profile, err = parseProfile(profPath)
		if err != nil {
			errText = append(errText, fmt.Sprintf("error: failed read profile %s. %s", profPath, err.Error()))
		}
	}

	opts.token = strGetFirstOne(*optToken, envToken, profile.Token)
	if opts.token == "" {
		errText = append(errText, "error: slack token is required")
	}

	opts.postOpts.channel = strGetFirstOne(*optChannel, profile.Channel)
	if opts.postOpts.channel == "" {
		errText = append(errText, "error: channel is required")
	}

	if _, err := Do(opts); err != nil {
		fmt.Println(err.Error())
		if *noFail {
			os.Exit(0)
		}
		os.Exit(1)
	}

	return
}

func Do(opts *Options) (*CliOutput, error) {
	var output *CliOutput
	log.Println("Do")
	return output, nil
}
