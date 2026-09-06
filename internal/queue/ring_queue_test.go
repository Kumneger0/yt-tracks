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
