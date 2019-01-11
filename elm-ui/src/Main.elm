module Main exposing (..)

import Browser
import Css exposing (..)
import Html
import Html.Styled exposing (..)
import Html.Styled.Attributes as A

main : Program () Model Msg
main =
    -- TODO: Switch to Browser.application
    Browser.document
        { init = init
        , view = \model -> {
            title = "Rol-o-decks",
            body = [(view >> toUnstyled) model]
        }
        , update = update
        , subscriptions = \_ -> Sub.none
        }

type alias Model =
    { searchTerm : String
    , searchResults : List SearchResult
    }

type alias SearchResult =
    { path : String
    , slide : Int
    , thumbnailBase64 : String
    , match : SearchMatch
    }

type alias SearchMatch =
    { text : String
    , start : Int
    , length : Int
    }

emptyModel : Model
emptyModel =
    { searchTerm = ""
    , searchResults = [ ]
    }

init : () -> ( Model, Cmd Msg )
init _ =
  ( emptyModel
  , Cmd.none
  )

type Msg
    = NoOp
    | UpdatedSearchTerm
    | FetchedSearchResults

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NoOp ->
            ( model, Cmd.none )
        UpdatedSearchTerm ->
            ( model, Cmd.none )
        FetchedSearchResults ->
            ( model, Cmd.none )

view : Model -> Html Msg
view model =
    div [ A.css [ fontFamilies ["Arial"] ] ]
        [ searchBoxView model ]

searchBoxView : Model -> Html Msg
searchBoxView model =
    div [ A.css [ displayFlex
              , flexDirection column
              , maxWidth (rem 60)
              , marginLeft auto
              , marginRight auto ] ]
        [ input [ A.placeholder "Search for presentation..."
                , A.type_ "text"
                , A.css
                    [ border3 (px 1) solid (hex "ddd")
                    , greyShadow
                    , fontSize (rem 1.8)
                    , padding (rem 0.5)
                    , marginBottom (rem 1)
                    , outline none
                    ]
                ]
                [ ]
        , div [] (List.map searchResultView model.searchResults)]

searchResultView : SearchResult -> Html Msg
searchResultView searchResult =
    div [ searchResultStyle ]
        [ div [ thumbnailStyle ]
            [ img [ thumbnailImgStyle, A.height 200, A.src (base64dataImage searchResult.thumbnailBase64) ] [] ] 
        , div [ matchContentStyle ]
            [
                div [ matchDetailsStyle ] 
                    [ p [ presPathStyle ] [ text searchResult.path ] 
                    , p [ slideNumStyle ] [ text ("Slide " ++ String.fromInt searchResult.slide)]
                    ]
                , formattedMatch searchResult.match
            ]            
        ]

searchResultStyle : Html.Styled.Attribute msg
searchResultStyle =
    A.css [ color (hex "333"), displayFlex, border3 (px 1) solid (hex "ddd"), greyShadow, marginBottom (rem 0.5), cursor pointer]

thumbnailStyle : Html.Styled.Attribute msg
thumbnailStyle =
    A.css [ margin (px 0), lineHeight (rem 0) ]

thumbnailImgStyle : Html.Styled.Attribute msg
thumbnailImgStyle =
    A.css [ margin (px 0), padding (px 0) ]

matchContentStyle : Html.Styled.Attribute msg
matchContentStyle =
    A.css [ displayFlex, flexGrow (num 1), flexDirection column ]

matchDetailsStyle : Html.Styled.Attribute msg
matchDetailsStyle =
    A.css [ paddingLeft (rem 0.4), margin (px 0), borderBottom3 (px 1) solid (hex "ddd"), color (hex "06b")]

presPathStyle : Html.Styled.Attribute msg
presPathStyle =
    A.css [ marginTop (rem 0.4), marginBottom (rem 0.2) ]

slideNumStyle : Html.Styled.Attribute msg
slideNumStyle =
    A.css [ marginTop (rem 0.2), fontSize (rem 0.8) ]

greyShadow : Css.Style
greyShadow = boxShadow5 (px 2) (px 2) (px 2) (px 0) (rgba 0 0 0 0.10)

formattedMatch : SearchMatch -> Html Msg
formattedMatch match =
    let
        end = match.start + match.length
        beforeText = String.slice 0 match.start match.text
        matchText = String.slice match.start end match.text
        afterText = String.slice end (String.length match.text) match.text
    in
        div [ formattedMatchStyles ] [ text beforeText
             , b [] [ text matchText]
             , text afterText
             ]

formattedMatchStyles : Html.Styled.Attribute msg
formattedMatchStyles =
    A.css [ padding (rem 0.5) ]

base64dataImage : String -> String
base64dataImage imageData =
    "data:image/png;base64," ++ imageData