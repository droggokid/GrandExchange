module Models.Item exposing (Item, PriceBox, SearchResponse, itemDecoder, searchResponseDecoder)

import Json.Decode as Decode exposing (Decoder, int, string, field, float)
import Json.Decode.Pipeline exposing (required, hardcoded)

-- DATA TYPES

type alias SearchResponse =
    { total : Int 
    , items : List Item
    }

type alias Item =
    { icon : String
    , iconLarge : String
    , id : Int
    , itemType : String
    , typeIcon : String
    , name : String
    , description : String
    , current : PriceBox
    , today : PriceBox
    , members : String
    }

type alias PriceBox = 
    { trend : String
    , price : String
    }

-- JSON DECODERS

-- Decode price that can be either a string or a number
priceDecoder : Decoder String
priceDecoder =
    Decode.oneOf
        [ string
        , Decode.map String.fromInt int
        , Decode.map String.fromFloat Decode.float
        ]


priceBoxDecoder : Decoder PriceBox
priceBoxDecoder =
    Decode.succeed PriceBox
         |> required "trend" string
         |> required "price" priceDecoder
 
 
itemDecoder : Decoder Item
itemDecoder =
    Decode.map8
        (\icon iconLarge id itemType typeIcon name description current ->
            \today members ->
                { icon = icon
                , iconLarge = iconLarge
                , id = id
                , itemType = itemType
                , typeIcon = typeIcon
                , name = name
                , description = description
                , current = current
                , today = today
                , members = members
                }
        )
        (field "icon" string)
        (field "icon_large" string)
        (field "id" int)
        (field "type" string)
        (field "typeIcon" string)
        (field "name" string)
        (field "description" string)
        (field "current" priceBoxDecoder)
        |> Decode.andThen (\fn ->
            Decode.map2 fn
                (field "today" priceBoxDecoder)
                (field "members" string)
        )
 
 
searchResponseDecoder : Decoder SearchResponse
searchResponseDecoder =
    Decode.succeed SearchResponse
         |> required "total" int
         |> required "items" (Decode.list itemDecoder)
