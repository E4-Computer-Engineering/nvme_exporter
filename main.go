package main

// Export nvme smart-log metrics in prometheus format

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"os/user"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tidwall/gjson"
)

var labels = []string{"device"}

type nvmeCollector struct {
	nvmeCriticalWarning *prometheus.Desc
	nvmeTemperature *prometheus.Desc
	nvmeAvailSpare *prometheus.Desc
	nvmeSpareThresh *prometheus.Desc
	nvmePercentUsed *prometheus.Desc
	nvmeEnduranceGrpCriticalWarningSummary *prometheus.Desc
	nvmeDataUnitsRead *prometheus.Desc
	nvmeDataUnitsWritten *prometheus.Desc
	nvmeHostReadCommands *prometheus.Desc
	nvmeHostWriteCommands *prometheus.Desc
	nvmeControllerBusyTime *prometheus.Desc
	nvmePowerCycles *prometheus.Desc
	nvmePowerOnHours *prometheus.Desc
	nvmeUnsafeShutdowns *prometheus.Desc
	nvmeMediaErrors *prometheus.Desc
	nvmeNumErrLogEntries *prometheus.Desc
	nvmeWarningTempTime *prometheus.Desc
	nvmeCriticalCompTime *prometheus.Desc
	nvmeThmTemp1TransCount *prometheus.Desc
	nvmeThmTemp2TransCount *prometheus.Desc
	nvmeThmTemp1TotalTime *prometheus.Desc
	nvmeThmTemp2TotalTime *prometheus.Desc
    nvmePhysicalMediaUnitsWrittenHi *prometheus.Desc
    nvmePhysicalMediaUnitsWrittenLo *prometheus.Desc
    nvmePhysicalMediaUnitsReadHi *prometheus.Desc
    nvmePhysicalMediaUnitsReadLo *prometheus.Desc
}

// nvme smart-log field descriptions can be found on page 180 of:
// https://nvmexpress.org/wp-content/uploads/NVM-Express-Base-Specification-2_0-2021.06.02-Ratified-5.pdf

