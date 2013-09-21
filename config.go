package main

import "encoding/json"
import "io"
import "io/ioutil"
import "os"
import "reflect"

const DEFAULT_PATH      string = "/etc/goldilocks.conf"
const ENV_VARIABLE_NAME string = "GOLDILOCKS_CONFIG"

type GLConfig struct {
    Meta map[string]interface{}           `json:"meta"`
    RPC  map[string]string                `json:"rpc_alias"`
    Services  map[string]GLConfigService  `json:"services"`
    Schedules map[string]GLConfigSchedule `json:"schedules"`
    Templates map[string]GLConfigTemplate `json:"templates"`
}

type GLConfigService struct {
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
    From      string `json:"from"`
    To        string `json:"to"`
    Amount    string `json:"amount"`
    Frequency string `json:"frequency"`
    RPC       string `json:"rpc_alias"`
}

type GLConfigTemplate struct {
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

func GetConfig(paths []string) (config GLConfig, err error) {
    paths = append(
        paths,
        os.Getenv(ENV_VARIABLE_NAME),
        DEFAULT_PATH,
    )

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
    
    _, ok := config.RPC["default"]
    if ! ok { 
        err = &GLConfigValidationError{
            "GLConfig.RPC",
            "No default RPC address",
        }
        return
    }

    for name, service := range config.Services {
        if service.RPC == "" { 
            service.RPC = "default"
            config.Services[name] = service
        }

        if name == "" {
            err = &GLConfigValidationError{
                "GLConfig.Services",
                "Blank name",
            }
            return
        }
        err = ValidateConfStruct(service)
        if err != nil { return }
    }

    for name, schedule := range config.Schedules {
        if schedule.RPC == "" { 
            schedule.RPC = "default"
            config.Schedules[name] = schedule
        }

        if name == "" {
            err = &GLConfigValidationError{
                "GLConfig.Schedules",
                "Blank name",
            }
            return
        }
        err = ValidateConfStruct(schedule)
        if err != nil { return }
    }

    for name, template := range config.Templates {
        if name == "" {
            err = &GLConfigValidationError{
                "GLConfig.Templates",
                "Blank name",
            }
            return
        }

        err = ValidateConfStruct(template)
        if err != nil { return }
    }
    return
}

func ConfigSanitize(config *GLConfig) {
    config.RPC = make(map[string]string)
}

func ConfigDump(config *GLConfig, w io.Writer) (error) {
    // Serialized, sanitized, and ready to write.
    ConfigSanitize(config)

    data, err := json.MarshalIndent(config, "", "    ")
    if err != nil { return err }

    // Gotta be a better way to do this...
    data = append(data, []byte("\n")...)

    sent_total  := 0
    data_length := len(data)
    for sent_total < data_length {
        sent_now, err := w.Write(data[sent_total:])
        if err != nil { return err }

        sent_total += sent_now
    }

    return nil
}
