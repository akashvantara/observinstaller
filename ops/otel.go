package ops

import (
	"fmt"
	"os"

	_ "embed"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed embeds/otel-config.yaml
	otelConfig []byte
)

func ListOtelOptions() {
	otelData := make(map[interface{}]interface{})
	unmarshalErr := yaml.Unmarshal(otelConfig, &otelData)

	if unmarshalErr != nil {
		fmt.Fprintf(os.Stderr, "Failed to unmarshal YAML, err: %v\n", unmarshalErr)
	}

	fmt.Fprintf(os.Stdin, "Supported entries:\n")
	for ke, ve := range otelData {
		fmt.Fprintf(os.Stdin, "  %s\n", ke)
		for kee, _ := range ve.(map[string]interface{}) {
			fmt.Fprintf(os.Stdin, "    %v\n", kee)
		}
	}
}