func newNvmeCollector() prometheus.Collector {
	return &nvmeCollector{
		nvmeCriticalWarning: prometheus.NewDesc(
			"nvme_critical_warning",
			"Critical warnings for the state of the controller",
			labels,
			nil,
		),
		nvmeTemperature: prometheus.NewDesc(
			"nvme_temperature",
			"Temperature in degrees fahrenheit",
			labels,
			nil,
		),
		nvmeAvailSpare: prometheus.NewDesc(
			"nvme_avail_spare",
			"Normalized percentage of remaining spare capacity available",
			labels,
			nil,
		),
		nvmeSpareThresh: prometheus.NewDesc(
			"nvme_spare_thresh",
			"Async event completion may occur when avail spare < threshold",
			labels,
			nil,
		),
		nvmePercentUsed: prometheus.NewDesc(
			"nvme_percent_used",
			"Vendor specific estimate of the percentage of life used",
			labels,
			nil,
		),
		nvmeEnduranceGrpCriticalWarningSummary: prometheus.NewDesc(
			"nvme_endurance_grp_critical_warning_summary",
			"Critical warnings for the state of endurance groups",
			labels,
			nil,
		),
		nvmeDataUnitsRead: prometheus.NewDesc(
			"nvme_data_units_read",
			"Number of 512 byte data units host has read",
			labels,
			nil,
		),
		nvmeDataUnitsWritten: prometheus.NewDesc(
			"nvme_data_units_written",
			"Number of 512 byte data units the host has written",
			labels,
			nil,
		),
		nvmeHostReadCommands: prometheus.NewDesc(
			"nvme_host_read_commands",
			"Number of read commands completed",
			labels,
			nil,
		),
		nvmeHostWriteCommands: prometheus.NewDesc(
			"nvme_host_write_commands",
			"Number of write commands completed",
			labels,
			nil,
		),
		nvmeControllerBusyTime: prometheus.NewDesc(
			"nvme_controller_busy_time",
			"Amount of time in minutes controller busy with IO commands",
			labels,
			nil,
		),
		nvmePowerCycles: prometheus.NewDesc(
			"nvme_power_cycles",
			"Number of power cycles",
			labels,
			nil,
		),
		nvmePowerOnHours: prometheus.NewDesc(
			"nvme_power_on_hours",
			"Number of power on hours",
			labels,
			nil,
		),
		nvmeUnsafeShutdowns: prometheus.NewDesc(
			"nvme_unsafe_shutdowns",
			"Number of unsafe shutdowns",
			labels,
			nil,
		),
		nvmeMediaErrors: prometheus.NewDesc(
			"nvme_media_errors",
			"Number of unrecovered data integrity errors",
			labels,
			nil,
		),
		nvmeNumErrLogEntries: prometheus.NewDesc(
			"nvme_num_err_log_entries",
			"Lifetime number of error log entries",
			labels,
			nil,
		),
		nvmeWarningTempTime: prometheus.NewDesc(
			"nvme_warning_temp_time",
			"Amount of time in minutes temperature > warning threshold",
			labels,
			nil,
		),
		nvmeCriticalCompTime: prometheus.NewDesc(
			"nvme_critical_comp_time",
			"Amount of time in minutes temperature > critical threshold",
			labels,
			nil,
		),
		nvmeThmTemp1TransCount: prometheus.NewDesc(
			"nvme_thm_temp1_trans_count",
			"Number of times controller transitioned to lower power",
			labels,
			nil,
		),
		nvmeThmTemp2TransCount: prometheus.NewDesc(
			"nvme_thm_temp2_trans_count",
			"Number of times controller transitioned to lower power",
			labels,
			nil,
		),
		nvmeThmTemp1TotalTime: prometheus.NewDesc(
			"nvme_thm_temp1_trans_time",
			"Total number of seconds controller transitioned to lower power",
			labels,
			nil,
		),
		nvmeThmTemp2TotalTime: prometheus.NewDesc(
			"nvme_thm_temp2_trans_time",
			"Total number of seconds controller transitioned to lower power",
			labels,
			nil,
		),
        nvmePhysicalMediaUnitsWrittenHi: prometheus.NewDesc(
			"nvme_physical_media_units_written_hi",
			"Physical meda units written high",
			labels,
			nil,
		),
        nvmePhysicalMediaUnitsWrittenLo: prometheus.NewDesc(
			"nvme_physical_media_units_written_lo",
			"Physical meda units written low",
			labels,
			nil,
		),
        nvmePhysicalMediaUnitsReadHi: prometheus.NewDesc(
			"nvme_physical_media_units_read_hi",
			"Physical meda units read high",
			labels,
			nil,
		),
        nvmePhysicalMediaUnitsReadLo: prometheus.NewDesc(
			"nvme_physical_media_units_read_lo",
			"Physical meda units read low",
			labels,
			nil,
		),
	}
}

func (c *nvmeCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.nvmeCriticalWarning
	ch <- c.nvmeTemperature
	ch <- c.nvmeAvailSpare
	ch <- c.nvmeSpareThresh
	ch <- c.nvmePercentUsed
	ch <- c.nvmeEnduranceGrpCriticalWarningSummary
	ch <- c.nvmeDataUnitsRead
	ch <- c.nvmeDataUnitsWritten
	ch <- c.nvmeHostReadCommands
	ch <- c.nvmeHostWriteCommands
	ch <- c.nvmeControllerBusyTime
	ch <- c.nvmePowerCycles
	ch <- c.nvmePowerOnHours
	ch <- c.nvmeUnsafeShutdowns
	ch <- c.nvmeMediaErrors
	ch <- c.nvmeNumErrLogEntries
	ch <- c.nvmeWarningTempTime
	ch <- c.nvmeCriticalCompTime
	ch <- c.nvmeThmTemp1TransCount
	ch <- c.nvmeThmTemp2TransCount
	ch <- c.nvmeThmTemp1TotalTime
	ch <- c.nvmeThmTemp2TotalTime
    ch <- c.nvmePhysicalMediaUnitsWrittenHi
    ch <- c.nvmePhysicalMediaUnitsWrittenLo
    ch <- c.nvmePhysicalMediaUnitsReadHi
    ch <- c.nvmePhysicalMediaUnitsReadLo
}

