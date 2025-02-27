package phlaredb

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/oklog/ulid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"

	ingestv1 "github.com/grafana/pyroscope/api/gen/proto/go/ingester/v1"
	typesv1 "github.com/grafana/pyroscope/api/gen/proto/go/types/v1"
	"github.com/grafana/pyroscope/pkg/iter"
	"github.com/grafana/pyroscope/pkg/model"
	"github.com/grafana/pyroscope/pkg/objstore/providers/filesystem"
	"github.com/grafana/pyroscope/pkg/phlaredb/block"
	"github.com/grafana/pyroscope/pkg/phlaredb/tsdb/index"
)

func TestQuerierBlockEviction(t *testing.T) {
	type testCase struct {
		blocks     []string
		expected   []string
		notEvicted bool
	}

	blockToEvict := "01H002D4Z9PKWSS17Q3XY1VEM9"
	testCases := []testCase{
		{
			notEvicted: true,
		},
		{
			blocks:     []string{"01H002D4Z9ES0DHMMSD18H5J5M"},
			expected:   []string{"01H002D4Z9ES0DHMMSD18H5J5M"},
			notEvicted: true,
		},
		{
			blocks:   []string{blockToEvict},
			expected: []string{},
		},
		{
			blocks:   []string{blockToEvict, "01H002D4Z9ES0DHMMSD18H5J5M"},
			expected: []string{"01H002D4Z9ES0DHMMSD18H5J5M"},
		},
		{
			blocks:   []string{"01H002D4Z9ES0DHMMSD18H5J5M", blockToEvict},
			expected: []string{"01H002D4Z9ES0DHMMSD18H5J5M"},
		},
		{
			blocks:   []string{"01H002D4Z9ES0DHMMSD18H5J5M", blockToEvict, "01H003A2QTY5JF30Z441CDQE70"},
			expected: []string{"01H002D4Z9ES0DHMMSD18H5J5M", "01H003A2QTY5JF30Z441CDQE70"},
		},
		{
			blocks:   []string{"01H003A2QTY5JF30Z441CDQE70", blockToEvict, "01H002D4Z9ES0DHMMSD18H5J5M"},
			expected: []string{"01H003A2QTY5JF30Z441CDQE70", "01H002D4Z9ES0DHMMSD18H5J5M"},
		},
	}

	for _, tc := range testCases {
		q := BlockQuerier{queriers: make([]*singleBlockQuerier, len(tc.blocks))}
		for i, b := range tc.blocks {
			q.queriers[i] = &singleBlockQuerier{
				meta:    &block.Meta{ULID: ulid.MustParse(b)},
				metrics: newBlocksMetrics(nil),
			}
		}

		evicted, err := q.evict(ulid.MustParse(blockToEvict))
		require.NoError(t, err)
		require.Equal(t, !tc.notEvicted, evicted)

		actual := make([]string, 0, len(tc.expected))
		for _, b := range q.queriers {
			actual = append(actual, b.meta.ULID.String())
		}

		require.ElementsMatch(t, tc.expected, actual)
	}
}

type profileCounter struct {
	iter.Iterator[Profile]
	count int
}

func (p *profileCounter) Next() bool {
	r := p.Iterator.Next()
	if r {
		p.count++
	}

	return r
}

func TestBlockCompatability(t *testing.T) {
	path := "./block/testdata/"
	bucket, err := filesystem.NewBucket(path)
	require.NoError(t, err)

	ctx := context.Background()
	metas, err := NewBlockQuerier(ctx, bucket).BlockMetas(ctx)
	require.NoError(t, err)

	for _, meta := range metas {
		t.Run(fmt.Sprintf("block-v%d-%s", meta.Version, meta.ULID.String()), func(t *testing.T) {
			q := NewSingleBlockQuerierFromMeta(ctx, bucket, meta)
			require.NoError(t, q.Open(ctx))

			profilesTypes, err := q.index.LabelValues("__profile_type__")
			require.NoError(t, err)

			profileCount := 0

			for _, profileType := range profilesTypes {
				t.Log(profileType)
				profileTypeParts := strings.Split(profileType, ":")

				it, err := q.SelectMatchingProfiles(ctx, &ingestv1.SelectProfilesRequest{
					LabelSelector: "{}",
					Start:         0,
					End:           time.Now().UnixMilli(),
					Type: &typesv1.ProfileType{
						Name:       profileTypeParts[0],
						SampleType: profileTypeParts[1],
						SampleUnit: profileTypeParts[2],
						PeriodType: profileTypeParts[3],
						PeriodUnit: profileTypeParts[4],
					},
				})
				require.NoError(t, err)

				pcIt := &profileCounter{Iterator: it}

				// TODO: It would be nice actually comparing the whole profile, but at present the result is not deterministic.
				_, err = q.MergePprof(ctx, pcIt)
				require.NoError(t, err)

				profileCount += pcIt.count
			}

			require.Equal(t, int(meta.Stats.NumProfiles), profileCount)
		})
	}
}

