package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const (
	noArgsError = 1
	errorCode   = 2
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var exitError *exec.ExitError
	if len(cmd) == 0 {
		return noArgsError
	}
	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec
	currentEnv := os.Environ()
	commandEnvMap, err := parseCommandMap(currentEnv)
	if err != nil {
		return errorCode
	}
	filterEnvironment(env, commandEnvMap)

	command.Env = buildCommand(commandEnvMap)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	filterEnvironment(env, commandEnvMap)

	if err := command.Start(); err != nil {
		log.Println(err)
	}
	if err := command.Wait(); err != nil {
		if errors.As(err, &exitError) {
			returnCode = exitError.ExitCode()
		} else {
			log.Println(err)
		}
	}

	return returnCode
}

// buildCommand make command strings name=value.
func buildCommand(commandEnvMap map[string]string) []string {
	commandStrings := make([]string, 0, len(commandEnvMap))
	bufferString := strings.Builder{}

	for name, value := range commandEnvMap {
		bufferString.WriteString(name)
		bufferString.WriteString("=")
		bufferString.WriteString(value)
		commandStrings = append(commandStrings, bufferString.String())
		bufferString.Reset()
	}
	sort.Strings(commandStrings)

	return commandStrings
}

// parseCommandMap parse envs string.
func parseCommandMap(envs []string) (map[string]string, error) {
	environmentMap := make(map[string]string, len(envs))
	for _, item := range envs {
		value := strings.SplitN(item, "=", 2)
		if len(value) != 2 {
			return nil, errors.New("error parse")
		}
		environmentMap[value[0]] = value[1]
	}

	return environmentMap, nil
}

// filterEnvironment filter value from commandEnvMap.
func filterEnvironment(env Environment, commandEnvMap map[string]string) {
	for name, value := range env {
		switch value.NeedRemove {
		case true:
			delete(commandEnvMap, name)
		case false:
			commandEnvMap[name] = value.Value
		}
	}
}
