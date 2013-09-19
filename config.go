package main

import "encoding/json"
import "io"
import "io/ioutil"
import "os"
import "reflect"

const DEFAULT_PATH      string = "/etc/goldilocks.conf"
const ENV_VARIABLE_NAME string = "GOLDILOCKS_CONFIG"

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

func GetConfigReader(paths []string) (r io.Reader, err error) {
    // Get first working reader from paths slice
    for _, path := range paths {
        r, err = os.Open(path)
        if err == nil { return }
    }
    return
}

func GetConfig(path string) (config GLConfig, err error) {
    paths := []string{
        path,
        os.Getenv(ENV_VARIABLE_NAME),
        DEFAULT_PATH,
    }

    r, err := GetConfigReader(paths)
    if err != nil {
        return
    }
    config, err = GetConfigFromReader(r)
    if err != nil {
        return
    }
    err = ValidateConfig(&config)
    if err != nil {
        return
    }
    return
}

type GLConfigValidationError struct {
    Location string
    Problem string
}

func (e *GLConfigValidationError) Error() string {
    return "Failure to validate " + e.Location + ": " + e.Problem
}

func ValidateConfStruct(s interface{}) (err error) {
    structtype := reflect.TypeOf(s)
    structval  := reflect.ValueOf(s)
    structname := structtype.Name()

    for i := 0; i < structtype.NumField(); i++ {
        field := structval.Field(i)
        fieldname := structtype.Field(i).Name

        if field.Kind() == reflect.String && field.String() == "" {
            err = &GLConfigValidationError{
                structname + "." + fieldname,
                "was blank string",
            }
            return
        } else if field.Kind() == reflect.Struct {
            err = ValidateConfStruct(field.Interface())
            if err != nil { return }
        }
    }
    return
}

func ValidateConfig(config *GLConfig) (err error) {
    // Ensure that all necessary data for operation is present
    
    default_rpc, ok := config.RPC["default"]
    if ! ok { 
        err = &GLConfigValidationError{
            "GLConfig.RPC",
            "No default RPC address",
        }
        return
    }

    for _, service := range config.Services {
        if service.RPC == "" { 
            service.RPC = default_rpc
        }

        err = ValidateConfStruct(service)
        if err != nil { return }
    }

    for _, schedule := range config.Schedules {
        if schedule.RPC == "" { 
            schedule.RPC = default_rpc
        }

        err = ValidateConfStruct(schedule)
        if err != nil { return }
    }

    for _, template := range config.Templates {
        err = ValidateConfStruct(template)
        if err != nil { return }
    }
    return
}
