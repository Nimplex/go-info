package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"time"

	color "github.com/fatih/color"
	tablewriter "github.com/olekukonko/tablewriter"
	cpu "github.com/shirou/gopsutil/v3/cpu"
	hostd "github.com/shirou/gopsutil/v3/host"
	mem "github.com/shirou/gopsutil/v3/mem"
)

func Round2(x, unit float64) float64 {
	if x > 0 {
		return float64(int64(x/unit+0.5)) * unit
	}
	return float64(int64(x/unit-0.5)) * unit
}

func BtoG(bytes uint64) float64 {
	size := float64(bytes) * (9.3 * math.Pow(10, -10))
	return Round2(size, 0.001)
}

func main() {
	displayRaw := flag.Bool("raw", false, "display raw data?")
	flag.Parse()
	
	cores := [][]string{}
	cpus := [][]string{}
	memo := [][]string{}
	host := [][]string{}
	
	memInfo, _ := mem.VirtualMemory()
	cpuInfo, _ := cpu.Info()
	hostInfo, _ := hostd.Info()
	cpuPercentage, _ := cpu.Percent(time.Second*2, true)
	
	if *displayRaw == false {
		for index, elem := range cpuPercentage {
			cores = append(cores, []string{color.YellowString(fmt.Sprintf("%v", index)), fmt.Sprintf("%v%v", color.YellowString(fmt.Sprintf("%v", Round2(elem, 0.01))), "%")})
		}

		for index, elem := range cpuInfo {
			cpus = append(cpus,
				[]string{
					color.YellowString(fmt.Sprintf("%v", index)),
					(elem.ModelName),
					color.YellowString(fmt.Sprintf("%v", elem.Cores)),
					color.YellowString(fmt.Sprintf("%v", elem.Mhz)),
				},
			)
		}

		memo = append(memo, []string{
			fmt.Sprintf("%v GB (%v bytes)", color.YellowString(fmt.Sprintf("%v", BtoG(memInfo.Total))), color.YellowString(fmt.Sprintf("%v", memInfo.Total))),
			fmt.Sprintf("%v GB (%v bytes)", color.YellowString(fmt.Sprintf("%v", BtoG(memInfo.Free))), color.YellowString(fmt.Sprintf("%v", memInfo.Free))),
			fmt.Sprintf("%v GB (%v bytes)", color.YellowString(fmt.Sprintf("%v", BtoG(memInfo.SwapTotal))), color.YellowString(fmt.Sprintf("%v", memInfo.SwapTotal))),
			fmt.Sprintf("%v GB (%v bytes)", color.YellowString(fmt.Sprintf("%v", BtoG(memInfo.SwapFree))), color.YellowString(fmt.Sprintf("%v", memInfo.SwapFree))),
		})

		host = append(host, []string{
			(hostInfo.Hostname),
			color.YellowString(fmt.Sprintf("%v", hostInfo.Uptime)),
			color.YellowString(fmt.Sprintf("%v", hostInfo.BootTime)),
			(hostInfo.KernelArch),
			(hostInfo.Platform),
			color.YellowString(fmt.Sprintf("%v", hostInfo.Procs)),
			(hostInfo.HostID),
		})

		Cores := tablewriter.NewWriter(os.Stdout)
		Cpus := tablewriter.NewWriter(os.Stdout)
		Memo := tablewriter.NewWriter(os.Stdout)
		Host := tablewriter.NewWriter(os.Stdout)

		Host.SetHeader([]string{"Hostname", "Uptime", "Boot time", "Arch", "Platform", "Procs", "Host ID"})
		Host.SetAlignment(tablewriter.ALIGN_CENTER)
		Host.SetAutoWrapText(false)
		Host.AppendBulk(host)
		Host.Render()

		Memo.SetHeader([]string{"Total memory", "Free memory", "Total swap", "Free swap"})
		Memo.SetAlignment(tablewriter.ALIGN_CENTER)
		Memo.SetAutoWrapText(false)
		Memo.AppendBulk(memo)
		Memo.Render()

		Cpus.SetHeader([]string{"CPU number", "CPU model", "CPU cores", "CPU clock speed (MHz)"})
		Cpus.SetAlignment(tablewriter.ALIGN_CENTER)
		Cpus.SetAutoWrapText(false)
		Cpus.AppendBulk(cpus)
		Cpus.Render()

		Cores.SetHeader([]string{"Core number", "Core usage"})
		Cores.SetAlignment(tablewriter.ALIGN_CENTER)
		Cores.SetAutoWrapText(false)
		Cores.AppendBulk(cores)
		Cores.Render()
	} else {
		fmt.Println("Memory")
		fmt.Println(memInfo)
		fmt.Println("CPU")
		fmt.Println(cpuInfo)
		fmt.Println("CPU Usage")
		fmt.Println(cpuPercentage)
		fmt.Println("Host")
		fmt.Println(hostInfo)
	}
}