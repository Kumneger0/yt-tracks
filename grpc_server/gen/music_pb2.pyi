from google.protobuf.internal import containers as _containers
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from collections.abc import Iterable as _Iterable, Mapping as _Mapping
from typing import ClassVar as _ClassVar, Optional as _Optional, Union as _Union

DESCRIPTOR: _descriptor.FileDescriptor

class HealthCheckRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class HealthCheckResponse(_message.Message):
    __slots__ = ("ok",)
    OK_FIELD_NUMBER: _ClassVar[int]
    ok: bool
    def __init__(self, ok: _Optional[bool] = ...) -> None: ...

class Song(_message.Message):
    __slots__ = ("video_id", "title", "artists", "album", "album_id", "duration_seconds", "liked", "thumbnails", "is_explicit", "url", "set_video_id")
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    ARTISTS_FIELD_NUMBER: _ClassVar[int]
    ALBUM_FIELD_NUMBER: _ClassVar[int]
    ALBUM_ID_FIELD_NUMBER: _ClassVar[int]
    DURATION_SECONDS_FIELD_NUMBER: _ClassVar[int]
    LIKED_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    IS_EXPLICIT_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    SET_VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    video_id: str
    title: str
    artists: _containers.RepeatedCompositeFieldContainer[Artist]
    album: str
    album_id: str
    duration_seconds: int
    liked: bool
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    is_explicit: bool
    url: str
    set_video_id: str
    def __init__(self, video_id: _Optional[str] = ..., title: _Optional[str] = ..., artists: _Optional[_Iterable[_Union[Artist, _Mapping]]] = ..., album: _Optional[str] = ..., album_id: _Optional[str] = ..., duration_seconds: _Optional[int] = ..., liked: _Optional[bool] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ..., is_explicit: _Optional[bool] = ..., url: _Optional[str] = ..., set_video_id: _Optional[str] = ...) -> None: ...

class Artist(_message.Message):
    __slots__ = ("id", "name")
    ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    id: str
    name: str
    def __init__(self, id: _Optional[str] = ..., name: _Optional[str] = ...) -> None: ...

class Thumbnail(_message.Message):
    __slots__ = ("url", "width", "height")
    URL_FIELD_NUMBER: _ClassVar[int]
    WIDTH_FIELD_NUMBER: _ClassVar[int]
    HEIGHT_FIELD_NUMBER: _ClassVar[int]
    url: str
    width: int
    height: int
    def __init__(self, url: _Optional[str] = ..., width: _Optional[int] = ..., height: _Optional[int] = ...) -> None: ...

class Album(_message.Message):
    __slots__ = ("browse_id", "title", "artists", "year", "thumbnails", "is_explicit", "type")
    BROWSE_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    ARTISTS_FIELD_NUMBER: _ClassVar[int]
    YEAR_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    IS_EXPLICIT_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    browse_id: str
    title: str
    artists: _containers.RepeatedCompositeFieldContainer[Artist]
    year: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    is_explicit: bool
    type: str
    def __init__(self, browse_id: _Optional[str] = ..., title: _Optional[str] = ..., artists: _Optional[_Iterable[_Union[Artist, _Mapping]]] = ..., year: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ..., is_explicit: _Optional[bool] = ..., type: _Optional[str] = ...) -> None: ...

class LibraryChannel(_message.Message):
    __slots__ = ("browse_id", "name", "subscribers", "thumbnails")
    BROWSE_ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    SUBSCRIBERS_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    browse_id: str
    name: str
    subscribers: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    def __init__(self, browse_id: _Optional[str] = ..., name: _Optional[str] = ..., subscribers: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ...) -> None: ...

class Playlist(_message.Message):
    __slots__ = ("playlist_id", "title", "description", "count", "thumbnails", "author")
    PLAYLIST_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    COUNT_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    playlist_id: str
    title: str
    description: str
    count: int
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    author: str
    def __init__(self, playlist_id: _Optional[str] = ..., title: _Optional[str] = ..., description: _Optional[str] = ..., count: _Optional[int] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ..., author: _Optional[str] = ...) -> None: ...

class UserProfile(_message.Message):
    __slots__ = ("channel_id", "name", "thumbnails")
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    channel_id: str
    name: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    def __init__(self, channel_id: _Optional[str] = ..., name: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ...) -> None: ...

