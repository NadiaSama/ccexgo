package spot

import (
	"testing"
	"time"
)

func TestMBPDS(t *testing.T) {
	t.Run("test depth not consistent", func(t *testing.T) {
		depth := []Depth{
			{SeqNum: 1, PrevSeqNum: 0},
			{PrevSeqNum: 1, SeqNum: 2},
			{PrevSeqNum: 3, SeqNum: 4},
		}

		ds := NewMBPDepthDS(nil)
		for i, d := range depth {
			_, err := ds.Push(&d, time.Now())
			if i != 2 {
				if err != nil {
					t.Fatalf("init fail error=%s", err.Error())
				}
			} else {
				if err == nil {
					t.Fatalf("expect error happen for i==2")
				}
			}
		}
	})

	t.Run("test with refresh message", func(t *testing.T) {
		ds := NewMBPDepthDS(nil)

		for i := 1; i < 10; i++ {
			d := Depth{
				SeqNum:     int64(i),
				PrevSeqNum: int64(i - 1),
				Bids:       [][2]float64{{float64(i) * 1.0, float64(i) * 2.0}},
				Asks:       [][2]float64{{float64(i) * 1.0, float64(i) * 2.0}},
			}

			ds.Push(&d, time.Now())
		}

		ds.AddRefresh(&Depth{
			SeqNum: 5,
			Bids:   [][2]float64{{101.0, 102.0}, {102.0, 103.0}},
			Asks:   [][2]float64{{35.0, 36.0}, {37.0, 38.0}},
		})
		t.Logf("tree bids.String=%s asks.String=%s inited=%+v", ds.bids.String(), ds.asks.String(), ds.inited)

		tree := ds.OrderBook(1)

		if elem := tree.Bids[0]; elem.Price != 102.0 || elem.Amount != 103.0 {
			t.Errorf("invalid bids %+v", elem)
		}

		if elem := tree.Asks[0]; elem.Price != 6.0 || elem.Amount != 12.0 {
			t.Errorf("invalid asks %+v", elem)
		}

		tree = ds.OrderBook(-1)
		if len(tree.Bids) != 6 || tree.Bids[5].Price != 6.0 {
			t.Errorf("invalid bids %+v", tree.Bids)
		}

		if len(tree.Asks) != 6 || tree.Asks[5].Price != 37.0 {
			t.Errorf("invalid asks %+v", tree.Asks)
		}

	})
}
