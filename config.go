package main

import "encoding/json"
import "io"
import "io/ioutil"
import "os"

const DEFAULT_PATH string = "/etc/goldilocks.conf"

type GLConfig struct {
    Meta map[string]interface{}  `json:"meta"`
    RPC  map[string]string       `json:"rpc_alias"`
    Services  []GLConfigService  `json:"services"`
    Schedules []GLConfigSchedule `json:"schedules"`
    Templates []GLConfigTemplate `json:"templates"`
}

type GLConfigService struct {
    Name        string           `json:"name"`
    Description string           `json:"description"`
    Address     string           `json:"address"`
    Threshold   string           `json:"threshold"`
    Commands    GLConfigCommands `json:"commands"`
    RPC         string           `json:"rpc_alias"`
}

type GLConfigCommands struct {
    Start  string `json:"start"`
    Stop   string `json:"stop"`
    Status string `json:"status"`
}

type GLConfigSchedule struct {
    Name      string `json:"name"`
    From      string `json:"from"`
    To        string `json:"to"`
    Amount    string `json:"amount"`
    Frequency string `json:"frequency"`
    RPC       string `json:"rpc_alias"`
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

func ValidateConfig(config *GLConfig) (bool) {
    // Ensure that all necessary data for operation is present
    
    default_rpc, ok := config.RPC["default"]
    if ! ok { return false }

    for _, service := range config.Services {
        if service.Name    == "" { return false }
        if service.Address == "" { return false }
        if service.Commands.Start  == "" { return false }
        if service.Commands.Stop   == "" { return false }
        if service.Commands.Status == "" { return false }

        if service.RPC == "" { 
            service.RPC = default_rpc
        }
    }

    for _, schedule := range config.Schedules {
        if schedule.Name      == "" { return false }
        if schedule.From      == "" { return false }
        if schedule.To        == "" { return false }
        if schedule.Amount    == "" { return false }
        if schedule.Frequency == "" { return false }

        if schedule.RPC == "" { 
            schedule.RPC = default_rpc
        }
    }

    for _, template := range config.Templates {
        if template.Name   == "" { return false }
        if template.Source == "" { return false }
        if template.Output == "" { return false }
    }

    return true
}
