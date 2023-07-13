package main

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"os"
	"path"
	"time"

	"github.com/VaiTon/uniboappelli/appelli"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

type ConfigEntry struct {
	Degree          string `toml:"degree"`
	TelegramChannel string `toml:"channel,omitempty"`
}

type Config struct {
	Degrees  []ConfigEntry `toml:"degrees"`
	BotToken string        `toml:"bot_token"`
}

var (
	rootCmd = &cobra.Command{
		Use:   "appelli",
		Short: "",
		Long:  "",
		Run:   run,
	}

	debugFlag  *bool
	jsonFlag   *bool
	configFile *string

	config Config
)

func main() {
	debugFlag = rootCmd.Flags().Bool("debug", false, "sets log level to debug")
	jsonFlag = rootCmd.Flags().Bool("json", false, "output json to console instead of pretty printed text")
	configFile = rootCmd.Flags().StringP("config", "c", "config.toml", "config file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(*cobra.Command, []string) {

	// Set global log level to debugFlag or info
	if *debugFlag {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Set log output to json or pretty printed text
	if !*jsonFlag {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, PartsExclude: []string{"time"}})
	}

	// Read config file
	readConfig()

	// Start analysis for each degree
	for _, course := range config.Degrees {
		log.Debug().Str("degree", course.Degree).Msg("starting analysis")
		analyzeDegree(course)
	}
}

func readConfig() {
	configData, err := os.ReadFile(*configFile)
	if err != nil {
		if os.IsNotExist(err) {
			log.Error().Msgf("config file '%s' does not exist", *configFile)
			log.Fatal().Msg("create a config file or run with --help for more info")
		} else {
			log.Fatal().Err(err).Msg("could not read config file")
		}
	}

	err = toml.Unmarshal(configData, &config)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse config file. run with --help for more info")
	}
}

func analyzeDegree(entry ConfigEntry) {
	degree := entry.Degree

	// Create degree specific logger
	degreeLog := log.With().Str("degree", degree).Logger()

	// Get published exams
	newExams, err := appelli.GetExams(degree)
	if err != nil {
		degreeLog.Err(err).Msg("could not get exams")
		return
	}
	degreeLog.Info().Int("exams", len(newExams)).Msg("exams found")

	// Create data directory if it does not exist
	dataPath := "data"
	err = os.MkdirAll(dataPath, 0755)
	if err != nil {
		degreeLog.Fatal().Err(err).Msg("could not create data directory")
	}

	dataFileName := fmt.Sprintf("%s.json", degree)
	dataFileName = path.Join(dataPath, dataFileName)

	// Read old exams from data file if it exists
	timedOld := TimedExams{}
	data, err := readDataFile(dataFileName)
	if err != nil {
		degreeLog.Fatal().Err(err).Msg("could not read data file")
	} else if data != nil {
		timedOld = data
	}

	diff := appelli.Diff(newExams, timedOld.ToExams())
	if len(diff) == 0 {
		degreeLog.Info().Msg("no new exams found")
		return
	}

	degreeLog.Info().Int("exams", len(diff)).Msg("new exams found")

	timedDiff := NewTimedExams(diff, time.Now())

	createFeedFile(degree, appelli.GetExamsUrl(degree), timedDiff, timedOld, degreeLog)

	timedOld = append(timedOld, timedDiff...)
	err = saveDataFile(dataFileName, timedOld)
	if err != nil {
		log.Fatal().Err(err).Msg("could not save data file")
	}

	// Post diff to telegram channel (if configured)
	postToTelegram(timedDiff, entry.TelegramChannel, degreeLog)
}
