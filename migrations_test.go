package main

import (
    "math/rand"
    "os"
    "strconv"
    "testing"
)

func TestGenerateRevision(t *testing.T) {
    configPath := "settings/development.yml"
    err := loadConfig(configPath)
    if err != nil {
        t.Fatal(err)
    }

    generatedFiles, err := generateRevision("test" + strconv.Itoa(rand.Int()))
    if err != nil {
        t.Fatal(err)
    }

    for _, file := range generatedFiles {
        err := os.Remove(file)
        if err != nil {
            t.Fatal(err)
        }
    }
}