func executeCommand(cmd string, args ...string) ([]byte, error) {
    command := exec.Command(cmd, args...)
    output, err := command.CombinedOutput()
    if err != nil {
        log.Printf("error running %s command: %s, output: %s", cmd, err, string(output))
        return nil, fmt.Errorf("error running %s command: %s, output: %s", cmd, err, string(output))
    }
    if !gjson.Valid(string(output)) {
        log.Printf("invalid JSON output from %s command: %s", cmd, string(output))
        return nil, fmt.Errorf("invalid JSON output from %s command: %s", cmd, string(output))
    }
    return output, nil
}

func (c *nvmeCollector) Collect(ch chan<- prometheus.Metric) {
    nvmeDeviceList, err := c.getDeviceList()
    if err != nil {
        log.Fatalf("Error running nvme command: %s\n", err)
    }
    for _, nvmeDevice := range nvmeDeviceList {
        c.collectSmartLogMetrics(ch, nvmeDevice)
        c.collectOcpSmartLogMetrics(ch, nvmeDevice)
    }
}

func (c *nvmeCollector) getDeviceList() ([]gjson.Result, error) {
    nvmeDeviceCmd, err := executeCommand("nvme", "list", "-o", "json")
    if err != nil {
        return nil, err
    }
    return gjson.Get(string(nvmeDeviceCmd), "Devices.#.DevicePath").Array(), nil
}

func (c *nvmeCollector) collectSmartLogMetrics(ch chan<- prometheus.Metric, device gjson.Result) {
    nvmeSmartLog, err := executeCommand("nvme", "smart-log", device.String(), "-o", "json")
    if err != nil {
        log.Fatalf("Error running nvme smart-log command for device %s: %s\n", device.String(), err)
    }
    nvmeSmartLogMetrics := gjson.GetMany(string(nvmeSmartLog),
										 "critical_warning",
										 "temperature",
										 "avail_spare",
										 "spare_thresh",
										 "percent_used",
										 "endurance_grp_critical_warning_summary",
										 "data_units_read",
										 "data_units_written",
										 "host_read_commands",
										 "host_write_commands",
										 "controller_busy_time",
										 "power_cycles",
										 "power_on_hours",
										 "unsafe_shutdowns",
										 "media_errors",
										 "num_err_log_entries",
										 "warning_temp_time",
										 "critical_comp_time",
										 "thm_temp1_trans_count",
										 "thm_temp2_trans_count",
										 "thm_temp1_total_time",
										 "thm_temp2_total_time",)
    c.sendSmartLogMetrics(ch, nvmeSmartLogMetrics, device.String())
}

func (c *nvmeCollector) collectOcpSmartLogMetrics(ch chan<- prometheus.Metric, device gjson.Result) {
    nvmeOcpSmartLog, err := executeCommand("nvme", "ocp", "smart-add-log", device.String(), "-o", "json")
    if err != nil {
        log.Fatalf("Error running nvme ocp smart-add-log command for device %s: %s\n", device.String(), err)
    }
    nvmeOcpSmartLogMetrics := gjson.GetMany(string(nvmeOcpSmartLog),
											"Physical media units written.hi",
											"Physical media units written.lo",
											"Physical media units read.hi",
											"Physical media units read.lo",
											"Bad user nand blocks - Raw",
											"Bad user nand blocks - Normalized",
											"Bad system nand blocks - Raw",
											"Bad system nand blocks - Normalized",
											"XOR recovery count",
											"Uncorrectable read error count",
											"Soft ecc error count",
											"End to end detected errors",
											"End to end corrected errors",
											"System data percent used",
											"Refresh counts",
											"Max User data erase counts",
											"Min User data erase counts",
											"Number of Thermal throttling events",
											"Current throttling status",
											"PCIe correctable error count",
											"Incomplete shutdowns",
											"Percent free blocks",
											"Capacitor health",
											"Unaligned I/O",
											"Security Version Number",
											"NUSE - Namespace utilization",
											"PLP start count",
											"Endurance estimate",
											"Log page version",
											"Log page GUID",
											"Errata Version Field",
											"Point Version Field",
											"Minor Version Field",
											"Major Version Field",
											"NVMe Errata Version",
											"PCIe Link Retraining Count",
                                            "Power State Change Count",)
    c.sendOcpSmartLogMetrics(ch, nvmeOcpSmartLogMetrics, device.String())
}

