package queue

import (
	"testing"

	musicpb "github.com/kumneger0/yt-tracks/gen"
	"github.com/kumneger0/yt-tracks/internal/types"
)

func TestRingQueue_Empty(t *testing.T) {
	ringQueue := NewRingQueue()
	if ringQueue.Len() != 0 {
		t.Fatalf("expected len 0, got %d", ringQueue.Len())
	}
	if ringQueue.Current() != nil {
		t.Fatal("expected nil current track")
	}
	if ringQueue.Next() != nil {
		t.Fatal("expected nil next track")
	}
	if ringQueue.Prev() != nil {
		t.Fatal("expected nil prev track")
	}
}

func TestRingQueue_AddAndNavigate(t *testing.T) {
	ringQueue := NewRingQueue()
	t1 := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "v1", Title: "Song 1"}}
	t2 := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "v2", Title: "Song 2"}}

	ringQueue.AddTrack(t1)
	ringQueue.AddTrack(t2)

	if ringQueue.Len() != 2 {
		t.Fatalf("expected len 2, got %d", ringQueue.Len())
	}
	if ringQueue.Current().Track.VideoId != "v1" {
		t.Fatalf("expected current v1, got %s", ringQueue.Current().Track.VideoId)
	}

	next := ringQueue.Next()
	if next.Track.VideoId != "v2" {
		t.Fatalf("expected next v2, got %s", next.Track.VideoId)
	}

	nextWrap := ringQueue.Next()
	if nextWrap.Track.VideoId != "v1" {
		t.Fatalf("expected wrap v1, got %s", nextWrap.Track.VideoId)
	}

	prev := ringQueue.Prev()
	if prev.Track.VideoId != "v2" {
		t.Fatalf("expected prev v2, got %s", prev.Track.VideoId)
	}
}

func TestRingQueue_PlayNextTrack(t *testing.T) {
	ringQueue := NewRingQueue()
	t1 := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "v1"}}
	t2 := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "v2"}}
	t3 := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "v3"}}

	ringQueue.AddTrack(t1)
	ringQueue.AddTrack(t2)
	ringQueue.PlayNextTrack(t3)

	if ringQueue.Len() != 3 {
		t.Fatalf("expected len 3, got %d", ringQueue.Len())
	}

	if ringQueue.Current().Track.VideoId != "v1" {
		t.Fatalf("expected current v1, got %s", ringQueue.Current().Track.VideoId)
	}
	if ringQueue.Next().Track.VideoId != "v3" {
		t.Fatalf("expected next v3, got %s", ringQueue.Current().Track.VideoId)
	}
}

func TestRingQueue_SetAndRemove(t *testing.T) {
	ringQueue := NewRingQueue()
	tracks := []*types.PlaylistTrackObject{
		{Track: &musicpb.Song{VideoId: "v1"}},
		{Track: &musicpb.Song{VideoId: "v2"}},
	}
	ringQueue.SetTracks(tracks)
	if ringQueue.Len() != 2 {
		t.Fatalf("expected len 2, got %d", ringQueue.Len())
	}

	ringQueue.RemoveCurrent()
	if ringQueue.Len() != 1 {
		t.Fatalf("expected len 1 after remove, got %d", ringQueue.Len())
	}
	if ringQueue.Current().Track.VideoId != "v2" {
		t.Fatalf("expected current v2 after remove, got %s", ringQueue.Current().Track.VideoId)
	}
}

func TestRingQueue_UpdateTrack(t *testing.T) {
	ringQueue := NewRingQueue()

	// Update on empty queue
	track1 := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "v1", Title: "Original 1"}}
	if ringQueue.UpdateTrack(track1) {
		t.Fatal("expected update to fail on empty queue")
	}

	// Update with invalid inputs
	ringQueue.AddTrack(track1)
	if ringQueue.UpdateTrack(nil) {
		t.Fatal("expected update to fail with nil track")
	}
	if ringQueue.UpdateTrack(&types.PlaylistTrackObject{}) {
		t.Fatal("expected update to fail with track without musicpb Song")
	}

	// Update single element
	updatedT1 := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "v1", Title: "Updated 1"}}
	if !ringQueue.UpdateTrack(updatedT1) {
		t.Fatal("expected update to succeed for v1")
	}
	if ringQueue.Current().Track.Title != "Updated 1" {
		t.Fatalf("expected title 'Updated 1', got '%s'", ringQueue.Current().Track.Title)
	}

	// Update with multiple elements
	track2 := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "v2", Title: "Original 2"}}
	track3 := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "v3", Title: "Original 3"}}
	ringQueue.AddTrack(track2)
	ringQueue.AddTrack(track3)

	updatedT2 := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "v2", Title: "Updated 2"}}
	if !ringQueue.UpdateTrack(updatedT2) {
		t.Fatal("expected update to succeed for v2")
	}

	all := ringQueue.AllTracks()
	if len(all) != 3 {
		t.Fatalf("expected 3 tracks, got %d", len(all))
	}
	if all[1].Track.Title != "Updated 2" {
		t.Fatalf("expected title 'Updated 2' at index 1, got '%s'", all[1].Track.Title)
	}

	// Non-existent videoId
	tUnknown := &types.PlaylistTrackObject{Track: &musicpb.Song{VideoId: "v99", Title: "Unknown"}}
	if ringQueue.UpdateTrack(tUnknown) {
		t.Fatal("expected update to fail for non-existent videoId")
	}
}
