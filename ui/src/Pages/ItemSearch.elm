module Pages.ItemSearch exposing (Model, Msg, init, update, view)

import Html exposing (Html, button, div, h1, input, p, text)
import Html.Attributes exposing (placeholder, style, value)
import Html.Events exposing (onClick, onInput)
import Http
import Models.Item exposing (Item, SearchResponse)
import Api.Items
import View.ItemCard


-- MODEL


type alias Model =
    { searchText : String
    , searchState : SearchState
    }


type SearchState
    = NotSearched
    | Searching
    | Success SearchResponse
    | Failure Http.Error


init : Model
init =
    { searchText = ""
    , searchState = NotSearched
    }


-- UPDATE


type Msg
    = UpdateSearchText String
    | SearchItems
    | GotSearchResponse (Result Http.Error SearchResponse)


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        UpdateSearchText newText ->
            ( { model | searchText = newText }
            , Cmd.none
            )

        SearchItems ->
            ( { model | searchState = Searching }
            , Api.Items.searchItems model.searchText GotSearchResponse
            )

        GotSearchResponse result ->
            case result of
                Ok searchResponse ->
                    ( { model | searchState = Success searchResponse }
                    , Cmd.none
                    )

                Err error ->
                    ( { model | searchState = Failure error }
                    , Cmd.none
                    )


-- VIEW


view : Model -> Html Msg
view model =
    div 
        [ style "max-width" "900px"
        , style "margin" "0 auto"
        , style "padding" "40px 20px"
        , style "font-family" "system-ui, -apple-system, sans-serif"
        ]
        [ h1 
            [ style "color" "#2c3e50"
            , style "margin-bottom" "30px"
            , style "font-size" "2.5rem"
            ] 
            [ text "🗡️ RuneScape Item Search" ]
        , viewSearchBar model
        , viewSearchResults model.searchState
        ]


viewSearchBar : Model -> Html Msg
viewSearchBar model =
    div 
        [ style "display" "flex"
        , style "gap" "10px"
        , style "margin-bottom" "30px"
        ]
        [ input
            [ placeholder "Enter item name (e.g., 'rune', 'dragon')..."
            , value model.searchText
            , onInput UpdateSearchText
            , style "flex" "1"
            , style "padding" "12px 16px"
            , style "font-size" "16px"
            , style "border" "2px solid #e0e0e0"
            , style "border-radius" "8px"
            , style "outline" "none"
            , style "transition" "border-color 0.2s"
            ]
            []
        , button 
            [ onClick SearchItems
            , style "padding" "12px 30px"
            , style "font-size" "16px"
            , style "font-weight" "600"
            , style "background" "#3498db"
            , style "color" "white"
            , style "border" "none"
            , style "border-radius" "8px"
            , style "cursor" "pointer"
            , style "transition" "background 0.2s"
            ] 
            [ text "Search" ]
        ]


viewSearchResults : SearchState -> Html Msg
viewSearchResults searchState =
    case searchState of
        NotSearched ->
            div 
                [ style "text-align" "center"
                , style "padding" "60px 20px"
                , style "color" "#7f8c8d"
                ]
                [ text "Search for items to get started" ]

        Searching ->
            div 
                [ style "text-align" "center"
                , style "padding" "60px 20px"
                , style "color" "#3498db"
                , style "font-size" "18px"
                ]
                [ text "⚔️ Searching..." ]

        Success response ->
            div []
                [ p 
                    [ style "color" "#7f8c8d"
                    , style "margin-bottom" "20px"
                    , style "font-size" "14px"
                    ] 
                    [ text ("Found " ++ String.fromInt response.total ++ " items") ]
                , div 
                    [ style "display" "grid"
                    , style "grid-template-columns" "repeat(auto-fill, minmax(300px, 1fr))"
                    , style "gap" "20px"
                    ] 
                    (List.map viewItemCard response.items)
                ]

        Failure error ->
            div 
                [ style "background" "#fee"
                , style "color" "#c33"
                , style "padding" "20px"
                , style "border-radius" "8px"
                , style "border-left" "4px solid #c33"
                ]
                [ text ("Error: " ++ httpErrorToString error) ]


viewItemCard : Item -> Html Msg
viewItemCard item =
    Html.map (\_ -> UpdateSearchText "") (View.ItemCard.view item)


httpErrorToString : Http.Error -> String
httpErrorToString error =
    case error of
        Http.BadUrl url ->
            "Bad URL: " ++ url

        Http.Timeout ->
            "Request timeout"

        Http.NetworkError ->
            "Network error"

        Http.BadStatus status ->
            "Bad status: " ++ String.fromInt status

        Http.BadBody body ->
            "Bad body: " ++ body