class LoginRequest(_message.Message):
    __slots__ = ("auth_json",)
    AUTH_JSON_FIELD_NUMBER: _ClassVar[int]
    auth_json: str
    def __init__(self, auth_json: _Optional[str] = ...) -> None: ...

class LoginResponse(_message.Message):
    __slots__ = ("authenticated", "error", "user_name")
    AUTHENTICATED_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    USER_NAME_FIELD_NUMBER: _ClassVar[int]
    authenticated: bool
    error: str
    user_name: str
    def __init__(self, authenticated: _Optional[bool] = ..., error: _Optional[str] = ..., user_name: _Optional[str] = ...) -> None: ...

class GetLibraryRequest(_message.Message):
    __slots__ = ("limit", "continuation")
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    CONTINUATION_FIELD_NUMBER: _ClassVar[int]
    limit: int
    continuation: str
    def __init__(self, limit: _Optional[int] = ..., continuation: _Optional[str] = ...) -> None: ...

class Podcast(_message.Message):
    __slots__ = ("podcast_id", "title", "author", "thumbnails")
    PODCAST_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    podcast_id: str
    title: str
    author: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    def __init__(self, podcast_id: _Optional[str] = ..., title: _Optional[str] = ..., author: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ...) -> None: ...

class GetLibraryResponse(_message.Message):
    __slots__ = ("songs", "albums", "playlists", "channels", "artists", "podcasts")
    SONGS_FIELD_NUMBER: _ClassVar[int]
    ALBUMS_FIELD_NUMBER: _ClassVar[int]
    PLAYLISTS_FIELD_NUMBER: _ClassVar[int]
    CHANNELS_FIELD_NUMBER: _ClassVar[int]
    ARTISTS_FIELD_NUMBER: _ClassVar[int]
    PODCASTS_FIELD_NUMBER: _ClassVar[int]
    songs: _containers.RepeatedCompositeFieldContainer[Song]
    albums: _containers.RepeatedCompositeFieldContainer[Album]
    playlists: _containers.RepeatedCompositeFieldContainer[Playlist]
    channels: _containers.RepeatedCompositeFieldContainer[LibraryChannel]
    artists: _containers.RepeatedCompositeFieldContainer[FollowedArtist]
    podcasts: _containers.RepeatedCompositeFieldContainer[Podcast]
    def __init__(self, songs: _Optional[_Iterable[_Union[Song, _Mapping]]] = ..., albums: _Optional[_Iterable[_Union[Album, _Mapping]]] = ..., playlists: _Optional[_Iterable[_Union[Playlist, _Mapping]]] = ..., channels: _Optional[_Iterable[_Union[LibraryChannel, _Mapping]]] = ..., artists: _Optional[_Iterable[_Union[FollowedArtist, _Mapping]]] = ..., podcasts: _Optional[_Iterable[_Union[Podcast, _Mapping]]] = ...) -> None: ...

class GetUserSavedTracksRequest(_message.Message):
    __slots__ = ("limit",)
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    limit: int
    def __init__(self, limit: _Optional[int] = ...) -> None: ...

class GetUserSavedTracksResponse(_message.Message):
    __slots__ = ("tracks", "total")
    TRACKS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    tracks: _containers.RepeatedCompositeFieldContainer[Song]
    total: int
    def __init__(self, tracks: _Optional[_Iterable[_Union[Song, _Mapping]]] = ..., total: _Optional[int] = ...) -> None: ...

class GetUserSavedAlbumsRequest(_message.Message):
    __slots__ = ("limit",)
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    limit: int
    def __init__(self, limit: _Optional[int] = ...) -> None: ...

class GetUserSavedAlbumsResponse(_message.Message):
    __slots__ = ("albums", "total")
    ALBUMS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    albums: _containers.RepeatedCompositeFieldContainer[Album]
    total: int
    def __init__(self, albums: _Optional[_Iterable[_Union[Album, _Mapping]]] = ..., total: _Optional[int] = ...) -> None: ...

class GetUserPlaylistsRequest(_message.Message):
    __slots__ = ("limit",)
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    limit: int
    def __init__(self, limit: _Optional[int] = ...) -> None: ...

