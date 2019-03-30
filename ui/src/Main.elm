module Main exposing (main)

import Browser
import Css exposing (..)
import Html.Entity
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
    , settings : Settings
    , previousSettings : Settings
    , showSettings: Bool
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

type alias Settings =
    { searchPaths : List String
    }


emptyModel : Model
emptyModel =
    { searchTerm = ""
    , searchResults = []
    , settings = emptySettings
    , previousSettings = emptySettings
    , showSettings = False
    }

emptySettings : Settings
emptySettings =
    { searchPaths = []
    }


init : () -> ( Model, Cmd Msg )
init _ =
    ( emptyModel
    , fetchSettings
    )


type Msg
    = NoOp
    | UpdatedSearchTerm String
    | FetchedSearchResults (List SearchResult)
    | OpenSlide String
    | OpenedSlide (Result Http.Error ())
    | FetchedSettings Settings
    | ShowSettingsModal
    | HideSettingsModal
    | CancelSettings
    | SaveSettings
    | AddSearchPath
    | UpdateSearchPathAtIndex Int String
    | RemoveSearchPathAtIndex Int


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

        FetchedSettings settings ->
            ( { model | settings = settings }, Cmd.none )

        ShowSettingsModal ->
            let
                p = model.settings
            in
                ( { model |
                    showSettings = True,
                    previousSettings = p }, Cmd.none )

        HideSettingsModal ->
            ( { model | showSettings = False }, Cmd.none )

        CancelSettings ->
            let
                revertedSettings = { model | settings = model.previousSettings }
            in
                revertedSettings |> update HideSettingsModal

        SaveSettings ->
            model |> update HideSettingsModal

        AddSearchPath ->
            let
                s = model.settings
                updatedSearchPaths = List.append s.searchPaths [""]
                updatedSettings = { s | searchPaths = updatedSearchPaths }
            in
                ( { model | settings = updatedSettings } , Cmd.none )

        RemoveSearchPathAtIndex index ->
            let
                s = model.settings
                updatedSearchPaths = removeValueAtIndex index s.searchPaths
                updatedSettings = { s | searchPaths = updatedSearchPaths }
            in
                ( { model | settings = updatedSettings }, Cmd.none)

        UpdateSearchPathAtIndex index updatedValue ->
            let
                s = model.settings
                updatedSearchPaths = updateValueAtIndex index updatedValue s.searchPaths
                updatedSettings = { s | searchPaths = updatedSearchPaths }
            in
                ( { model | settings = updatedSettings }, Cmd.none )

removeValueAtIndex : Int -> List a -> List a
removeValueAtIndex index list =
    let
        remover itemIndex value =
            if itemIndex == index then
                []
            else
                [value]
    in
        List.concat (List.indexedMap remover list)

updateValueAtIndex : Int -> a -> List a -> List a
updateValueAtIndex index newValue values =
    let
        updater itemIndex item =
            if itemIndex == index then
                newValue
            else
                item
    in
        List.indexedMap updater values

view : Model -> Html Msg
view model =
    div [ A.css [ fontFamilies [ "Arial" ] ] ]
        (searchBoxView model :: settingsDialog model)


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
        [ div
            [ A.css
                [ displayFlex
                , flexDirection row
                , alignItems baseline
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
                    , flexGrow (num 1)
                    ]
                , E.onInput UpdatedSearchTerm
                ]
                [ ]
            , span
                [ A.css
                    [ flexGrow (num 0)
                    , fontSize (px 30)
                    , fontWeight bold
                    , width (px 20)
                    , textAlign right
                    , hover undecoratedPointer
                    , focus undecoratedPointer
                    ]
                , E.onClick ShowSettingsModal
                ]
                [ text Html.Entity.vellip ]
            ]
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


settingsDialog : Model -> List (Html Msg)
settingsDialog model =
    if model.showSettings then
        [div [ modalBackgroundStyle ]
            [ modalContent model ] ]
    else
        []

modalBackgroundStyle : Html.Styled.Attribute msg
modalBackgroundStyle =
    A.css [ position fixed
    , left (px 0)
    , top (px 0)
    , width (pct 100)
    , height (pct 100)
    , backgroundColor (rgba 0 0 0 0.4)]

modalContent : Model -> Html Msg
modalContent model =
    div [modalContentStyle]
        [ span [ closeStyle, E.onClick CancelSettings ] [ text Html.Entity.times ]
        , h3 [] [ text "Settings" ]
        , h4 [] [ text "Search Paths" ]
        , div [] (List.indexedMap searchPathControl model.settings.searchPaths)
        , addSearchPath
        , settingsDialogButtons model.settings
        ]

modalContentStyle : Html.Styled.Attribute msg
modalContentStyle =
    A.css [ backgroundColor (rgb 230 230 230)
        , color (rgb 40 40 40)
        , marginTop (rem 5)
        , marginLeft auto
        , marginRight auto
        , maxWidth (rem 48)
        , padding (rem 1)]

closeStyle : Html.Styled.Attribute msg
closeStyle =
    A.css [ float right
        , fontWeight bold
        , fontSize (px 24)
        , hover undecoratedPointer
        , focus undecoratedPointer
        ]

undecoratedPointer : List Style
undecoratedPointer =
    [ cursor pointer
    , textDecoration none
    ]

searchPathControl : Int -> String -> Html Msg
searchPathControl index path =
    div [ A.css [ displayFlex ] ]
        [ input
            [ A.css [ fontSize (px 20), flexGrow (num 1) ]
            , E.onInput (UpdateSearchPathAtIndex index)
            , A.width 100
            , A.value path ]
            [ text path ]
        , button
            [ E.onClick (RemoveSearchPathAtIndex index)
            , removeSearchPathStyle ]
            [ text Html.Entity.times ]
        ]

removeSearchPathStyle : Html.Styled.Attribute msg
removeSearchPathStyle =
    A.css
        [ fontSize (px 20) ]

addSearchPath : Html Msg
addSearchPath =
    button
        [ E.onClick AddSearchPath ]
        [ text "Add Search Path" ]

settingsDialogButtons : Settings -> Html Msg
settingsDialogButtons settings =
    div []
        [ button
            [ E.onClick CancelSettings ]
            [ text "Cancel" ]
        , button
            [ E.onClick SaveSettings ]
            [ text "Save" ]
        ]

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


fetchSettings : Cmd Msg
fetchSettings =
    Http.get
        { url = "/api/settings"
        , expect = Http.expectJson extractSettings settingsDecoder
        }


extractSettings : Result Http.Error Settings -> Msg
extractSettings result =
    case result of
        Ok settings ->
            FetchedSettings settings

        Err e ->
            -- TODO: Handle errors properly
            NoOp


settingsDecoder : D.Decoder Settings
settingsDecoder =
    D.field "indexPaths"
        (D.map Settings
            (D.list D.string))