func TestBlockCompatability_SelectMergeSpans(t *testing.T) {
	path := "./block/testdata/"
	bucket, err := filesystem.NewBucket(path)
	require.NoError(t, err)

	ctx := context.Background()
	metas, err := NewBlockQuerier(ctx, bucket).BlockMetas(ctx)
	require.NoError(t, err)

	for _, meta := range metas {
		t.Run(fmt.Sprintf("block-v%d-%s", meta.Version, meta.ULID.String()), func(t *testing.T) {
			q := NewSingleBlockQuerierFromMeta(ctx, bucket, meta)
			require.NoError(t, q.Open(ctx))

			profilesTypes, err := q.index.LabelValues("__profile_type__")
			require.NoError(t, err)

			profileCount := 0

			for _, profileType := range profilesTypes {
				t.Log(profileType)
				profileTypeParts := strings.Split(profileType, ":")

				it, err := q.SelectMatchingProfiles(ctx, &ingestv1.SelectProfilesRequest{
					LabelSelector: "{}",
					Start:         0,
					End:           time.Now().UnixMilli(),
					Type: &typesv1.ProfileType{
						Name:       profileTypeParts[0],
						SampleType: profileTypeParts[1],
						SampleUnit: profileTypeParts[2],
						PeriodType: profileTypeParts[3],
						PeriodUnit: profileTypeParts[4],
					},
				})
				require.NoError(t, err)

				pcIt := &profileCounter{Iterator: it}

				spanSelector, err := model.NewSpanSelector([]string{})
				require.NoError(t, err)
				resp, err := q.MergeBySpans(ctx, pcIt, spanSelector)
				require.NoError(t, err)

				require.Zero(t, resp.Total())
				profileCount += pcIt.count
			}

			require.Zero(t, profileCount)
		})
	}
}

type fakeQuerier struct {
	Querier
	doErr bool
}

func (f *fakeQuerier) SelectMatchingProfiles(ctx context.Context, params *ingestv1.SelectProfilesRequest) (iter.Iterator[Profile], error) {
	// add some jitter
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	if f.doErr {
		return nil, fmt.Errorf("fake error")
	}
	profiles := []Profile{}
	for i := 0; i < 100000; i++ {
		profiles = append(profiles, BlockProfile{})
	}
	return iter.NewSliceIterator(profiles), nil
}

func TestSelectMatchingProfilesCleanUp(t *testing.T) {
	defer goleak.VerifyNone(t, goleak.IgnoreCurrent())

	_, err := SelectMatchingProfiles(context.Background(), &ingestv1.SelectProfilesRequest{}, Queriers{
		&fakeQuerier{},
		&fakeQuerier{},
		&fakeQuerier{},
		&fakeQuerier{},
		&fakeQuerier{doErr: true},
	})
	require.Error(t, err)
}