class GetUserPlaylistsResponse(_message.Message):
    __slots__ = ("playlists", "total")
    PLAYLISTS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    playlists: _containers.RepeatedCompositeFieldContainer[Playlist]
    total: int
    def __init__(self, playlists: _Optional[_Iterable[_Union[Playlist, _Mapping]]] = ..., total: _Optional[int] = ...) -> None: ...

class GetTrackRequest(_message.Message):
    __slots__ = ("video_id",)
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    video_id: str
    def __init__(self, video_id: _Optional[str] = ...) -> None: ...

class GetTrackResponse(_message.Message):
    __slots__ = ("track",)
    TRACK_FIELD_NUMBER: _ClassVar[int]
    track: Song
    def __init__(self, track: _Optional[_Union[Song, _Mapping]] = ...) -> None: ...

class GetAlbumTracksRequest(_message.Message):
    __slots__ = ("browse_id",)
    BROWSE_ID_FIELD_NUMBER: _ClassVar[int]
    browse_id: str
    def __init__(self, browse_id: _Optional[str] = ...) -> None: ...

class GetAlbumTracksResponse(_message.Message):
    __slots__ = ("title", "artists", "year", "thumbnails", "tracks", "total", "description")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    ARTISTS_FIELD_NUMBER: _ClassVar[int]
    YEAR_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    TRACKS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    title: str
    artists: _containers.RepeatedCompositeFieldContainer[Artist]
    year: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    tracks: _containers.RepeatedCompositeFieldContainer[Song]
    total: int
    description: str
    def __init__(self, title: _Optional[str] = ..., artists: _Optional[_Iterable[_Union[Artist, _Mapping]]] = ..., year: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ..., tracks: _Optional[_Iterable[_Union[Song, _Mapping]]] = ..., total: _Optional[int] = ..., description: _Optional[str] = ...) -> None: ...

class GetPlaylistItemsRequest(_message.Message):
    __slots__ = ("playlist_id", "limit")
    PLAYLIST_ID_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    playlist_id: str
    limit: int
    def __init__(self, playlist_id: _Optional[str] = ..., limit: _Optional[int] = ...) -> None: ...

class GetWatchPlaylistItemsRequest(_message.Message):
    __slots__ = ("video_id", "limit")
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    video_id: str
    limit: int
    def __init__(self, video_id: _Optional[str] = ..., limit: _Optional[int] = ...) -> None: ...

class GetWatchPlaylistItemsResponse(_message.Message):
    __slots__ = ("playlist_id", "tracks")
    PLAYLIST_ID_FIELD_NUMBER: _ClassVar[int]
    TRACKS_FIELD_NUMBER: _ClassVar[int]
    playlist_id: str
    tracks: _containers.RepeatedCompositeFieldContainer[Song]
    def __init__(self, playlist_id: _Optional[str] = ..., tracks: _Optional[_Iterable[_Union[Song, _Mapping]]] = ...) -> None: ...

class GetPlaylistItemsResponse(_message.Message):
    __slots__ = ("title", "description", "author", "year", "track_count", "thumbnails", "tracks")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    YEAR_FIELD_NUMBER: _ClassVar[int]
    TRACK_COUNT_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    TRACKS_FIELD_NUMBER: _ClassVar[int]
    title: str
    description: str
    author: str
    year: str
    track_count: int
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    tracks: _containers.RepeatedCompositeFieldContainer[Song]
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., author: _Optional[str] = ..., year: _Optional[str] = ..., track_count: _Optional[int] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ..., tracks: _Optional[_Iterable[_Union[Song, _Mapping]]] = ...) -> None: ...

class GetSearchResultsRequest(_message.Message):
    __slots__ = ("query", "filter", "limit")
    QUERY_FIELD_NUMBER: _ClassVar[int]
    FILTER_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    query: str
    filter: str
    limit: int
    def __init__(self, query: _Optional[str] = ..., filter: _Optional[str] = ..., limit: _Optional[int] = ...) -> None: ...

