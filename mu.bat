for /L %%n in (1 4000 1) do (
    mu http://__RADARR__HOST__:__RADARR__PORT__/radarr/api/movie/%%n?apikey=__RADARR__API__KEY__
)

pause
