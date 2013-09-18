package main

import "encoding/json"
import "io"
import "io/ioutil"
import "os"

const DEFAULT_PATH string = "/etc/goldilocks.conf"

type GLConfig struct {
    Meta map[string]interface{}  `json:"meta"`
    RPC map[string]string        `json:"rpc_alias"`
    Services []GLConfigService   `json:"services"`
    Schedules []GLConfigSchedule `json:"schedules"`
    Templates []GLConfigTemplate `json:"templates"`
}

type GLConfigService struct {
    Name string                `json:"name"`
    Description string         `json:"description"`
    Address string             `json:"address"`
    Threshold interface{}      `json:"threshold"`
    Commands GLConfigCommands  `json:"commands"`
}

type GLConfigCommands struct {
    Start  string `json:"start"`
    Stop   string `json:"stop"`
    Status string `json:"status"`
}

type GLConfigSchedule struct {
    From string         `json:"from"`
    To string           `json:"to"`
    Amount interface{}  `json:"amount"`
    Frequency string    `json:"frequency"`
}

type GLConfigTemplate struct {
    Name string   `json:"name"`
    Source string `json:"source"`
    Output string `json:"output"`
}

func GetConfigFromReader(r io.Reader) (config GLConfig, err error) {
    // Read full file into memory, because json unmarshal does not
    // accept io.Reader as a possible input.
    data, err := ioutil.ReadAll(r)
    if err != nil {
        return
    }
    err = json.Unmarshal(data, &config)
    return
}

func GetConfigReader(path string) (r io.Reader, err error) {
    // Attempt to load from specific path
    r, err = os.Open(path)
    if err != nil {
        r, err = os.Open(DEFAULT_PATH)
    }
    return
}

func GetConfig(path string) (config GLConfig, err error) {
    r, err := GetConfigReader(path)
    if err != nil {
        return
    }
    config, err = GetConfigFromReader(r)
    return
}
