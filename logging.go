package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/gookit/color"
)

func LogErr(m string) {
	color.Printf("<fg=white>[</><fg=red;op=bold>ERROR</><fg=white>]</> » %s\n", m)
}

//d
func LogInvalid(m string) {
	color.Printf("<fg=white>[</><fg=red;op=bold>INVALID</><fg=white>]</> » %s\n", m)
}

//d
func LogInfo(m string) {
	color.Printf("<fg=white>[</><fg=cyan;op=bold>INFO</><fg=white>]</> » %s\n", m)
}

//d
func LogSuccess(m string) {
	color.Printf("<fg=white>[</><fg=green;op=bold>SUCCESS</><fg=white>]</> » %s\n", m)
}

//d
func LogWarn(m string) {
	color.Printf("<fg=white>[</><fg=yellow;op=bold>WARN</><fg=white>]</> » %s\n", m)
}

//d
func LogFatal(m string) {
	color.Printf("<fg=white>[</><fg=red;op=bold>FATAL</><fg=white>]</> » %s\n", m)
	color.Printf("<fg=white>[</><fg=red;op=bold>LOG</><fg=white>]</> » exitting\n")
	os.Exit(0)
}

func Clear() {
	cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}

//
func UserInput(m string) string {
	reader := bufio.NewReader(os.Stdin)
	var out string
	color.Printf("<fg=white>[</><fg=cyan;op=bold>input</><fg=white>]</> %s » ", m)
	out, err := reader.ReadString('\n')
	if err != nil {
		LogFatal(err.Error())
	}
	out = strings.TrimSuffix(out, "\r\n")
	out = strings.TrimSuffix(out, "\n")
	return out
}

//
func PrintLogo() {
	header := `
██████╗ ██╗███████╗ ██████╗ ██████╗ ██████╗ ██████╗       ███████╗███╗   ██╗██╗██████╗ ███████╗██████╗        ██████╗  ██████╗ 
██╔══██╗██║██╔════╝██╔════╝██╔═══██╗██╔══██╗██╔══██╗      ██╔════╝████╗  ██║██║██╔══██╗██╔════╝██╔══██╗      ██╔════╝ ██╔═══██╗
██║  ██║██║███████╗██║     ██║   ██║██████╔╝██║  ██║█████╗███████╗██╔██╗ ██║██║██████╔╝█████╗  ██████╔╝█████╗██║  ███╗██║   ██║
██║  ██║██║╚════██║██║     ██║   ██║██╔══██╗██║  ██║╚════╝╚════██║██║╚██╗██║██║██╔═══╝ ██╔══╝  ██╔══██╗╚════╝██║   ██║██║   ██║
██████╔╝██║███████║╚██████╗╚██████╔╝██║  ██║██████╔╝      ███████║██║ ╚████║██║██║     ███████╗██║  ██║      ╚██████╔╝╚██████╔╝
╚═════╝ ╚═╝╚══════╝ ╚═════╝ ╚═════╝ ╚═╝  ╚═╝╚═════╝       ╚══════╝╚═╝  ╚═══╝╚═╝╚═╝     ╚══════╝╚═╝  ╚═╝       ╚═════╝  ╚═════╝ 																						
<fg=cyan>[version 1.0]</> 
	`

	for _, char := range []string{"╗", "║", "╝", "╔", "═"} {
		header = strings.ReplaceAll(header, char, fmt.Sprintf("<fg=white>%v</>", char))
	}

	header = strings.ReplaceAll(header, "█", "<fg=blue>█</>")

	color.Printf(header)
	fmt.Println("")

}
