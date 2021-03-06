package aws

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestNewSSEParsedConfig(t *testing.T) {
	kmsKeyID := "test"
	kmsEncryptionContext := `{"a": "bc", "b": "cd"}`
	// compact form of kmsEncryptionContext
	parsedKMSEncryptionContext := "eyJhIjoiYmMiLCJiIjoiY2QifQ=="

	tests := []struct {
		name        string
		params      SSEConfig
		expected    *SSEParsedConfig
		expectedErr error
	}{
		{
			name: "Test SSE encryption with SSES3 type",
			params: SSEConfig{
				Type: SSES3,
			},
			expected: &SSEParsedConfig{
				ServerSideEncryption: sseS3Type,
			},
		},
		{
			name: "Test SSE encryption with SSEKMS type without context",
			params: SSEConfig{
				Type:     SSEKMS,
				KMSKeyID: kmsKeyID,
			},
			expected: &SSEParsedConfig{
				ServerSideEncryption: sseKMSType,
				KMSKeyID:             &kmsKeyID,
			},
		},
		{
			name: "Test SSE encryption with SSEKMS type with context",
			params: SSEConfig{
				Type:                 SSEKMS,
				KMSKeyID:             kmsKeyID,
				KMSEncryptionContext: kmsEncryptionContext,
			},
			expected: &SSEParsedConfig{
				ServerSideEncryption: sseKMSType,
				KMSKeyID:             &kmsKeyID,
				KMSEncryptionContext: &parsedKMSEncryptionContext,
			},
		},
		{
			name: "Test invalid SSE type",
			params: SSEConfig{
				Type: "invalid",
			},
			expectedErr: errors.New("SSE type is empty or invalid"),
		},
		{
			name: "Test SSE encryption with SSEKMS type without KMS Key ID",
			params: SSEConfig{
				Type:     SSEKMS,
				KMSKeyID: "",
			},
			expectedErr: errors.New("KMS key id must be passed when SSE-KMS encryption is selected"),
		},
		{
			name: "Test SSE with invalid KMS encryption context JSON",
			params: SSEConfig{
				Type:                 SSEKMS,
				KMSKeyID:             kmsKeyID,
				KMSEncryptionContext: `INVALID_JSON`,
			},
			expectedErr: errors.New("failed to parse KMS encryption context: failed to marshal KMS encryption context: json: error calling MarshalJSON for type json.RawMessage: invalid character 'I' looking for beginning of value"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewSSEParsedConfig(tt.params)
			if tt.expectedErr != nil {
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			}
			assert.Equal(t, tt.expected, result)
		})
	}
}
