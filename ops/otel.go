package ops

import (
	"fmt"
	"os"

	_ "embed"

	"github.com/hv/akash.chandra/observinstaller/conf"
	"gopkg.in/yaml.v3"
)

var (
	//go:embed embeds/otel-config.yaml
	otelConfig []byte
)

const (
	OTEL_RECEIVERS  = "receivers"
	OTEL_PROCESSORS = "processors"
	OTEL_EXPORTERS  = "exporters"
	OTEL_EXTENSIONS = "extensions"
	OTEL_SERVICE    = "service"
	OTEL_PIPELINES  = "pipelines"
)

func ListOtelOptions() {
	otelData := make(map[interface{}]interface{})
	unmarshalErr := yaml.Unmarshal(otelConfig, &otelData)

	if unmarshalErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal YAML, err: %v\n", unmarshalErr)
	}

	fmt.Fprintf(os.Stdin, "Supported entries:\n")
	for ke, ve := range otelData {
		if ke == OTEL_SERVICE {
			continue
		}
		fmt.Fprintf(os.Stdin, "  %s\n", ke)
		for kee := range ve.(map[string]interface{}) {
			fmt.Fprintf(os.Stdin, "    %v\n", kee)
		}
	}
}

func PrepareOtelCfgFile(fileConfig *conf.FileConfig, writeLoc string) bool {
	cfgMap := make(map[interface{}]interface{})           // Data to be prepared
	embeddedOtelData := make(map[interface{}]interface{}) // Data from the embedded file

	unmarshalErr := yaml.Unmarshal(otelConfig, &embeddedOtelData)
	if unmarshalErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal YAML, err: %v\n", unmarshalErr)
		return false
	}

	cfgMap[OTEL_SERVICE] = embeddedOtelData[OTEL_SERVICE] // Copy the service struct completely
	cfgMap[OTEL_EXTENSIONS] = make(map[string]interface{})
	cfgMap[OTEL_RECEIVERS] = make(map[string]interface{})
	cfgMap[OTEL_PROCESSORS] = make(map[string]interface{})
	cfgMap[OTEL_EXPORTERS] = make(map[string]interface{})

	receiversBlock := cfgMap[OTEL_RECEIVERS].(map[string]interface{})
	processorsBlock := cfgMap[OTEL_PROCESSORS].(map[string]interface{})
	exportersBlock := cfgMap[OTEL_EXPORTERS].(map[string]interface{})
	serviceBlock := cfgMap[OTEL_SERVICE].(map[string]interface{})
	extensionsBlock := cfgMap[OTEL_EXTENSIONS].(map[string]interface{})

	embeddedReceiversBlock := embeddedOtelData[OTEL_RECEIVERS].(map[string]interface{})
	embeddedProcessorsBlock := embeddedOtelData[OTEL_PROCESSORS].(map[string]interface{})
	embeddedExportersBlock := embeddedOtelData[OTEL_EXPORTERS].(map[string]interface{})

	serviceExtensionBlock := make([]string, 0, 4)
	embeddedExtensionBlock := embeddedOtelData[OTEL_EXTENSIONS].(map[string]interface{})

	for _, extension := range fileConfig.BaseOtelConfig.Extensions {
		if embeddedExtensionBlock[extension] != nil {
			serviceExtensionBlock = append(serviceExtensionBlock, extension)
			extensionsBlock[extension] = embeddedExtensionBlock[extension]
		} else {
			fmt.Fprintf(os.Stderr, "Failed to find %s extension in our config\n", extension)
		}
	}
	serviceBlock[OTEL_EXTENSIONS] = serviceExtensionBlock

	servicePipelineBlock := make(map[string]interface{})
	for pipelineName, pipelineIf := range fileConfig.BaseOtelConfig.Pipelines {
		servicePipelineBlock[pipelineName] = make(map[string]interface{})
		servicePipelineBlockEntry := servicePipelineBlock[pipelineName].(map[string]interface{})

		entryArr := make([]string, 0, 4)
		for _, receiverName := range pipelineIf[OTEL_RECEIVERS] {
			if embeddedReceiversBlock[receiverName] != nil {
				receiversBlock[receiverName] = embeddedReceiversBlock[receiverName]
				entryArr = append(entryArr, receiverName)
			} else {
				fmt.Fprintf(os.Stderr, "Failed to find %s receiver in our config\n", receiverName)
			}
		}
		servicePipelineBlockEntry[OTEL_RECEIVERS] = entryArr

		entryArr = make([]string, 0, 4)
		for _, processorName := range pipelineIf[OTEL_PROCESSORS] {
			if embeddedProcessorsBlock[processorName] != nil {
				processorsBlock[processorName] = embeddedProcessorsBlock[processorName]
				entryArr = append(entryArr, processorName)
			} else {
				fmt.Fprintf(os.Stderr, "Failed to find %s processor in our config\n", processorName)
			}
		}
		servicePipelineBlockEntry[OTEL_PROCESSORS] = entryArr

		entryArr = make([]string, 0, 4)
		for _, exporterName := range pipelineIf[OTEL_EXPORTERS] {
			if embeddedExportersBlock[exporterName] != nil {
				exportersBlock[exporterName] = embeddedExportersBlock[exporterName]
				entryArr = append(entryArr, exporterName)
			} else {
				fmt.Fprintf(os.Stderr, "Failed to find %s exporter in our config\n", exporterName)
			}
		}
		servicePipelineBlockEntry[OTEL_EXPORTERS] = entryArr
	}

	for _, pkg := range fileConfig.Pkg {
		if cfgMap[pkg.PkgOtelConfig.Type] != nil {
			servicePipelineBlockEntry := servicePipelineBlock[pkg.PkgOtelConfig.Pipeline].(map[string]interface{})
			servicePipelineBlockTypeEntry := servicePipelineBlockEntry[pkg.PkgOtelConfig.Type].([]string)
			var unmarshalData map[string]interface{} = make(map[string]interface{})
			if err := yaml.Unmarshal([]byte(pkg.PkgOtelConfig.Config), &unmarshalData); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to parse the otel config inside %s\n", pkg.Name)
			}

			genBlock := cfgMap[pkg.PkgOtelConfig.Type].(map[string]interface{})
			for k, v := range unmarshalData {
				genBlock[k] = v

				similarKFound := false
				for _, similarK := range servicePipelineBlockEntry[pkg.PkgOtelConfig.Type].([]string) {
					if similarK == k {
						similarKFound = true
						break
					}
				}

				if !similarKFound {
					servicePipelineBlockEntry[pkg.PkgOtelConfig.Type] = append(servicePipelineBlockTypeEntry, k)
				}
			}
		}
	}
	serviceBlock[OTEL_PIPELINES] = servicePipelineBlock

	configFile, err := os.OpenFile(writeLoc, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open %s file, err: %v\n", writeLoc, err)
		return false
	}
	defer configFile.Close()

	configData, err := yaml.Marshal(cfgMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to marshal config to prepare file, err: %v\n", err)
		configFile.Close()
		return false
	}

	_, err = configFile.Write(configData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to write to file, err: %v\n", err)
		configFile.Close()
		return false
	}

	fmt.Fprintf(os.Stdin, "Config written to file: %s\n", writeLoc)
	return true
}