func (c *nvmeCollector) sendSmartLogMetrics(ch chan<- prometheus.Metric, metrics []gjson.Result, device string) {
    ch <- prometheus.MustNewConstMetric(c.nvmeCriticalWarning, prometheus.GaugeValue, metrics[0].Float(), device)
	// convert kelvin to fahrenheit
	ch <- prometheus.MustNewConstMetric(c.nvmeTemperature, prometheus.GaugeValue, (metrics[1].Float() - 273.15) * 9/5 + 32, device)
	ch <- prometheus.MustNewConstMetric(c.nvmeAvailSpare, prometheus.GaugeValue, metrics[2].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeSpareThresh, prometheus.GaugeValue, metrics[3].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmePercentUsed, prometheus.GaugeValue, metrics[4].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeEnduranceGrpCriticalWarningSummary, prometheus.GaugeValue, metrics[5].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeDataUnitsRead, prometheus.CounterValue, metrics[6].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeDataUnitsWritten, prometheus.CounterValue, metrics[7].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeHostReadCommands, prometheus.CounterValue, metrics[8].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeHostWriteCommands, prometheus.CounterValue, metrics[9].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeControllerBusyTime, prometheus.CounterValue, metrics[10].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmePowerCycles, prometheus.CounterValue, metrics[11].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmePowerOnHours, prometheus.CounterValue, metrics[12].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeUnsafeShutdowns, prometheus.CounterValue, metrics[13].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeMediaErrors, prometheus.CounterValue, metrics[14].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeNumErrLogEntries, prometheus.CounterValue, metrics[15].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeWarningTempTime, prometheus.CounterValue, metrics[16].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeCriticalCompTime, prometheus.CounterValue, metrics[17].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeThmTemp1TransCount, prometheus.CounterValue, metrics[18].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeThmTemp2TransCount, prometheus.CounterValue, metrics[19].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeThmTemp1TotalTime, prometheus.CounterValue, metrics[20].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmeThmTemp2TotalTime, prometheus.CounterValue, metrics[21].Float(), device)
}

func (c *nvmeCollector) sendOcpSmartLogMetrics(ch chan<- prometheus.Metric, metrics []gjson.Result, device string) {
    ch <- prometheus.MustNewConstMetric(c.nvmePhysicalMediaUnitsWrittenHi, prometheus.CounterValue, metrics[0].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmePhysicalMediaUnitsWrittenLo, prometheus.CounterValue, metrics[1].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmePhysicalMediaUnitsReadHi, prometheus.CounterValue, metrics[2].Float(), device)
	ch <- prometheus.MustNewConstMetric(c.nvmePhysicalMediaUnitsReadLo, prometheus.CounterValue, metrics[3].Float(), device)
}

func main() {
	port := flag.String("port", "9998", "port to listen on")
	flag.Parse()
	// check user
	currentUser, err := user.Current()
	if err != nil {
		log.Fatalf("Error getting current user  %s\n", err)
	}
	if currentUser.Username != "root" {
		log.Fatalln("Error: you must be root to use nvme-cli")
	}
	// check for nvme-cli executable
	_, err = exec.LookPath("nvme")
	if err != nil {
		log.Fatalf("Cannot find nvme command in path: %s\n", err)
	}
	prometheus.MustRegister(newNvmeCollector())
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
