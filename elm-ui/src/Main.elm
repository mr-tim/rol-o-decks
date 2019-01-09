module Main exposing (..)

import Browser
import Browser.Dom as Dom
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (..)
import Html.Keyed as Keyed
import Html.Lazy exposing (lazy, lazy2)
import Task

main : Program (Maybe Model) Model Msg
main =
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
    {}

emptyModel : Model
emptyModel =
    {}

init : Maybe Model -> ( Model, Cmd Msg )
init maybeModel =
  ( Maybe.withDefault emptyModel maybeModel
  , Cmd.none
  )

type Msg
    = NoOp

update : Msg -> Model -> ( Model, Cmd Msg )
update msg model =
    case msg of
        NoOp ->
            ( model, Cmd.none )

view : Model -> Html Msg
view model =
    div []
        [ text "Rol-o-decks" ]