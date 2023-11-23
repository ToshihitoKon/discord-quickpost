package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/bwmarrin/discordgo"
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
	session  *discordgo.Session
	token    string
	postOpts *PostOptions
}

type PostOptions struct {
	channel   string
	text      string
	filePaths []string
}

func main() {
	var (
		optText = flag.String("text", "", "post text")
		optFile = flag.String("file", "", "post file path")

		envToken = os.Getenv("DISCORD_TOKEN")
		optToken = flag.String("token", "", "discord app token")

		optChannel = flag.String("channel", "", "post slack channel id")

		// envProfile = os.Getenv("DISCORD_QUICKPOST_PROFILE")
		// optProfile = flag.String("profile", "", "discord quickpost profile name")

		noFail       = flag.Bool("nofail", false, "always return success code(0)")
		printVersion = flag.Bool("version", false, "print version")

		errText []string
	)
	flag.Parse()

	if *printVersion {
		fmt.Println(version())
		os.Exit(0)
	}

	opts := &Options{
		postOpts: &PostOptions{},
	}

	opts.postOpts.text = *optText

	// var profile = &profile{}
	// profileName := strGetFirstOne(*optProfile, envProfile)
	// if profileName != "" {
	// 	var err error
	// 	profile, err = getProfile(profileName)
	// 	if err != nil {
	// 		errText = append(errText, fmt.Sprintf("error: failed read profile %s. %s", profPath, err.Error()))
	// 	}
	// }

	// opts.token = strGetFirstOne(*optToken, envToken, profile.Token)
	opts.token = strGetFirstOne(*optToken, envToken)
	if opts.token == "" {
		errText = append(errText, "error: slack token is required")
	}

	// opts.postOpts.channel = strGetFirstOne(*optChannel, profile.Channel)
	opts.postOpts.channel = strGetFirstOne(*optChannel)
	if opts.postOpts.channel == "" {
		errText = append(errText, "error: channel is required")
	}

	if *optFile != "" {
		opts.postOpts.filePaths = []string{*optFile}
	}

	// Discord initialize
	session, err := discordgo.New("Bot " + opts.token)
	if err != nil {
		fmt.Println(err.Error())
		if *noFail {
			os.Exit(0)
		}
		os.Exit(1)
	}
	opts.session = session

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
	var sendParameter = &discordgo.MessageSend{}
	if opts.postOpts.text != "" {
		sendParameter.Content = opts.postOpts.text
	}

	// TODO: 2つ以上に対応する
	var f *os.File
	if len(opts.postOpts.filePaths) > 0 {
		var err error
		if f, err = os.Open(opts.postOpts.filePaths[0]); err != nil {
			return nil, err
		}
		defer f.Close()

		file := &discordgo.File{
			Name:        f.Name(),
			ContentType: "text/plain",
			Reader:      f,
		}

		sendParameter.Files = []*discordgo.File{file}
	}

	if _, err := opts.session.ChannelMessageSendComplex(
		opts.postOpts.channel,
		sendParameter,
	); err != nil {
		return nil, err
	}

	return output, nil
}
