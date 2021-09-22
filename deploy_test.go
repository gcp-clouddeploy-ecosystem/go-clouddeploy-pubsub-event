package deploy

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const (
	testDataDir = "test"
)

func TestParseEvent(t *testing.T) {
	type testcase struct {
		filename string
		want     *PubSubMessage
	}

	publishTime, _ := time.Parse(time.RFC3339, "2021-05-17T21:31:25.143Z")

	tcs := map[string]*testcase{
		"approval required": {
			filename: filepath.Join(testDataDir, "approval-required.json"),
			want: &PubSubMessage{
				AckID: "ackId",
				Message: &Message{
					Attributes: &Attributes{
						Action:  Required,
						Rollout: "projects/120123456789/locations/us-central1/deliveryPipelines/etest/releases/f2/rollouts/rollout-123",
					},
					MessageID:   "messageId",
					PublishTime: publishTime,
				},
			},
		},
		"delivery pipeline create": {
			filename: filepath.Join(testDataDir, "delivery-pipeline-create.json"),
			want: &PubSubMessage{
				AckID: "ackId",
				Message: &Message{
					Attributes: &Attributes{
						Action:       Create,
						Resource:     "projects/120123456789/locations/us-central1/deliveryPipelines/etest",
						ResourceType: DeliveryPipeline,
					},
					MessageID:   "messageId",
					PublishTime: publishTime,
				},
			},
		},
		"render start": {
			filename: filepath.Join(testDataDir, "render-start.json"),
			want: &PubSubMessage{
				AckID: "ackId",
				Message: &Message{
					Attributes: &Attributes{
						Action:       Start,
						Resource:     "projects/120123456789/locations/us-central1/deliveryPipelines/etest/releases/f2",
						ResourceType: Release,
					},
					MessageID:   "messageId",
					PublishTime: publishTime,
				},
			},
		},
	}

	for n, tc := range tcs {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			f, err := os.Open(tc.filename)
			if err != nil {
				t.Fatalf("failed to load testdata: %v", err)
			}
			defer f.Close()

			b, _ := ioutil.ReadAll(f)

			got, err := ParseEvent(b)
			if err != nil {
				t.Fatalf("failed: %v", err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("mismatch (-got +want):%s\n", diff)
			}
		})
	}
}
