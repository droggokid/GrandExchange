module View.ItemCard exposing (view)

import Html exposing (Html, div, h3, img, p, text)
import Html.Attributes exposing (src, style)
import Models.Item exposing (Item)


view : Item -> Html msg
view item =
    div 
        [ style "background" "white"
        , style "border" "1px solid #e0e0e0"
        , style "border-radius" "12px"
        , style "padding" "20px"
        , style "transition" "all 0.2s"
        , style "box-shadow" "0 2px 4px rgba(0,0,0,0.05)"
        , style "display" "flex"
        , style "flex-direction" "column"
        , style "gap" "12px"
        ]
        [ div 
            [ style "display" "flex"
            , style "align-items" "center"
            , style "gap" "12px"
            ]
            [ img 
                [ src item.icon
                , style "width" "36px"
                , style "height" "36px"
                ] 
                []
            , h3 
                [ style "margin" "0"
                , style "color" "#2c3e50"
                , style "font-size" "18px"
                ] 
                [ text item.name ]
            ]
        , p 
            [ style "color" "#7f8c8d"
            , style "font-size" "14px"
            , style "margin" "0"
            , style "line-height" "1.5"
            ] 
            [ text item.description ]
        , viewPriceInfo item
        ]


viewPriceInfo : Item -> Html msg
viewPriceInfo item =
    div 
        [ style "display" "flex"
        , style "justify-content" "space-between"
        , style "align-items" "center"
        , style "padding-top" "12px"
        , style "border-top" "1px solid #ecf0f1"
        ]
        [ div []
            [ div 
                [ style "font-size" "12px"
                , style "color" "#95a5a6"
                , style "margin-bottom" "4px"
                ] 
                [ text "Current Price" ]
            , div 
                [ style "font-size" "20px"
                , style "font-weight" "700"
                , style "color" "#27ae60"
                ] 
                [ text item.current.price ]
            ]
        , viewTrendBadge item.current.trend
        ]


viewTrendBadge : String -> Html msg
viewTrendBadge trend =
    let
        (bgColor, textColor) =
            if trend == "positive" then
                ("#d5f4e6", "#27ae60")
            else if trend == "negative" then
                ("#fee", "#c33")
            else
                ("#f0f0f0", "#7f8c8d")
    in
    div 
        [ style "font-size" "12px"
        , style "padding" "4px 8px"
        , style "background" bgColor
        , style "color" textColor
        , style "border-radius" "4px"
        , style "font-weight" "600"
        ] 
        [ text ("↗ " ++ trend) ]
