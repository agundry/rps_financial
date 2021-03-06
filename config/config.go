package config

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// Config struct for webapp config
type Config struct {
	Server struct {
		// Port is the local machine TCP Port to bind the HTTP Server to
		Port    string `yaml:"port"`
	} `yaml:"server"`

	Db struct {
		Username    string `yaml:"username"`
		Password	string `yaml:"password"`
		Server 		string `yaml:"server"`
		DbName		string `yaml:"db_name"`
	} `yaml:"db"`
}

func Configure() (*Config) {
	// Process config
	cfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := NewConfig(cfgPath)

	// In prod we'll use environment variables for the password
	dbPassword, ok := os.LookupEnv("RPS_DB_PASSWORD")
	if ok {
		cfg.Db.Password = dbPassword
	}

	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}
