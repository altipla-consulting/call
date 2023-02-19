package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/TylerBrock/colorjson"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

func sendValidArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if _, err := os.Stat("protos"); os.IsNotExist(err) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var complete []string
	fn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Trace(err)
		}
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".proto" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return errors.Trace(err)
		}
		var pkg, svc string
		var methods []string
		for _, line := range strings.Split(string(content), "\n") {
			line = strings.TrimSpace(line)

			if strings.HasPrefix(line, "package ") {
				pkg = line
				pkg = strings.TrimPrefix(pkg, "package ")
				pkg = strings.TrimSuffix(pkg, ";")
				pkg = strings.TrimSpace(pkg)
			}
			if strings.HasPrefix(line, "service ") {
				svc = line
				svc = strings.TrimPrefix(svc, "service ")
				svc = strings.TrimSuffix(svc, "{")
				svc = strings.TrimSpace(svc)
			}
			if strings.HasPrefix(line, "rpc ") {
				methodName := line
				methodName = strings.TrimPrefix(methodName, "rpc ")
				methodName = strings.Split(methodName, "(")[0]
				methodName = strings.TrimSpace(methodName)
				methods = append(methods, methodName)
			}
		}

		for _, method := range methods {
			sug := fmt.Sprintf("%s.%s/%s", pkg, svc, method)
			if strings.HasPrefix(sug, toComplete) {
				complete = append(complete, sug)
			}
		}

		return nil
	}
	if err := filepath.Walk("protos", fn); err != nil {
		cobra.CompErrorln(err.Error())
		return nil, cobra.ShellCompDirectiveError
	}

	return complete, cobra.ShellCompDirectiveNoFileComp
}

func sendRequest(ctx context.Context, remote string, args []string) error {
	parts := strings.Split(args[0], "/")
	if len(parts) != 2 {
		return errors.Errorf("Invalid method name %q. It should be like: package.subpackage.FooService/FooMethod", args[0])
	}

	data := map[string]any{}
	for _, arg := range args[1:] {
		if !strings.Contains(arg, "=") {
			continue
		}

		parts := strings.SplitN(arg, "=", 2)
		key := strings.TrimSpace(parts[0])

		if key == "" {
			return errors.Errorf("Invalid argument %q", arg)
		}

		var value any = strings.TrimSpace(parts[1])
		if strings.HasSuffix(key, ":") {
			key = strings.TrimSpace(strings.TrimSuffix(key, ":"))

			switch value {
			case "true":
				value = true
			case "false":
				value = false
			case "null":
				value = nil
			default:
				var err error
				value, err = strconv.ParseInt(value.(string), 10, 64)
				if err != nil {
					return errors.Errorf("Invalid argument %q", arg)
				}
			}
		}

		if err := setPath(data, key, value); err != nil {
			return errors.Trace(err)
		}
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return errors.Trace(err)
	}

	target := fmt.Sprintf("https://%s/%s/%s", remote, parts[0], parts[1])
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, target, &buf)
	if err != nil {
		return errors.Trace(err)
	}
	req.Header.Set("Content-Type", "application/json")

	for _, arg := range args {
		if !strings.Contains(arg, ":") {
			continue
		}

		parts := strings.SplitN(arg, ":", 2)
		req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}

	if err := printRequest(req, buf); err != nil {
		return errors.Trace(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Trace(err)
	}
	defer resp.Body.Close()

	if err := printResponse(resp); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func printRequest(req *http.Request, body bytes.Buffer) error {
	green := color.New(color.FgGreen)
	cyan := color.New(color.FgCyan)
	blue := color.New(color.FgBlue)

	green.Print(req.Method)
	fmt.Print(" ")
	cyan.Print(req.URL.Path)
	fmt.Print(" ")
	blue.Print("HTTP/1.1")
	fmt.Println()

	blue.Print("Host")
	fmt.Println(": ", req.URL.Host)

	for key, values := range req.Header {
		for _, value := range values {
			blue.Print(key)
			fmt.Println(": ", value)
		}
	}

	fmt.Println()

	if err := printJSON(body.Bytes()); err != nil {
		return errors.Trace(err)
	}
	fmt.Println()
	fmt.Println()

	return nil
}

func printResponse(resp *http.Response) error {
	cyan := color.New(color.FgCyan)
	blue := color.New(color.FgBlue)

	blue.Print("HTTP/1.1")
	fmt.Print(" ")
	blue.Print(resp.StatusCode)
	fmt.Print(" ")
	cyan.Print(http.StatusText(resp.StatusCode))
	fmt.Println()

	for key, values := range resp.Header {
		for _, value := range values {
			blue.Print(key)
			fmt.Println(": ", value)
		}
	}

	fmt.Println()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.Trace(err)
	}
	if resp.Header.Get("Content-Type") == "application/json" {
		if err := printJSON(body); err != nil {
			return errors.Trace(err)
		}
	} else {
		fmt.Println(string(body))
	}

	return nil
}

func printJSON(b []byte) error {
	var obj map[string]any
	if err := json.Unmarshal(b, &obj); err != nil {
		return errors.Trace(err)
	}

	f := colorjson.NewFormatter()
	f.Indent = 2
	pretty, err := f.Marshal(obj)
	if err != nil {
		return errors.Trace(err)
	}
	fmt.Println(string(pretty))

	return nil
}

func setPath(object map[string]any, path string, value any) error {
	it := object
	parts := strings.Split(path, ".")
	for _, p := range parts[:len(parts)-1] {
		switch v := it[p].(type) {
		case map[string]any:
			it = v
		case nil:
			m := map[string]any{}
			it[p] = m
			it = m
		default:
			return errors.Errorf("Invalid path %q already has value %#v", path, v)
		}
	}

	if it == nil {
		return errors.Errorf("Invalid path %q", path)
	}

	it[parts[len(parts)-1]] = value

	return nil
}
