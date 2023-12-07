package configurationreader

import (
	"context"
	"os"
	"reflect"
	"strconv"

	awsssm "github.com/oslokommune/common-lib-go/aws/awsparameterstore"
	"github.com/rs/zerolog/log"
)

// Read a ParameterStore variable and/or environment variables into a config struct.
//
// The struct should be annotated with a `json` tags and optionally `env` tags.
// to indicate which JSON properties and environment variable to read from.
// Struct fields should be of type string, bool or int.
//
// If the environment variable is set, it overrides the value from ParameterStore.
func ReadConfiguration[T any](ctx context.Context, client awsssm.GetParameterAPI, name string) (*T, error) {
	var cfg T

	// Read and fill with parameterstore values
	if err := awsssm.GetParameterStoreParameter(ctx, client, name, &cfg); err != nil {
		log.Info().Err(err).Msgf("Failed to read parameterstore configuration for parameter %s", name)
	}

	// Override with env variabler and panic if no value is set in either way
	cfgType := reflect.TypeOf(cfg)
	cfgValue := reflect.ValueOf(&cfg).Elem()

	// Iterate over the fields of the struct
	for i := 0; i < cfgType.NumField(); i++ {
		// Get the field type and value
		fieldType := cfgType.Field(i)
		fieldValue := cfgValue.Field(i)

		// Determine what env var to read from
		tag := fieldType.Tag.Get("env")
		if tag == "" {
			tag = fieldType.Tag.Get("json")
		}

		envValue, found := os.LookupEnv(tag)
		if !found && (fieldValue.IsZero() && fieldValue.Kind() != reflect.Bool) {
			log.Panic().Msgf("Either parameterstore value for parameter %s or environment variable is not set.", tag)
		}

		// Set the field value using reflection
		if found && fieldValue.CanSet() {
			switch fieldValue.Kind() {
			case reflect.String:
				fieldValue.SetString(envValue)
			case reflect.Bool:
				{
					i, err := strconv.ParseBool(envValue)
					if err != nil {
						log.Panic().Err(err).Msg("Failed to convert string to bool")
					}
					fieldValue.SetBool(i)
				}
			case reflect.Int:
				{
					i, err := strconv.Atoi(envValue)
					if err != nil {
						log.Panic().Err(err).Msg("Failed to convert string to int")
					}
					fieldValue.SetInt(int64(i))
				}
			}
		}
	}

	return &cfg, nil
}
