package main

import "testing"

func assert(success bool, t *testing.T, msg string) {
    if success == false {
        t.Errorf(msg)
    }
}

func TestRead(t *testing.T) {
    path := "example.conf.json"
    config, err := GetConfig(path)
    if err != nil {
        t.Errorf("Got error %v", err)
    }
    if len(config.Meta) != 2 {
        t.Errorf("Expected 2 items of metadata")
    }

    // Compare services ====================================
    if len(config.Services) != 1 {
        t.Errorf("Expected 1 services")
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
    }

    // Compare schedules ===================================
    if len(config.Schedules) != 1 {
        t.Errorf("Expected 1 schedule")
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
    }

    // Compare templates ===================================
    if len(config.Templates) != 2 {
        t.Errorf("Expected 2 templates")
    }
    template := config.Templates[0]
    expected_template := GLConfigTemplate{
        "overview",
        "/srv/goldilocks/templates/overview",
        "/srv/www/gl/index.html",
    }
    if template != expected_template {
        t.Errorf("Incorrect template data for %s", expected_template.Name)
    }
    template = config.Templates[1]
    expected_template = GLConfigTemplate{
        "global json dump",
        "core.json",
        "/srv/www/gl/core.json",
    }
    if template != expected_template {
        t.Errorf("Incorrect template data for %s", expected_template.Name)
    }
}

func TestValidateFromRead(t *testing.T) {
    path := "example.conf.json"
    config, err := GetConfig(path)
    if err != nil {
        t.Errorf("Got error %v", err)
    }

    ok := ValidateConfig(&config)
    if ! ok {
        t.Errorf("Validation failed")
    }
}

func TestValidateNoDefaultRPC(t *testing.T) {
    path := "example.conf.json"
    config, err := GetConfig(path)
    if err != nil {
        t.Errorf("Got error %v", err)
    }

    delete(config.RPC, "default")
    ok := ValidateConfig(&config)
    if ok {
        t.Errorf("Validation should have failed")
    }
}