func Test_singleBlockQuerier_Series(t *testing.T) {
	ctx := context.Background()
	reader, err := index.NewFileReader("testdata/01HA2V3CPSZ9E0HMQNNHH89WSS/index.tsdb")
	assert.NoError(t, err)

	q := &singleBlockQuerier{
		metrics: newBlocksMetrics(nil),
		meta:    &block.Meta{ULID: ulid.MustParse("01HA2V3CPSZ9E0HMQNNHH89WSS")},
		opened:  true, // Skip trying to open the block.
		index:   reader,
	}

	t.Run("get all names", func(t *testing.T) {
		want := []string{
			"__delta__",
			"__name__",
			"__period_type__",
			"__period_unit__",
			"__profile_type__",
			"__service_name__",
			"__type__",
			"__unit__",
			"foo",
			"function",
			"pyroscope_spy",
			"service_name",
			"target",
			"version",
		}
		got, err := q.index.LabelNames()
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("get label", func(t *testing.T) {
		want := []*typesv1.Labels{
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "block"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "goroutine"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "mutex"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "process_cpu"},
			}},
		}
		got, err := q.Series(ctx, &ingestv1.SeriesRequest{
			LabelNames: []string{
				"__name__",
			},
		})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("get label with matcher", func(t *testing.T) {
		want := []*typesv1.Labels{
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "block"},
			}},
		}
		got, err := q.Series(ctx, &ingestv1.SeriesRequest{
			Matchers:   []string{`{__name__="block"}`},
			LabelNames: []string{"__name__"},
		})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("get multiple labels", func(t *testing.T) {
		want := []*typesv1.Labels{
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "block"},
				{Name: "__type__", Value: "contentions"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "block"},
				{Name: "__type__", Value: "delay"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "goroutine"},
				{Name: "__type__", Value: "goroutines"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
				{Name: "__type__", Value: "alloc_objects"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
				{Name: "__type__", Value: "alloc_space"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
				{Name: "__type__", Value: "inuse_objects"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
				{Name: "__type__", Value: "inuse_space"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "mutex"},
				{Name: "__type__", Value: "contentions"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "mutex"},
				{Name: "__type__", Value: "delay"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "process_cpu"},
				{Name: "__type__", Value: "cpu"},
			}},
		}
		got, err := q.Series(ctx, &ingestv1.SeriesRequest{
			LabelNames: []string{"__name__", "__type__"},
		})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("get multiple labels with matcher", func(t *testing.T) {
		want := []*typesv1.Labels{
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
				{Name: "__type__", Value: "alloc_objects"},
			}},
		}
		got, err := q.Series(ctx, &ingestv1.SeriesRequest{
			Matchers:   []string{`{__name__="memory",__type__="alloc_objects"}`},
			LabelNames: []string{"__name__", "__type__"},
		})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("empty labels and empty matcher", func(t *testing.T) {
		want := []*typesv1.Labels{
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "block"},
				{Name: "__profile_type__", Value: "block:contentions:count::"},
				{Name: "__service_name__", Value: "pyroscope"},
				{Name: "__type__", Value: "contentions"},
				{Name: "__unit__", Value: "count"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "pyroscope"},
				{Name: "target", Value: "all"},
				{Name: "version", Value: "label-names-store-gateway-0e430f1e-WIP"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "block"},
				{Name: "__profile_type__", Value: "block:delay:nanoseconds::"},
				{Name: "__service_name__", Value: "pyroscope"},
				{Name: "__type__", Value: "delay"},
				{Name: "__unit__", Value: "nanoseconds"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "pyroscope"},
				{Name: "target", Value: "all"},
				{Name: "version", Value: "label-names-store-gateway-0e430f1e-WIP"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "goroutine"},
				{Name: "__profile_type__", Value: "goroutine:goroutines:count::"},
				{Name: "__service_name__", Value: "pyroscope"},
				{Name: "__type__", Value: "goroutines"},
				{Name: "__unit__", Value: "count"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "pyroscope"},
				{Name: "target", Value: "all"},
				{Name: "version", Value: "label-names-store-gateway-0e430f1e-WIP"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:alloc_objects:count::"},
				{Name: "__service_name__", Value: "pyroscope"},
				{Name: "__type__", Value: "alloc_objects"},
				{Name: "__unit__", Value: "count"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "pyroscope"},
				{Name: "target", Value: "all"},
				{Name: "version", Value: "label-names-store-gateway-0e430f1e-WIP"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:alloc_objects:count::"},
				{Name: "__service_name__", Value: "simple.golang.app"},
				{Name: "__type__", Value: "alloc_objects"},
				{Name: "__unit__", Value: "count"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "simple.golang.app"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:alloc_space:bytes::"},
				{Name: "__service_name__", Value: "pyroscope"},
				{Name: "__type__", Value: "alloc_space"},
				{Name: "__unit__", Value: "bytes"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "pyroscope"},
				{Name: "target", Value: "all"},
				{Name: "version", Value: "label-names-store-gateway-0e430f1e-WIP"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:alloc_space:bytes::"},
				{Name: "__service_name__", Value: "simple.golang.app"},
				{Name: "__type__", Value: "alloc_space"},
				{Name: "__unit__", Value: "bytes"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "simple.golang.app"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:inuse_objects:count::"},
				{Name: "__service_name__", Value: "pyroscope"},
				{Name: "__type__", Value: "inuse_objects"},
				{Name: "__unit__", Value: "count"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "pyroscope"},
				{Name: "target", Value: "all"},
				{Name: "version", Value: "label-names-store-gateway-0e430f1e-WIP"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:inuse_objects:count::"},
				{Name: "__service_name__", Value: "simple.golang.app"},
				{Name: "__type__", Value: "inuse_objects"},
				{Name: "__unit__", Value: "count"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "simple.golang.app"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:inuse_space:bytes::"},
				{Name: "__service_name__", Value: "pyroscope"},
				{Name: "__type__", Value: "inuse_space"},
				{Name: "__unit__", Value: "bytes"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "pyroscope"},
				{Name: "target", Value: "all"},
				{Name: "version", Value: "label-names-store-gateway-0e430f1e-WIP"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:inuse_space:bytes::"},
				{Name: "__service_name__", Value: "simple.golang.app"},
				{Name: "__type__", Value: "inuse_space"},
				{Name: "__unit__", Value: "bytes"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "simple.golang.app"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "mutex"},
				{Name: "__profile_type__", Value: "mutex:contentions:count::"},
				{Name: "__service_name__", Value: "pyroscope"},
				{Name: "__type__", Value: "contentions"},
				{Name: "__unit__", Value: "count"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "pyroscope"},
				{Name: "target", Value: "all"},
				{Name: "version", Value: "label-names-store-gateway-0e430f1e-WIP"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "mutex"},
				{Name: "__profile_type__", Value: "mutex:delay:nanoseconds::"},
				{Name: "__service_name__", Value: "pyroscope"},
				{Name: "__type__", Value: "delay"},
				{Name: "__unit__", Value: "nanoseconds"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "pyroscope"},
				{Name: "target", Value: "all"},
				{Name: "version", Value: "label-names-store-gateway-0e430f1e-WIP"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "process_cpu"},
				{Name: "__period_type__", Value: "cpu"},
				{Name: "__period_unit__", Value: "nanoseconds"},
				{Name: "__profile_type__", Value: "process_cpu:cpu:nanoseconds:cpu:nanoseconds"},
				{Name: "__service_name__", Value: "pyroscope"},
				{Name: "__type__", Value: "cpu"},
				{Name: "__unit__", Value: "nanoseconds"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "pyroscope"},
				{Name: "target", Value: "all"},
				{Name: "version", Value: "label-names-store-gateway-0e430f1e-WIP"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "process_cpu"},
				{Name: "__period_type__", Value: "cpu"},
				{Name: "__period_unit__", Value: "nanoseconds"},
				{Name: "__profile_type__", Value: "process_cpu:cpu:nanoseconds:cpu:nanoseconds"},
				{Name: "__service_name__", Value: "simple.golang.app"},
				{Name: "__type__", Value: "cpu"},
				{Name: "__unit__", Value: "nanoseconds"},
				{Name: "foo", Value: "bar"},
				{Name: "function", Value: "fast"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "simple.golang.app"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "process_cpu"},
				{Name: "__period_type__", Value: "cpu"},
				{Name: "__period_unit__", Value: "nanoseconds"},
				{Name: "__profile_type__", Value: "process_cpu:cpu:nanoseconds:cpu:nanoseconds"},
				{Name: "__service_name__", Value: "simple.golang.app"},
				{Name: "__type__", Value: "cpu"},
				{Name: "__unit__", Value: "nanoseconds"},
				{Name: "foo", Value: "bar"},
				{Name: "function", Value: "slow"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "simple.golang.app"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "process_cpu"},
				{Name: "__period_type__", Value: "cpu"},
				{Name: "__period_unit__", Value: "nanoseconds"},
				{Name: "__profile_type__", Value: "process_cpu:cpu:nanoseconds:cpu:nanoseconds"},
				{Name: "__service_name__", Value: "simple.golang.app"},
				{Name: "__type__", Value: "cpu"},
				{Name: "__unit__", Value: "nanoseconds"},
				{Name: "foo", Value: "bar"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "simple.golang.app"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__delta__", Value: "false"},
				{Name: "__name__", Value: "process_cpu"},
				{Name: "__period_type__", Value: "cpu"},
				{Name: "__period_unit__", Value: "nanoseconds"},
				{Name: "__profile_type__", Value: "process_cpu:cpu:nanoseconds:cpu:nanoseconds"},
				{Name: "__service_name__", Value: "simple.golang.app"},
				{Name: "__type__", Value: "cpu"},
				{Name: "__unit__", Value: "nanoseconds"},
				{Name: "pyroscope_spy", Value: "gospy"},
				{Name: "service_name", Value: "simple.golang.app"},
			}},
		}
		got, err := q.Series(ctx, &ingestv1.SeriesRequest{
			Matchers:   []string{},
			LabelNames: []string{},
		})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("ui plugin", func(t *testing.T) {
		want := []*typesv1.Labels{
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "block"},
				{Name: "__profile_type__", Value: "block:contentions:count::"},
				{Name: "__type__", Value: "contentions"},
				{Name: "service_name", Value: "pyroscope"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "block"},
				{Name: "__profile_type__", Value: "block:delay:nanoseconds::"},
				{Name: "__type__", Value: "delay"},
				{Name: "service_name", Value: "pyroscope"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "goroutine"},
				{Name: "__profile_type__", Value: "goroutine:goroutines:count::"},
				{Name: "__type__", Value: "goroutines"},
				{Name: "service_name", Value: "pyroscope"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:alloc_objects:count::"},
				{Name: "__type__", Value: "alloc_objects"},
				{Name: "service_name", Value: "pyroscope"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:alloc_objects:count::"},
				{Name: "__type__", Value: "alloc_objects"},
				{Name: "service_name", Value: "simple.golang.app"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:alloc_space:bytes::"},
				{Name: "__type__", Value: "alloc_space"},
				{Name: "service_name", Value: "pyroscope"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:alloc_space:bytes::"},
				{Name: "__type__", Value: "alloc_space"},
				{Name: "service_name", Value: "simple.golang.app"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:inuse_objects:count::"},
				{Name: "__type__", Value: "inuse_objects"},
				{Name: "service_name", Value: "pyroscope"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:inuse_objects:count::"},
				{Name: "__type__", Value: "inuse_objects"},
				{Name: "service_name", Value: "simple.golang.app"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:inuse_space:bytes::"},
				{Name: "__type__", Value: "inuse_space"},
				{Name: "service_name", Value: "pyroscope"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "memory"},
				{Name: "__profile_type__", Value: "memory:inuse_space:bytes::"},
				{Name: "__type__", Value: "inuse_space"},
				{Name: "service_name", Value: "simple.golang.app"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "mutex"},
				{Name: "__profile_type__", Value: "mutex:contentions:count::"},
				{Name: "__type__", Value: "contentions"},
				{Name: "service_name", Value: "pyroscope"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "mutex"},
				{Name: "__profile_type__", Value: "mutex:delay:nanoseconds::"},
				{Name: "__type__", Value: "delay"},
				{Name: "service_name", Value: "pyroscope"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "process_cpu"},
				{Name: "__profile_type__", Value: "process_cpu:cpu:nanoseconds:cpu:nanoseconds"},
				{Name: "__type__", Value: "cpu"},
				{Name: "service_name", Value: "pyroscope"},
			}},
			{Labels: []*typesv1.LabelPair{
				{Name: "__name__", Value: "process_cpu"},
				{Name: "__profile_type__", Value: "process_cpu:cpu:nanoseconds:cpu:nanoseconds"},
				{Name: "__type__", Value: "cpu"},
				{Name: "service_name", Value: "simple.golang.app"},
			}},
		}
		got, err := q.Series(ctx, &ingestv1.SeriesRequest{
			Matchers: []string{},
			LabelNames: []string{
				"pyroscope_app",
				"service_name",
				"__profile_type__",
				"__type__",
				"__name__",
			},
		})

		assert.NoError(t, err)
		assert.Equal(t, want, got)
	})
}

func Benchmark_singleBlockQuerier_Series(b *testing.B) {
	ctx := context.Background()
	reader, err := index.NewFileReader("testdata/01HA2V3CPSZ9E0HMQNNHH89WSS/index.tsdb")
	assert.NoError(b, err)

	q := &singleBlockQuerier{
		metrics: newBlocksMetrics(nil),
		meta:    &block.Meta{ULID: ulid.MustParse("01HA2V3CPSZ9E0HMQNNHH89WSS")},
		opened:  true, // Skip trying to open the block.
		index:   reader,
	}

	b.Run("multiple labels", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			q.Series(ctx, &ingestv1.SeriesRequest{ //nolint:errcheck
				Matchers:   []string{`{__name__="block"}`},
				LabelNames: []string{"__name__"},
			})
		}
	})

	b.Run("multiple labels with matcher", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			q.Series(ctx, &ingestv1.SeriesRequest{ //nolint:errcheck
				Matchers:   []string{`{__name__="memory",__type__="alloc_objects"}`},
				LabelNames: []string{"__name__", "__type__"},
			})
		}
	})
}