class SearchResultSong(_message.Message):
    __slots__ = ("video_id", "title", "artists", "album", "album_id", "duration_seconds", "liked", "thumbnails", "is_explicit", "url")
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    ARTISTS_FIELD_NUMBER: _ClassVar[int]
    ALBUM_FIELD_NUMBER: _ClassVar[int]
    ALBUM_ID_FIELD_NUMBER: _ClassVar[int]
    DURATION_SECONDS_FIELD_NUMBER: _ClassVar[int]
    LIKED_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    IS_EXPLICIT_FIELD_NUMBER: _ClassVar[int]
    URL_FIELD_NUMBER: _ClassVar[int]
    video_id: str
    title: str
    artists: _containers.RepeatedCompositeFieldContainer[Artist]
    album: str
    album_id: str
    duration_seconds: int
    liked: bool
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    is_explicit: bool
    url: str
    def __init__(self, video_id: _Optional[str] = ..., title: _Optional[str] = ..., artists: _Optional[_Iterable[_Union[Artist, _Mapping]]] = ..., album: _Optional[str] = ..., album_id: _Optional[str] = ..., duration_seconds: _Optional[int] = ..., liked: _Optional[bool] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ..., is_explicit: _Optional[bool] = ..., url: _Optional[str] = ...) -> None: ...

class SearchResultAlbum(_message.Message):
    __slots__ = ("browse_id", "title", "artists", "year", "type", "thumbnails", "is_explicit")
    BROWSE_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    ARTISTS_FIELD_NUMBER: _ClassVar[int]
    YEAR_FIELD_NUMBER: _ClassVar[int]
    TYPE_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    IS_EXPLICIT_FIELD_NUMBER: _ClassVar[int]
    browse_id: str
    title: str
    artists: _containers.RepeatedCompositeFieldContainer[Artist]
    year: str
    type: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    is_explicit: bool
    def __init__(self, browse_id: _Optional[str] = ..., title: _Optional[str] = ..., artists: _Optional[_Iterable[_Union[Artist, _Mapping]]] = ..., year: _Optional[str] = ..., type: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ..., is_explicit: _Optional[bool] = ...) -> None: ...

class SearchResultArtist(_message.Message):
    __slots__ = ("browse_id", "name", "subscribers", "thumbnails")
    BROWSE_ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    SUBSCRIBERS_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    browse_id: str
    name: str
    subscribers: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    def __init__(self, browse_id: _Optional[str] = ..., name: _Optional[str] = ..., subscribers: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ...) -> None: ...

class SearchResultPlaylist(_message.Message):
    __slots__ = ("browse_id", "title", "author", "item_count", "thumbnails")
    BROWSE_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    ITEM_COUNT_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    browse_id: str
    title: str
    author: str
    item_count: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    def __init__(self, browse_id: _Optional[str] = ..., title: _Optional[str] = ..., author: _Optional[str] = ..., item_count: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ...) -> None: ...

class SearchResultPodcast(_message.Message):
    __slots__ = ("browse_id", "title", "author", "thumbnails")
    BROWSE_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    AUTHOR_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    browse_id: str
    title: str
    author: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    def __init__(self, browse_id: _Optional[str] = ..., title: _Optional[str] = ..., author: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ...) -> None: ...

class SearchResultEpisode(_message.Message):
    __slots__ = ("video_id", "title", "podcast_name", "podcast_id", "date", "thumbnails")
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    TITLE_FIELD_NUMBER: _ClassVar[int]
    PODCAST_NAME_FIELD_NUMBER: _ClassVar[int]
    PODCAST_ID_FIELD_NUMBER: _ClassVar[int]
    DATE_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    video_id: str
    title: str
    podcast_name: str
    podcast_id: str
    date: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    def __init__(self, video_id: _Optional[str] = ..., title: _Optional[str] = ..., podcast_name: _Optional[str] = ..., podcast_id: _Optional[str] = ..., date: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ...) -> None: ...

