//go:build windows

package main

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

func getSteamPath() string {
	reg, err := registry.OpenKey(registry.CURRENT_USER, `Software\Valve\Steam`, registry.QUERY_VALUE)
	if err != nil {
		fmt.Printf("打开注册表失败！\n")
		return "./"
	}
	defer reg.Close()
	path, _, err := reg.GetStringValue("SteamPath")
	if err != nil {
		fmt.Printf("读取注册表失败！\n")
		return "./"
	}
	return path + "/"
}
