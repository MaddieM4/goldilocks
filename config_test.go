package main

import "testing"

func assert(success bool, t *testing.T, msg string) {
    if success == false {
        t.Errorf(msg)
        return
    }
}

func TestRead(t *testing.T) {
    path := "example.conf.json"
    config, err := GetConfig(path)
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
    service := config.Services[0]
    expected_service := GLConfigService{
        "nginx",
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
    schedule := config.Schedules[0]
    expected_schedule := GLConfigSchedule{
        "Daily retrieval",
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
    template := config.Templates[0]
    expected_template := GLConfigTemplate{
        "overview",
        "/srv/goldilocks/templates/overview",
        "/srv/www/gl/index.html",
    }
    if template != expected_template {
        t.Errorf("Incorrect template data for %s", expected_template.Name)
        return
    }
    template = config.Templates[1]
    expected_template = GLConfigTemplate{
        "global json dump",
        "core.json",
        "/srv/www/gl/core.json",
    }
    if template != expected_template {
        t.Errorf("Incorrect template data for %s", expected_template.Name)
        return
    }
}

func TestValidateFromRead(t *testing.T) {
    path := "example.conf.json"
    config, err := GetConfig(path)
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
    path := "example.conf.json"
    config, err := GetConfig(path)
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
    path := "example.conf.json"
    config, err := GetConfig(path)
    if err != nil {
        t.Errorf("Got error %v", err)
    }

    config.Services[0].Commands.Stop = ""
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