class GetSearchResultsResponse(_message.Message):
    __slots__ = ("songs", "albums", "artists", "playlists", "podcasts", "episodes")
    SONGS_FIELD_NUMBER: _ClassVar[int]
    ALBUMS_FIELD_NUMBER: _ClassVar[int]
    ARTISTS_FIELD_NUMBER: _ClassVar[int]
    PLAYLISTS_FIELD_NUMBER: _ClassVar[int]
    PODCASTS_FIELD_NUMBER: _ClassVar[int]
    EPISODES_FIELD_NUMBER: _ClassVar[int]
    songs: _containers.RepeatedCompositeFieldContainer[SearchResultSong]
    albums: _containers.RepeatedCompositeFieldContainer[SearchResultAlbum]
    artists: _containers.RepeatedCompositeFieldContainer[SearchResultArtist]
    playlists: _containers.RepeatedCompositeFieldContainer[SearchResultPlaylist]
    podcasts: _containers.RepeatedCompositeFieldContainer[SearchResultPodcast]
    episodes: _containers.RepeatedCompositeFieldContainer[SearchResultEpisode]
    def __init__(self, songs: _Optional[_Iterable[_Union[SearchResultSong, _Mapping]]] = ..., albums: _Optional[_Iterable[_Union[SearchResultAlbum, _Mapping]]] = ..., artists: _Optional[_Iterable[_Union[SearchResultArtist, _Mapping]]] = ..., playlists: _Optional[_Iterable[_Union[SearchResultPlaylist, _Mapping]]] = ..., podcasts: _Optional[_Iterable[_Union[SearchResultPodcast, _Mapping]]] = ..., episodes: _Optional[_Iterable[_Union[SearchResultEpisode, _Mapping]]] = ...) -> None: ...

class GetArtistTopTracksRequest(_message.Message):
    __slots__ = ("channel_id",)
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    channel_id: str
    def __init__(self, channel_id: _Optional[str] = ...) -> None: ...

class GetArtistTopTracksResponse(_message.Message):
    __slots__ = ("name", "thumbnails", "subscribers", "tracks")
    NAME_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    SUBSCRIBERS_FIELD_NUMBER: _ClassVar[int]
    TRACKS_FIELD_NUMBER: _ClassVar[int]
    name: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    subscribers: str
    tracks: _containers.RepeatedCompositeFieldContainer[Song]
    def __init__(self, name: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ..., subscribers: _Optional[str] = ..., tracks: _Optional[_Iterable[_Union[Song, _Mapping]]] = ...) -> None: ...

class GetFollowedArtistsRequest(_message.Message):
    __slots__ = ("limit",)
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    limit: int
    def __init__(self, limit: _Optional[int] = ...) -> None: ...

class FollowedArtist(_message.Message):
    __slots__ = ("channel_id", "name", "subscribers", "thumbnails")
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    SUBSCRIBERS_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    channel_id: str
    name: str
    subscribers: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    def __init__(self, channel_id: _Optional[str] = ..., name: _Optional[str] = ..., subscribers: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ...) -> None: ...

class GetFollowedArtistsResponse(_message.Message):
    __slots__ = ("artists", "total")
    ARTISTS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    artists: _containers.RepeatedCompositeFieldContainer[FollowedArtist]
    total: int
    def __init__(self, artists: _Optional[_Iterable[_Union[FollowedArtist, _Mapping]]] = ..., total: _Optional[int] = ...) -> None: ...

class GetUserProfileRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class GetUserProfileResponse(_message.Message):
    __slots__ = ("channel_id", "name", "thumbnails")
    CHANNEL_ID_FIELD_NUMBER: _ClassVar[int]
    NAME_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    channel_id: str
    name: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    def __init__(self, channel_id: _Optional[str] = ..., name: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ...) -> None: ...

class GetUserTopItemsRequest(_message.Message):
    __slots__ = ("limit",)
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    limit: int
    def __init__(self, limit: _Optional[int] = ...) -> None: ...

class GetUserTopItemsResponse(_message.Message):
    __slots__ = ("tracks", "total")
    TRACKS_FIELD_NUMBER: _ClassVar[int]
    TOTAL_FIELD_NUMBER: _ClassVar[int]
    tracks: _containers.RepeatedCompositeFieldContainer[Song]
    total: int
    def __init__(self, tracks: _Optional[_Iterable[_Union[Song, _Mapping]]] = ..., total: _Optional[int] = ...) -> None: ...

class CheckUserSavedTrackRequest(_message.Message):
    __slots__ = ("video_id",)
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    video_id: str
    def __init__(self, video_id: _Optional[str] = ...) -> None: ...

class CheckUserSavedTrackResponse(_message.Message):
    __slots__ = ("is_saved",)
    IS_SAVED_FIELD_NUMBER: _ClassVar[int]
    is_saved: bool
    def __init__(self, is_saved: _Optional[bool] = ...) -> None: ...

