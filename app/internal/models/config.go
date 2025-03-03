package models

import (
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"slices"
	"strings"
)

var detectorMethods = []string{"exclude", "include"}

type section interface {
	validate() error
}

// ================
// ==== CONFIG ====
// ================

type Config struct {
	Present bool

	Detectors DetectorsConf `yaml:"detectors"`
}

func (c *Config) Read(path string) error {
	file, err := os.ReadFile(path)

	if err != nil {
		return nil
	}

	err = yaml.Unmarshal(file, &c)

	if err != nil {
		return err
	}

	err = c.validate()

	if err != nil {
		return err
	}

	c.Present = true

	return nil
}

func (c *Config) validate() error {
	if !c.Present {
		return nil
	}

	if err := c.Detectors.validate(); err != nil {
		return err
	}

	log.Print(c.Detectors.Names)

	return nil
}

// ===================
// ==== DETECTORS ====
// ===================

type DetectorsConf struct {
	section
	Method string   `yaml:"method"`
	Names  []string `yaml:"names"`
}

func (d *DetectorsConf) validate() error {
	if d.Method != "" || !slices.Contains(detectorMethods, d.Method) {
		return errors.New("invalid method, must be one of: " + strings.Join(detectorMethods, ", "))
	}

	if d.Names == nil {
		return errors.New("no detectors provided")
	}

	return nil
}
