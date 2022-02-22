package service

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	types "github.com/jack-hughes/ports/internal"
	"github.com/jack-hughes/ports/pkg/apis/ports"
	"github.com/jack-hughes/ports/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"testing"
)

var someErr = fmt.Errorf("some-error")

func TestService_Update(t *testing.T) {
	type table struct {
		name     string
		ctx      context.Context
		expCalls func(tt table, m *mocks.MockPorts_UpdateClient)
		req      *ports.Port
		err      error
	}
	tests := []table{
		{
			name: "successfully send a port update on the stream",
			req: &ports.Port{
				ID: "test-port-id",
			},
			expCalls: func(tt table, m *mocks.MockPorts_UpdateClient) {
				m.EXPECT().Send(tt.req).Return(nil)
			},
			err: nil,
		},
		{
			name: "send fails, expect error",
			req: &ports.Port{
				ID: "test-port-id",
			},
			expCalls: func(tt table, m *mocks.MockPorts_UpdateClient) {
				m.EXPECT().Send(tt.req).Return(someErr)
			},
			err: someErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			stream := mocks.NewMockPorts_UpdateClient(ctrl)
			client := mocks.NewMockPortsClient(ctrl)
			svc := Service{
				log:    zap.NewNop(),
				stream: stream,
				client: client,
			}
			if tt.expCalls != nil {
				tt.expCalls(tt, stream)
			}

			err := svc.Update(tt.ctx, tt.req)
			if tt.err != nil {
				require.Error(t, err)
				assert.Equal(t, tt.err, err)
			}

			ctrl.Finish()
		})
	}
}

func TestService_Get(t *testing.T) {
	type table struct {
		name     string
		ctx      context.Context
		expCalls func(tt table, m *mocks.MockPortsClient)
		stringID string
		id       *ports.GetPortRequest
		err      error
	}
	tests := []table{
		{
			name:     "successfully retrieve a port from the client",
			stringID: "test-port-id",
			id:       &ports.GetPortRequest{ID: "test-port-id"},
			expCalls: func(tt table, m *mocks.MockPortsClient) {
				m.EXPECT().Get(tt.ctx, tt.id).Return(getTestGRPCPort(), nil)
			},
			err: nil,
		},
		{
			name:     "get fails, expect error",
			stringID: "test-port-id",
			id:       &ports.GetPortRequest{ID: "test-port-id"},
			expCalls: func(tt table, m *mocks.MockPortsClient) {
				m.EXPECT().Get(tt.ctx, tt.id).Return(getTestGRPCPort(), someErr)
			},
			err: someErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			stream := mocks.NewMockPorts_UpdateClient(ctrl)
			client := mocks.NewMockPortsClient(ctrl)
			svc := Service{
				log:    zap.NewNop(),
				stream: stream,
				client: client,
			}
			if tt.expCalls != nil {
				tt.expCalls(tt, client)
			}

			port, err := svc.Get(tt.ctx, tt.stringID)
			if tt.err != nil {
				require.Error(t, err)
				assert.Equal(t, tt.err, err)
			} else {
				assert.IsType(t, types.Port{}, port)
				assert.Equal(t, getTestInternalPort(), port)
			}

			ctrl.Finish()
		})
	}
}

func TestService_List(t *testing.T) {
	type table struct {
		name     string
		ctx      context.Context
		expCalls func(tt table, m *mocks.MockPortsClient, listClient *mocks.MockPorts_ListClient)
		port     *ports.Port
		err      error
	}
	tests := []table{
		{
			name: "successfully retrieve two ports",
			ctx:  context.TODO(),
			port: getTestGRPCPort(),
			expCalls: func(tt table, c *mocks.MockPortsClient, listClient *mocks.MockPorts_ListClient) {
				c.EXPECT().List(tt.ctx, gomock.Any()).Return(listClient, nil)
				listClient.EXPECT().Recv().Return(tt.port, nil)
				listClient.EXPECT().Recv().Return(tt.port, nil)
				listClient.EXPECT().Recv().Return(tt.port, io.EOF)
			},
			err: nil,
		},
		{
			name: "fail on client listing",
			ctx:  context.TODO(),
			port: getTestGRPCPort(),
			expCalls: func(tt table, c *mocks.MockPortsClient, listClient *mocks.MockPorts_ListClient) {
				c.EXPECT().List(tt.ctx, gomock.Any()).Return(listClient, someErr)
			},
			err: someErr,
		},
		{
			name: "fail on recv",
			ctx:  context.TODO(),
			port: getTestGRPCPort(),
			expCalls: func(tt table, c *mocks.MockPortsClient, listClient *mocks.MockPorts_ListClient) {
				c.EXPECT().List(tt.ctx, gomock.Any()).Return(listClient, nil)
				listClient.EXPECT().Recv().Return(tt.port, someErr)
				listClient.EXPECT().Recv().Return(tt.port, nil)
				listClient.EXPECT().Recv().Return(tt.port, io.EOF)
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			listClient := mocks.NewMockPorts_ListClient(ctrl)
			client := mocks.NewMockPortsClient(ctrl)
			svc := Service{
				log:    zap.NewNop(),
				client: client,
			}
			if tt.expCalls != nil {
				tt.expCalls(tt, client, listClient)
			}

			port, err := svc.List(tt.ctx)
			if tt.err != nil {
				require.Error(t, err)
				assert.Equal(t, tt.err, err)
			} else {
				assert.IsType(t, []types.Port{}, port)
				assert.Equal(t, getMultipleInternalPorts(), port)
			}

			ctrl.Finish()
		})
	}
}

func getTestGRPCPort() *ports.Port {
	return &ports.Port{
		ID:          "aaa",
		Name:        "bbb",
		City:        "ccc",
		Country:     "ddd",
		Alias:       []string{"1", "2", "3"},
		Regions:     []string{"4", "5", "6"},
		Coordinates: []float32{1.2, 1.3, 1.4},
		Province:    "eee",
		Timezone:    "fff",
		Unlocs:      []string{"7", "8", "9"},
		Code:        "ggg",
	}
}

func getTestInternalPort() types.Port {
	return types.Port{
		ID:          "aaa",
		Name:        "bbb",
		City:        "ccc",
		Country:     "ddd",
		Alias:       []string{"1", "2", "3"},
		Regions:     []string{"4", "5", "6"},
		Coordinates: []float32{1.2, 1.3, 1.4},
		Province:    "eee",
		Timezone:    "fff",
		Unlocs:      []string{"7", "8", "9"},
		Code:        "ggg",
	}
}

func getMultipleInternalPorts() []types.Port {
	return []types.Port{
		{
			ID:          "aaa",
			Name:        "bbb",
			City:        "ccc",
			Country:     "ddd",
			Alias:       []string{"1", "2", "3"},
			Regions:     []string{"4", "5", "6"},
			Coordinates: []float32{1.2, 1.3, 1.4},
			Province:    "eee",
			Timezone:    "fff",
			Unlocs:      []string{"7", "8", "9"},
			Code:        "ggg",
		},
		{
			ID:          "aaa",
			Name:        "bbb",
			City:        "ccc",
			Country:     "ddd",
			Alias:       []string{"1", "2", "3"},
			Regions:     []string{"4", "5", "6"},
			Coordinates: []float32{1.2, 1.3, 1.4},
			Province:    "eee",
			Timezone:    "fff",
			Unlocs:      []string{"7", "8", "9"},
			Code:        "ggg",
		},
	}
}
