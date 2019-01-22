module Main exposing (main)

import Browser
import Css exposing (..)
import Html.Styled exposing (..)
import Html.Styled.Attributes as A
import Html.Styled.Events as E
import Http
import Json.Decode as D


main : Program () Model Msg
main =
    -- TODO: Switch to Browser.application
    Browser.document
        { init = init
        , view =
            \model ->
                { title = "Rol-o-decks"
                , body = [ (view >> toUnstyled) model ]
                }
        , update = update
        , subscriptions = \_ -> Sub.none
        }


type alias Model =
    { searchTerm : String
    , searchResults : List SearchResult
    }


type alias SearchResult =
    { slideId : String
    , path : String
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
    , searchResults = []
    }


init : () -> ( Model, Cmd Msg )
init _ =
    ( emptyModel
    , Cmd.none
    )


type Msg
    = NoOp
    | UpdatedSearchTerm String
    | FetchedSearchResults (List SearchResult)
    | OpenSlide String
    | OpenedSlide (Result Http.Error ())


update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NoOp ->
            ( model, Cmd.none )

        UpdatedSearchTerm term ->
            ( { model | searchTerm = term }, getSearchResults term )

        FetchedSearchResults results ->
            ( { model | searchResults = results }, Cmd.none )

        OpenSlide slideId ->
            ( model, openSlide slideId )

        OpenedSlide _ ->
            ( model, Cmd.none )


view : Model -> Html Msg
view model =
    div [ A.css [ fontFamilies [ "Arial" ] ] ]
        [ searchBoxView model ]


searchBoxView : Model -> Html Msg
searchBoxView model =
    div
        [ A.css
            [ displayFlex
            , flexDirection column
            , maxWidth (rem 60)
            , marginLeft auto
            , marginRight auto
            ]
        ]
        [ input
            [ A.placeholder "Search for presentation..."
            , A.type_ "text"
            , A.css
                [ border3 (px 1) solid (hex "ddd")
                , greyShadow
                , fontSize (rem 1.8)
                , padding (rem 0.5)
                , marginBottom (rem 1)
                , outline none
                ]
            , E.onInput UpdatedSearchTerm
            ]
            []
        , div [] (List.map searchResultView model.searchResults)
        ]


searchResultView : SearchResult -> Html Msg
searchResultView searchResult =
    div [ searchResultStyle ]
        [ div [ thumbnailStyle ]
            [ img
                [ thumbnailImgStyle
                , A.height 200
                , A.src (base64dataImage searchResult.thumbnailBase64)
                , E.onClick (OpenSlide searchResult.slideId)
                ]
                []
            ]
        , div [ matchContentStyle, E.onClick (OpenSlide searchResult.slideId) ]
            [ div [ matchDetailsStyle ]
                [ p [ presPathStyle ] [ text searchResult.path ]
                , p [ slideNumStyle ] [ text ("Slide " ++ String.fromInt searchResult.slide) ]
                ]
            , formattedMatch searchResult.match
            ]
        ]


searchResultStyle : Html.Styled.Attribute msg
searchResultStyle =
    A.css [ color (hex "333"), displayFlex, border3 (px 1) solid (hex "ddd"), greyShadow, marginBottom (rem 0.5), cursor pointer ]


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
    A.css [ paddingLeft (rem 0.4), margin (px 0), borderBottom3 (px 1) solid (hex "ddd"), color (hex "06b") ]


presPathStyle : Html.Styled.Attribute msg
presPathStyle =
    A.css [ marginTop (rem 0.4), marginBottom (rem 0.2) ]


slideNumStyle : Html.Styled.Attribute msg
slideNumStyle =
    A.css [ marginTop (rem 0.2), fontSize (rem 0.8) ]


greyShadow : Css.Style
greyShadow =
    boxShadow5 (px 2) (px 2) (px 2) (px 0) (rgba 0 0 0 0.1)


formattedMatch : SearchMatch -> Html Msg
formattedMatch match =
    let
        end =
            match.start + match.length

        beforeText =
            String.slice 0 match.start match.text

        matchText =
            String.slice match.start end match.text

        afterText =
            String.slice end (String.length match.text) match.text
    in
    div [ formattedMatchStyles ]
        [ text beforeText
        , b [] [ text matchText ]
        , text afterText
        ]


formattedMatchStyles : Html.Styled.Attribute msg
formattedMatchStyles =
    A.css [ padding (rem 0.5) ]


base64dataImage : String -> String
base64dataImage imageData =
    "data:image/png;base64," ++ imageData


getSearchResults : String -> Cmd Msg
getSearchResults term =
    Http.get
        { url = "/api/search?q=" ++ term
        , expect = Http.expectJson extractSearchResults searchResponseDecoder
        }


extractSearchResults : Result Http.Error (List SearchResult) -> Msg
extractSearchResults result =
    case result of
        Ok results ->
            FetchedSearchResults results

        Err e ->
            -- TODO: Handle errors properly
            NoOp


searchResponseDecoder : D.Decoder (List SearchResult)
searchResponseDecoder =
    D.field "results"
        (D.list searchResultDecoder)


searchResultDecoder : D.Decoder SearchResult
searchResultDecoder =
    D.map5 SearchResult
        (D.field "slideId" D.string)
        (D.field "path" D.string)
        (D.field "slide" D.int)
        (D.field "thumbnail" D.string)
        (D.field "match" searchMatchDecoder)


searchMatchDecoder : D.Decoder SearchMatch
searchMatchDecoder =
    D.map3 SearchMatch
        (D.field "text" D.string)
        (D.field "start" D.int)
        (D.field "length" D.int)


openSlide : String -> Cmd Msg
openSlide slideId =
    Http.get
        { url = "/api/open/" ++ slideId
        , expect = Http.expectWhatever OpenedSlide
        }
