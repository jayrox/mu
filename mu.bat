for /L %%n in (1 1 1) do (
    sm http://__RADARR__HOST__:__RADARR__PORT__/radarr/api/movie/%%n?apikey=__RADARR__API__KEY__
)

pause
