package controllers

import (
	"testing"

	"github.com/stretchr/testify/require"

	gpuv1 "github.com/NVIDIA/gpu-operator/api/nvidia/v1"
)

func TestGetDRADriverAdditionalNamespaces(t *testing.T) {
	testCases := []struct {
		description string
		draSpec     gpuv1.DRADriverSpec
		expected    []string
	}{
		{
			description: "additional namespaces variable not defined",
			draSpec: gpuv1.DRADriverSpec{
				ComputeDomains: gpuv1.DRADriverComputeDomains{
					Controller: gpuv1.DRADriverController{
						Env: nil,
					},
				},
			},
			expected: nil,
		},
		{
			description: "additional namespaces variable empty",
			draSpec: gpuv1.DRADriverSpec{
				ComputeDomains: gpuv1.DRADriverComputeDomains{
					Controller: gpuv1.DRADriverController{
						Env: []gpuv1.EnvVar{
							{Name: "ADDITIONAL_NAMESPACES", Value: ""},
						},
					},
				},
			},
			expected: nil,
		},
		{
			description: "additional namespaces variable defined",
			draSpec: gpuv1.DRADriverSpec{
				ComputeDomains: gpuv1.DRADriverComputeDomains{
					Controller: gpuv1.DRADriverController{
						Env: []gpuv1.EnvVar{
							{Name: "ADDITIONAL_NAMESPACES", Value: "foo,bar"},
						},
					},
				},
			},
			expected: []string{"foo", "bar"},
		},
		{
			description: "additional namespaces variable defined, white space trimmed",
			draSpec: gpuv1.DRADriverSpec{
				ComputeDomains: gpuv1.DRADriverComputeDomains{
					Controller: gpuv1.DRADriverController{
						Env: []gpuv1.EnvVar{
							{Name: "ADDITIONAL_NAMESPACES", Value: " foo,bar "},
						},
					},
				},
			},
			expected: []string{"foo", "bar"},
		},
	}

	for _, tc := range testCases {
		namespaces := getDRADriverAdditionalNamespaces(tc.draSpec)
		require.Equal(t, tc.expected, namespaces)
	}
}
