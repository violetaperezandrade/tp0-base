package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/bet"
	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/common"
)

// InitConfig Function that uses viper library to parse configuration parameters.
// Viper is configured to read variables from both environment variables and the
// config file ./config.yaml. Environment variables takes precedence over parameters
// defined in the configuration file. If some of the variables cannot be parsed,
// an error is returned
func InitConfig() (*viper.Viper, error) {
	v := viper.New()

	// Configure viper to read env variables with the CLI_ prefix
	v.AutomaticEnv()
	v.SetEnvPrefix("cli")
	// Use a replacer to replace env variables underscores with points. This let us
	// use nested configurations in the config file and at the same time define
	// env variables for the nested configurations
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Add env variables supported
	v.BindEnv("id")
	v.BindEnv("server", "address")
	v.BindEnv("loop", "period")
	v.BindEnv("loop", "lapse")
	v.BindEnv("log", "level")

	// Try to read configuration from config file. If config file
	// does not exists then ReadInConfig will fail but configuration
	// can be loaded from the environment variables so we shouldn't
	// return an error in that case
	v.SetConfigFile("./config.yaml")
	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("Configuration could not be read from config file. Using env variables instead")
	}

	// Parse time.Duration variables and return an error if those variables cannot be parsed
	if _, err := time.ParseDuration(v.GetString("loop.lapse")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse CLI_LOOP_LAPSE env var as time.Duration.")
	}

	if _, err := time.ParseDuration(v.GetString("loop.period")); err != nil {
		return nil, errors.Wrapf(err, "Could not parse CLI_LOOP_PERIOD env var as time.Duration.")
	}

	return v, nil
}

// InitLogger Receives the log level to be set in logrus as a string. This method
// parses the string and set the level to the logger. If the level string is not
// valid an error is returned
func InitLogger(logLevel string) error {
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return err
	}

	customFormatter := &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   false,
	}
	logrus.SetFormatter(customFormatter)
	logrus.SetLevel(level)
	return nil
}

// PrintConfig Print all the configuration parameters of the program.
// For debugging purposes only
func PrintConfig(v *viper.Viper) {
	logrus.Infof("action: config | result: success | client_id: %s | server_address: %s | loop_lapse: %v | loop_period: %v | log_level: %s",
		v.GetString("id"),
		v.GetString("server.address"),
		v.GetDuration("loop.lapse"),
		v.GetDuration("loop.period"),
		v.GetString("log.level"),
	)
}

func betFromRegister(register []string, agency string) bet.Bet {

	firstName := register[0]
	lastName := register[1]
	dni, _ := strconv.Atoi(register[2])
	number, _ := strconv.Atoi(register[4])
	birthDate, _ := time.Parse("2006-01-02", register[3])
	agencyId, _ := strconv.Atoi(agency)

	return bet.Bet{
		Agency:    agencyId,
		Dni:       dni,
		Number:    number,
		Year:      birthDate.Year(),
		Month:     int(birthDate.Month()),
		Day:       birthDate.Day(),
		FirstName: firstName,
		LastName:  lastName,
	}
}

func main() {
	v, err := InitConfig()
	if err != nil {
		log.Fatalf("%s", err)
	}

	if err := InitLogger(v.GetString("log.level")); err != nil {
		log.Fatalf("%s", err)
	}

	// Print program config with debugging purposes
	PrintConfig(v)

	batchSize, _ := strconv.Atoi(os.Getenv("CHUNK_SIZE"))

	clientConfig := common.ClientConfig{
		ServerAddress: v.GetString("server.address"),
		ID:            v.GetString("id"),
		LoopLapse:     v.GetDuration("loop.lapse"),
		LoopPeriod:    v.GetDuration("loop.period"),
		Agency:        os.Getenv("CLI_ID"),
		BatchSize:     batchSize,
	}
	filePath := fmt.Sprintf("/.data/dataset/agency-%s.csv", os.Getenv("CLI_ID"))
	log.Infof("FILEPATH: %s", filePath)

	client := common.NewClient(clientConfig)
	log.Infof("Client created succesfully")
	var betsList []bet.Bet

	betsParsed := 0
	batches := 0

	readFile, err := os.Open(filePath)

	if err != nil {
		log.Error("Error creating reader: %s", err)
	}
	fileScanner := bufio.NewScanner(readFile)

	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		if betsParsed == clientConfig.BatchSize {
			batches += 1
			//log.Debugf("Sending batch number %v for agency %v", batches, clientConfig.Agency)
			client.SendBets(betsList, batches)
			betsParsed = 0
			betsList = []bet.Bet{}
		}
		line := fileScanner.Text()
		register := strings.Split(line, ",")
		bet := betFromRegister(register, clientConfig.Agency)
		betsList = append(betsList, bet)
		betsParsed += 1
	}

	if len(betsList) != 0 {
		//log.Debugf("Sending batch number %v for agency %v", batches, clientConfig.Agency)
		client.SendBets(betsList, batches)

	}

	readFile.Close()
}
