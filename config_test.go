package main

import "testing"

const PATH_EXAMPLE  = "resources/conf/example.json"
const PATH_TEMPLATE = "resources/conf/tmpl.json"

func assert(success bool, t *testing.T, msg string) {
    if success == false {
        t.Errorf(msg)
        return
    }
}

func TestRead(t *testing.T) {
    config, err := GetConfig([]string{PATH_EXAMPLE})
    if err != nil {
        t.Errorf("Got error %v", err)
        return
    }
    if len(config.Meta) != 2 {
        t.Errorf("Expected 2 items of metadata")
        return
    }

    // Compare services ====================================
    if len(config.Services) != 1 {
        t.Errorf("Expected 1 services")
        return
    }
    service := config.Services["nginx"]
    expected_service := GLConfigService{
        "Turn nginx on and off",
        "<some long bitcoin addr>",
        "0 BTC",
        GLConfigCommands{
            "sudo /bin/goldilocks_start",
            "sudo /bin/goldilocks_stop",
            "pgrep nginx",
        },
        "",
    }
    if service != expected_service {
        t.Errorf("Incorrect service data")
        return
    }

    // Compare schedules ===================================
    if len(config.Schedules) != 1 {
        t.Errorf("Expected 1 schedule")
        return
    }
    schedule := config.Schedules["daily_retrieval"]
    expected_schedule := GLConfigSchedule{
        "<same bt addr as earlier>",
        "<personal bt addr>",
        "0.002 BTC",
        "0 5 * * *",
        "",
    }
    if schedule != expected_schedule {
        t.Errorf("Incorrect schedule data")
        return
    }

    // Compare templates ===================================
    if len(config.Templates) != 2 {
        t.Errorf("Expected 2 templates")
        return
    }
    expected_templates := map[string]GLConfigTemplate{
        "overview": GLConfigTemplate{
            "/srv/goldilocks/templates/overview",
            "/srv/www/gl/index.html",
        },
        "global_json_dump": GLConfigTemplate{
            "core.json",
            "/srv/www/gl/core.json",
        },
    }
    for name, expected := range expected_templates {
        template := config.Templates[name]
        if template != expected {
            t.Errorf("Incorrect template data for '%s'", name)
            return
        }
    }
}

func TestValidateFromRead(t *testing.T) {
    config, err := GetConfig([]string{PATH_EXAMPLE})
    if err != nil {
        t.Errorf("Got error %v", err)
        return
    }

    err = ValidateConfig(&config)
    if err != nil {
        t.Errorf("Validation failed")
        return
    }
}

func TestValidateNoDefaultRPC(t *testing.T) {
    config, err := GetConfig([]string{PATH_EXAMPLE})
    if err != nil {
        t.Errorf("Got error %v", err)
        return
    }

    delete(config.RPC, "default")
    err = ValidateConfig(&config)
    if err == nil {
        t.Errorf("Validation should have failed")
        return
    }
    expected_error := "Failure to validate GLConfig.RPC: No default RPC address"
    if err.Error() != expected_error {
        t.Errorf("Expected '%s', got '%s'", expected_error, err.Error())
    }
}

func TestValidateNoStopCommand(t *testing.T) {
    config, err := GetConfig([]string{PATH_EXAMPLE})
    if err != nil {
        t.Errorf("Got error %v", err)
    }

    service_name := "nginx"
    service, ok  := config.Services[service_name]
    if ! ok {
        t.Errorf("No service '%s'", service_name)
    }
    service.Commands.Stop = ""
    config.Services[service_name] = service

    err = ValidateConfig(&config)
    if err == nil {
        t.Errorf("Validation should have failed")
        return
    }
    expected_error := "Failure to validate GLConfigCommands.Stop: was blank string"
    if err.Error() != expected_error {
        t.Errorf("Expected '%s', got '%s'", expected_error, err.Error())
        return
    }
}
