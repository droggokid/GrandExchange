module Api.Items exposing (searchItems)

import Http
import Config
import Models.Item exposing (SearchResponse, searchResponseDecoder)
import Url.Builder

-- API FUNCTIONS

searchItems : String -> (Result Http.Error SearchResponse -> msg) -> Cmd msg
searchItems itemName toMsg =
    let
        url =
            Url.Builder.crossOrigin Config.apiUrl
                [ "search-item", itemName ]
                []
    in
    Http.get
        { url = url
        , expect = Http.expectJson toMsg searchResponseDecoder
        }