class SaveRemoveTrackRequest(_message.Message):
    __slots__ = ("video_ids", "is_remove")
    VIDEO_IDS_FIELD_NUMBER: _ClassVar[int]
    IS_REMOVE_FIELD_NUMBER: _ClassVar[int]
    video_ids: _containers.RepeatedScalarFieldContainer[str]
    is_remove: bool
    def __init__(self, video_ids: _Optional[_Iterable[str]] = ..., is_remove: _Optional[bool] = ...) -> None: ...

class SaveRemoveTrackResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class LikeSongRequest(_message.Message):
    __slots__ = ("video_id",)
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    video_id: str
    def __init__(self, video_id: _Optional[str] = ...) -> None: ...

class LikeSongResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class UnlikeSongRequest(_message.Message):
    __slots__ = ("video_id",)
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    video_id: str
    def __init__(self, video_id: _Optional[str] = ...) -> None: ...

class UnlikeSongResponse(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class SearchSongsRequest(_message.Message):
    __slots__ = ("query", "limit")
    QUERY_FIELD_NUMBER: _ClassVar[int]
    LIMIT_FIELD_NUMBER: _ClassVar[int]
    query: str
    limit: int
    def __init__(self, query: _Optional[str] = ..., limit: _Optional[int] = ...) -> None: ...

class SearchSongsResponse(_message.Message):
    __slots__ = ("songs",)
    SONGS_FIELD_NUMBER: _ClassVar[int]
    songs: _containers.RepeatedCompositeFieldContainer[Song]
    def __init__(self, songs: _Optional[_Iterable[_Union[Song, _Mapping]]] = ...) -> None: ...

class GetHomePageRequest(_message.Message):
    __slots__ = ()
    def __init__(self) -> None: ...

class GetHomePageResponse(_message.Message):
    __slots__ = ("sections",)
    SECTIONS_FIELD_NUMBER: _ClassVar[int]
    sections: _containers.RepeatedCompositeFieldContainer[HomePageSection]
    def __init__(self, sections: _Optional[_Iterable[_Union[HomePageSection, _Mapping]]] = ...) -> None: ...

class HomePageSection(_message.Message):
    __slots__ = ("title", "contents")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    CONTENTS_FIELD_NUMBER: _ClassVar[int]
    title: str
    contents: _containers.RepeatedCompositeFieldContainer[HomePageContent]
    def __init__(self, title: _Optional[str] = ..., contents: _Optional[_Iterable[_Union[HomePageContent, _Mapping]]] = ...) -> None: ...

class HomePageContent(_message.Message):
    __slots__ = ("title", "playlist_id", "thumbnails", "description", "video_id", "browse_id", "content_type", "artists", "duration_seconds")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    PLAYLIST_ID_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    BROWSE_ID_FIELD_NUMBER: _ClassVar[int]
    CONTENT_TYPE_FIELD_NUMBER: _ClassVar[int]
    ARTISTS_FIELD_NUMBER: _ClassVar[int]
    DURATION_SECONDS_FIELD_NUMBER: _ClassVar[int]
    title: str
    playlist_id: str
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    description: str
    video_id: str
    browse_id: str
    content_type: str
    artists: _containers.RepeatedCompositeFieldContainer[Artist]
    duration_seconds: int
    def __init__(self, title: _Optional[str] = ..., playlist_id: _Optional[str] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ..., description: _Optional[str] = ..., video_id: _Optional[str] = ..., browse_id: _Optional[str] = ..., content_type: _Optional[str] = ..., artists: _Optional[_Iterable[_Union[Artist, _Mapping]]] = ..., duration_seconds: _Optional[int] = ...) -> None: ...

class GetVideoStreamURLAndDurationRequest(_message.Message):
    __slots__ = ("video_id",)
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    video_id: str
    def __init__(self, video_id: _Optional[str] = ...) -> None: ...

class GetVideoStreamURLAndDurationResponse(_message.Message):
    __slots__ = ("url", "duration")
    URL_FIELD_NUMBER: _ClassVar[int]
    DURATION_FIELD_NUMBER: _ClassVar[int]
    url: str
    duration: str
    def __init__(self, url: _Optional[str] = ..., duration: _Optional[str] = ...) -> None: ...

class GetSongRelatedRequest(_message.Message):
    __slots__ = ("video_id", "browse_id")
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    BROWSE_ID_FIELD_NUMBER: _ClassVar[int]
    video_id: str
    browse_id: str
    def __init__(self, video_id: _Optional[str] = ..., browse_id: _Optional[str] = ...) -> None: ...

class GetSongRelatedResponse(_message.Message):
    __slots__ = ("sections",)
    SECTIONS_FIELD_NUMBER: _ClassVar[int]
    sections: _containers.RepeatedCompositeFieldContainer[SongRelatedSection]
    def __init__(self, sections: _Optional[_Iterable[_Union[SongRelatedSection, _Mapping]]] = ...) -> None: ...

class SongRelatedSection(_message.Message):
    __slots__ = ("title", "contents", "text_content")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    CONTENTS_FIELD_NUMBER: _ClassVar[int]
    TEXT_CONTENT_FIELD_NUMBER: _ClassVar[int]
    title: str
    contents: _containers.RepeatedCompositeFieldContainer[SongRelatedContent]
    text_content: str
    def __init__(self, title: _Optional[str] = ..., contents: _Optional[_Iterable[_Union[SongRelatedContent, _Mapping]]] = ..., text_content: _Optional[str] = ...) -> None: ...

class SongRelatedContent(_message.Message):
    __slots__ = ("title", "video_id", "playlist_id", "browse_id", "artists", "thumbnails", "is_explicit", "album", "album_id", "description", "subscribers", "year", "content_type")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    PLAYLIST_ID_FIELD_NUMBER: _ClassVar[int]
    BROWSE_ID_FIELD_NUMBER: _ClassVar[int]
    ARTISTS_FIELD_NUMBER: _ClassVar[int]
    THUMBNAILS_FIELD_NUMBER: _ClassVar[int]
    IS_EXPLICIT_FIELD_NUMBER: _ClassVar[int]
    ALBUM_FIELD_NUMBER: _ClassVar[int]
    ALBUM_ID_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    SUBSCRIBERS_FIELD_NUMBER: _ClassVar[int]
    YEAR_FIELD_NUMBER: _ClassVar[int]
    CONTENT_TYPE_FIELD_NUMBER: _ClassVar[int]
    title: str
    video_id: str
    playlist_id: str
    browse_id: str
    artists: _containers.RepeatedCompositeFieldContainer[Artist]
    thumbnails: _containers.RepeatedCompositeFieldContainer[Thumbnail]
    is_explicit: bool
    album: str
    album_id: str
    description: str
    subscribers: str
    year: str
    content_type: str
    def __init__(self, title: _Optional[str] = ..., video_id: _Optional[str] = ..., playlist_id: _Optional[str] = ..., browse_id: _Optional[str] = ..., artists: _Optional[_Iterable[_Union[Artist, _Mapping]]] = ..., thumbnails: _Optional[_Iterable[_Union[Thumbnail, _Mapping]]] = ..., is_explicit: _Optional[bool] = ..., album: _Optional[str] = ..., album_id: _Optional[str] = ..., description: _Optional[str] = ..., subscribers: _Optional[str] = ..., year: _Optional[str] = ..., content_type: _Optional[str] = ...) -> None: ...

class GetLyricsRequest(_message.Message):
    __slots__ = ("video_id", "timestamps")
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    TIMESTAMPS_FIELD_NUMBER: _ClassVar[int]
    video_id: str
    timestamps: bool
    def __init__(self, video_id: _Optional[str] = ..., timestamps: _Optional[bool] = ...) -> None: ...

class LyricLine(_message.Message):
    __slots__ = ("text", "start_time", "end_time", "id")
    TEXT_FIELD_NUMBER: _ClassVar[int]
    START_TIME_FIELD_NUMBER: _ClassVar[int]
    END_TIME_FIELD_NUMBER: _ClassVar[int]
    ID_FIELD_NUMBER: _ClassVar[int]
    text: str
    start_time: int
    end_time: int
    id: int
    def __init__(self, text: _Optional[str] = ..., start_time: _Optional[int] = ..., end_time: _Optional[int] = ..., id: _Optional[int] = ...) -> None: ...

class GetLyricsResponse(_message.Message):
    __slots__ = ("lyrics", "source", "has_timestamps", "lines")
    LYRICS_FIELD_NUMBER: _ClassVar[int]
    SOURCE_FIELD_NUMBER: _ClassVar[int]
    HAS_TIMESTAMPS_FIELD_NUMBER: _ClassVar[int]
    LINES_FIELD_NUMBER: _ClassVar[int]
    lyrics: str
    source: str
    has_timestamps: bool
    lines: _containers.RepeatedCompositeFieldContainer[LyricLine]
    def __init__(self, lyrics: _Optional[str] = ..., source: _Optional[str] = ..., has_timestamps: _Optional[bool] = ..., lines: _Optional[_Iterable[_Union[LyricLine, _Mapping]]] = ...) -> None: ...

class CreatePlaylistRequest(_message.Message):
    __slots__ = ("title", "description", "privacy_status", "video_ids", "source_playlist")
    TITLE_FIELD_NUMBER: _ClassVar[int]
    DESCRIPTION_FIELD_NUMBER: _ClassVar[int]
    PRIVACY_STATUS_FIELD_NUMBER: _ClassVar[int]
    VIDEO_IDS_FIELD_NUMBER: _ClassVar[int]
    SOURCE_PLAYLIST_FIELD_NUMBER: _ClassVar[int]
    title: str
    description: str
    privacy_status: str
    video_ids: _containers.RepeatedScalarFieldContainer[str]
    source_playlist: str
    def __init__(self, title: _Optional[str] = ..., description: _Optional[str] = ..., privacy_status: _Optional[str] = ..., video_ids: _Optional[_Iterable[str]] = ..., source_playlist: _Optional[str] = ...) -> None: ...

class CreatePlaylistResponse(_message.Message):
    __slots__ = ("playlist_id", "success", "error")
    PLAYLIST_ID_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    playlist_id: str
    success: bool
    error: str
    def __init__(self, playlist_id: _Optional[str] = ..., success: _Optional[bool] = ..., error: _Optional[str] = ...) -> None: ...

class AddPlaylistItemsRequest(_message.Message):
    __slots__ = ("playlist_id", "video_ids", "source_playlist", "duplicates")
    PLAYLIST_ID_FIELD_NUMBER: _ClassVar[int]
    VIDEO_IDS_FIELD_NUMBER: _ClassVar[int]
    SOURCE_PLAYLIST_FIELD_NUMBER: _ClassVar[int]
    DUPLICATES_FIELD_NUMBER: _ClassVar[int]
    playlist_id: str
    video_ids: _containers.RepeatedScalarFieldContainer[str]
    source_playlist: str
    duplicates: bool
    def __init__(self, playlist_id: _Optional[str] = ..., video_ids: _Optional[_Iterable[str]] = ..., source_playlist: _Optional[str] = ..., duplicates: _Optional[bool] = ...) -> None: ...

class AddPlaylistItemsResponse(_message.Message):
    __slots__ = ("status", "success", "error")
    STATUS_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    status: str
    success: bool
    error: str
    def __init__(self, status: _Optional[str] = ..., success: _Optional[bool] = ..., error: _Optional[str] = ...) -> None: ...

class PlaylistItemRef(_message.Message):
    __slots__ = ("video_id", "set_video_id")
    VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    SET_VIDEO_ID_FIELD_NUMBER: _ClassVar[int]
    video_id: str
    set_video_id: str
    def __init__(self, video_id: _Optional[str] = ..., set_video_id: _Optional[str] = ...) -> None: ...

class RemovePlaylistItemsRequest(_message.Message):
    __slots__ = ("playlist_id", "videos")
    PLAYLIST_ID_FIELD_NUMBER: _ClassVar[int]
    VIDEOS_FIELD_NUMBER: _ClassVar[int]
    playlist_id: str
    videos: _containers.RepeatedCompositeFieldContainer[PlaylistItemRef]
    def __init__(self, playlist_id: _Optional[str] = ..., videos: _Optional[_Iterable[_Union[PlaylistItemRef, _Mapping]]] = ...) -> None: ...

class RemovePlaylistItemsResponse(_message.Message):
    __slots__ = ("status", "success", "error")
    STATUS_FIELD_NUMBER: _ClassVar[int]
    SUCCESS_FIELD_NUMBER: _ClassVar[int]
    ERROR_FIELD_NUMBER: _ClassVar[int]
    status: str
    success: bool
    error: str
    def __init__(self, status: _Optional[str] = ..., success: _Optional[bool] = ..., error: _Optional[str] = ...) -> None: ...
