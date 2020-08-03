package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

type Config struct {
	Ke struct {
		YumPath string `yaml:"yum_path"`
		ToPath string `yaml:"to_path"`
		AssetsPath string `yaml:"assets_path"`
	}
}

func getConfig() Config {
	file, err := ioutil.ReadFile("./pubspec.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(500)
	}
	config := Config{}
	err = yaml.Unmarshal(file, &config)
	if (err != nil) {
		fmt.Println(err)
		os.Exit(500)
	}

	return config
}

func main()  {
	config := getConfig()
	println("start...")
	Move(config.Ke.YumPath, config.Ke.ToPath)
	if config.Ke.AssetsPath != "" {
		GenerateAssetsClass(config.Ke.ToPath, config.Ke.AssetsPath)
	}
	println("start success!")
}
