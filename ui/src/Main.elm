module Main exposing (main)

import Browser
import Html exposing (Html)
import Pages.ItemSearch


-- MAIN

main : Program () Pages.ItemSearch.Model Pages.ItemSearch.Msg
main =
    Browser.element
        { init = init
        , update = Pages.ItemSearch.update
        , subscriptions = subscriptions
        , view = Pages.ItemSearch.view
        }


init : () -> ( Pages.ItemSearch.Model, Cmd Pages.ItemSearch.Msg )
init _ =
    ( Pages.ItemSearch.init, Cmd.none )


subscriptions : Pages.ItemSearch.Model -> Sub Pages.ItemSearch.Msg
subscriptions _ =
    Sub.none
