package stream

import (
	"fmt"
	types "github.com/jack-hughes/ports/internal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestStream_Start(t *testing.T) {
	t.Run("cannot open file", func(t *testing.T) {
		s := NewJSONStream(zap.NewNop())
		go func() {
			for data := range s.Watch() {
				if data.Error != nil {
					require.Error(t, data.Error)
					assert.Equal(t, fmt.Errorf("open file: open broken-file-path: no such file or directory"), data.Error)
				}
			}
		}()

		s.Start("broken-file-path")
	})
	t.Run("can open file, but json cannot be decoded", func(t *testing.T) {
		s := NewJSONStream(zap.NewNop())
		go func() {
			for data := range s.Watch() {
				if data.Error != nil {
					require.Error(t, data.Error)
					assert.Equal(t, fmt.Errorf("decode opening delimiter: invalid character 'h' looking for beginning of value"), data.Error)
				}
			}
		}()

		s.Start("../../../test/testdata/broken.json")
	})
	t.Run("fail to decode second line", func(t *testing.T) {
		s := NewJSONStream(zap.NewNop())
		go func() {
			for data := range s.Watch() {
				if data.Error != nil {
					require.Error(t, data.Error)
					assert.Equal(t, fmt.Errorf("decode opening delimiter: invalid character 'z'"), data.Error)
				}
			}
		}()

		s.Start("../../../test/testdata/broken-1.json")
	})
	t.Run("fail to decode port", func(t *testing.T) {
		s := NewJSONStream(zap.NewNop())
		go func() {
			for data := range s.Watch() {
				if data.Error != nil {
					require.Error(t, data.Error)
					assert.Equal(t, fmt.Errorf("decode line failure json: cannot unmarshal string into Go value of type types.Port"), data.Error)
				}
			}
		}()

		s.Start("../../../test/testdata/broken-2.json")
	})
	t.Run("working port - success", func(t *testing.T) {
		s := NewJSONStream(zap.NewNop())
		go func() {
			for data := range s.Watch() {
				assert.Equal(t, types.Port{
					Name:        "Ajman",
					City:        "Ajman",
					Country:     "United Arab Emirates",
					Alias:       []string{},
					Regions:     []string{},
					Coordinates: []float32{
						55.5136433,
						25.4052165,
					},
					Province:    "Ajman",
					Timezone:    "Asia/Dubai",
					Unlocs:      []string{
						"AEAJM",
					},
					Code:        "52000",
				}, data.Port)
			}
		}()

		s.Start("../../../test/testdata/working-port.json")
	})
	t.Run("broken closing delimiter", func(t *testing.T) {
		s := NewJSONStream(zap.NewNop())
		go func() {
			for data := range s.Watch() {
				if data.Error != nil {
					require.Error(t, data.Error)
					assert.Equal(t, fmt.Errorf("decode closing delimiter: EOF"), data.Error)
				}
			}
		}()

		s.Start("../../../test/testdata/broken-3.json")
	})
}
