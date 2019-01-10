module Main exposing (..)

import Browser
import Html exposing (..)

main : Program () Model Msg
main =
    -- TODO: Switch to Browser.application
    Browser.document
        { init = init
        , view = \model -> {
            title = "Rol-o-decks",
            body = [view model]
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
    , end : Int
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
    div []
        [ text "Rol-o-decks" ]