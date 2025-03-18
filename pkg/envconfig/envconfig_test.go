package envconfig

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestConfig struct {
	Simple      string `envconfig:"SIMPLE" required:"true"`
	WithDefault string `envconfig:"WITH_DEFAULT" default:"default_value"`

	Nested struct {
		Required string `envconfig:"REQUIRED" required:"true"`
		Optional string `envconfig:"OPTIONAL"`
	} `envconfig:"NESTED"`

	PointerNested *struct {
		Required string `envconfig:"REQUIRED" required:"true"`
	} `envconfig:"PTR"`
}

func TestProcess(t *testing.T) {
	tests := []struct {
		name          string
		envVars       map[string]string
		wantErr       bool
		expectedError string
	}{
		{
			name: "all required variables set",
			envVars: map[string]string{
				"SIMPLE":          "value",
				"NESTED_REQUIRED": "nested_value",
				"PTR_REQUIRED":    "ptr_value",
			},
			wantErr: false,
		},
		{
			name: "missing simple required variable",
			envVars: map[string]string{
				"NESTED_REQUIRED": "nested_value",
				"PTR_REQUIRED":    "ptr_value",
			},
			wantErr:       true,
			expectedError: "configuration error:\n\tmissing required environment variable \"SIMPLE\"",
		},
		{
			name: "missing nested required variable",
			envVars: map[string]string{
				"SIMPLE":       "value",
				"PTR_REQUIRED": "ptr_value",
			},
			wantErr:       true,
			expectedError: "configuration error:\n\tmissing required environment variable \"NESTED_REQUIRED\"",
		},
		{
			name: "missing pointer nested required variable",
			envVars: map[string]string{
				"SIMPLE":          "value",
				"NESTED_REQUIRED": "nested_value",
			},
			wantErr:       true,
			expectedError: "configuration error:\n\tmissing required environment variable \"PTR_REQUIRED\"",
		},
		{
			name: "missing multiple required variables",
			envVars: map[string]string{
				"WITH_DEFAULT": "custom_default", // This is not required due to default value
			},
			wantErr:       true,
			expectedError: "configuration error:\n\tmissing required environment variable \"SIMPLE\"\n\tmissing required environment variable \"NESTED_REQUIRED\"\n\tmissing required environment variable \"PTR_REQUIRED\"",
		},
		{
			name: "default value is used",
			envVars: map[string]string{
				"SIMPLE":          "value",
				"NESTED_REQUIRED": "nested_value",
				"PTR_REQUIRED":    "ptr_value",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment before each test
			os.Clearenv()

			// Set environment variables for the test
			for k, v := range tt.envVars {
				require.NoError(t, os.Setenv(k, v))
			}

			var cfg TestConfig
			err := Process("", &cfg)

			if tt.wantErr {
				require.Error(t, err)
				if tt.expectedError != "" {
					assert.Equal(t, tt.expectedError, err.Error())
				}
			} else {
				require.NoError(t, err)
				// Verify values are set correctly
				if val, exists := tt.envVars["SIMPLE"]; exists {
					assert.Equal(t, val, cfg.Simple)
				}
				if val, exists := tt.envVars["NESTED_REQUIRED"]; exists {
					assert.Equal(t, val, cfg.Nested.Required)
				}
				if val, exists := tt.envVars["PTR_REQUIRED"]; exists {
					require.NotNil(t, cfg.PointerNested)
					assert.Equal(t, val, cfg.PointerNested.Required)
				}
				// Check default value
				assert.Equal(t, "default_value", cfg.WithDefault)
			}
		})
	}
}

func TestMissingEnvError(t *testing.T) {
	err := &MissingEnvError{
		StructField: "Host",
		EnvKey:      "DB_HOST",
	}

	assert.Equal(t, `missing required environment variable "DB_HOST"`, err.Error())
}

func TestProcessEdgeCases(t *testing.T) {
	t.Run("nil pointer", func(t *testing.T) {
		err := Process("", nil)
		assert.Error(t, err)
	})

	t.Run("non-pointer value", func(t *testing.T) {
		var cfg TestConfig
		err := Process("", cfg)
		assert.Error(t, err)
	})

	t.Run("pointer to non-struct", func(t *testing.T) {
		var str string
		err := Process("", &str)
		assert.Error(t, err)
	})
}
