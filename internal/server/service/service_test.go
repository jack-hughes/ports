package service

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	types "github.com/jack-hughes/ports/internal"
	"github.com/jack-hughes/ports/pkg/apis/ports"
	"github.com/jack-hughes/ports/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"testing"
)

var errFoo = fmt.Errorf("some-error")
func TestService_New(t *testing.T) {
	t.Run("test new successfully", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		store := mocks.NewMockStorage(ctrl)

		res := New(store, zap.NewNop())
		assert.IsType(t, &PortsServer{}, res)

		ctrl.Finish()
	})
}

func TestService_Update(t *testing.T) {
	type table struct {
		name     string
		expCalls func(tt table, st *mocks.MockStorage, ms *mocks.MockPorts_UpdateServer)
		err      error
	}
	tests := []table{
		{
			name: "successfully update a record in storage",
			expCalls: func(tt table, st *mocks.MockStorage, ms *mocks.MockPorts_UpdateServer) {
				ms.EXPECT().Recv().Return(getTestGRPCPort(), nil)
				ms.EXPECT().Recv().Return(getTestGRPCPort(), io.EOF)
				ms.EXPECT().SendAndClose(gomock.Any()).Return(nil)
				st.EXPECT().Update(types.Clone(getTestGRPCPort()))
			},
			err: nil,
		},
		{
			name: "fail to recv a record",
			expCalls: func(tt table, st *mocks.MockStorage, ms *mocks.MockPorts_UpdateServer) {
				ms.EXPECT().Recv().Return(getTestGRPCPort(), errFoo)
			},
			err: errFoo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mocks.NewMockStorage(ctrl)
			stream := mocks.NewMockPorts_UpdateServer(ctrl)
			ps := PortsServer{
				log:   zap.NewNop(),
				store: store,
			}
			if tt.expCalls != nil {
				tt.expCalls(tt, store, stream)
			}

			err := ps.Update(stream)
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
		portReq  *ports.GetPortRequest
		expCalls func(tt table, st *mocks.MockStorage)
		err      error
	}
	tests := []table{
		{
			name: "successfully retrieve a record from storage",
			portReq: &ports.GetPortRequest{ID: "test-port-id"},
			expCalls: func(tt table, st *mocks.MockStorage) {
				st.EXPECT().Get(tt.portReq.ID).Return(getTestInternalPort(), nil)
			},
			err: nil,
		},
		{
			name: "fail to retrieve a record from storage",
			portReq: &ports.GetPortRequest{ID: "test-port-id"},
			expCalls: func(tt table, st *mocks.MockStorage) {
				st.EXPECT().Get(tt.portReq.ID).Return(getTestInternalPort(), errFoo)
			},
			err: errFoo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mocks.NewMockStorage(ctrl)
			ps := PortsServer{
				log:   zap.NewNop(),
				store: store,
			}
			if tt.expCalls != nil {
				tt.expCalls(tt, store)
			}

			port, err := ps.Get(tt.ctx, tt.portReq)
			if tt.err != nil {
				require.Error(t, err)
				assert.Equal(t, tt.err, err)
			}  else {
				assert.IsType(t, &ports.Port{}, port)
				assert.Equal(t, getTestGRPCPort(), port)
			}

			ctrl.Finish()
		})
	}
}

func TestService_List(t *testing.T) {
	type table struct {
		name     string
		expCalls func(tt table, st *mocks.MockStorage, ms *mocks.MockPorts_ListServer)
		err      error
	}
	tests := []table{
		{
			name: "successfully list all records in storage",
			expCalls: func(tt table, st *mocks.MockStorage, ls *mocks.MockPorts_ListServer) {
				st.EXPECT().List().Return(getMultipleInternalPorts())
				ls.EXPECT().Send(gomock.Any()).Return(nil)
				ls.EXPECT().Send(gomock.Any()).Return(nil)
			},
			err: nil,
		},
		{
			name: "fail on second srv send",
			expCalls: func(tt table, st *mocks.MockStorage, ls *mocks.MockPorts_ListServer) {
				st.EXPECT().List().Return(getMultipleInternalPorts())
				ls.EXPECT().Send(gomock.Any()).Return(nil)
				ls.EXPECT().Send(gomock.Any()).Return(errFoo)
			},
			err: errFoo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mocks.NewMockStorage(ctrl)
			stream := mocks.NewMockPorts_ListServer(ctrl)
			ps := PortsServer{
				log:   zap.NewNop(),
				store: store,
			}
			if tt.expCalls != nil {
				tt.expCalls(tt, store, stream)
			}

			err := ps.List(&empty.Empty{}, stream)
			if tt.err != nil {
				require.Error(t, err)
				assert.Equal(t, tt.err, err)
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
